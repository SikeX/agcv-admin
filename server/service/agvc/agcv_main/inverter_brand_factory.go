package agcv_main

import (
	"fmt"
	"sync"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/service/agvc/cons"
	"go.uber.org/zap"
)

// InverterBrandFactory 逆变器品牌工厂
type InverterBrandFactory struct {
	brands map[int]InverterBrandInterface
	mu     sync.RWMutex
}

var BrandFactory = &InverterBrandFactory{
	brands: make(map[int]InverterBrandInterface),
}

// RegisterBrand 注册品牌实现
func (f *InverterBrandFactory) RegisterBrand(brand InverterBrandInterface) {
	f.mu.Lock()
	defer f.mu.Unlock()

	brandCode := brand.GetBrandCode()
	f.brands[brandCode] = brand

	//global.GVA_LOG.Info("注册逆变器品牌",
	//	zap.Int("brandCode", brandCode),
	//	zap.String("brandName", brand.GetBrandName()))
}

// GetBrand 获取品牌实现
func (f *InverterBrandFactory) GetBrand(brandCode int) (InverterBrandInterface, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	brand, exists := f.brands[brandCode]
	if !exists {
		return nil, fmt.Errorf("未找到品牌代码 %d 的实现", brandCode)
	}

	return brand, nil
}

// GetBrandName 获取品牌名称
func (f *InverterBrandFactory) GetBrandName(brandCode int) string {
	brand, err := f.GetBrand(brandCode)
	if err != nil {
		return "未知品牌"
	}
	return brand.GetBrandName()
}

// InitializeAllBrands 初始化所有已注册品牌的点位映射
func (f *InverterBrandFactory) InitializeAllBrands() error {
	f.mu.RLock()
	defer f.mu.RUnlock()

	global.GVA_LOG.Info("开始初始化所有品牌点位映射",
		zap.Int("品牌数量", len(f.brands)))

	for brandCode, brand := range f.brands {
		if err := brand.InitializePoints(); err != nil {
			global.GVA_LOG.Warn("初始化品牌点位失败",
				zap.Int("brandCode", brandCode),
				zap.String("brandName", brand.GetBrandName()),
				zap.Error(err))
			continue
		}
		global.GVA_LOG.Info("初始化品牌点位成功",
			zap.Int("brandCode", brandCode),
			zap.String("brandName", brand.GetBrandName()))
	}

	global.GVA_LOG.Info("所有品牌点位映射初始化完成")
	return nil
}

// GetAllRegisteredBrands 获取所有已注册的品牌
func (f *InverterBrandFactory) GetAllRegisteredBrands() []InverterBrandInterface {
	f.mu.RLock()
	defer f.mu.RUnlock()

	brands := make([]InverterBrandInterface, 0, len(f.brands))
	for _, brand := range f.brands {
		brands = append(brands, brand)
	}
	return brands
}

// init 初始化时自动注册所有品牌
func init() {
	// 注册华为品牌
	BrandFactory.RegisterBrand(&HuaweiInverterBrand{})

	// 未来可以在这里注册其他品牌
	// BrandFactory.RegisterBrand(&SungrowInverterBrand{})
	// BrandFactory.RegisterBrand(&GoodweInverterBrand{})

	//global.GVA_LOG.Info("逆变器品牌工厂初始化完成")
}

// GetDefaultBrand 获取默认品牌（华为）
func (f *InverterBrandFactory) GetDefaultBrand() InverterBrandInterface {
	brand, err := f.GetBrand(cons.INVERTER_BRAND_HUAWEI)
	if err != nil {
		global.GVA_LOG.Warn("获取默认品牌失败，返回nil", zap.Error(err))
		return nil
	}
	return brand
}
