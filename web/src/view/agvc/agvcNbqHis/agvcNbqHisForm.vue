
<template>
  <div>
    <div class="gva-form-box">
      <el-form :model="formData" ref="elFormRef" label-position="right" :rules="rule" label-width="80px">
        <el-form-item label="逆变器编号:" prop="inverterNo">
    <el-input v-model.number="formData.inverterNo" :clearable="true" placeholder="请输入逆变器编号" />
</el-form-item>
        <el-form-item label="逆变器名称:" prop="name">
    <el-input v-model="formData.name" :clearable="true" placeholder="请输入逆变器名称" />
</el-form-item>
        <el-form-item label="逆变器运行状态:" prop="status">
    <el-input v-model="formData.status" :clearable="true" placeholder="请输入逆变器运行状态" />
</el-form-item>
        <el-form-item label="是否调节:" prop="isParticipateAdjust">
    <el-switch v-model="formData.isParticipateAdjust" active-color="#13ce66" inactive-color="#ff4949" active-text="是" inactive-text="否" clearable ></el-switch>
</el-form-item>
        <el-form-item label="额定功率(kW):" prop="ratedPower">
    <el-input-number v-model="formData.ratedPower" style="width:100%" :precision="2" :clearable="true" />
</el-form-item>
        <el-form-item label="有功功率:" prop="activePower">
    <el-input-number v-model="formData.activePower" style="width:100%" :precision="2" :clearable="true" />
</el-form-item>
        <el-form-item label="有功目标(kW):" prop="apTargetValue">
    <el-input-number v-model="formData.apTargetValue" style="width:100%" :precision="2" :clearable="true" />
</el-form-item>
        <el-form-item label="无功功率(kVar):" prop="reactivePower">
    <el-input-number v-model="formData.reactivePower" style="width:100%" :precision="2" :clearable="true" />
</el-form-item>
        <el-form-item label="无功目标(kVar):" prop="rpTargetValue">
    <el-input-number v-model="formData.rpTargetValue" style="width:100%" :precision="2" :clearable="true" />
</el-form-item>
        <el-form-item label="功率因数:" prop="powerFactor">
    <el-input-number v-model="formData.powerFactor" style="width:100%" :precision="2" :clearable="true" />
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
  createAgvcNbqHis,
  updateAgvcNbqHis,
  findAgvcNbqHis
} from '@/api/agvc/agvcNbqHis'

defineOptions({
    name: 'AgvcNbqHisForm'
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
            inverterNo: undefined,
            name: '',
            status: '',
            isParticipateAdjust: false,
            ratedPower: 0,
            activePower: 0,
            apTargetValue: 0,
            reactivePower: 0,
            rpTargetValue: 0,
            powerFactor: 0,
        })
// 验证规则
const rule = reactive({
})

const elFormRef = ref()

// 初始化方法
const init = async () => {
 // 建议通过url传参获取目标数据ID 调用 find方法进行查询数据操作 从而决定本页面是create还是update 以下为id作为url参数示例
    if (route.query.id) {
      const res = await findAgvcNbqHis({ ID: route.query.id })
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
               res = await createAgvcNbqHis(formData.value)
               break
             case 'update':
               res = await updateAgvcNbqHis(formData.value)
               break
             default:
               res = await createAgvcNbqHis(formData.value)
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
