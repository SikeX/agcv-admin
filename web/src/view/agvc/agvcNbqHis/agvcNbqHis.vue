
<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" class="demo-form-inline" @keyup.enter="onSubmit">
        <el-form-item label="逆变器编号">
          <el-input v-model="searchInfo.inverter_no" placeholder="请输入逆变器编号" clearable />
        </el-form-item>
        <el-form-item label="逆变器名称">
          <el-input v-model="searchInfo.name" placeholder="请输入逆变器名称" clearable />
        </el-form-item>
        
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    
    <div class="gva-table-box">
      <el-table
        ref="multipleTable"
        style="width: 100%"
        tooltip-effect="dark"
        :data="tableData"
        row-key="inverterNo"
        border
        :max-height="600"
      >
        <!-- 固定列 -->
        <el-table-column align="center" label="逆变器编号" prop="inverterNo" width="120" fixed="left" />
        <el-table-column align="center" label="逆变器名称" prop="name" width="150" fixed="left" />
        
        <!-- 动态生成数据列 -->
        <el-table-column
          v-for="field in dataFields"
          :key="field.prop"
          :label="field.label"
          :prop="field.prop"
          :width="field.width || 150"
          align="center"
        >
          <template #default="scope">
            <div class="field-cell">
              <span class="field-value">{{ formatValue(scope.row[field.prop]) }}</span>
              <el-button
                v-if="field.point"
                type="primary"
                link
                size="small"
                icon="TrendCharts"
                @click="showHistoryChart(scope.row, field)"
                class="history-btn"
              >
                历史
              </el-button>
            </div>
          </template>
        </el-table-column>
      </el-table>
      
      <div class="gva-pagination">
        <el-pagination
          layout="total, sizes, prev, pager, next, jumper"
          :current-page="page"
          :page-size="pageSize"
          :page-sizes="[10, 30, 50, 100]"
          :total="total"
          @current-change="handleCurrentChange"
          @size-change="handleSizeChange"
        />
      </div>
    </div>

    <!-- 历史数据弹窗 -->
    <el-dialog 
      v-model="historyDialogVisible" 
      :title="`${currentDevice?.name || ''} - ${currentField?.label || ''}历史数据`" 
      width="80%" 
      top="5vh"
    >
      <el-row :gutter="20">
        <!-- 日期选择区域 -->
        <el-col :span="6">
          <el-card>
            <template #header>
              <span>时间选择</span>
            </template>
            <el-date-picker
              v-model="historyDateRange"
              type="datetimerange"
              range-separator="至"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
              value-format="YYYY-MM-DD"
              style="width: 100%; margin-bottom: 20px;"
            />
            <el-button type="primary" @click="loadHistoryData" :loading="loading" style="width: 100%;">查询</el-button>
          </el-card>
        </el-col>
        
        <!-- 图表展示区域 -->
        <el-col :span="18">
          <el-card>
            <div ref="chartContainer" style="width: 100%; height: 500px;"></div>
          </el-card>
        </el-col>
      </el-row>
      
      <template #footer>
        <el-button @click="historyDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import {
  getAgvcNbqHisList,
  getAgvcNbqHistory
} from '@/api/agvc/agvcNbqHis'

import { ElMessage } from 'element-plus'
import { ref, reactive, nextTick, onMounted, onUnmounted } from 'vue'
import * as echarts from 'echarts'

defineOptions({
  name: 'AgvcNbqHis'
})

// 定义数据字段配置（根据AgvcNbqHis结构体的point标签生成）
const dataFields = ref([
  // 发电量相关
  { prop: 'totalPowerGeneration', label: '总发电量(kWh)', point: '1', width: 150 },
  { prop: 'dailyPowerGeneration', label: '日发电量(kWh)', point: '2', width: 150 },
  { prop: 'monthlyPowerGeneration', label: '月发电量(kWh)', point: '3', width: 150 },
  { prop: 'annualPowerGeneration', label: '年发电量(kWh)', point: '4', width: 150 },
  
  // 功率相关
  { prop: 'acPower', label: '交流功率(kW)', point: '10', width: 130 },
  { prop: 'dcPower', label: '直流功率(kW)', point: '202', width: 130 },
  { prop: 'reactivePower', label: '无功功率(kVar)', point: '27', width: 140 },
  { prop: 'apparentPower', label: '视在功率(kVa)', point: '28', width: 140 },
  { prop: 'gridFrequency', label: '电网频率(Hz)', point: '15', width: 130 },
  { prop: 'totalDCCurrent', label: '总直流电流(A)', point: '16', width: 130 },
  
  // 相电压
  { prop: 'phaseAVoltage', label: 'A相电压Ua(V)', point: '94', width: 140 },
  { prop: 'phaseBVoltage', label: 'B相电压Ub(V)', point: '95', width: 140 },
  { prop: 'phaseCVoltage', label: 'C相电压Uc(V)', point: '96', width: 140 },
  
  // 线电压
  { prop: 'lineABVoltage', label: 'AB线电压Uab(V)', point: '7', width: 150 },
  { prop: 'lineBCVoltage', label: 'BC线电压Ubc(V)', point: '8', width: 150 },
  { prop: 'lineCAVoltage', label: 'CA线电压Uca(V)', point: '9', width: 150 },
  
  // 相电流
  { prop: 'phaseACurrent', label: 'A相电流Ia(A)', point: '12', width: 130 },
  { prop: 'phaseBCurrent', label: 'B相电流Ib(A)', point: '13', width: 130 },
  { prop: 'phaseCCurrent', label: 'C相电流Ic(A)', point: '14', width: 130 },
  
  // 功率因数和温度
  { prop: 'powerFactor', label: '功率因数', point: '21', width: 110 },
  { prop: 'deviceTemperature', label: '设备温度(℃)', point: '22', width: 130 },
  { prop: 'conversionEfficiency', label: '转换效率(%)', point: '102', width: 130 },
  
  // 支路电压 (20路)
  ...Array.from({ length: 20 }, (_, i) => ({
    prop: `branchVoltage${i + 1}`,
    label: `支路电压${i + 1}(V)`,
    point: String(301 + i),
    width: 140
  })),
  
  // 支路电流 (20路)
  ...Array.from({ length: 20 }, (_, i) => ({
    prop: `branchCurrent${i + 1}`,
    label: `支路电流${i + 1}(A)`,
    point: String(36 + i),
    width: 140
  })),
  
  // MPPT
  { prop: 'mppt1', label: 'mppt1', point: '401', width: 110 },
  { prop: 'mppt2', label: 'mppt2', point: '402', width: 110 },
  { prop: 'mppt3', label: 'mppt3', point: '403', width: 110 },
  { prop: 'mppt4', label: 'mppt4', point: '404', width: 110 },
  
  // 设备状态和其他
  { prop: 'deviceStatusCode', label: '设备状态码', point: '501', width: 120 },
  { prop: 'insulationImpedance', label: '绝缘阻抗', point: '502', width: 110 },
  { prop: 'peakActivePower', label: '当天峰值有功功率(kW)', point: '503', width: 180 },
  { prop: 'startupTime', label: '开机时间', point: '511', width: 120 },
  { prop: 'shutdownTime', label: '关机时间', point: '512', width: 120 },
  { prop: 'phaseATemperature', label: 'A相温度', point: '251', width: 110 },
  { prop: 'phaseBTemperature', label: 'B相温度', point: '252', width: 110 },
  { prop: 'phaseCTemperature', label: 'C相温度', point: '253', width: 110 }
])

const elSearchFormRef = ref()

// 表格控制部分
const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({})
const loading = ref(false)

// 格式化数值显示
const formatValue = (value) => {
  if (value === null || value === undefined) {
    return '-'
  }
  if (typeof value === 'number') {
    return value.toFixed(2)
  }
  return value
}

// 重置
const onReset = () => {
  searchInfo.value = {}
  getTableData()
}

// 搜索
const onSubmit = () => {
  elSearchFormRef.value?.validate(async(valid) => {
    if (!valid) return
    page.value = 1
    getTableData()
  })
}

// 分页
const handleSizeChange = (val) => {
  pageSize.value = val
  getTableData()
}

// 修改页面容量
const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

// 查询
const getTableData = async() => {
  loading.value = true
  try {
    const table = await getAgvcNbqHisList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
    if (table.code === 0) {
      tableData.value = table.data.list
      total.value = table.data.total
      page.value = table.data.page
      pageSize.value = table.data.pageSize
    }
  } catch (error) {
    ElMessage.error('获取数据失败: ' + error.message)
  } finally {
    loading.value = false
  }
}

getTableData()

// 历史数据相关变量
const historyDialogVisible = ref(false)
const historyData = ref([])
const historyDateRange = ref([])
const chartContainer = ref(null)
const currentDevice = ref(null)
const currentField = ref(null)

// 图表实例
let chartInstance = null

// 显示历史数据图表
const showHistoryChart = async (row, field) => {
  // 保存当前设备和字段信息
  currentDevice.value = row
  currentField.value = field
  
  // 设置默认日期范围为最近7天
  const end = new Date()
  const start = new Date()
  start.setTime(start.getTime() - 3600 * 1000 * 24 * 7)
  historyDateRange.value = [start, end]
  
  // 显示弹窗
  historyDialogVisible.value = true
  
  // 在下次DOM更新后初始化图表并加载数据
  nextTick(() => {
    initChart()
    loadHistoryData()
  })
}

// 加载历史数据
const loadHistoryData = async () => {
  if (!currentDevice.value || !currentField.value) return
  
  loading.value = true
  try {
    const params = {
      eqid: currentDevice.value.inverterNo,
      point: currentField.value.point,
      startTime: historyDateRange.value[0],
      endTime: historyDateRange.value[1]
    }
    
    const res = await getAgvcNbqHistory(params)
    if (res.code === 0) {
      // 处理从 InfluxDB获取的历史数据
      historyData.value = res.data.map(item => {
        return {
          timestamp: new Date(item.time),
          value: item.value,
        }
      })
      // 按时间排序
      historyData.value.sort((a, b) => a.timestamp - b.timestamp)
      
      // 更新图表
      updateChart()
    } else {
      ElMessage.error('获取历史数据失败: ' + res.msg)
    }
  } catch (error) {
    ElMessage.error('获取历史数据失败: ' + error.message)
  } finally {
    loading.value = false
  }
}

// 初始化图表
const initChart = () => {
  if (!chartContainer.value) return
  
  // 使用 echarts 绘制图表
  if (!chartInstance) {
    chartInstance = echarts.init(chartContainer.value)
  }
}

// 更新图表
const updateChart = () => {
  if (!chartInstance || !historyData.value.length) {
    if (chartInstance) {
      chartInstance.setOption({
        title: {
          text: '暂无数据',
          left: 'center',
          top: 'center',
          textStyle: {
            fontSize: 20,
            color: '#999'
          }
        }
      })
    }
    return
  }
  
  // 准备图表数据
  const xData = historyData.value.map(item => 
    item.timestamp.toLocaleString('zh-CN', {
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit'
    })
  )
  const yData = historyData.value.map(item => item.value)
  
  // 配置图表选项
  const option = {
    title: {
      text: `${currentField.value?.label || ''} 历史趋势`,
      left: 'center',
      top: 10
    },
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'cross'
      },
      formatter: (params) => {
        if (params && params.length > 0) {
          const param = params[0]
          return `${param.name}<br/>${currentField.value?.label}: ${param.value}`
        }
        return ''
      }
    },
    grid: {
      left: '3%',
      right: '4%',
      top: 60,
      bottom: 80,
      containLabel: true
    },
    xAxis: {
      type: 'category',
      data: xData,
      boundaryGap: false,
      axisLabel: {
        rotate: 45,
        fontSize: 11
      }
    },
    yAxis: {
      type: 'value',
      name: currentField.value?.label || '',
      axisLabel: {
        fontSize: 11
      }
    },
    series: [{
      name: currentField.value?.label,
      type: 'line',
      data: yData,
      smooth: true,
      lineStyle: {
        width: 2,
        color: '#5470c6'
      },
      areaStyle: {
        color: {
          type: 'linear',
          x: 0,
          y: 0,
          x2: 0,
          y2: 1,
          colorStops: [{
            offset: 0, color: 'rgba(84, 112, 198, 0.4)'
          }, {
            offset: 1, color: 'rgba(84, 112, 198, 0.1)'
          }]
        }
      }
    }],
    dataZoom: [
      {
        type: 'inside',
        start: 0,
        end: 100
      },
      {
        type: 'slider',
        start: 0,
        end: 100,
        height: 25,
        bottom: 15,
        handleSize: '110%',
        handleStyle: {
          color: '#5470c6'
        },
        textStyle: {
          fontSize: 11
        }
      }
    ]
  }
  
  // 设置图表选项
  chartInstance.setOption(option, true)
}

// 监听窗口大小变化，重置图表大小
const handleResize = () => {
  if (chartInstance) {
    chartInstance.resize()
  }
}

// 在组件挂载时添加窗口大小监听
onMounted(() => {
  window.addEventListener('resize', handleResize)
})

// 在组件卸载时移除监听
onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  if (chartInstance) {
    chartInstance.dispose()
    chartInstance = null
  }
})
</script>

<style scoped>
.field-cell {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}

.field-value {
  flex: 1;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.history-btn {
  flex-shrink: 0;
  padding: 2px 4px;
  font-size: 12px;
}
</style>
