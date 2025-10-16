
<template>
  <div>
    <div class="gva-form-box">
      <el-form :model="formData" ref="elFormRef" label-position="right" :rules="rule" label-width="80px">
        <el-form-item label="入库点名:" prop="inStorageName">
    <el-input v-model="formData.inStorageName" :clearable="true" placeholder="请输入入库点名" />
</el-form-item>
        <el-form-item label="测点名称:" prop="pointName">
    <el-input v-model="formData.pointName" :clearable="true" placeholder="请输入测点名称" />
</el-form-item>
        <el-form-item label="测点类型:" prop="pointType">
    <el-select v-model="formData.pointType" placeholder="请选择测点类型" style="width:100%" filterable :clearable="true">
        <el-option v-for="(item,key) in DataTypeOptions" :key="key" :label="item.label" :value="item.value" />
    </el-select>
</el-form-item>
        <el-form-item label="设备类型:" prop="deviceType">
    <el-select v-model="formData.deviceType" placeholder="请选择设备类型" style="width:100%" filterable :clearable="true">
        <el-option v-for="(item,key) in DeviceTypeOptions" :key="key" :label="item.label" :value="item.value" />
    </el-select>
</el-form-item>
        <el-form-item label="是否全栈闭锁:" prop="isAllClose">
    <el-input v-model="formData.isAllClose" :clearable="true" placeholder="请输入是否全栈闭锁" />
</el-form-item>
        <el-form-item label="是否设备闭锁:" prop="isDeviceClose">
    <el-input v-model="formData.isDeviceClose" :clearable="true" placeholder="请输入是否设备闭锁" />
</el-form-item>
        <el-form-item label="闭锁位置:" prop="closePosition">
    <el-input v-model.number="formData.closePosition" :clearable="true" placeholder="请输入闭锁位置" />
</el-form-item>
        <el-form-item label="闭锁上限:" prop="closeRangeMax">
    <el-input-number v-model="formData.closeRangeMax" style="width:100%" :precision="2" :clearable="true" />
</el-form-item>
        <el-form-item label="闭锁下限:" prop="closeRangeMin">
    <el-input-number v-model="formData.closeRangeMin" style="width:100%" :precision="2" :clearable="true" />
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
  createSysTestPoint,
  updateSysTestPoint,
  findSysTestPoint
} from '@/api/system/sysTestPoint'

defineOptions({
    name: 'SysTestPointForm'
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
const DataTypeOptions = ref([])
const DeviceTypeOptions = ref([])
const formData = ref({
            inStorageName: '',
            pointName: '',
            pointType: '',
            deviceType: '',
            isAllClose: '',
            isDeviceClose: '',
            closePosition: undefined,
            closeRangeMax: 0,
            closeRangeMin: 0,
        })
// 验证规则
const rule = reactive({
})

const elFormRef = ref()

// 初始化方法
const init = async () => {
 // 建议通过url传参获取目标数据ID 调用 find方法进行查询数据操作 从而决定本页面是create还是update 以下为id作为url参数示例
    if (route.query.id) {
      const res = await findSysTestPoint({ ID: route.query.id })
      if (res.code === 0) {
        formData.value = res.data
        type.value = 'update'
      }
    } else {
      type.value = 'create'
    }
    DataTypeOptions.value = await getDictFunc('DataType')
    DeviceTypeOptions.value = await getDictFunc('DeviceType')
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
               res = await createSysTestPoint(formData.value)
               break
             case 'update':
               res = await updateSysTestPoint(formData.value)
               break
             default:
               res = await createSysTestPoint(formData.value)
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
