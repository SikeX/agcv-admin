package agvc

import (
	"context"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/agvc"
	agvcReq "github.com/flipped-aurora/gin-vue-admin/server/model/agvc/request"
)

type AgvcEventHisService struct{}

// CreateAgvcEventHis 创建历史事件记录
// Author [yourname](https://github.com/yourname)
func (agvcEventHisService *AgvcEventHisService) CreateAgvcEventHis(ctx context.Context, agvcEventHis *agvc.AgvcEventHis) (err error) {
	err = global.GVA_DB.Create(agvcEventHis).Error
	return err
}

// DeleteAgvcEventHis 删除历史事件记录
// Author [yourname](https://github.com/yourname)
func (agvcEventHisService *AgvcEventHisService) DeleteAgvcEventHis(ctx context.Context, ID string) (err error) {
	err = global.GVA_DB.Delete(&agvc.AgvcEventHis{}, "id = ?", ID).Error
	return err
}

// DeleteAgvcEventHisByIds 批量删除历史事件记录
// Author [yourname](https://github.com/yourname)
func (agvcEventHisService *AgvcEventHisService) DeleteAgvcEventHisByIds(ctx context.Context, IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]agvc.AgvcEventHis{}, "id in ?", IDs).Error
	return err
}

// UpdateAgvcEventHis 更新历史事件记录
// Author [yourname](https://github.com/yourname)
func (agvcEventHisService *AgvcEventHisService) UpdateAgvcEventHis(ctx context.Context, agvcEventHis agvc.AgvcEventHis) (err error) {
	err = global.GVA_DB.Model(&agvc.AgvcEventHis{}).Where("id = ?", agvcEventHis.ID).Updates(&agvcEventHis).Error
	return err
}

// GetAgvcEventHis 根据ID获取历史事件记录
// Author [yourname](https://github.com/yourname)
func (agvcEventHisService *AgvcEventHisService) GetAgvcEventHis(ctx context.Context, ID string) (agvcEventHis agvc.AgvcEventHis, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&agvcEventHis).Error
	return
}

// GetAgvcEventHisInfoList 分页获取历史事件记录
// Author [yourname](https://github.com/yourname)
func (agvcEventHisService *AgvcEventHisService) GetAgvcEventHisInfoList(ctx context.Context, info agvcReq.AgvcEventHisSearch) (list []agvc.AgvcEventHis, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&agvc.AgvcEventHis{})
	var agvcEventHiss []agvc.AgvcEventHis
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.Psid != "" {
		db = db.Where("psid = ?", info.Psid)
	}
	if info.EqType != "" {
		db = db.Where("eq_type = ?", info.EqType)
	}
	if info.Eqid != "" {
		db = db.Where("eqid LIKE ?", "%"+info.Eqid+"%")
	}
	if len(info.RecordTimeRange) == 2 {
		db = db.Where("record_time BETWEEN ? AND ?", info.RecordTimeRange[0], info.RecordTimeRange[1])
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Find(&agvcEventHiss).Error
	return agvcEventHiss, total, err
}
func (agvcEventHisService *AgvcEventHisService) GetAgvcEventHisPublic(ctx context.Context) {
	// 此方法为获取数据源定义的数据
	// 请自行实现
}
