
<template>
  <div>
    <div class="gva-form-box">
      <el-form :model="formData" ref="elFormRef" label-position="right" :rules="rule" label-width="80px">
        <el-form-item label="逆变器编号:" prop="inverterNo">
    <el-input v-model="formData.inverterNo" :clearable="true" placeholder="请输入逆变器编号" />
</el-form-item>
        <el-form-item label="并网点编号:" prop="bwdNo">
    <el-input v-model="formData.bwdNo" :clearable="true" placeholder="请输入并网点编号" />
</el-form-item>
        <el-form-item label="逆变器名称:" prop="name">
    <el-input v-model="formData.name" :clearable="true" placeholder="请输入逆变器名称" />
</el-form-item>
        <el-form-item label="额定有功功率:" prop="ratedActivePower">
    <el-input-number v-model="formData.ratedActivePower" style="width:100%" :precision="2" :clearable="true" />
</el-form-item>
        <el-form-item label="额定无功功率:" prop="ratedReactivePower">
    <el-input-number v-model="formData.ratedReactivePower" style="width:100%" :precision="2" :clearable="true" />
</el-form-item>
        <el-form-item label="功率抖动区间:" prop="jitterRange">
    <el-input-number v-model="formData.jitterRange" style="width:100%" :precision="2" :clearable="true" />
</el-form-item>
        <el-form-item label="功率死区区间:" prop="deadbandRange">
    <el-input-number v-model="formData.deadbandRange" style="width:100%" :precision="2" :clearable="true" />
</el-form-item>
        <el-form-item label="升额优先级:" prop="upgradePriority">
    <el-input-number v-model="formData.upgradePriority" style="width:100%" :precision="2" :clearable="true" />
</el-form-item>
        <el-form-item label="降级优先级:" prop="downgradePriority">
    <el-input-number v-model="formData.downgradePriority" style="width:100%" :precision="2" :clearable="true" />
</el-form-item>
        <el-form-item label="参与调节:" prop="isParticipateAdjust">
    <el-switch v-model="formData.isParticipateAdjust" active-color="#13ce66" inactive-color="#ff4949" active-text="是" inactive-text="否" clearable ></el-switch>
</el-form-item>
        <el-form-item label="标杆逆变器:" prop="isBenchmarkInverter">
    <el-switch v-model="formData.isBenchmarkInverter" active-color="#13ce66" inactive-color="#ff4949" active-text="是" inactive-text="否" clearable ></el-switch>
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
  createAgvcNbqSetting,
  updateAgvcNbqSetting,
  findAgvcNbqSetting
} from '@/api/agvc/agvcNbqSetting'

defineOptions({
    name: 'AgvcNbqSettingForm'
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
            inverterNo: '',
            bwdNo: '',
            name: '',
            ratedActivePower: 0,
            ratedReactivePower: 0,
            jitterRange: 0,
            deadbandRange: 0,
            upgradePriority: 0,
            downgradePriority: 0,
            isParticipateAdjust: false,
            isBenchmarkInverter: false,
        })
// 验证规则
const rule = reactive({
  inverterNo: [
    { required: true, message: '请输入逆变器编号', trigger: 'blur' },
    { pattern: /^\d+$/, message: '逆变器编号只能输入数字', trigger: 'blur' }
  ],
  bwdNo: [
    { required: true, message: '请输入并网点编号', trigger: 'blur' },
    { pattern: /^\d+$/, message: '并网点编号只能输入数字', trigger: 'blur' }
  ]
})

const elFormRef = ref()

// 初始化方法
const init = async () => {
 // 建议通过url传参获取目标数据ID 调用 find方法进行查询数据操作 从而决定本页面是create还是update 以下为id作为url参数示例
    if (route.query.id) {
      const res = await findAgvcNbqSetting({ ID: route.query.id })
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
               res = await createAgvcNbqSetting(formData.value)
               break
             case 'update':
               res = await updateAgvcNbqSetting(formData.value)
               break
             default:
               res = await createAgvcNbqSetting(formData.value)
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
