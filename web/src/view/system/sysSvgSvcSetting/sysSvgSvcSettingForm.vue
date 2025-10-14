
<template>
  <div>
    <div class="gva-form-box">
      <el-form :model="formData" ref="elFormRef" label-position="right" :rules="rule" label-width="80px">
        <el-form-item label="无功补偿装置名称:" prop="wugongName">
    <el-input v-model="formData.wugongName" :clearable="true" placeholder="请输入无功补偿装置名称" />
</el-form-item>
        <el-form-item label="无功容量:" prop="wugongCapacity">
    <el-input-number v-model="formData.wugongCapacity" style="width:100%" :precision="2" :clearable="true" />
</el-form-item>
        <el-form-item label="容量系数:" prop="capacityCoefficient">
    <el-input-number v-model="formData.capacityCoefficient" style="width:100%" :precision="2" :clearable="true" />
</el-form-item>
        <el-form-item label="参与调节:" prop="isAdjustment">
    <el-input v-model="formData.isAdjustment" :clearable="true" placeholder="请输入参与调节" />
</el-form-item>
        <el-form-item label="备注:" prop="note">
    <el-input v-model="formData.note" :clearable="true" placeholder="请输入备注" />
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
  createSysSvgSvcSetting,
  updateSysSvgSvcSetting,
  findSysSvgSvcSetting
} from '@/api/system/sysSvgSvcSetting'

defineOptions({
    name: 'SysSvgSvcSettingForm'
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
            wugongName: '',
            wugongCapacity: 0,
            capacityCoefficient: 0,
            isAdjustment: '',
            note: '',
        })
// 验证规则
const rule = reactive({
})

const elFormRef = ref()

// 初始化方法
const init = async () => {
 // 建议通过url传参获取目标数据ID 调用 find方法进行查询数据操作 从而决定本页面是create还是update 以下为id作为url参数示例
    if (route.query.id) {
      const res = await findSysSvgSvcSetting({ ID: route.query.id })
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
               res = await createSysSvgSvcSetting(formData.value)
               break
             case 'update':
               res = await updateSysSvgSvcSetting(formData.value)
               break
             default:
               res = await createSysSvgSvcSetting(formData.value)
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
