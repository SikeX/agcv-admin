package internal

import (
    "os"
    "path/filepath"
    "sync"
    "time"
)

// Cutter 实现 io.Writer 接口
// 用于日志切割, strings.Join([]string{director,layout, formats..., level+".log"}, os.PathSeparator)
type Cutter struct {
    level        string        // 日志级别(debug, info, warn, error, dpanic, panic, fatal)
    layout       string        // 时间格式 2006-01-02 15:04:05
    formats      []string      // 自定义参数([]string{Director,"2006-01-02", "business"(此参数可不写), level+".log"}
    director     string        // 日志文件夹
    retentionDay int           // 日志保留天数
    maxSize      int64         // 单个日志文件最大大小(字节)
    file         *os.File      // 文件句柄
    mutex        *sync.RWMutex // 读写锁
}

type CutterOption func(*Cutter)

// CutterWithLayout 时间格式
func CutterWithLayout(layout string) CutterOption {
    return func(c *Cutter) {
        c.layout = layout
    }
}

// CutterWithFormats 格式化参数
func CutterWithFormats(format ...string) CutterOption {
    return func(c *Cutter) {
        if len(format) > 0 {
            c.formats = format
        }
    }
}

// CutterWithMaxSize 设置最大文件大小(MB)
func CutterWithMaxSize(maxSizeMB int) CutterOption {
    return func(c *Cutter) {
        c.maxSize = int64(maxSizeMB) * 1024 * 1024
    }
}

func NewCutter(director string, level string, retentionDay int, options ...CutterOption) *Cutter {
    rotate := &Cutter{
        level:        level,
        director:     director,
        retentionDay: retentionDay,
        mutex:        new(sync.RWMutex),
    }
    for i := 0; i < len(options); i++ {
        options[i](rotate)
    }
    return rotate
}

// Write satisfies the io.Writer interface. It writes to the
// appropriate file handle that is currently being used.
// If we have reached rotation time, the target file gets
// automatically rotated, and also purged if necessary.
func (c *Cutter) Write(bytes []byte) (n int, err error) {
    c.mutex.Lock()
    defer func() {
        if c.file != nil {
            _ = c.file.Close()
            c.file = nil
        }
        c.mutex.Unlock()
    }()
    length := len(c.formats)
    values := make([]string, 0, 3+length)
    values = append(values, c.director)
    if c.layout != "" {
        values = append(values, time.Now().Format(c.layout))
    }
    for i := 0; i < length; i++ {
        values = append(values, c.formats[i])
    }
    values = append(values, c.level+".log")
    filename := filepath.Join(values...)
    
    // 如果设置了maxSize，检查文件大小并进行切割
    if c.maxSize > 0 {
        filename = c.getRotatedFilename(filename)
    }
    
    director := filepath.Dir(filename)
    err = os.MkdirAll(director, os.ModePerm)
    if err != nil {
        return 0, err
    }
    err = removeNDaysFolders(c.director, c.retentionDay)
    if err != nil {
        return 0, err
    }
    c.file, err = os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
    if err != nil {
        return 0, err
    }
    return c.file.Write(bytes)
}

func (c *Cutter) Sync() error {
    c.mutex.Lock()
    defer c.mutex.Unlock()

    if c.file != nil {
        return c.file.Sync()
    }
    return nil
}

// getRotatedFilename 根据文件大小获取切割后的文件名
func (c *Cutter) getRotatedFilename(baseFilename string) string {
    // 检查文件是否存在以及大小
    info, err := os.Stat(baseFilename)
    if err != nil || info.Size() < c.maxSize {
        // 文件不存在或未达到最大大小，使用原文件名
        return baseFilename
    }
    
    // 文件达到最大大小，删除旧文件并重新创建
    // 这样可以确保只有一个日志文件，超过大小后覆盖旧内容
    os.Remove(baseFilename)
    return baseFilename
}

// 增加日志目录文件清理 小于等于零的值默认忽略不再处理
func removeNDaysFolders(dir string, days int) error {
    if days <= 0 {
        return nil
    }
    cutoff := time.Now().AddDate(0, 0, -days)
    return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if info.IsDir() && info.ModTime().Before(cutoff) && path != dir {
            err = os.RemoveAll(path)
            if err != nil {
                return err
            }
        }
        return nil
    })
}
