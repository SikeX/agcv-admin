<template>
  <el-button type="primary" icon="download" @click="exportExcelFunc"
    >导出历史数据</el-button
  >
</template>

<script setup>
import { ElMessage } from 'element-plus'

const props = defineProps({
  apiUrl: {
    type: String,
    required: true
  },
  condition: {
    type: Object,
    default: () => ({})
  },
  exportFunc: {
    type: Function,
    required: true
  }
})

const exportExcelFunc = async () => {
  if (!props.apiUrl || !props.exportFunc) {
    ElMessage.error('导出组件配置错误')
    return
  }

  try {
    let baseUrl = import.meta.env.VITE_BASE_API
    if (baseUrl === "/"){
      baseUrl = ""
    }

    // 调用传入的API函数获取导出token
    const res = await props.exportFunc(props.condition)

    if(res.code === 0){
      ElMessage.success('创建导出任务成功，开始下载')
      const url = `${baseUrl}${res.data}`
      window.open(url, '_blank')
    } else {
      ElMessage.error(res.msg || '创建导出任务失败')
    }
  } catch (error) {
    ElMessage.error('导出失败: ' + error.message)
  }
}
</script>
