package initialize

import (
    "github.com/flipped-aurora/gin-vue-admin/server/global"
    "github.com/flipped-aurora/gin-vue-admin/server/model/agvc"
)

func bizModel() error {
    db := global.GVA_DB
    err := db.AutoMigrate(
        agvc.AgvcBwdSetting{},
        agvc.AgvcQxySetting{},
        agvc.AgvcNbqSetting{},
        agvc.AgvcEventHis{},
        agvc.AgvcNbqHis{},
        agvc.AgvcBwdHis{},
        agvc.AgvcQxyHis{},
        agvc.AgvcScheduleCurve{}, // 新增：计划曲线表
    )
    if err != nil {
        return err
    }
    return nil
}
