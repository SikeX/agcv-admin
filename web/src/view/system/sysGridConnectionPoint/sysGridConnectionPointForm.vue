
<template>
  <div>
    <div class="gva-form-box">
      <el-form :model="formData" ref="elFormRef" label-position="right" :rules="rule" label-width="80px">
        <el-form-item label="并网点名称:" prop="name">
    <el-input v-model="formData.name" :clearable="true" placeholder="请输入并网点名称" />
</el-form-item>
        <el-form-item label="电压等级(kV):" prop="voltageLevel">
    <el-input-number v-model="formData.voltageLevel" style="width:100%" :precision="2" :clearable="true" />
</el-form-item>
        <el-form-item label="AGC功能退出模式:" prop="agcFunctionExit">
    <el-input v-model.number="formData.agcFunctionExit" :clearable="true" placeholder="请输入AGC功能退出模式" />
</el-form-item>
        <el-form-item label="AGC调节步长(kW):" prop="agcStepSize">
    <el-input-number v-model="formData.agcStepSize" style="width:100%" :precision="2" :clearable="true" />
</el-form-item>
        <el-form-item label="AGC步长周期(秒):" prop="agcStepPeriod">
    <el-input v-model.number="formData.agcStepPeriod" :clearable="true" placeholder="请输入AGC步长周期(秒)" />
</el-form-item>
        <el-form-item label="AGC抖动区间(kW):" prop="agcVibrationRange">
    <el-input-number v-model="formData.agcVibrationRange" style="width:100%" :precision="2" :clearable="true" />
</el-form-item>
        <el-form-item label="AGC调控周期(秒):" prop="agcControlPeriod">
    <el-input v-model.number="formData.agcControlPeriod" :clearable="true" placeholder="请输入AGC调控周期(秒)" />
</el-form-item>
        <el-form-item label="AGC微调系数:" prop="agcMicroAdjustmentCoefficient">
    <el-input-number v-model="formData.agcMicroAdjustmentCoefficient" style="width:100%" :precision="2" :clearable="true" />
</el-form-item>
        <el-form-item label="AVC调节步长(kV):" prop="avcStepSize">
    <el-input-number v-model="formData.avcStepSize" style="width:100%" :precision="2" :clearable="true" />
</el-form-item>
        <el-form-item label="AVC步长周期(秒):" prop="avcStepPeriod">
    <el-input v-model.number="formData.avcStepPeriod" :clearable="true" placeholder="请输入AVC步长周期(秒)" />
</el-form-item>
        <el-form-item label="AVC抖动区间(kV):" prop="avcVibrationRange">
    <el-input-number v-model="formData.avcVibrationRange" style="width:100%" :precision="2" :clearable="true" />
</el-form-item>
        <el-form-item label="AVC调控周期(秒):" prop="avcControlPeriod">
    <el-input v-model.number="formData.avcControlPeriod" :clearable="true" placeholder="请输入AVC调控周期(秒)" />
</el-form-item>
        <el-form-item label="AVC系统阻抗:" prop="avcSystemImpedance">
    <el-input-number v-model="formData.avcSystemImpedance" style="width:100%" :precision="2" :clearable="true" />
</el-form-item>
        <el-form-item label="AVC调节范围最小值(kV):" prop="avcAdjustmentRangeMin">
    <el-input-number v-model="formData.avcAdjustmentRangeMin" style="width:100%" :precision="2" :clearable="true" />
</el-form-item>
        <el-form-item label="AVC调节范围最大值(kV):" prop="avcAdjustmentRangeMax">
    <el-input-number v-model="formData.avcAdjustmentRangeMax" style="width:100%" :precision="2" :clearable="true" />
</el-form-item>
        <el-form-item>
          <el-button :loading="btnLoading" type="primary" @click="save">保存</el-button>
          <el-button type="primary" @click="back">返回</el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup>
import {
  createSysGridConnectionPoint,
  updateSysGridConnectionPoint,
  findSysGridConnectionPoint
} from '@/api/system/sysGridConnectionPoint'

defineOptions({
    name: 'SysGridConnectionPointForm'
})

// 自动获取字典
import { getDictFunc } from '@/utils/format'
import { useRoute, useRouter } from "vue-router"
import { ElMessage } from 'element-plus'
import { ref, reactive } from 'vue'


const route = useRoute()
const router = useRouter()

// 提交按钮loading
const btnLoading = ref(false)

const type = ref('')
const formData = ref({
            name: '',
            voltageLevel: 0,
            agcFunctionExit: undefined,
            agcStepSize: 0,
            agcStepPeriod: undefined,
            agcVibrationRange: 0,
            agcControlPeriod: undefined,
            agcMicroAdjustmentCoefficient: 0,
            avcStepSize: 0,
            avcStepPeriod: undefined,
            avcVibrationRange: 0,
            avcControlPeriod: undefined,
            avcSystemImpedance: 0,
            avcAdjustmentRangeMin: 0,
            avcAdjustmentRangeMax: 0,
        })
// 验证规则
const rule = reactive({
})

const elFormRef = ref()

// 初始化方法
const init = async () => {
 // 建议通过url传参获取目标数据ID 调用 find方法进行查询数据操作 从而决定本页面是create还是update 以下为id作为url参数示例
    if (route.query.id) {
      const res = await findSysGridConnectionPoint({ ID: route.query.id })
      if (res.code === 0) {
        formData.value = res.data
        type.value = 'update'
      }
    } else {
      type.value = 'create'
    }
}

init()
// 保存按钮
const save = async() => {
      btnLoading.value = true
      elFormRef.value?.validate( async (valid) => {
         if (!valid) return btnLoading.value = false
            let res
           switch (type.value) {
             case 'create':
               res = await createSysGridConnectionPoint(formData.value)
               break
             case 'update':
               res = await updateSysGridConnectionPoint(formData.value)
               break
             default:
               res = await createSysGridConnectionPoint(formData.value)
               break
           }
           btnLoading.value = false
           if (res.code === 0) {
             ElMessage({
               type: 'success',
               message: '创建/更改成功'
             })
           }
       })
}

// 返回按钮
const back = () => {
    router.go(-1)
}

</script>

<style>
</style>
