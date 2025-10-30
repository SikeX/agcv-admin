
<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" class="demo-form-inline" @keyup.enter="onSubmit">
      <el-form-item label="气象仪名称" prop="name">
        <el-input v-model="searchInfo.name" placeholder="请输入气象仪名称" />
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
        <el-table-column align="left" label="气象仪名称" prop="name" width="200">
          <template #default="scope">
            {{ scope.row.name || '未命名气象仪' }}
          </template>
        </el-table-column>

        <el-table-column align="left" label="操作" fixed="right" :min-width="appStore.operateMinWith">
            <template #default="scope">
            <el-button  type="primary" link class="table-button" @click="viewHistoryData(scope.row)"><el-icon style="margin-right: 5px"><InfoFilled /></el-icon>查看历史数据</el-button>
            <el-button  type="success" link class="table-button" @click="generateTestDataForQixiangyi(scope.row)"><el-icon style="margin-right: 5px"><DataAnalysis /></el-icon>生成测试数据</el-button>
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
    <el-dialog v-model="historyDialogVisible" title="气象仪历史数据" width="80%" destroy-on-close>
      <div v-if="currentQixiangyi">
        <div class="mb-4">
          <span class="text-lg font-semibold">{{ currentQixiangyi.name }} - 历史数据</span>
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
  getSysQixiangyiHistoryList,
  getQixiangyiHistoryData,
  generateTestData
} from '@/api/monitor/sysQixiangyiHistory'

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
const currentQixiangyi = ref(null)
const historyTimeRange = ref([])
const selectedPointID = ref('')
const pointOptions = ref([])
const historyChart = ref(null)
const allHistoryData = ref([]) // 存储所有历史数据

// 点标识配置 - 根据气象仪数据规范，与后端保持一致
const pointIdentifiers = {
  qixiangyi1: [
    { pointID: 2, pointId: 2, pointName: '环境温度(℃)', unit: '℃' },
    { pointID: 3, pointId: 3, pointName: '露点温度(℃)', unit: '℃' },
    { pointID: 4, pointId: 4, pointName: '风速(m/s)', unit: 'm/s' },
    { pointID: 5, pointId: 5, pointName: '2分风速(m/s)', unit: 'm/s' },
    { pointID: 6, pointId: 6, pointName: '10分风速(m/s)', unit: 'm/s' },
    { pointID: 7, pointId: 7, pointName: '风向(°)', unit: '°' },
    { pointID: 8, pointId: 8, pointName: '水平辐射强度(W/㎡)', unit: 'W/㎡' },
    { pointID: 9, pointId: 9, pointName: '倾斜辐射强度(W/㎡)', unit: 'W/㎡' },
    { pointID: 11, pointId: 11, pointName: '组件温度(℃)', unit: '℃' },
    { pointID: 12, pointId: 12, pointName: '雨量(mm)', unit: 'mm' }
  ],
  qixiangyi2: [
    { pointID: 21, pointId: 21, pointName: '散射辐射强度(W/㎡)', unit: 'W/㎡' },
    { pointID: 22, pointId: 22, pointName: '直射辐射强度(W/㎡)', unit: 'W/㎡' },
    { pointID: 39, pointId: 39, pointName: '水平总辐射量(MJ/㎡)', unit: 'MJ/㎡' },
    { pointID: 40, pointId: 40, pointName: '倾斜总辐射量(MJ/㎡)', unit: 'MJ/㎡' },
    { pointID: 46, pointId: 46, pointName: '环境湿度(%RH)', unit: '%RH' },
    { pointID: 51, pointId: 51, pointName: '日照小时数(h)', unit: 'h' },
    { pointID: 52, pointId: 52, pointName: '气压(Pa)', unit: 'Pa' }
  ]
}

// 根据气象仪ID获取对应的点标识列表
const getPointIdentifiersByQixiangyiId = (qixiangyiId) => {
  console.log('获取点标识，气象仪ID:', qixiangyiId)
  // 根据具体的气象仪ID返回对应的点标识，而不是根据ID范围
  // 这里需要根据实际的气象仪配置来确定哪个ID对应哪组点标识
  // 假设气象仪ID 1 对应 qixiangyi1，气象仪ID 2 对应 qixiangyi2
  let result
  if (qixiangyiId === 1) {
    result = pointIdentifiers.qixiangyi1
    console.log('返回气象仪1的点标识:', result)
  } else if (qixiangyiId === 2) {
    result = pointIdentifiers.qixiangyi2
    console.log('返回气象仪2的点标识:', result)
  } else {
    // 如果是其他ID，可以根据实际情况添加更多条件
    // 或者使用默认的点标识
    result = []
    console.log('未知气象仪ID，返回空数组:', qixiangyiId)
  }
  return result
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
    const table = await getSysQixiangyiHistoryList()
    if (table.code === 0) {
      // 后端直接返回数组，不是分页格式
      tableData.value = table.data.list || []
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
  console.log('查看历史数据，行数据:', row)
  currentQixiangyi.value = row
  historyDialogVisible.value = true
  // 初始化时间范围为最近24小时
  const endTime = new Date()
  const startTime = new Date(endTime.getTime() - 24 * 60 * 60 * 1000)
  historyTimeRange.value = [startTime.toISOString(), endTime.toISOString()]
  
  // 根据气象仪ID获取对应的点标识列表
  // 兼容不同的字段名称：qixiangyiID, QixiangyiID, ID
  const qixiangyiId = row.qixiangyiID || row.QixiangyiID || row.ID
  console.log('使用的气象仪ID:', qixiangyiId)
  pointOptions.value = getPointIdentifiersByQixiangyiId(qixiangyiId)
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
  if (!currentQixiangyi.value || historyTimeRange.value.length !== 2) {
    ElMessage.warning('请选择有效的时间范围')
    return
  }
  
  const startTime = historyTimeRange.value[0]
  const endTime = historyTimeRange.value[1]
  
  // 兼容不同的字段名称：qixiangyiID, QixiangyiID, ID
  const qixiangyiId = currentQixiangyi.value.qixiangyiID || currentQixiangyi.value.QixiangyiID || currentQixiangyi.value.ID
  
  console.log('发送请求参数:', {
    qixiangyiId: qixiangyiId,
    startTime: startTime,
    endTime: endTime
  }) // 添加调试日志
  
  try {
    const res = await getQixiangyiHistoryData({
      qixiangyiId: qixiangyiId,
      startTime: startTime,
      endTime: endTime
    })
    
    console.log('服务器返回数据:', res) // 添加调试日志
    
    if (res.code === 0) {
      allHistoryData.value = res.data.data || []
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

    // 修复字段名匹配问题：服务器返回pointId，前端使用pointID
    const seriesData = responseData.find(series => 
      series.pointID === selectedPointID.value || series.pointId === selectedPointID.value
    )
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
    const xAxisData = seriesData.data.map(item => {
      const date = new Date(item.time)
      return date.toLocaleString('zh-CN', {
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit'
      })
    })
    
    const seriesValues = seriesData.data.map(item => item.value)
    
    // 获取当前选中的数据项名称，修复字段名匹配问题
    const pointId = seriesData.pointID || seriesData.pointId
    const currentPoint = pointOptions.value.find(p => p.pointID === pointId || p.pointId === pointId)
    const chartTitle = currentPoint ? currentPoint.pointName : (seriesData.pointName || `点标识${pointId}`)
    
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
          if (!allTimes.includes(item.time)) {
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
      // 获取数据项名称，修复字段名匹配问题
      const pointId = series.pointID || series.pointId
      const pointInfo = pointOptions.value.find(p => p.pointID === pointId || p.pointId === pointId)
      const seriesName = pointInfo ? pointInfo.pointName : (series.pointName || `点标识${pointId}`)
      
      // 创建数据数组，缺失数据用null填充
      const seriesData = allTimes.map(time => {
        const dataPoint = series.data.find(item => item.time === time)
        return dataPoint ? dataPoint.value : null
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
        text: '气象仪历史数据',
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
const generateTestDataForQixiangyi = async(row) => {
  console.log('生成测试数据，行数据:', row)
  ElMessageBox.confirm('确定要为该气象仪生成测试数据吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async() => {
    try {
      // 兼容不同的字段名称：qixiangyiID, QixiangyiID, ID
      const qixiangyiId = row.qixiangyiID || row.QixiangyiID || row.ID
      console.log('生成测试数据使用的气象仪ID:', qixiangyiId)
      const res = await generateTestData({ qixiangyiId: qixiangyiId })
      if (res.code === 0) {
        ElMessage.success('测试数据生成成功')
        // 如果当前正在查看该气象仪的历史数据，则刷新图表
        const currentId = currentQixiangyi.value && (currentQixiangyi.value.qixiangyiID || currentQixiangyi.value.QixiangyiID || currentQixiangyi.value.ID)
        if (currentId === qixiangyiId) {
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
  if (!currentQixiangyi.value) return
  
  ElMessageBox.confirm('确定要为当前气象仪生成测试数据吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async() => {
    try {
      // 兼容不同的字段名称：qixiangyiID, QixiangyiID, ID
      const qixiangyiId = currentQixiangyi.value.qixiangyiID || currentQixiangyi.value.QixiangyiID || currentQixiangyi.value.ID
      console.log('为当前点生成测试数据使用的气象仪ID:', qixiangyiId)
      const res = await generateTestData({ qixiangyiId: qixiangyiId })
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

<style>
</style>
