package agcv_main

import (
	"fmt"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/agvc"
	"github.com/flipped-aurora/gin-vue-admin/server/model/agvc/agvc_main"
	"github.com/flipped-aurora/gin-vue-admin/server/model/agvc/agvc_main/request"
)

type device struct{}

var Device = new(device)

// CreateDevice 创建设备
func (s *device) CreateDevice(dev *agvc_main.Device) error {
	// 生成设备编号
	dev.DeviceCode = fmt.Sprintf("%s%s%s%s", dev.PSID, dev.EQID, dev.DataType, dev.EQType)
	return global.GVA_DB.Create(dev).Error
}

// DeleteDevice 删除设备
func (s *device) DeleteDevice(id uint) error {
	return global.GVA_DB.Delete(&agvc_main.Device{}, id).Error
}

// UpdateDevice 更新设备
func (s *device) UpdateDevice(dev *agvc_main.Device) error {
	// 更新设备编号
	dev.DeviceCode = fmt.Sprintf("%s%s%s%s", dev.PSID, dev.EQID, dev.DataType, dev.EQType)
	return global.GVA_DB.Model(&agvc_main.Device{}).Where("id = ?", dev.ID).Updates(dev).Error
}

// GetDevice 获取设备详情
func (s *device) GetDevice(id uint) (agvc_main.Device, error) {
	var dev agvc_main.Device
	err := global.GVA_DB.Where("id = ?", id).First(&dev).Error
	return dev, err
}

// GetDeviceByCode 根据设备编号获取设备
func (s *device) GetDeviceByCode(deviceCode string) (agvc_main.Device, error) {
	var dev agvc_main.Device
	err := global.GVA_DB.Where("device_code = ?", deviceCode).First(&dev).Error
	return dev, err
}

// GetDeviceByPSIDAndEQID 根据PSID和EQID获取设备
func (s *device) GetDeviceByPSIDAndEQID(psid, eqid, eqType string) (agvc_main.Device, error) {
	var dev agvc_main.Device
	err := global.GVA_DB.Where("psid = ? AND eqid = ? AND eq_type = ?", psid, eqid, eqType).First(&dev).Error
	return dev, err
}

// GetDeviceList 获取设备列表
func (s *device) GetDeviceList(req request.DeviceSearch) ([]agvc_main.Device, int64, error) {
	var devices []agvc_main.Device
	var total int64

	db := global.GVA_DB.Model(&agvc_main.Device{})

	// 条件过滤
	if req.PSID != "" {
		db = db.Where("psid = ?", req.PSID)
	}
	if req.EQID != "" {
		db = db.Where("eqid = ?", req.EQID)
	}
	if req.EQType != "" {
		db = db.Where("eq_type = ?", req.EQType)
	}
	if req.DataType != "" {
		db = db.Where("data_type = ?", req.DataType)
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}
	if req.Keyword != "" {
		db = db.Where("name LIKE ? OR device_code LIKE ?", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	if req.PageSize > 0 {
		offset := (req.Page - 1) * req.PageSize
		db = db.Limit(req.PageSize).Offset(offset)
	}

	err = db.Order("created_at DESC").Find(&devices).Error
	return devices, total, err
}

// GetInvertersByPSID 获取电站的所有逆变器
func (s *device) GetInvertersByPSID(psid string) ([]agvc_main.Device, error) {
	var devices []agvc_main.Device
	err := global.GVA_DB.Where("psid = ? AND eq_type = ?", psid, "01").Find(&devices).Error
	return devices, err
}

// GetOnlineInvertersByBwdNo 获取电站的所有在线逆变器
func (s *device) GetOnlineInvertersByBwdNo(bwdNo int) ([]agvc.AgvcNbqSetting, error) {
	var devices []agvc.AgvcNbqSetting
	// status := 1
	err := global.GVA_DB.Where("bwdNo = ? AND isParticipateAdjust = ? ", bwdNo, true).Find(&devices).Error
	return devices, err
}

// UpdateDeviceStatus 更新设备状态
func (s *device) UpdateDeviceStatus(psid, eqid, eqType string, status int) error {
	return global.GVA_DB.Model(&agvc_main.Device{}).
		Where("psid = ? AND eqid = ? AND eq_type = ?", psid, eqid, eqType).
		Update("status", status).Error
}

// GetPointMapping 获取测点映射
func (s *device) GetPointMapping(eqType, dataType int, point string) (agvc_main.PointMapping, error) {
	var mapping agvc_main.PointMapping
	err := global.GVA_DB.Where("eq_type = ? AND data_type = ? AND point = ?", eqType, dataType, point).
		First(&mapping).Error
	return mapping, err
}

// GetPointMappingsByCategory 根据分类获取测点映射
func (s *device) GetPointMappingsByCategory(category string) ([]agvc_main.PointMapping, error) {
	var mappings []agvc_main.PointMapping
	err := global.GVA_DB.Where("category = ?", category).Find(&mappings).Error
	return mappings, err
}

// GetPointMappingsByDevice 获取设备的所有测点映射
func (s *device) GetPointMappingsByDevice(eqType, dataType string) ([]agvc_main.PointMapping, error) {
	var mappings []agvc_main.PointMapping
	err := global.GVA_DB.Where("eq_type = ? AND data_type = ?", eqType, dataType).Find(&mappings).Error
	return mappings, err
}

// CreatePointMapping 创建测点映射
func (s *device) CreatePointMapping(mapping *agvc_main.PointMapping) error {
	return global.GVA_DB.Create(mapping).Error
}

// BatchCreatePointMappings 批量创建测点映射
func (s *device) BatchCreatePointMappings(mappings []agvc_main.PointMapping) error {
	return global.GVA_DB.CreateInBatches(mappings, 100).Error
}

// GetDeviceRealtimeData 获取设备实时数据（结合实时数据和测点映射）
func (s *device) GetDeviceRealtimeData(psid, eqid, eqType string) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	// 获取设备基本信息
	dev, err := s.GetDeviceByPSIDAndEQID(psid, eqid, eqType)
	if err != nil {
		return nil, fmt.Errorf("设备不存在: %v", err)
	}

	result["device"] = dev

	// 获取所有实时数据
	allData := DataStorage.GetDeviceAllData(psid, eqid, eqType)

	// 按数据类型组织数据
	dataByType := make(map[int]map[string]interface{})
	for key, data := range allData {
		if dataByType[data.DataType] == nil {
			dataByType[data.DataType] = make(map[string]interface{})
		}

		// 获取测点映射信息
		mapping, err := s.GetPointMapping(data.EQType, data.DataType, data.Point)
		if err == nil {
			dataByType[data.DataType][data.Point] = map[string]interface{}{
				"value":       data.Value,
				"pointName":   mapping.PointName,
				"unit":        mapping.Unit,
				"description": mapping.Description,
				"timestamp":   data.Timestamp,
			}
		} else {
			// 如果没有映射信息，只返回原始数据
			dataByType[data.DataType][data.Point] = map[string]interface{}{
				"value":     data.Value,
				"timestamp": data.Timestamp,
			}
		}
		_ = key // unused variable
	}

	result["realtimeData"] = dataByType

	return result, nil
}
