<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" class="demo-form-inline" @keyup.enter="onSubmit">
      <el-form-item label="并网点名称" prop="name">
        <el-input v-model="searchInfo.name" placeholder="请输入并网点名称" />
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
        <el-table
        ref="multipleTable"
        style="width: 100%"
        tooltip-effect="dark"
        :data="tableData"
        row-key="ID"
        >
        <el-table-column align="left" label="并网点名称" prop="name" width="200">
          <template #default="scope">
            {{ scope.row.name || '未命名并网点' }}
          </template>
        </el-table-column>

        <el-table-column align="left" label="操作" fixed="right" :min-width="appStore.operateMinWith">
            <template #default="scope">
            <el-button  type="primary" link class="table-button" @click="viewHistoryData(scope.row)"><el-icon style="margin-right: 5px"><InfoFilled /></el-icon>查看历史数据</el-button>
            <el-button  type="success" link class="table-button" @click="generateTestDataForGridPoint(scope.row)"><el-icon style="margin-right: 5px"><DataAnalysis /></el-icon>生成测试数据</el-button>
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
    <el-dialog v-model="historyDialogVisible" title="并网点历史数据" width="80%" destroy-on-close>
      <div v-if="currentGridPoint">
        <div class="mb-4">
          <span class="text-lg font-semibold">{{ currentGridPoint.name }} - 历史数据</span>
        </div>
        
        <!-- 时间范围选择 -->
        <div class="flex items-center mb-4">
          <span class="mr-2">时间范围:</span>
          <el-date-picker
            v-model="historyTimeRange"
            type="datetimerange"
            range-separator="至"
            start-placeholder="开始时间"
            end-placeholder="结束时间"
            value-format="YYYY-MM-DDTHH:mm:ss.000Z"
            class="!w-380px"
          />
          <el-button type="primary" icon="search" @click="loadHistoryData" class="ml-2">查询</el-button>
          <el-button type="success" icon="data-analysis" @click="generateTestDataForCurrentPoint" class="ml-2">生成测试数据</el-button>
        </div>
        
        <!-- 数据选择下拉框 -->
        <div class="flex items-center mb-4">
          <span class="mr-2">选择数据项:</span>
          <el-select v-model="selectedPointID" placeholder="请选择数据项" class="!w-200px">
            <el-option
              v-for="point in pointOptions"
              :key="point.pointID"
              :label="point.pointName"
              :value="point.pointID"
            />
          </el-select>
        </div>
        
        <!-- 折线图 -->
        <div id="historyChart" style="width: 100%; height: 400px;"></div>
      </div>
    </el-dialog>

  </div>
</template>

<script setup>
import {
  getSysGridConnectionPointHistoryList,
  getGridPointHistoryData,
  generateTestData
} from '@/api/monitor/sysGridConnectionPointHistory'

import { useAppStore } from '@/pinia/modules/app'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, onMounted, watch } from 'vue'
import * as echarts from 'echarts'

const appStore = useAppStore()

// 分页相关
const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({})

// 查询相关
const showAllQuery = ref(false)

// 历史数据弹窗相关
const historyDialogVisible = ref(false)
const currentGridPoint = ref(null)
const historyTimeRange = ref([])
const selectedPointID = ref('')
const pointOptions = ref([])
const historyChart = ref(null)
const allHistoryData = ref([]) // 存储所有历史数据

// 点标识配置 - 根据开关柜（保护装置）数据规范
const pointIdentifiers = {
  gridPoint1: [
    { pointID: 1, pointName: 'AB线电压Uab' },
    { pointID: 2, pointName: 'BC线电压Ubc' },
    { pointID: 3, pointName: 'CA线电压Uca' },
    { pointID: 4, pointName: 'A相电流Ia' },
    { pointID: 5, pointName: 'B相电流Ib' },
    { pointID: 6, pointName: 'C相电流Ic' },
    { pointID: 7, pointName: '有功功率P' },
    { pointID: 8, pointName: '无功功率Q' },
    { pointID: 9, pointName: 'PF(COS功率因数)' },
    { pointID: 10, pointName: 'F（频率）' },
    { pointID: 11, pointName: 'S（视在功率）' },
    { pointID: 12, pointName: 'A相电压Ua' },
    { pointID: 13, pointName: 'B相电压Ub' },
    { pointID: 14, pointName: 'C相电压Uc' },
    { pointID: 15, pointName: '零序电流（A）' },
    { pointID: 16, pointName: 'PA' },
    { pointID: 17, pointName: 'PB' },
    { pointID: 18, pointName: 'PC' },
    { pointID: 19, pointName: 'QA' },
    { pointID: 20, pointName: 'QB' }
  ],
  gridPoint2: [
    { pointID: 21, pointName: 'QC' },
    { pointID: 22, pointName: 'SA' },
    { pointID: 23, pointName: 'SB' },
    { pointID: 24, pointName: 'SC' },
    { pointID: 25, pointName: '零序电压（V）' },
    { pointID: 40, pointName: '正向有功' },
    { pointID: 41, pointName: '正向无功' },
    { pointID: 42, pointName: '反向有功' },
    { pointID: 43, pointName: '反向无功' }
  ]
}

// 根据并网点ID获取对应的点标识列表
const getPointIdentifiersByGridPointId = (gridPointId) => {
  // 确保与后端定义一致
  if (gridPointId === 1) {
    return pointIdentifiers.gridPoint1
  } else if (gridPointId === 2) {
    return pointIdentifiers.gridPoint2
  }
  return []
}

// 查询
const onSubmit = () => {
  page.value = 1
  pageSize.value = 10
  getTableData()
}

// 重置
const onReset = () => {
  searchInfo.value = {}
  showAllQuery.value = false
  onSubmit()
}

// 分页
const handleSizeChange = (val) => {
  pageSize.value = val
  getTableData()
}

const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

// 获取表格数据
const getTableData = async() => {
  try {
    const table = await getSysGridConnectionPointHistoryList()
    if (table.code === 0) {
      // 后端直接返回数组，不是分页格式
      tableData.value = table.data || []
      total.value = tableData.value.length
      // 重置分页信息
      if (tableData.value.length === 0) {
        page.value = 1
      }
    } else {
      ElMessage.error(table.msg || '获取数据失败')
      tableData.value = []
      total.value = 0
    }
  } catch (error) {
    console.error('获取表格数据失败:', error)
    ElMessage.error('获取数据失败，请检查网络连接')
    tableData.value = []
    total.value = 0
  }
}

// 查看历史数据
const viewHistoryData = (row) => {
  currentGridPoint.value = row
  historyDialogVisible.value = true
  // 初始化时间范围为最近24小时
  const endTime = new Date()
  const startTime = new Date(endTime.getTime() - 24 * 60 * 60 * 1000)
  historyTimeRange.value = [startTime.toISOString(), endTime.toISOString()]
  
  // 根据并网点ID获取对应的点标识列表
  pointOptions.value = getPointIdentifiersByGridPointId(row.ID)
  // 默认选中第一个数据项
  if (pointOptions.value.length > 0) {
    selectedPointID.value = pointOptions.value[0].pointID
  }
  
  // 延迟加载图表，确保DOM已渲染
  setTimeout(() => {
    initChart()
    loadHistoryData()
  }, 100)
}

// 初始化图表
const initChart = () => {
  const chartDom = document.getElementById('historyChart')
  if (chartDom) {
    // 销毁之前的实例
    if (historyChart.value) {
      historyChart.value.dispose()
    }
    // 创建新的实例
    historyChart.value = echarts.init(chartDom)
  }
}

// 加载历史数据
const loadHistoryData = async() => {
  if (!currentGridPoint.value || historyTimeRange.value.length !== 2) {
    ElMessage.warning('请选择有效的时间范围')
    return
  }
  
  const startTime = historyTimeRange.value[0]
  const endTime = historyTimeRange.value[1]
  
  console.log('发送请求参数:', {
    GridPointID: currentGridPoint.value.ID,
    StartTime: startTime,
    EndTime: endTime
  }) // 添加调试日志
  
  try {
    const res = await getGridPointHistoryData({
      GridPointID: currentGridPoint.value.ID,  // 修正参数名称，使用大写G
      StartTime: startTime,
      EndTime: endTime
    })
    
    console.log('服务器返回数据:', res) // 添加调试日志
    
    if (res.code === 0) {
      allHistoryData.value = res.data || []
      console.log('处理后的数据:', allHistoryData.value) // 添加调试日志
      
      if (allHistoryData.value.length === 0) {
        ElMessage.info('该时间范围内暂无数据，您可以生成测试数据进行查看')
      }
      // 渲染图表
      renderChart(allHistoryData.value)
    } else {
      ElMessage.error(res.msg || '获取历史数据失败')
    }
  } catch (error) {
    console.error('获取历史数据失败:', error)
    ElMessage.error('获取历史数据失败，请检查网络连接')
  }
}

// 渲染图表
const renderChart = (responseData) => {
  if (!historyChart.value) {
    ElMessage.warning('图表初始化失败，请重新打开弹窗')
    return
  }
  
  if (!responseData || responseData.length === 0) {
    // 显示空数据提示
    const option = {
      title: {
        text: '暂无数据',
        left: 'center',
        top: 'center',
        textStyle: {
          color: '#999',
          fontSize: 18
        }
      }
    }
    historyChart.value.setOption(option)
    return
  }
  
  // 如果只选择了一个数据项，则只显示该数据项的图表
  if (selectedPointID.value) {
    console.log('选中的点标识ID:', selectedPointID.value)
    console.log('响应数据结构:', responseData)

    const seriesData = responseData.find(series => series.pointID === selectedPointID.value)
    console.log('找到的系列数据:', seriesData)
    
    if (!seriesData || !seriesData.data || seriesData.data.length === 0) {
      console.log('数据为空，原因:', {
        seriesData: !!seriesData,
        hasData: seriesData ? !!seriesData.data : false,
        dataLength: seriesData && seriesData.data ? seriesData.data.length : 0
      }) // 添加调试日志
      
      // 显示空数据提示
      const option = {
        title: {
          text: '选中的数据项暂无数据',
          left: 'center',
          top: 'center',
          textStyle: {
            color: '#999',
            fontSize: 18
          }
        }
      }
      historyChart.value.setOption(option)
      return
    }
    
    // 准备图表数据
    const xAxisData = seriesData.data.map(item => { //time
      const date = new Date(item.time)
      return date.toLocaleString('zh-CN', {
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit'
      })
    })
    
    const seriesValues = seriesData.data.map(item => item.value) //value
    
    // 获取当前选中的数据项名称
    const currentPoint = pointOptions.value.find(p => p.pointID === selectedPointID.value)
    const chartTitle = currentPoint ? currentPoint.pointName : seriesData.pointName //pointName
    
    // 配置图表选项
    const option = {
      title: {
        text: chartTitle,
        left: 'center'
      },
      tooltip: {
        trigger: 'axis',
        formatter: function(params) {
          const date = new Date(seriesData.data[params[0].dataIndex].time)
          return `${date.toLocaleString('zh-CN')}<br/>${params[0].marker}${params[0].seriesName}: ${params[0].value}`
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
        data: xAxisData,
        axisLabel: {
          rotate: 45
        }
      },
      yAxis: {
        type: 'value'
      },
      series: [{
        name: chartTitle,
        type: 'line',
        data: seriesValues,
        smooth: true,
        lineStyle: {
          width: 2
        },
        itemStyle: {
          color: '#409EFF'
        }
      }]
    }
    
    // 渲染图表
    historyChart.value.setOption(option)
  } else {
    // 处理多个数据系列的情况（初始加载时）
    console.log('渲染多个数据系列:', responseData) // 添加调试日志
    
    // 提取所有时间点
    const allTimes = []
    responseData.forEach(series => {
      if (series.data && series.data.length > 0) {
        series.data.forEach(item => {
          if (!allTimes.includes(item.time)) { // 修复字段名：time 而不是 Time
            allTimes.push(item.time)
          }
        })
      }
    })
    allTimes.sort()
    
    if (allTimes.length === 0) {
      // 显示空数据提示
      const option = {
        title: {
          text: '所有数据项都暂无数据',
          left: 'center',
          top: 'center',
          textStyle: {
            color: '#999',
            fontSize: 18
          }
        }
      }
      historyChart.value.setOption(option)
      return
    }
    
    // 准备图表数据
    const xAxisData = allTimes.map(time => {
      const date = new Date(time)
      return date.toLocaleString('zh-CN', {
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit'
      })
    })
    
    // 准备多个系列的数据
    const seriesList = responseData.filter(series => series.data && series.data.length > 0).map(series => {
      // 获取数据项名称
      const pointInfo = pointOptions.value.find(p => p.pointID === series.pointID) // 修复字段名：pointID 而不是 PointID
      const seriesName = pointInfo ? pointInfo.pointName : series.pointName // 修复字段名：pointName 而不是 PointName
      
      // 创建数据数组，缺失数据用null填充
      const seriesData = allTimes.map(time => {
        const dataPoint = series.data.find(item => item.time === time) // 修复字段名：time 而不是 Time
        return dataPoint ? dataPoint.value : null // 修复字段名：value 而不是 Value
      })
      
      return {
        name: seriesName,
        type: 'line',
        data: seriesData,
        smooth: true,
        lineStyle: {
          width: 2
        },
        itemStyle: {
          color: '#409EFF'
        }
      }
    })
    
    if (seriesList.length === 0) {
      ElMessage.info('所有数据系列都暂无数据')
      return
    }
    
    // 配置图表选项
    const option = {
      title: {
        text: '并网点历史数据',
        left: 'center'
      },
      tooltip: {
        trigger: 'axis',
        formatter: function(params) {
          const date = new Date(allTimes[params[0].dataIndex])
          let tooltip = `${date.toLocaleString('zh-CN')}<br/>`
          params.forEach(param => {
            if (param.value !== null) {
              tooltip += `${param.marker}${param.seriesName}: ${param.value}<br/>`
            }
          })
          return tooltip
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
        data: xAxisData,
        axisLabel: {
          rotate: 45
        }
      },
      yAxis: {
        type: 'value'
      },
      series: seriesList
    }
    
    // 渲染图表
    historyChart.value.setOption(option)
  }
}

// 生成测试数据
const generateTestDataForGridPoint = async(row) => {
  ElMessageBox.confirm('确定要为该并网点生成测试数据吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async() => {
    try {
      const res = await generateTestData({ GridPointID: row.ID })  // 修正参数名称
      if (res.code === 0) {
        ElMessage.success('测试数据生成成功')
        // 如果当前正在查看该并网点的历史数据，则刷新图表
        if (currentGridPoint.value && currentGridPoint.value.ID === row.ID) {
          loadHistoryData()
        }
      } else {
        ElMessage.error(res.msg || '测试数据生成失败')
      }
    } catch (error) {
      console.error('生成测试数据失败:', error)
      ElMessage.error('生成测试数据失败，请检查网络连接')
    }
  }).catch(() => {
    // 用户取消操作
  })
}

// 为当前点生成测试数据
const generateTestDataForCurrentPoint = async() => {
  if (!currentGridPoint.value) return
  
  ElMessageBox.confirm('确定要为当前并网点生成测试数据吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async() => {
    try {
      const res = await generateTestData({ GridPointID: currentGridPoint.value.ID })  // 修正参数名称
      if (res.code === 0) {
        ElMessage.success('测试数据生成成功')
        // 刷新图表
        loadHistoryData()
      } else {
        ElMessage.error(res.msg || '测试数据生成失败')
      }
    } catch (error) {
      console.error('生成测试数据失败:', error)
      ElMessage.error('生成测试数据失败，请检查网络连接')
    }
  }).catch(() => {
    // 用户取消操作
  })
}

// 监听选中数据项的变化，重新渲染图表
watch(selectedPointID, () => {
  if (allHistoryData.value && allHistoryData.value.length > 0) {
    renderChart(allHistoryData.value)
  }
})

onMounted(() => {
  getTableData()
})
</script>

