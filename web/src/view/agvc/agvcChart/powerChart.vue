<template>
  <div>
    <div class="gva-table-box">
      <el-card>
        <div style="height: 400px">
          <div ref="chartRef" style="width: 100%; height: 100%"></div>
        </div>
      </el-card>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, nextTick, watch, computed } from 'vue'
import { ElMessage } from 'element-plus'
import * as echarts from 'echarts'
import { getPowerChartData } from '@/api/agvc/agvcChart'
import { useAppStore } from '@/pinia'

defineOptions({
  name: 'PowerChart'
})

const appStore = useAppStore()
const isDark = computed(() => appStore.isDark)

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
let refreshTimer = null

const tooltipFormatter = (params) =>{
  const value = params[0].axisValue
  const date = new Date(value)
  const dateStr = `${date.getHours().toString().padStart(2, '0')}:${date.getMinutes().toString().padStart(2, '0')}:${date.getSeconds().toString().padStart(2, '0')}`
  const dateStyle =`时间:${dateStr}<br>`
  let allStyle = ''
  for (const item of params) {
   allStyle += ` ${item.marker} ${item.seriesName}: ${item.value || 0}<br>`
  }
  allStyle = dateStyle + allStyle

  return allStyle
}

// 初始化图表
const initChart = () => {
  if (!chartRef.value) return
  
  if (chartInstance) {
    chartInstance.dispose()
  }
  
  chartInstance = echarts.init(chartRef.value)
  
  // 根据暗黑模式设置颜色
  const backgroundColor = isDark.value ? '#19202D' : 'transparent'
  const textColor = isDark.value ? 'white' : '#333'
  const gridLineColor = isDark.value ? '#303642' : '#e0e6f1'
  
  const option = {
    backgroundColor: backgroundColor,
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'cross'
      },
      formatter: (params) => {
        return tooltipFormatter(params)
      }
    },
    legend: {
      data: ['当前有功功率', '目标有功功率'],
      top: 30,
      textStyle: {
        color: textColor
      }
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
        color: textColor,
        formatter: (value) => {
          const date = new Date(value)
          return `${date.getHours().toString().padStart(2, '0')}:${date.getMinutes().toString().padStart(2, '0')}:${date.getSeconds().toString().padStart(2, '0')}`
        }
      },
      axisLine: {
        lineStyle: {
          color: gridLineColor
        }
      },
      splitLine: {
        lineStyle: {
          color: gridLineColor
        }
      }
    },
    yAxis: {
      type: 'value',
      name: '功率 (kW)',
      nameTextStyle: {
        color: textColor
      },
      axisLabel: {
        color: textColor,
        formatter: '{value}'
      },
      axisLine: {
        lineStyle: {
          color: gridLineColor
        }
      },
      splitLine: {
        lineStyle: {
          color: gridLineColor
        }
      },
      // interval: 0.2,
      min: function (value) {
        return Math.floor(value.min - 1);
      },
      max: function (value) {
        return Math.ceil(value.max + 1);
      },
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
  console.log("props.psid, props.eqid", props.psid, props.eqid)
  if (props.psid === "" || props.eqid === "") {
    return
  }
  
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
  if (props.psid !== "" && props.eqid !== "") {
    // 重新初始化图表以确保完全重新渲染
    initChart()
    loadChartData()
  }
}, { immediate: false })

// 监听暗黑模式变化，重新初始化图表
watch(isDark, () => {
  if (chartInstance) {
    initChart()
    loadChartData()
  }
})

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
