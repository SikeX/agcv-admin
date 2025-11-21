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
      <el-row :gutter="20">
        <el-col :span="12">
          <el-card>
            <template #header>
              <div class="card-header">
                <span>电压监控</span>
                <el-tag v-if="autoRefresh" type="success" effect="dark">自动刷新中...</el-tag>
              </div>
            </template>
            <div v-loading="loading" style="height: 400px">
              <div ref="voltageChartRef" style="width: 100%; height: 100%"></div>
            </div>
          </el-card>
        </el-col>
        
        <el-col :span="12">
          <el-card>
            <template #header>
              <div class="card-header">
                <span>无功监控</span>
                <el-tag v-if="autoRefresh" type="success" effect="dark">自动刷新中...</el-tag>
              </div>
            </template>
            <div v-loading="loading" style="height: 400px">
              <div ref="reactiveChartRef" style="width: 100%; height: 100%"></div>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import * as echarts from 'echarts'
import { getVoltageReactiveChartData } from '@/api/agvc/agvcChart'

defineOptions({
  name: 'VoltageReactiveChart'
})

const searchFormRef = ref()
const voltageChartRef = ref()
const reactiveChartRef = ref()
let voltageChartInstance = null
let reactiveChartInstance = null
const loading = ref(false)
const autoRefresh = ref(false)
let refreshTimer = null

const searchInfo = ref({
  psid: 1,
  eqid: 1
})

// 初始化电压图表
const initVoltageChart = () => {
  if (!voltageChartRef.value) return
  
  if (voltageChartInstance) {
    voltageChartInstance.dispose()
  }
  
  voltageChartInstance = echarts.init(voltageChartRef.value)
  
  const option = {
    title: {
      text: '电压监控',
      left: 'center',
      textStyle: {
        fontSize: 14
      }
    },
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'cross'
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
      }
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

// 初始化无功图表
const initReactiveChart = () => {
  if (!reactiveChartRef.value) return
  
  if (reactiveChartInstance) {
    reactiveChartInstance.dispose()
  }
  
  reactiveChartInstance = echarts.init(reactiveChartRef.value)
  
  const option = {
    title: {
      text: '无功监控',
      left: 'center',
      textStyle: {
        fontSize: 14
      }
    },
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'cross'
      }
    },
    legend: {
      data: ['当前无功', '目标无功'],
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
      name: '无功 (kVar)',
      axisLabel: {
        formatter: '{value}'
      }
    },
    series: [
      {
        name: '当前无功',
        type: 'line',
        data: [],
        smooth: true,
        itemStyle: {
          color: '#ee6666'
        }
      },
      {
        name: '目标无功',
        type: 'line',
        data: [],
        smooth: true,
        itemStyle: {
          color: '#fac858'
        }
      }
    ]
  }
  
  reactiveChartInstance.setOption(option)
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
    
    const res = await getVoltageReactiveChartData(params)
    if (res.code === 0 && res.data) {
      updateCharts(res.data)
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
const updateCharts = (data) => {
  if (!data || data.length === 0) {
    ElMessage.warning('暂无数据')
    return
  }
  
  const times = data.map(item => item.time)
  const currentVoltages = data.map(item => item.currentVoltage || null)
  const targetVoltages = data.map(item => item.targetVoltage || null)
  const currentReactivePowers = data.map(item => item.currentReactivePower || null)
  const targetReactivePowers = data.map(item => item.targetReactivePower || null)
  
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
  
  // 更新无功图表
  if (reactiveChartInstance) {
    reactiveChartInstance.setOption({
      xAxis: {
        data: times
      },
      series: [
        {
          data: currentReactivePowers
        },
        {
          data: targetReactivePowers
        }
      ]
    })
  }
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
  if (voltageChartInstance) {
    voltageChartInstance.resize()
  }
  if (reactiveChartInstance) {
    reactiveChartInstance.resize()
  }
}

onMounted(async() => {
  await nextTick()
  initVoltageChart()
  initReactiveChart()
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  stopAutoRefresh()
  window.removeEventListener('resize', handleResize)
  if (voltageChartInstance) {
    voltageChartInstance.dispose()
  }
  if (reactiveChartInstance) {
    reactiveChartInstance.dispose()
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
