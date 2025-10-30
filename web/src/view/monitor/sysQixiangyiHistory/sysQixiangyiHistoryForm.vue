
<template>
  <div>
    <div class="gva-form-box">
      <el-form :model="formData" ref="elFormRef" label-position="right" :rules="rule" label-width="80px">
        <el-form-item label="气象仪名称:" prop="name">
    <el-input v-model="formData.name" :clearable="true" placeholder="请输入气象仪名称" />
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
  getSysQixiangyiHistoryList,
  getQixiangyiHistoryData,
  generateTestData
} from '@/api/monitor/sysQixiangyiHistory'

defineOptions({
    name: 'SysQixiangyiHistoryForm'
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
        })
// 验证规则
const rule = reactive({
})

const elFormRef = ref()

// 初始化方法
const init = async () => {
 // 气象仪历史数据页面不需要CRUD操作，这里可以获取气象仪列表用于选择
    if (route.query.id) {
      // 如果有ID参数，可以用于指定特定的气象仪
      type.value = 'view'
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
           // 气象仪历史数据页面主要用于查看，这里可以实现生成测试数据等功能
           res = await generateTestData({ qixiangyiId: formData.value.id || 1 })
           btnLoading.value = false
           if (res.code === 0) {
             ElMessage({
               type: 'success',
               message: '操作成功'
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
