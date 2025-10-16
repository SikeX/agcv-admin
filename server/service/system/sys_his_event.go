package system

import (
	"context"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	systemReq "github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
)

type SysHisEventService struct{}

// CreateSysHisEvent 创建历史事件记录
// Author [yourname](https://github.com/yourname)
func (sysHisEventService *SysHisEventService) CreateSysHisEvent(ctx context.Context, sysHisEvent *system.SysHisEvent) (err error) {
	err = global.GVA_DB.Create(sysHisEvent).Error
	return err
}

// DeleteSysHisEvent 删除历史事件记录
// Author [yourname](https://github.com/yourname)
func (sysHisEventService *SysHisEventService) DeleteSysHisEvent(ctx context.Context, ID string) (err error) {
	err = global.GVA_DB.Delete(&system.SysHisEvent{}, "id = ?", ID).Error
	return err
}

// DeleteSysHisEventByIds 批量删除历史事件记录
// Author [yourname](https://github.com/yourname)
func (sysHisEventService *SysHisEventService) DeleteSysHisEventByIds(ctx context.Context, IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]system.SysHisEvent{}, "id in ?", IDs).Error
	return err
}

// UpdateSysHisEvent 更新历史事件记录
// Author [yourname](https://github.com/yourname)
func (sysHisEventService *SysHisEventService) UpdateSysHisEvent(ctx context.Context, sysHisEvent system.SysHisEvent) (err error) {
	err = global.GVA_DB.Model(&system.SysHisEvent{}).Where("id = ?", sysHisEvent.ID).Updates(&sysHisEvent).Error
	return err
}

// GetSysHisEvent 根据ID获取历史事件记录
// Author [yourname](https://github.com/yourname)
func (sysHisEventService *SysHisEventService) GetSysHisEvent(ctx context.Context, ID string) (sysHisEvent system.SysHisEvent, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&sysHisEvent).Error
	return
}

// GetSysHisEventInfoList 分页获取历史事件记录
// Author [yourname](https://github.com/yourname)
func (sysHisEventService *SysHisEventService) GetSysHisEventInfoList(ctx context.Context, info systemReq.SysHisEventSearch) (list []system.SysHisEvent, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&system.SysHisEvent{})
	var sysHisEvents []system.SysHisEvent
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.IsRead != "" {
		db = db.Where("is_read = ?", info.IsRead)
	}
	if info.DataType != "" {
		db = db.Where("data_type = ?", info.DataType)
	}
	if len(info.HappenTimeRange) == 2 {
		db = db.Where("happen_time BETWEEN ? AND ?", info.HappenTimeRange[0], info.HappenTimeRange[1])
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Find(&sysHisEvents).Error
	return sysHisEvents, total, err
}
func (sysHisEventService *SysHisEventService) GetSysHisEventPublic(ctx context.Context) {
	// 此方法为获取数据源定义的数据
	// 请自行实现
}
