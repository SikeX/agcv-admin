
<template>
  <div>
    <div class="gva-form-box">
      <el-form :model="formData" ref="elFormRef" label-position="right" :rules="rule" label-width="80px">
        <el-form-item label="设备名称:" prop="deviceName">
    <el-input v-model="formData.deviceName" :clearable="true" placeholder="请输入设备名称" />
</el-form-item>
        <el-form-item label="设备类型:" prop="deviceType">
    <el-input v-model="formData.deviceType" :clearable="true" placeholder="请输入设备类型" />
</el-form-item>
        <el-form-item label="设备位置:" prop="devicePosition">
    <el-input v-model="formData.devicePosition" :clearable="true" placeholder="请输入设备位置" />
</el-form-item>
        <el-form-item label="设备厂家:" prop="deviceFactory">
    <el-input v-model="formData.deviceFactory" :clearable="true" placeholder="请输入设备厂家" />
</el-form-item>
        <el-form-item label="设备型号:" prop="deviceModel">
    <el-input v-model.number="formData.deviceModel" :clearable="true" placeholder="请输入设备型号" />
</el-form-item>
        <el-form-item label="并网点:" prop="gcpName">
    <el-input v-model="formData.gcpName" :clearable="true" placeholder="请输入并网点" />
</el-form-item>
        <el-form-item label="安装角度:" prop="installAngle">
    <el-input-number v-model="formData.installAngle" style="width:100%" :precision="2" :clearable="true" />
</el-form-item>
        <el-form-item label="是否主气象仪:" prop="isMaster">
    <el-input v-model="formData.isMaster" :clearable="true" placeholder="请输入是否主气象仪" />
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
  createSysQixiangyiSetting,
  updateSysQixiangyiSetting,
  findSysQixiangyiSetting
} from '@/api/system/sysQixiangyiSetting'

defineOptions({
    name: 'SysQixiangyiSettingForm'
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
            deviceName: '',
            deviceType: '',
            devicePosition: '',
            deviceFactory: '',
            deviceModel: undefined,
            gcpName: '',
            installAngle: 0,
            isMaster: '',
        })
// 验证规则
const rule = reactive({
})

const elFormRef = ref()

// 初始化方法
const init = async () => {
 // 建议通过url传参获取目标数据ID 调用 find方法进行查询数据操作 从而决定本页面是create还是update 以下为id作为url参数示例
    if (route.query.id) {
      const res = await findSysQixiangyiSetting({ ID: route.query.id })
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
               res = await createSysQixiangyiSetting(formData.value)
               break
             case 'update':
               res = await updateSysQixiangyiSetting(formData.value)
               break
             default:
               res = await createSysQixiangyiSetting(formData.value)
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
