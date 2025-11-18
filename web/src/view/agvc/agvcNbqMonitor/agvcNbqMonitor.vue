<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" class="demo-form-inline" @keyup.enter="onSubmit">
        <el-form-item label="电站编号">
          <el-input v-model.number="searchInfo.psid" placeholder="请输入电站编号" clearable />
        </el-form-item>
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
        <el-table-column align="center" label="电站编号" prop="psid" width="100" fixed="left" />
        <el-table-column align="center" label="逆变器编号" prop="inverterNo" width="120" fixed="left" />
        <el-table-column align="center" label="逆变器名称" prop="name" width="150" fixed="left" />
        
        <!-- 主要数据列 -->
        <el-table-column align="center" label="日发电量(kWh)" prop="dailyPowerGeneration" width="150">
          <template #default="scope">
            {{ formatValue(scope.row.dailyPowerGeneration) }}
          </template>
        </el-table-column>
        <el-table-column align="center" label="交流功率(kW)" prop="acPower" width="130">
          <template #default="scope">
            {{ formatValue(scope.row.acPower) }}
          </template>
        </el-table-column>
        <el-table-column align="center" label="采集时间" prop="ctime" width="180">
          <template #default="scope">
            {{ formatTime(scope.row.ctime) }}
          </template>
        </el-table-column>
        
        <!-- 操作列 -->
        <el-table-column align="center" label="操作" width="120" fixed="right">
          <template #default="scope">
            <el-button
              type="primary"
              link
              icon="View"
              @click="showMoreData(scope.row)"
            >
              更多
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

    <!-- 更多数据弹窗 -->
    <el-dialog 
      v-model="moreDialogVisible" 
      :title="`${currentDevice?.name || ''} - 全部实时数据`" 
      width="90%" 
      top="5vh"
    >
      <div v-loading="moreDataLoading">
        <el-descriptions :column="3" border>
          <el-descriptions-item label="采集时间" :span="3">
            {{ formatTime(realData?.ctime) }}
          </el-descriptions-item>
          
          <!-- 发电量相关 -->
          <el-descriptions-item label="总发电量(kWh)">{{ formatValue(realData?.totalPowerGeneration) }}</el-descriptions-item>
          <el-descriptions-item label="日发电量(kWh)">{{ formatValue(realData?.dailyPowerGeneration) }}</el-descriptions-item>
          <el-descriptions-item label="月发电量(kWh)">{{ formatValue(realData?.monthlyPowerGeneration) }}</el-descriptions-item>
          <el-descriptions-item label="年发电量(kWh)">{{ formatValue(realData?.annualPowerGeneration) }}</el-descriptions-item>
          
          <!-- 功率相关 -->
          <el-descriptions-item label="交流功率(kW)">{{ formatValue(realData?.acPower) }}</el-descriptions-item>
          <el-descriptions-item label="直流功率(kW)">{{ formatValue(realData?.dcPower) }}</el-descriptions-item>
          <el-descriptions-item label="无功功率(kVar)">{{ formatValue(realData?.reactivePower) }}</el-descriptions-item>
          <el-descriptions-item label="视在功率(kVa)">{{ formatValue(realData?.apparentPower) }}</el-descriptions-item>
          <el-descriptions-item label="电网频率(Hz)">{{ formatValue(realData?.gridFrequency) }}</el-descriptions-item>
          <el-descriptions-item label="总直流电流(A)">{{ formatValue(realData?.totalDCCurrent) }}</el-descriptions-item>
          <el-descriptions-item label="当天峰值有功功率(kW)">{{ formatValue(realData?.peakActivePower) }}</el-descriptions-item>
          
          <!-- 相电压 -->
          <el-descriptions-item label="A相电压(V)">{{ formatValue(realData?.phaseAVoltage) }}</el-descriptions-item>
          <el-descriptions-item label="B相电压(V)">{{ formatValue(realData?.phaseBVoltage) }}</el-descriptions-item>
          <el-descriptions-item label="C相电压(V)">{{ formatValue(realData?.phaseCVoltage) }}</el-descriptions-item>
          
          <!-- 线电压 -->
          <el-descriptions-item label="AB线电压(V)">{{ formatValue(realData?.lineABVoltage) }}</el-descriptions-item>
          <el-descriptions-item label="BC线电压(V)">{{ formatValue(realData?.lineBCVoltage) }}</el-descriptions-item>
          <el-descriptions-item label="CA线电压(V)">{{ formatValue(realData?.lineCAVoltage) }}</el-descriptions-item>
          
          <!-- 相电流 -->
          <el-descriptions-item label="A相电流(A)">{{ formatValue(realData?.phaseACurrent) }}</el-descriptions-item>
          <el-descriptions-item label="B相电流(A)">{{ formatValue(realData?.phaseBCurrent) }}</el-descriptions-item>
          <el-descriptions-item label="C相电流(A)">{{ formatValue(realData?.phaseCCurrent) }}</el-descriptions-item>
          
          <!-- 功率因数和温度 -->
          <el-descriptions-item label="功率因数">{{ formatValue(realData?.powerFactor) }}</el-descriptions-item>
          <el-descriptions-item label="设备温度(℃)">{{ formatValue(realData?.deviceTemperature) }}</el-descriptions-item>
          <el-descriptions-item label="转换效率(%)">{{ formatValue(realData?.conversionEfficiency) }}</el-descriptions-item>
          
          <!-- 支路电压 -->
          <template v-for="i in 20" :key="`voltage-${i}`">
            <el-descriptions-item :label="`支路电压${i}(V)`">
              {{ formatValue(realData?.[`branchVoltage${i}`]) }}
            </el-descriptions-item>
          </template>
          
          <!-- 支路电流 -->
          <template v-for="i in 20" :key="`current-${i}`">
            <el-descriptions-item :label="`支路电流${i}(A)`">
              {{ formatValue(realData?.[`branchCurrent${i}`]) }}
            </el-descriptions-item>
          </template>
          
          <!-- MPPT -->
          <el-descriptions-item label="MPPT1">{{ formatValue(realData?.mppt1) }}</el-descriptions-item>
          <el-descriptions-item label="MPPT2">{{ formatValue(realData?.mppt2) }}</el-descriptions-item>
          <el-descriptions-item label="MPPT3">{{ formatValue(realData?.mppt3) }}</el-descriptions-item>
          <el-descriptions-item label="MPPT4">{{ formatValue(realData?.mppt4) }}</el-descriptions-item>
          
          <!-- 设备状态和其他 -->
          <el-descriptions-item label="设备状态码">{{ formatValue(realData?.deviceStatusCode) }}</el-descriptions-item>
          <el-descriptions-item label="绝缘阻抗">{{ formatValue(realData?.insulationImpedance) }}</el-descriptions-item>
          <el-descriptions-item label="开机时间">{{ formatValue(realData?.startupTime) }}</el-descriptions-item>
          <el-descriptions-item label="关机时间">{{ formatValue(realData?.shutdownTime) }}</el-descriptions-item>
          <el-descriptions-item label="A相温度">{{ formatValue(realData?.phaseATemperature) }}</el-descriptions-item>
          <el-descriptions-item label="B相温度">{{ formatValue(realData?.phaseBTemperature) }}</el-descriptions-item>
          <el-descriptions-item label="C相温度">{{ formatValue(realData?.phaseCTemperature) }}</el-descriptions-item>
        </el-descriptions>
      </div>
      
      <template #footer>
        <el-button @click="moreDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import {
  getAgvcNbqHisList,
  getNbqRealData
} from '@/api/agvc/agvcNbqHis'

import { ElMessage } from 'element-plus'
import { ref, onMounted, onUnmounted } from 'vue'

defineOptions({
  name: 'AgvcNbqMonitor'
})

const elSearchFormRef = ref()

// 表格控制部分
const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({})
const loading = ref(false)

// 实时刷新相关
const refreshInterval = ref(5000) // 默认5秒刷新
let refreshTimer = null

// 更多数据弹窗相关
const moreDialogVisible = ref(false)
const currentDevice = ref(null)
const realData = ref(null)
const moreDataLoading = ref(false)

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

// 显示更多数据
const showMoreData = async(row) => {
  currentDevice.value = row
  moreDialogVisible.value = true
  moreDataLoading.value = true
  
  try {
    const res = await getNbqRealData({ psid: row.psid, inverterNo: row.inverterNo })
    if (res.code === 0) {
      realData.value = res.data
    } else {
      ElMessage.error('获取实时数据失败: ' + res.msg)
    }
  } catch (error) {
    ElMessage.error('获取实时数据失败: ' + error.message)
  } finally {
    moreDataLoading.value = false
  }
}

// 启动自动刷新
const startAutoRefresh = () => {
  // 清除已存在的定时器
  if (refreshTimer) {
    clearInterval(refreshTimer)
  }
  // 设置新的定时器
  refreshTimer = setInterval(() => {
    getTableData()
  }, refreshInterval.value)
}

// 停止自动刷新
const stopAutoRefresh = () => {
  if (refreshTimer) {
    clearInterval(refreshTimer)
    refreshTimer = null
  }
}

onMounted(() => {
  getTableData()
  startAutoRefresh() // 启动自动刷新
})

// 组件卸载时停止刷新
onUnmounted(() => {
  stopAutoRefresh()
})
</script>

<style scoped>
.el-descriptions {
  margin-top: 20px;
}
</style>
