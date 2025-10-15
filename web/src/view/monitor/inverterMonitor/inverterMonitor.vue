
<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" class="demo-form-inline" @keyup.enter="onSubmit">
            <el-form-item label="逆变器编号" prop="inverter_no">
  <el-input v-model="searchInfo.inverter_no" placeholder="搜索条件" />
</el-form-item>
            
            <el-form-item label="逆变器名称" prop="name">
  <el-input v-model="searchInfo.name" placeholder="搜索条件" />
</el-form-item>
            
            <el-form-item label="运行状态" prop="status">
  <el-select v-model="searchInfo.status" clearable filterable placeholder="请选择" @clear="()=>{searchInfo.status=undefined}">
    <el-option v-for="(item,key) in nb_statusOptions" :key="key" :label="item.label" :value="item.value" />
  </el-select>
</el-form-item>
            

        <template v-if="showAllQuery">
          <!-- 将需要控制显示状态的查询条件添加到此范围内 -->
        </template>

        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
          <el-button link type="primary" icon="arrow-down" @click="showAllQuery=true" v-if="!showAllQuery">展开</el-button>
          <el-button link type="primary" icon="arrow-up" @click="showAllQuery=false" v-else>收起</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
        <div class="gva-btn-list">
            <el-button  type="primary" icon="plus" @click="openDialog()">新增</el-button>
            <el-button  icon="delete" style="margin-left: 10px;" :disabled="!multipleSelection.length" @click="onDelete">删除</el-button>
            
        </div>
        <el-table
        ref="multipleTable"
        style="width: 100%"
        tooltip-effect="dark"
        :data="tableData"
        row-key="ID"
        @selection-change="handleSelectionChange"
        >
        <el-table-column type="selection" width="55" />
        
            <el-table-column align="left" label="逆变器编号" prop="inverter_no" width="120" />

            <el-table-column align="left" label="逆变器名称" prop="name" width="120" />

            <el-table-column align="left" label="运行状态" prop="status" width="120">
    <template #default="scope">
    {{ filterDict(scope.row.status,nb_statusOptions) }}
    </template>
</el-table-column>
            <el-table-column align="left" label="是否调节" prop="isParticipateAdjust" width="120">
    <template #default="scope">{{ formatBoolean(scope.row.isParticipateAdjust) }}</template>
</el-table-column>
            <el-table-column align="left" label="额定功率(kW)" prop="ratedPower" width="120" />

            <el-table-column align="left" label="有功功率" prop="ygPreal" width="120" />

            <el-table-column align="left" label="有功目标(kW)" prop="apTargetValue" width="120" />

            <el-table-column align="left" label="无功功率(kVar)" prop="wgPreal" width="120" />

            <el-table-column align="left" label="无功目标(kVar)" prop="rpTargetValue" width="120" />

            <el-table-column align="left" label="功率因数" prop="pf" width="120" />

            <!-- 新增的历史数据列 -->
            <el-table-column align="left" label="历史数据" fixed="right" :min-width="appStore.operateMinWith">
              <template #default="scope">
                <el-button type="primary" link @click="showHistoryData(scope.row)">查询</el-button>
              </template>
            </el-table-column>

            <!-- <el-table-column align="left" label="操作" fixed="right" :min-width="appStore.operateMinWith">
            <template #default="scope">
            <el-button  type="primary" link class="table-button" @click="getDetails(scope.row)"><el-icon style="margin-right: 5px"><InfoFilled /></el-icon>查看</el-button>
            <el-button  type="primary" link icon="edit" class="table-button" @click="updateInverterMonitorFunc(scope.row)">编辑</el-button>
            <el-button   type="primary" link icon="delete" @click="deleteRow(scope.row)">删除</el-button>
            </template>
        </el-table-column> -->
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
    <el-drawer destroy-on-close :size="appStore.drawerSize" v-model="dialogFormVisible" :show-close="false" :before-close="closeDialog">
       <template #header>
              <div class="flex justify-between items-center">
                <span class="text-lg">{{type==='create'?'新增':'编辑'}}</span>
                <div>
                  <el-button :loading="btnLoading" type="primary" @click="enterDialog">确 定</el-button>
                  <el-button @click="closeDialog">取 消</el-button>
                </div>
              </div>
            </template>

          <el-form :model="formData" label-position="top" ref="elFormRef" :rules="rule" label-width="80px">
          </el-form>
    </el-drawer>

    <el-drawer destroy-on-close :size="appStore.drawerSize" v-model="detailShow" :show-close="true" :before-close="closeDetailShow" title="查看">
            <el-descriptions :column="1" border>
            </el-descriptions>
        </el-drawer>

    <!-- 历史数据弹窗 -->
    <el-dialog v-model="historyDialogVisible" title="历史数据" width="80%" top="5vh">
      <el-row :gutter="20">
        <!-- 日期选择区域 -->
        <el-col :span="6">
          <el-card>
            <div slot="header" class="clearfix">
              <span>时间选择</span>
            </div>
            <el-date-picker
              v-model="historyDateRange"
              type="daterange"
              range-separator="至"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
              value-format="YYYY-MM-DD"
              style="width: 100%; margin-bottom: 20px;"
            />
           <el-select v-model="selectedMetrics" multiple placeholder="请选择要显示的指标" style="width: 100%; margin-bottom: 20px;">
              <el-option label="有功功率" value="ygPreal" />
              <el-option label="有功目标(kW)" value="apTargetValue" />
              <el-option label="无功功率(kVar)" value="wgPreal" />
              <el-option label="无功目标(kVar)" value="rpTargetValue" />
              <el-option label="功率因数" value="pf" />
            </el-select>
            <el-button type="primary" @click="loadHistoryData" style="width: 100%;">查询</el-button>
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
        <span class="dialog-footer">
          <el-button @click="historyDialogVisible = false">关闭</el-button>
        </span>
      </template>
    </el-dialog>

  </div>
</template>

<script setup>
import {
  createInverterMonitor,
  deleteInverterMonitor,
  deleteInverterMonitorByIds,
  updateInverterMonitor,
  findInverterMonitor,
  getInverterMonitorList,
  getInverterHistory
} from '@/api/monitor/inverterMonitor'

import {
  findSysInverterSettingByInverterNo
} from '@/api/setting/sysInverterSetting'

// 全量引入格式化工具 请按需保留
import { getDictFunc, formatDate, formatBoolean, filterDict ,filterDataSource, returnArrImg, onDownloadFile } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive, nextTick, onMounted, onUnmounted, watch } from 'vue'
import { useAppStore } from "@/pinia"
import * as echarts from 'echarts'

defineOptions({
    name: 'InverterMonitor'
})

// 提交按钮loading
const btnLoading = ref(false)
const appStore = useAppStore()

// 控制更多查询条件显示/隐藏状态
const showAllQuery = ref(false)

// 自动化生成的字典（可能为空）以及字段
const nb_statusOptions = ref([])
const formData = ref({
        })

// 历史数据相关变量
const historyDialogVisible = ref(false)
const historyData = ref([])
const historyDateRange = ref([])
const selectedMetrics = ref(['ygPreal'])
const chartContainer = ref(null)
const currentInverter = ref(null)

// 图表实例
let chartInstance = null

// 验证规则
const rule = reactive({
})

const elFormRef = ref()
const elSearchFormRef = ref()

// =========== 表格控制部分 ===========
const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({})
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
    if (searchInfo.value.isParticipateAdjust === ""){
        searchInfo.value.isParticipateAdjust=null
    }
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
  const table = await getInverterMonitorList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
}

getTableData()

// ============== 表格控制部分结束 ===============

// 显示历史数据
const showHistoryData = async (row) => {
  // 保存当前逆变器信息
  currentInverter.value = row
  
  // 先根据逆变器编号获取设置信息
  try {
    const res = await findSysInverterSettingByInverterNo({ inverterNo: row.inverter_no })
    if (res.code === 0) {
      // 可以在这里使用设置信息
      console.log('逆变器设置:', res.data)
    }
  } catch (error) {
    ElMessage.warning('获取逆变器设置信息失败')
  }
  
  // 设置默认日期范围为最近7天
  const end = new Date()
  const start = new Date()
  start.setTime(start.getTime() - 3600 * 1000 * 24 * 7)
  historyDateRange.value = [start, end]
  
  // 加载历史数据
  loadHistoryData()
  
  // 显示弹窗
  historyDialogVisible.value = true
  
  // 在下次DOM更新后初始化图表
  nextTick(() => {
    initChart()
  })
}

// 加载历史数据
const loadHistoryData = async () => {
  if (!currentInverter.value) return
  
  try {
    const params = {
      inverterNo: currentInverter.value.inverter_no,
      startTime: historyDateRange.value[0],
      endTime: historyDateRange.value[1]
    }
    
    const res = await getInverterHistory(params)
    if (res.code === 0) {
      // 处理从InfluxDB获取的历史数据
      historyData.value = res.data.map(item => {
        return {
          timestamp: new Date(item._time),
          code: item.code,
          value: item.value,
          ygPreal: item.ygPreal || 0,
          wgPreal: item.wgPreal || 0,
          pf: item.pf || 0,
        }
      })
      //按时间排序
      historyData.value.sort((a, b) => a.timestamp - b.timestamp)
      
      // 更新图表
      updateChart()
    } else {
      ElMessage.error('获取历史数据失败: ' + res.msg)
    }
  } catch (error) {
    ElMessage.error('获取历史数据失败: ' + error.message)
  }
}

// 初始化图表
const initChart = () => {
  if (!chartContainer.value) return
  
  // 使用 echarts 绘制图表
  if (!chartInstance) {
    chartInstance = echarts.init(chartContainer.value)
  }
  
  updateChart()
}

// 更新图表
const updateChart = () => {
  if (!chartInstance || !historyData.value.length) return
  
  // 准备x轴数据（时间）
  const xData = historyData.value.map(item => 
    item.timestamp.toLocaleString('zh-CN', {
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit'
    })
  )
  
  // 根据选中的指标准备y轴数据
  const series = []
  const metricLabels = {
    ygPreal: '有功功率',
    wgPreal: '无功功率(kVar)',
    rpTargetValue: '无功目标(kVar)',
    pf: '功率因数'
  }
  
  selectedMetrics.value.forEach(metric => {
    series.push({
      name: metricLabels[metric],
      type: 'line',
      data: historyData.value.map(item => item[metric]),
      smooth: true
    })
  })
  
  // 配置图表选项
  const option = {
    title: {
      text: '历史数据趋势图'
    },
    tooltip: {
      trigger: 'axis'
    },
    legend: {
      data: selectedMetrics.value.map(m => metricLabels[m])
    },
    xAxis: {
      type: 'category',
      data: xData
    },
    yAxis: {
      type: 'value'
    },
    series: series,
    dataZoom: [
      {
        type: 'inside',
        start: 0,
        end: 100
      },
      {
        start: 0,
        end: 100
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

// 监听选中指标变化，更新图表
watch(selectedMetrics, () => {
  updateChart()
})

// 获取需要的字典 可能为空 按需保留
const setOptions = async () =>{
    nb_statusOptions.value = await getDictFunc('nb_status')
}

// 获取需要的字典 可能为空 按需保留
setOptions()


// 多选数据
const multipleSelection = ref([])
// 多选
const handleSelectionChange = (val) => {
    multipleSelection.value = val
}

// 删除行
const deleteRow = (row) => {
    ElMessageBox.confirm('确定要删除吗?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
    }).then(() => {
            deleteInverterMonitorFunc(row)
        })
    }

// 多选删除
const onDelete = async() => {
  ElMessageBox.confirm('确定要删除吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async() => {
      const IDs = []
      if (multipleSelection.value.length === 0) {
        ElMessage({
          type: 'warning',
          message: '请选择要删除的数据'
        })
        return
      }
      multipleSelection.value &&
        multipleSelection.value.map(item => {
          IDs.push(item.ID)
        })
      const res = await deleteInverterMonitorByIds({ IDs })
      if (res.code === 0) {
        ElMessage({
          type: 'success',
          message: '删除成功'
        })
        if (tableData.value.length === IDs.length && page.value > 1) {
          page.value--
        }
        getTableData()
      }
      })
    }

// 行为控制标记（弹窗内部需要增还是改）
const type = ref('')

// 更新行
const updateInverterMonitorFunc = async(row) => {
    const res = await findInverterMonitor({ ID: row.ID })
    type.value = 'update'
    if (res.code === 0) {
        formData.value = res.data
        dialogFormVisible.value = true
    }
}


// 删除行
const deleteInverterMonitorFunc = async (row) => {
    const res = await deleteInverterMonitor({ ID: row.ID })
    if (res.code === 0) {
        ElMessage({
                type: 'success',
                message: '删除成功'
            })
            if (tableData.value.length === 1 && page.value > 1) {
            page.value--
        }
        getTableData()
    }
}

// 弹窗控制标记
const dialogFormVisible = ref(false)

// 打开弹窗
const openDialog = () => {
    type.value = 'create'
    dialogFormVisible.value = true
}

// 关闭弹窗
const closeDialog = () => {
    dialogFormVisible.value = false
    formData.value = {
        }
}
// 弹窗确定
const enterDialog = async () => {
     btnLoading.value = true
     elFormRef.value?.validate( async (valid) => {
             if (!valid) return btnLoading.value = false
              let res
              switch (type.value) {
                case 'create':
                  res = await createInverterMonitor(formData.value)
                  break
                case 'update':
                  res = await updateInverterMonitor(formData.value)
                  break
                default:
                  res = await createInverterMonitor(formData.value)
                  break
              }
              btnLoading.value = false
              if (res.code === 0) {
                ElMessage({
                  type: 'success',
                  message: '创建/更改成功'
                })
                closeDialog()
                getTableData()
              }
      })
}

const detailForm = ref({})

// 查看详情控制标记
const detailShow = ref(false)


// 打开详情弹窗
const openDetailShow = () => {
  detailShow.value = true
}


// 打开详情
const getDetails = async (row) => {
  // 打开弹窗
  const res = await findInverterMonitor({ ID: row.ID })
  if (res.code === 0) {
    detailForm.value = res.data
    openDetailShow()
  }
}


// 关闭详情弹窗
const closeDetailShow = () => {
  detailShow.value = false
  detailForm.value = {}
}
</script>