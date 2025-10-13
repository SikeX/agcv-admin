
package system

import (
	"context"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	systemReq "github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
)

type SysGridConnectionPointService struct {}
// CreateSysGridConnectionPoint 创建并网点配置记录
// Author [yourname](https://github.com/yourname)
func (sysGridConnectionPointService *SysGridConnectionPointService) CreateSysGridConnectionPoint(ctx context.Context, sysGridConnectionPoint *system.SysGridConnectionPoint) (err error) {
	err = global.GVA_DB.Create(sysGridConnectionPoint).Error
	return err
}

// DeleteSysGridConnectionPoint 删除并网点配置记录
// Author [yourname](https://github.com/yourname)
func (sysGridConnectionPointService *SysGridConnectionPointService)DeleteSysGridConnectionPoint(ctx context.Context, ID string) (err error) {
	err = global.GVA_DB.Delete(&system.SysGridConnectionPoint{},"id = ?",ID).Error
	return err
}

// DeleteSysGridConnectionPointByIds 批量删除并网点配置记录
// Author [yourname](https://github.com/yourname)
func (sysGridConnectionPointService *SysGridConnectionPointService)DeleteSysGridConnectionPointByIds(ctx context.Context, IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]system.SysGridConnectionPoint{},"id in ?",IDs).Error
	return err
}

// UpdateSysGridConnectionPoint 更新并网点配置记录
// Author [yourname](https://github.com/yourname)
func (sysGridConnectionPointService *SysGridConnectionPointService)UpdateSysGridConnectionPoint(ctx context.Context, sysGridConnectionPoint system.SysGridConnectionPoint) (err error) {
	err = global.GVA_DB.Model(&system.SysGridConnectionPoint{}).Where("id = ?",sysGridConnectionPoint.ID).Updates(&sysGridConnectionPoint).Error
	return err
}

// GetSysGridConnectionPoint 根据ID获取并网点配置记录
// Author [yourname](https://github.com/yourname)
func (sysGridConnectionPointService *SysGridConnectionPointService)GetSysGridConnectionPoint(ctx context.Context, ID string) (sysGridConnectionPoint system.SysGridConnectionPoint, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&sysGridConnectionPoint).Error
	return
}
// GetSysGridConnectionPointInfoList 分页获取并网点配置记录
// Author [yourname](https://github.com/yourname)
func (sysGridConnectionPointService *SysGridConnectionPointService)GetSysGridConnectionPointInfoList(ctx context.Context, info systemReq.SysGridConnectionPointSearch) (list []system.SysGridConnectionPoint, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
    // 创建db
	db := global.GVA_DB.Model(&system.SysGridConnectionPoint{})
    var sysGridConnectionPoints []system.SysGridConnectionPoint
    // 如果有条件搜索 下方会自动创建搜索语句
    if info.Name != "" {
     db = db.Where("name LIKE ?", "%"+info.Name+"%")
    }
    
	err = db.Count(&total).Error
	if err!=nil {
    	return
    }

	if limit != 0 {
       db = db.Limit(limit).Offset(offset)
    }

	err = db.Find(&sysGridConnectionPoints).Error
	return  sysGridConnectionPoints, total, err
}
func (sysGridConnectionPointService *SysGridConnectionPointService)GetSysGridConnectionPointPublic(ctx context.Context) {
    // 此方法为获取数据源定义的数据
    // 请自行实现
}
