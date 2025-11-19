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
        <el-table-column align="center" label="设备编号" prop="inverterNo" width="120" fixed="left" />
        <el-table-column align="center" label="逆变器名称" prop="name" width="150" fixed="left" />
        
        <!-- 操作列 -->
        <el-table-column align="center" label="操作" width="120" fixed="right">
          <template #default="scope">
            <el-button
              type="primary"
              link
              icon="View"
              @click="showHistoryData(scope.row)"
            >
              查看历史数据
            </el-button>
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
      :title="`${currentDevice?.name || ''} - 历史数据`" 
      width="95%" 
      top="2vh"
      :close-on-click-modal="false"
    >
      <div>
        <!-- 时间选择器 -->
        <el-row :gutter="20" style="margin-bottom: 20px;">
          <el-col :span="18">
            <el-date-picker
              v-model="historyDateRange"
              type="datetimerange"
              range-separator="至"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
              style="width: 100%;"
            />
          </el-col>
          <el-col :span="6">
            <el-button type="primary" @click="loadHistoryData" :loading="loading" style="width: 100%;">查询历史数据</el-button>
          </el-col>
        </el-row>

        <!-- 历史数据表格 -->
        <el-table
          :data="historyData"
          border
          style="width: 100%"
          max-height="600"
          v-loading="loading"
        >
          <el-table-column align="center" label="采集时间" prop="ctime" width="180" fixed="left">
            <template #default="scope">
              {{ formatTime(scope.row.ctime) }}
            </template>
          </el-table-column>
          <el-table-column align="center" label="设备编号" prop="inverterNo" width="120" fixed="left" />
          <el-table-column align="center" label="设备名称" prop="name" width="150" fixed="left" />
          
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
      </div>
      
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
import { ref, onMounted } from 'vue'

defineOptions({
  name: 'AgvcNbqHistory'
})

const elSearchFormRef = ref()

// 表格控制部分
const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({})
const loading = ref(false)

// 历史数据相关变量
const historyDialogVisible = ref(false)
const historyData = ref([])
const historyDateRange = ref([])
const currentDevice = ref(null)

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

// 查询设备列表
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

// 显示历史数据
const showHistoryData = (row) => {
  currentDevice.value = row
  historyDialogVisible.value = true
  
  // 设置默认日期范围为最近7天
  const end = new Date()
  const start = new Date()
  start.setTime(start.getTime() - 3600 * 1000 * 24 * 7)
  historyDateRange.value = [start, end]
  
  // 清空历史数据
  historyData.value = []
}

// 加载历史数据
const loadHistoryData = async () => {
  if (!currentDevice.value) return
  
  // 检查时间范围是否有效
  if (!historyDateRange.value || historyDateRange.value.length !== 2) {
    ElMessage.warning('请选择时间范围')
    return
  }
  
  loading.value = true
  try {
    // 构建查询参数
    const agvcNbqHis = {
      psid: currentDevice.value.psid,
      inverterNo: currentDevice.value.inverterNo,
      name: currentDevice.value.name,
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
      startTime: formatToRFC3339(historyDateRange.value[0]),
      endTime: formatToRFC3339(historyDateRange.value[1])
    }
    
    const res = await getAgvcNbqHistory(requestData)
    if (res.code === 0) {
      if (res.data && res.data.length > 0) {
        // 按时间降序排列
        historyData.value = res.data.sort((a, b) => {
          return new Date(b.ctime) - new Date(a.ctime)
        })
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
</style>
