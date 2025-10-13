
<template>
  <div>
    <div class="gva-form-box">
      <el-form :model="formData" ref="elFormRef" label-position="right" :rules="rule" label-width="80px">
        <el-form-item label="电站名称:" prop="psid">
    <el-input v-model.number="formData.psid" :clearable="true" placeholder="请输入电站名称" />
</el-form-item>
        <el-form-item label="设备编号:" prop="eqid">
    <el-input v-model.number="formData.eqid" :clearable="true" placeholder="请输入设备编号" />
</el-form-item>
        <el-form-item label="设备类型:" prop="eqType">
    <el-input v-model.number="formData.eqType" :clearable="true" placeholder="请输入设备类型" />
</el-form-item>
        <el-form-item label="数据类型:" prop="dataType">
    <el-select v-model="formData.dataType" placeholder="请选择数据类型" style="width:100%" filterable :clearable="true">
        <el-option v-for="(item,key) in DataTypeOptions" :key="key" :label="item.label" :value="item.value" />
    </el-select>
</el-form-item>
        <el-form-item label="点号:" prop="datapoint">
    <el-input v-model="formData.datapoint" :clearable="true" placeholder="请输入点号" />
</el-form-item>
        <el-form-item label="发生时间:" prop="happenTime">
    <el-date-picker v-model="formData.happenTime" type="date" style="width:100%" placeholder="选择日期" :clearable="true" />
</el-form-item>
        <el-form-item label="数值:" prop="dataValue">
    <el-input v-model="formData.dataValue" :clearable="true" placeholder="请输入数值" />
</el-form-item>
        <el-form-item label="已读:" prop="isRead">
    <el-input v-model="formData.isRead" :clearable="true" placeholder="请输入已读" />
</el-form-item>
        <el-form-item label="名称:" prop="name">
    <el-input v-model="formData.name" :clearable="true" placeholder="请输入名称" />
</el-form-item>
        <el-form-item label="恢复时间:" prop="recoverTime">
    <el-date-picker v-model="formData.recoverTime" type="date" style="width:100%" placeholder="选择日期" :clearable="true" />
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
  createSysHisEvent,
  updateSysHisEvent,
  findSysHisEvent
} from '@/api/system/sysHisEvent'

defineOptions({
    name: 'SysHisEventForm'
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
const formData = ref({
            psid: undefined,
            eqid: undefined,
            eqType: undefined,
            dataType: '',
            datapoint: '',
            happenTime: new Date(),
            dataValue: '',
            isRead: '',
            name: '',
            recoverTime: new Date(),
        })
// 验证规则
const rule = reactive({
})

const elFormRef = ref()

// 初始化方法
const init = async () => {
 // 建议通过url传参获取目标数据ID 调用 find方法进行查询数据操作 从而决定本页面是create还是update 以下为id作为url参数示例
    if (route.query.id) {
      const res = await findSysHisEvent({ ID: route.query.id })
      if (res.code === 0) {
        formData.value = res.data
        type.value = 'update'
      }
    } else {
      type.value = 'create'
    }
    DataTypeOptions.value = await getDictFunc('DataType')
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
               res = await createSysHisEvent(formData.value)
               break
             case 'update':
               res = await updateSysHisEvent(formData.value)
               break
             default:
               res = await createSysHisEvent(formData.value)
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
