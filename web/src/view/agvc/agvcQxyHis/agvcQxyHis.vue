
<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" class="demo-form-inline" @keyup.enter="onSubmit">
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
            <el-button type="success" @click="generateTestData">生成测试数据</el-button>
        </div>
        <el-table
        ref="multipleTable"
        style="width: 100%"
        tooltip-effect="dark"
        :data="tableData"
        row-key="ID"
        >
        
            <el-table-column align="left" label="设备编号" prop="number" min-width="120" />

            <el-table-column align="left" label="设备名称" prop="name" min-width="150" />

            <!-- 历史数据列 -->
            <el-table-column align="left" label="历史数据" fixed="right" width="100">
              <template #default="scope">
                <el-button type="primary" link @click="showHistoryData(scope.row)">查询</el-button>
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
              type="datetimerange"
              range-separator="至"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
              style="width: 100%; margin-bottom: 20px;"
            />
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
  createAgvcQxyHis,
  deleteAgvcQxyHis,
  deleteAgvcQxyHisByIds,
  updateAgvcQxyHis,
  findAgvcQxyHis,
  getAgvcQxyHisList,
  getAgvcQxyHistory,
  generateAgvcQxyTestData
} from '@/api/agvc/agvcQxyHis'

// 全量引入格式化工具 请按需保留
import { getDictFunc, formatDate, formatBoolean, filterDict ,filterDataSource, returnArrImg, onDownloadFile } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive, nextTick, onMounted, onUnmounted, watch } from 'vue'
import { useAppStore } from "@/pinia"
import * as echarts from 'echarts'




defineOptions({
    name: 'AgvcQxyHis'
})

// 提交按钮loading
const btnLoading = ref(false)
const appStore = useAppStore()

// 控制更多查询条件显示/隐藏状态
const showAllQuery = ref(false)

// 自动化生成的字典(可能为空)以及字段
const formData = ref({
            name: '',
        })



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
  const table = await getAgvcQxyHisList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
}

getTableData()

// ============== 表格控制部分结束 ===============

// 获取需要的字典 可能为空 按需保留
const setOptions = async () =>{
}

// 获取需要的字典 可能为空 按需保留
setOptions()

// 历史数据相关变量
const historyDialogVisible = ref(false)
const historyData = ref([])
const historyDateRange = ref([])
const chartContainer = ref(null)
const currentDevice = ref(null)

// 图表实例
let chartInstance = null

// 生成测试数据
const generateTestData = async () => {
  ElMessageBox.prompt('请输入设备编号', '生成测试数据', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    inputPlaceholder: '请输入设备编号',
  }).then(async ({ value }) => {
    if (!value) {
      ElMessage.warning('请输入设备编号')
      return
    }
    const res = await generateAgvcQxyTestData({ eqid: value, count: 100 })
    if (res.code === 0) {
      ElMessage.success('测试数据生成成功')
    } else {
      ElMessage.error('生成失败: ' + res.msg)
    }
  }).catch(() => {})
}

// 显示历史数据
const showHistoryData = async (row) => {
  // 保存当前设备信息
  currentDevice.value = row
  
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
  if (!currentDevice.value) return
  
  // 检查时间范围是否有效
  if (!historyDateRange.value || historyDateRange.value.length !== 2) {
    ElMessage.warning('请选择时间范围')
    return
  }
  
  try {
    // 格式化时间为RFC3339格式
    const formatToRFC3339 = (date) => {
      if (typeof date === 'string') {
        date = new Date(date)
      }
      return date.toISOString()
    }
    
    const params = {
      eqid: currentDevice.value.number,  // 使用设备编号
      startTime: formatToRFC3339(historyDateRange.value[0]),
      endTime: formatToRFC3339(historyDateRange.value[1])
    }
    
    const res = await getAgvcQxyHistory(params)
    if (res.code === 0) {
      // 处理从 InfluxDB获取的历史数据
      if (res.data && res.data.length > 0) {
        historyData.value = res.data.map(item => {
          return {
            timestamp: new Date(item.time),
            point: item.point,
            pointName: item.pointName || `点号${item.point}`, // 使用后端返回的点位名称
            value: item.value,
          }
        })
        // 按时间排序
        historyData.value.sort((a, b) => a.timestamp - b.timestamp)
      } else {
        // 没有数据时清空
        historyData.value = []
        ElMessage.info('选择的时间范围内没有历史数据')
      }
      
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
  
  // 按点号分组数据
  const dataByPoint = {}
  const pointNames = {} // 存储点位名称
  
  historyData.value.forEach(item => {
    const point = item.point || 'unknown'
    if (!dataByPoint[point]) {
      dataByPoint[point] = []
      pointNames[point] = item.pointName || `点号${point}`
    }
    dataByPoint[point].push(item)
  })
  
  // 获取所有唯一的时间点(去重)
  const timeSet = new Set()
  historyData.value.forEach(item => {
    const timeStr = item.timestamp.toLocaleString('zh-CN', {
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit'
    })
    timeSet.add(timeStr)
  })
  const xData = Array.from(timeSet).sort()
  
  // 根据点号生成系列
  const series = []
  
  Object.keys(dataByPoint).forEach(point => {
    // 为每个点位创建完整的时间序列数据
    const pointData = dataByPoint[point]
    const dataMap = new Map()
    pointData.forEach(item => {
      const timeStr = item.timestamp.toLocaleString('zh-CN', {
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit'
      })
      dataMap.set(timeStr, item.value)
    })
    
    // 按x轴顺序填充数据
    const seriesData = xData.map(time => dataMap.get(time) ?? null)
    
    series.push({
      name: pointNames[point],
      type: 'line',
      data: seriesData,
      smooth: true,
      connectNulls: false // 不连接空值
    })
  })
  
  // 配置图表选项
  const option = {
    title: {
      text: '历史数据趋势图',
      left: 'center',
      top: 10
    },
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'cross'
      }
    },
    legend: {
      orient: 'vertical',  // 垂直方向
      left: 10,            // 左侧对齐
      top: 50,             // 从顶部50px开始
      type: 'scroll',      // 支持滚动
      data: Object.keys(dataByPoint).map(p => pointNames[p]),
      textStyle: {
        fontSize: 12
      },
      pageIconSize: 12,
      width: 180,          // 图例宽度
      height: 350          // 图例高度
    },
    grid: {
      left: 210,           // 留出图例的空间
      right: '4%',
      top: 60,
      bottom: 80,          // 留出dataZoom的空间
      containLabel: true
    },
    xAxis: {
      type: 'category',
      data: xData,
      boundaryGap: false,
      axisLabel: {
        rotate: 45,        // 旋转标签避免重叠
        fontSize: 11
      }
    },
    yAxis: {
      type: 'value',
      axisLabel: {
        fontSize: 11
      }
    },
    series: series,
    dataZoom: [
      {
        type: 'inside',    // 鼠标滚轮缩放
        start: 0,
        end: 100
      },
      {
        type: 'slider',    // 滑块缩放
        start: 0,
        end: 100,
        height: 25,        // 滑块高度
        bottom: 15,        // 距离底部距离
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

// 监听窗口大小变化,重置图表大小
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

<style>
</style>
