
package system

import (
	"context"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
    systemReq "github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
)

type SysTestPointService struct {}
// CreateSysTestPoint 创建测试管理记录
// Author [yourname](https://github.com/yourname)
func (sysTestPointService *SysTestPointService) CreateSysTestPoint(ctx context.Context, sysTestPoint *system.SysTestPoint) (err error) {
	err = global.GVA_DB.Create(sysTestPoint).Error
	return err
}

// DeleteSysTestPoint 删除测试管理记录
// Author [yourname](https://github.com/yourname)
func (sysTestPointService *SysTestPointService)DeleteSysTestPoint(ctx context.Context, ID string) (err error) {
	err = global.GVA_DB.Delete(&system.SysTestPoint{},"id = ?",ID).Error
	return err
}

// DeleteSysTestPointByIds 批量删除测试管理记录
// Author [yourname](https://github.com/yourname)
func (sysTestPointService *SysTestPointService)DeleteSysTestPointByIds(ctx context.Context, IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]system.SysTestPoint{},"id in ?",IDs).Error
	return err
}

// UpdateSysTestPoint 更新测试管理记录
// Author [yourname](https://github.com/yourname)
func (sysTestPointService *SysTestPointService)UpdateSysTestPoint(ctx context.Context, sysTestPoint system.SysTestPoint) (err error) {
	err = global.GVA_DB.Model(&system.SysTestPoint{}).Where("id = ?",sysTestPoint.ID).Updates(&sysTestPoint).Error
	return err
}

// GetSysTestPoint 根据ID获取测试管理记录
// Author [yourname](https://github.com/yourname)
func (sysTestPointService *SysTestPointService)GetSysTestPoint(ctx context.Context, ID string) (sysTestPoint system.SysTestPoint, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&sysTestPoint).Error
	return
}
// GetSysTestPointInfoList 分页获取测试管理记录
// Author [yourname](https://github.com/yourname)
func (sysTestPointService *SysTestPointService)GetSysTestPointInfoList(ctx context.Context, info systemReq.SysTestPointSearch) (list []system.SysTestPoint, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
    // 创建db
	db := global.GVA_DB.Model(&system.SysTestPoint{})
    var sysTestPoints []system.SysTestPoint
    // 如果有条件搜索 下方会自动创建搜索语句
    if info.InStorageName != "" {
     db = db.Where("in_storage_name LIKE ?", "%"+info.InStorageName+"%")
    }
    if info.PointName != "" {
     db = db.Where("point_name LIKE ?", "%"+info.PointName+"%")
    }
    if info.PointType != nil {
     db = db.Where("point_type = ?", *info.PointType)
    }
    if info.DeviceType != nil {
     db = db.Where("device_type = ?", *info.DeviceType)
    }
    if len(info.CreatedAtRange) == 2 {
     db = db.Where("created_at BETWEEN ? AND ?", info.CreatedAtRange[0], info.CreatedAtRange[1])
    }
    
	err = db.Count(&total).Error
	if err!=nil {
    	return
    }

	if limit != 0 {
       db = db.Limit(limit).Offset(offset)
    }

	err = db.Find(&sysTestPoints).Error
	return  sysTestPoints, total, err
}
func (sysTestPointService *SysTestPointService)GetSysTestPointPublic(ctx context.Context) {
    // 此方法为获取数据源定义的数据
    // 请自行实现
}
