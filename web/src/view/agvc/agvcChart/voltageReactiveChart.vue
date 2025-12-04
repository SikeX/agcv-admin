<template>
  <div>
    <div class="gva-table-box">
      <el-card>
        <div  style="height: 400px">
          <div ref="voltageChartRef" style="width: 100%; height: 100%"></div>
        </div>
      </el-card>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, nextTick, watch } from 'vue'
import { ElMessage } from 'element-plus'
import * as echarts from 'echarts'
import { getVoltageReactiveChartData } from '@/api/agvc/agvcChart'

defineOptions({
  name: 'VoltageReactiveChart'
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

const voltageChartRef = ref()
let voltageChartInstance = null
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

// 初始化电压图表
const initVoltageChart = () => {
  if (!voltageChartRef.value) return
  
  if (voltageChartInstance) {
    voltageChartInstance.dispose()
  }
  
  voltageChartInstance = echarts.init(voltageChartRef.value)
  
  const option = {
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
      data: ['当前电压', '目标电压'],
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
      name: '电压 (kV)',
      axisLabel: {
        formatter: '{value}'
      },
      min: function (value) {
        return Math.floor(value.min - 1);
      },
      max: function (value) {
        return Math.ceil(value.max + 1);
      },
    },
    series: [
      {
        name: '当前电压',
        type: 'line',
        data: [],
        smooth: true,
        itemStyle: {
          color: '#5470c6'
        }
      },
      {
        name: '目标电压',
        type: 'line',
        data: [],
        smooth: true,
        itemStyle: {
          color: '#91cc75'
        }
      }
    ]
  }
  
  voltageChartInstance.setOption(option)
}

// 加载图表数据
const loadChartData = async() => {
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
    
    const res = await getVoltageReactiveChartData(params)
    if (res.code === 0 && res.data) {
      updateCharts(res.data)
    }
  } catch (error) {
    console.error('获取数据失败:', error)
  } 
}

// 更新图表数据
const updateCharts = (data) => {
  if (!data || data.length === 0) {
    return
  }
  
  // 对数据按时间排序
  const sortedData = [...data].sort((a, b) => new Date(a.time) - new Date(b.time))
  
  const times = sortedData.map(item => item.time)
  const currentVoltages = sortedData.map(item => item.currentVoltage || null)
  const targetVoltages = sortedData.map(item => item.targetVoltage || null)
  
  // 更新电压图表
  if (voltageChartInstance) {
    voltageChartInstance.setOption({
      xAxis: {
        data: times
      },
      series: [
        {
          data: currentVoltages
        },
        {
          data: targetVoltages
        }
      ]
    })
  }
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
  if (voltageChartInstance) {
    voltageChartInstance.resize()
  }
}

// 监听props变化，重新加载数据
watch(() => [props.psid, props.eqid], () => {
  if (props.psid !== "" && props.eqid !== "") {
    // 重新初始化图表以确保完全重新渲染
    initVoltageChart()
    loadChartData()
  }
}, { immediate: false })

onMounted(async() => {
  await nextTick()
  initVoltageChart()
  startAutoRefresh() // 默认启动自动刷新
  handleResize
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  stopAutoRefresh()
  handleResize
  window.removeEventListener('resize', handleResize)
  if (voltageChartInstance) {
    voltageChartInstance.dispose()
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
