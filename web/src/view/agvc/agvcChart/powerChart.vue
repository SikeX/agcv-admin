<template>
  <div>
    <div class="gva-table-box">
      <el-card>
        <template #header>
          <div class="card-header">
            <span>电站出力监控</span>
            <el-tag type="success" effect="dark">自动刷新中(5s)</el-tag>
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
import { ref, onMounted, onUnmounted, nextTick, watch } from 'vue'
import { ElMessage } from 'element-plus'
import * as echarts from 'echarts'
import { getPowerChartData } from '@/api/agvc/agvcChart'

defineOptions({
  name: 'PowerChart'
})

// 接收父组件传递的props
const props = defineProps({
  psid: {
    type: Number,
    required: true
  },
  eqid: {
    type: Number,
    required: true
  }
})

const chartRef = ref()
let chartInstance = null
const loading = ref(false)
let refreshTimer = null

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
  if (!props.psid || !props.eqid) {
    return
  }
  
  loading.value = true
  try {
    // 查询最近5分钟的数据
    const now = new Date()
    const fiveMinutesAgo = new Date(now.getTime() - 5 * 60 * 1000)
    
    const params = {
      psid: props.psid,
      eqid: props.eqid,
      startTime: fiveMinutesAgo.toISOString(),
      endTime: now.toISOString()
    }
    
    const res = await getPowerChartData(params)
    if (res.code === 0 && res.data) {
      updateChart(res.data)
    }
  } catch (error) {
    console.error('获取数据失败:', error)
  } finally {
    loading.value = false
  }
}

// 更新图表数据
const updateChart = (data) => {
  if (!chartInstance || !data || data.length === 0) {
    return
  }
  
  // 对数据按时间排序
  const sortedData = [...data].sort((a, b) => new Date(a.time) - new Date(b.time))
  
  const times = sortedData.map(item => item.time)
  const currentPowers = sortedData.map(item => item.currentActivePower || null)
  const targetPowers = sortedData.map(item => item.targetActivePower || null)
  
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

// 启动自动刷新 - 默认5s
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

// 窗口大小改变时调整图表
const handleResize = () => {
  if (chartInstance) {
    chartInstance.resize()
  }
}

// 监听props变化，重新加载数据
watch(() => [props.psid, props.eqid], () => {
  if (props.psid && props.eqid) {
    loadChartData()
  }
}, { immediate: false })

onMounted(async() => {
  await nextTick()
  initChart()
  startAutoRefresh() // 默认启动自动刷新
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
