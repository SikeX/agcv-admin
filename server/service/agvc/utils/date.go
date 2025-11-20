package utils

import (
	"fmt"
	"time"
)

func FormatTimeToShanghai(t string) string {
	tTimeDate, _ := time.Parse(time.RFC3339, t)
	shanghaiLoc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		fmt.Printf("加载时区失败：%v\n", err)
		return ""
	}
	localTime := tTimeDate.In(shanghaiLoc)
	t = localTime.Format("2006-01-02 15:04:05")
	return t
}
