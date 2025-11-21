<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="searchFormRef" :inline="true" :model="searchInfo" class="demo-form-inline">
        <el-form-item label="电站编号">
          <el-input-number v-model.number="searchInfo.psid" placeholder="请输入电站编号" clearable />
        </el-form-item>
        <el-form-item label="设备编号">
          <el-input-number v-model.number="searchInfo.eqid" placeholder="请输入设备编号" clearable />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="loadChartData">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
          <el-button :type="autoRefresh ? 'danger' : 'success'" @click="toggleAutoRefresh">
            {{ autoRefresh ? '停止刷新' : '自动刷新(5s)' }}
          </el-button>
        </el-form-item>
      </el-form>
    </div>
    
    <div class="gva-table-box">
      <el-card>
        <template #header>
          <div class="card-header">
            <span>电站出力监控</span>
            <el-tag v-if="autoRefresh" type="success" effect="dark">自动刷新中...</el-tag>
          </div>
        </template>
        <div v-loading="loading" style="height: 500px">
          <div ref="chartRef" style="width: 100%; height: 100%"></div>
        </div>
      </el-card>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import * as echarts from 'echarts'
import { getPowerChartData } from '@/api/agvc/agvcChart'

defineOptions({
  name: 'PowerChart'
})

const searchFormRef = ref()
const chartRef = ref()
let chartInstance = null
const loading = ref(false)
const autoRefresh = ref(false)
let refreshTimer = null

const searchInfo = ref({
  psid: 1,
  eqid: 1
})

// 初始化图表
const initChart = () => {
  if (!chartRef.value) return
  
  if (chartInstance) {
    chartInstance.dispose()
  }
  
  chartInstance = echarts.init(chartRef.value)
  
  const option = {
    title: {
      text: '电站出力监控',
      left: 'center'
    },
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'cross'
      }
    },
    legend: {
      data: ['当前有功功率', '目标有功功率'],
      top: 30
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: [],
      axisLabel: {
        formatter: (value) => {
          const date = new Date(value)
          return `${date.getHours().toString().padStart(2, '0')}:${date.getMinutes().toString().padStart(2, '0')}:${date.getSeconds().toString().padStart(2, '0')}`
        }
      }
    },
    yAxis: {
      type: 'value',
      name: '功率 (kW)',
      axisLabel: {
        formatter: '{value}'
      }
    },
    series: [
      {
        name: '当前有功功率',
        type: 'line',
        data: [],
        smooth: true,
        itemStyle: {
          color: '#5470c6'
        }
      },
      {
        name: '目标有功功率',
        type: 'line',
        data: [],
        smooth: true,
        itemStyle: {
          color: '#91cc75'
        }
      }
    ]
  }
  
  chartInstance.setOption(option)
}

// 加载图表数据
const loadChartData = async() => {
  if (!searchInfo.value.psid || !searchInfo.value.eqid) {
    ElMessage.warning('请输入电站编号和设备编号')
    return
  }
  
  loading.value = true
  try {
    // 查询最近5分钟的数据
    const now = new Date()
    const fiveMinutesAgo = new Date(now.getTime() - 5 * 60 * 1000)
    
    const params = {
      psid: searchInfo.value.psid,
      eqid: searchInfo.value.eqid,
      startTime: fiveMinutesAgo.toISOString(),
      endTime: now.toISOString()
    }
    
    const res = await getPowerChartData(params)
    if (res.code === 0 && res.data) {
      updateChart(res.data)
    } else {
      ElMessage.error(res.msg || '获取数据失败')
    }
  } catch (error) {
    ElMessage.error('获取数据失败: ' + error.message)
  } finally {
    loading.value = false
  }
}

// 更新图表数据
const updateChart = (data) => {
  if (!chartInstance || !data || data.length === 0) {
    ElMessage.warning('暂无数据')
    return
  }
  
  const times = data.map(item => item.time)
  const currentPowers = data.map(item => item.currentActivePower || null)
  const targetPowers = data.map(item => item.targetActivePower || null)
  
  chartInstance.setOption({
    xAxis: {
      data: times
    },
    series: [
      {
        data: currentPowers
      },
      {
        data: targetPowers
      }
    ]
  })
}

// 切换自动刷新
const toggleAutoRefresh = () => {
  autoRefresh.value = !autoRefresh.value
  if (autoRefresh.value) {
    startAutoRefresh()
  } else {
    stopAutoRefresh()
  }
}

// 启动自动刷新
const startAutoRefresh = () => {
  if (refreshTimer) {
    clearInterval(refreshTimer)
  }
  loadChartData() // 立即加载一次
  refreshTimer = setInterval(() => {
    loadChartData()
  }, 5000)
}

// 停止自动刷新
const stopAutoRefresh = () => {
  if (refreshTimer) {
    clearInterval(refreshTimer)
    refreshTimer = null
  }
}

// 重置
const onReset = () => {
  searchInfo.value = {
    psid: 1,
    eqid: 1
  }
  stopAutoRefresh()
  autoRefresh.value = false
}

// 窗口大小改变时调整图表
const handleResize = () => {
  if (chartInstance) {
    chartInstance.resize()
  }
}

onMounted(async() => {
  await nextTick()
  initChart()
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  stopAutoRefresh()
  window.removeEventListener('resize', handleResize)
  if (chartInstance) {
    chartInstance.dispose()
  }
})
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
