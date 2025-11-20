<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" class="demo-form-inline" @keyup.enter="onSubmit">
        <el-form-item label="设备编号">
          <el-input v-model="searchInfo.inverter_no" placeholder="请输入设备编号" clearable />
        </el-form-item>
        <el-form-item label="逆变器名称">
          <el-input v-model="searchInfo.name" placeholder="请输入逆变器名称" clearable />
        </el-form-item>
        <el-form-item label="时间范围">
          <el-date-picker
            v-model="queryInfo.dateRange"
            type="datetimerange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            style="width: 400px;"
            :disabled-date="disabledDate"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    
    <!-- 历史数据表格 -->
    <div class="gva-table-box">
      <div class="gva-btn-list">
            <ExportNonDbExcel 
              api-url="/agvcNbqHis/exportAgvcNbqHistory" 
              :condition="searchInfo" 
              :export-func="exportAgvcNbqHistory"
            />
        </div>
      <el-table
        :data="tableData"
        border
        style="width: 100%"
        max-height="600"
        v-loading="loading"
      >
        <el-table-column align="center" label="设备编号" prop="inverterNo" width="120" fixed="left" />
        <el-table-column align="center" label="设备名称" prop="name" width="150" fixed="left" />
        <el-table-column align="center" label="采集时间" prop="ctime" width="180" fixed="left">
          <template #default="scope">
            {{ formatTime(scope.row.ctime) }}
          </template>
        </el-table-column>
        
        
        <!-- 发电量相关 -->
        <el-table-column align="center" label="总发电量(kWh)" prop="totalPowerGeneration" width="150">
          <template #default="scope">{{ formatValue(scope.row.totalPowerGeneration) }}</template>
        </el-table-column>
        <el-table-column align="center" label="日发电量(kWh)" prop="dailyPowerGeneration" width="150">
          <template #default="scope">{{ formatValue(scope.row.dailyPowerGeneration) }}</template>
        </el-table-column>
        <el-table-column align="center" label="月发电量(kWh)" prop="monthlyPowerGeneration" width="150">
          <template #default="scope">{{ formatValue(scope.row.monthlyPowerGeneration) }}</template>
        </el-table-column>
        <el-table-column align="center" label="年发电量(kWh)" prop="annualPowerGeneration" width="150">
          <template #default="scope">{{ formatValue(scope.row.annualPowerGeneration) }}</template>
        </el-table-column>
        
        <!-- 功率相关 -->
        <el-table-column align="center" label="交流功率(kW)" prop="acPower" width="130">
          <template #default="scope">{{ formatValue(scope.row.acPower) }}</template>
        </el-table-column>
        <el-table-column align="center" label="直流功率(kW)" prop="dcPower" width="130">
          <template #default="scope">{{ formatValue(scope.row.dcPower) }}</template>
        </el-table-column>
        <el-table-column align="center" label="无功功率(kVar)" prop="reactivePower" width="140">
          <template #default="scope">{{ formatValue(scope.row.reactivePower) }}</template>
        </el-table-column>
        <el-table-column align="center" label="视在功率(kVa)" prop="apparentPower" width="140">
          <template #default="scope">{{ formatValue(scope.row.apparentPower) }}</template>
        </el-table-column>
        
        <!-- 电压电流 -->
        <el-table-column align="center" label="A相电压(V)" prop="phaseAVoltage" width="130">
          <template #default="scope">{{ formatValue(scope.row.phaseAVoltage) }}</template>
        </el-table-column>
        <el-table-column align="center" label="B相电压(V)" prop="phaseBVoltage" width="130">
          <template #default="scope">{{ formatValue(scope.row.phaseBVoltage) }}</template>
        </el-table-column>
        <el-table-column align="center" label="C相电压(V)" prop="phaseCVoltage" width="130">
          <template #default="scope">{{ formatValue(scope.row.phaseCVoltage) }}</template>
        </el-table-column>
        <el-table-column align="center" label="A相电流(A)" prop="phaseACurrent" width="130">
          <template #default="scope">{{ formatValue(scope.row.phaseACurrent) }}</template>
        </el-table-column>
        <el-table-column align="center" label="B相电流(A)" prop="phaseBCurrent" width="130">
          <template #default="scope">{{ formatValue(scope.row.phaseBCurrent) }}</template>
        </el-table-column>
        <el-table-column align="center" label="C相电流(A)" prop="phaseCCurrent" width="130">
          <template #default="scope">{{ formatValue(scope.row.phaseCCurrent) }}</template>
        </el-table-column>
        
        <!-- 其他参数 -->
        <el-table-column align="center" label="功率因数" prop="powerFactor" width="110">
          <template #default="scope">{{ formatValue(scope.row.powerFactor) }}</template>
        </el-table-column>
        <el-table-column align="center" label="设备温度(℃)" prop="deviceTemperature" width="130">
          <template #default="scope">{{ formatValue(scope.row.deviceTemperature) }}</template>
        </el-table-column>
        <el-table-column align="center" label="转换效率(%)" prop="conversionEfficiency" width="130">
          <template #default="scope">{{ formatValue(scope.row.conversionEfficiency) }}</template>
        </el-table-column>
        <el-table-column align="center" label="电网频率(Hz)" prop="gridFrequency" width="130">
          <template #default="scope">{{ formatValue(scope.row.gridFrequency) }}</template>
        </el-table-column>
        <el-table-column align="center" label="设备状态码" prop="deviceStatusCode" width="120">
          <template #default="scope">{{ formatValue(scope.row.deviceStatusCode) }}</template>
        </el-table-column>
      </el-table>
      
      <!-- 分页组件 -->
      <div class="gva-pagination">
        <el-pagination
          layout="total, sizes, prev, pager, next, jumper"
          :current-page="currentPage"
          :page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          @current-change="handleCurrentChange"
          @size-change="handleSizeChange"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import {
  getAgvcNbqHisList,
  getAgvcNbqHistory,
  exportAgvcNbqHistory
} from '@/api/agvc/agvcNbqHis'

import { ElMessage } from 'element-plus'
import { start } from 'nprogress'
import { ref, onMounted } from 'vue'
// 导出组件
import ExportExcel from '@/components/exportExcel/exportExcel.vue'
import ExportTemplate from '@/components/exportExcel/exportTemplate.vue'
import ExportNonDbExcel from '@/components/exportExcel/exportNonDbExcel.vue'

defineOptions({
  name: 'AgvcNbqHistory'
})

const elSearchFormRef = ref()

// 设备列表相关变量
const tableData = ref([])
const searchInfo = ref({})

// 查询历史数据相关变量
const queryInfo = ref({
  selectedDevice: null,
  dateRange: []
})
const historyData = ref([])
const loading = ref(false)

// 分页相关变量
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

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

// 格式化时间显示
const formatTime = (timeStr) => {
  if (!timeStr) return '-'
  try {
    const date = new Date(timeStr)
    return date.toLocaleString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit'
    })
  } catch (e) {
    return timeStr
  }
}

// 禁用超过7天范围的日期
const disabledDate = (time) => {
  if (!queryInfo.value.dateRange || queryInfo.value.dateRange.length === 0) {
    return false
  }
  
  const selectedDate = queryInfo.value.dateRange[0]
  if (!selectedDate) {
    return false
  }
  
  // 计算选中日期前后7天的范围
  const minTime = new Date(selectedDate).getTime() - 7 * 24 * 60 * 60 * 1000
  const maxTime = new Date(selectedDate).getTime() + 7 * 24 * 60 * 60 * 1000
  
  return time.getTime() < minTime || time.getTime() > maxTime
}

// 重置搜索
const onReset = () => {
  searchInfo.value = {}
  queryInfo.value.dateRange = []
  currentPage.value = 1
  getTableData()
}

// 搜索设备
const onSubmit = () => {
  elSearchFormRef.value?.validate(async(valid) => {
    if (!valid) return
    
    // 验证时间范围不超过7天
    if (queryInfo.value.dateRange && queryInfo.value.dateRange.length === 2) {
      // 时间的格式“yyyy-MM-dd HH:mm:ss”
      const start = new Date(queryInfo.value.dateRange[0])
      const end = new Date(queryInfo.value.dateRange[1])
      const diffTime = Math.abs(end - start)
      const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24))
      
      if (diffDays > 7) {
        ElMessage.warning('查询时间范围不能超过7天')
        return
      }
    }
    
    currentPage.value = 1
    getTableData()
  })
}

// 分页相关
const handleSizeChange = (val) => {
  pageSize.value = val
  getTableData()
}

const handleCurrentChange = (val) => {
  currentPage.value = val
  getTableData()
}

// 查询设备列表
const getTableData = async() => {
  loading.value = true
  console.log(queryInfo.value.dateRange)
  if (queryInfo.value.dateRange == undefined || 
      queryInfo.value.dateRange.length === 0 
    ) { 
    //默认查询最近7天
    const end = new Date()
    const start = new Date()
    start.setTime(start.getTime() - 3600 * 1000 * 24 * 7)
    searchInfo.value.startTime = start
    searchInfo.value.endTime = end
  } else {
    searchInfo.value.startTime = queryInfo.value.dateRange[0]
    searchInfo.value.endTime = queryInfo.value.dateRange[1]
  }
  

  try {
    // 使用分页查询
    const table = await getAgvcNbqHisList({ 
      page: currentPage.value, 
      pageSize: pageSize.value, 
      ...searchInfo.value 
    })
    if (table.code === 0) {
      tableData.value = table.data.list || []
      total.value = table.data.total
      currentPage.value = table.data.page
      pageSize.value = table.data.pageSize
    }
  } catch (error) {
    ElMessage.error('获取设备列表失败: ' + error.message)
  } finally {
    loading.value = false
  }
}

// 设备选择变化
const onDeviceChange = () => {
  // 清空历史数据
  historyData.value = []
  
  // 设置默认日期范围为最近7天
  if (!queryInfo.value.dateRange || queryInfo.value.dateRange.length === 0) {
    const end = new Date()
    const start = new Date()
    start.setTime(start.getTime() - 3600 * 1000 * 24 * 7)
    queryInfo.value.dateRange = [start, end]
  }
}

// 加载历史数据
const loadHistoryData = async () => {
  // 检查是否选择了设备
  if (!queryInfo.value.selectedDevice) {
    ElMessage.warning('请先选择要查询的设备')
    return
  }
  
  // 检查时间范围是否有效
  if (!queryInfo.value.dateRange || queryInfo.value.dateRange.length !== 2) {
    ElMessage.warning('请选择时间范围')
    return
  }
  
  // 找到选中的设备
  const currentDevice = tableData.value.find(item => item.inverterNo === queryInfo.value.selectedDevice)
  if (!currentDevice) {
    ElMessage.error('未找到选中的设备')
    return
  }
  
  loading.value = true
  try {
    // 构建查询参数
    const agvcNbqHis = {
      psid: currentDevice.psid,
      inverterNo: currentDevice.inverterNo,
      name: currentDevice.name,
      // 设置一个字段的值以触发查询所有点位
      dailyPowerGeneration: 0
    }
    
    // 格式化时间为RFC3339格式
    const formatToRFC3339 = (date) => {
      if (typeof date === 'string') {
        date = new Date(date)
      }
      return date.toISOString()
    }
    
    const requestData = {
      agvcNbqHis: agvcNbqHis,
      startTime: formatToRFC3339(queryInfo.value.dateRange[0]),
      endTime: formatToRFC3339(queryInfo.value.dateRange[1])
    }
    
    const res = await getAgvcNbqHistory(requestData)
    if (res.code === 0) {
      if (res.data && res.data.length > 0) {
        // 按时间降序排列
        historyData.value = res.data.sort((a, b) => {
          return new Date(b.ctime) - new Date(a.ctime)
        })
        ElMessage.success(`成功加载 ${historyData.value.length} 条历史数据`)
      } else {
        historyData.value = []
        ElMessage.info('选择的时间范围内没有历史数据')
      }
    } else {
      ElMessage.error('获取历史数据失败: ' + res.msg)
    }
  } catch (error) {
    ElMessage.error('获取历史数据失败: ' + error.message)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  getTableData()
})
</script>

<style scoped>
.el-table {
  margin-top: 20px;
}

.gva-pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>
