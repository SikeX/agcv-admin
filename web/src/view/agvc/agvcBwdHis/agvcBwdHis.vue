
<template>
  <div>
    <!-- 设备选择卡片 -->
    <el-card class="device-selector-card" shadow="hover">
      <div class="device-info">
        <div class="device-name-row">
          <el-dropdown 
            trigger="click"
            size="large"
            @command="handleDeviceCommand"
          >
            <el-button type="primary" size="large" style="width: 300px;">
              {{ getCurrentDevice()?.name || '请选择并网点设备' }}
              <el-icon class="el-icon--right"><ArrowDown /></el-icon>
            </el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item
                  v-for="device in deviceList"
                  :key="device.number"
                  :command="device.number"
                >
                  {{ device.name }}
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
        <div class="device-details">
          <span class="detail-item">
            <span class="detail-label">并网点容量：</span>
            <span class="detail-value">{{ getCurrentDevice()?.capacity || '未设置' }} kW</span>
          </span>
          <span class="detail-item">
            <span class="detail-label">并网日期：</span>
            <span class="detail-value">{{ getCurrentDevice()?.gridConnectionDate || '未设置' }}</span>
          </span>
        </div>
      </div>
    </el-card>

    <!-- AGC/AVC控制标签页 -->
    <div class="main-tabs" style="margin-top: 20px;">
      <el-tabs v-model="controlActiveTab" type="border-card" style="position: relative;">
        <!-- 按钮组 - 绝对定位放在右上角 -->
        <div style="position: absolute; right: 20px; top: 8px; display: flex; gap: 10px; z-index: 100;" @click.stop>
          <!-- <el-button size="small" @click="showHistoryDataDialog">历史数据查看</el-button> -->
          <!-- <el-button size="small" @click="generateTestData">生成测试数据</el-button> -->
        </div>
        <!-- AGC控制标签页 -->
        <el-tab-pane label="AGC控制" name="agc">
          <div class="control-content">
            <!-- AGC折线图 - 放在最上面 -->
            <div class="chart-section chart-top">
              <div class="section-title">电站出力曲线图</div>
              <!-- PowerChart组件渲染 -->
              <div style="margin-top: 20px;">
                <PowerChart :psid="currentPsid" :eqid="mainActiveTab" />
              </div>
            </div>
                
                <!-- AGC控制面板和参数设置 - 左右布局 -->
                <div class="control-section">
                  <el-row :gutter="20">
                    <!-- 左侧：AGC控制面板 -->
                    <el-col :span="16">
                      <div class="section-title">AGC控制面板</div>
                      <div class="control-grid">
                        <!-- 电站AGC功能投退 -->
                        <div class="control-item">
                          <div class="item-label">电站AGC功能投退</div>
                          <div class="radio-group">
                            <el-radio-group v-model="agcFunctionState" size="small">
                              <el-radio :label="1">投入</el-radio>
                              <el-radio :label="0">退出</el-radio>
                            </el-radio-group>
                          </div>
                        </div>
                        
                        <!-- 电站AGC调节方式 -->
                        <div class="control-item">
                          <div class="item-label">电站AGC调节方式</div>
                          <div class="radio-group">
                            <el-radio-group v-model="agcControlMode" size="small">
                              <el-radio :label="1">闭环指导</el-radio>
                              <el-radio :label="0">开环指导</el-radio>
                            </el-radio-group>
                          </div>
                        </div>
                        
                        <!-- 电站AGC控制权限 -->
                        <div class="control-item">
                          <div class="item-label">电站AGC控制权限</div>
                          <div class="radio-group">
                            <el-radio-group v-model="agcControlAuthority" size="small">
                              <el-radio :label="1">调度控制</el-radio>
                              <el-radio :label="0">站内控制</el-radio>
                            </el-radio-group>
                          </div>
                        </div>
                      </div>
                    </el-col>
                    
                    <!-- 右侧：AGC参数设置 -->
                    <el-col :span="8">
                      <div class="section-title">AGC参数设置</div>
                      <div class="control-row" style="margin-top: 20px;">
                        <el-button size="default" @click="openAgcParametersDialog" style="width: 100%; margin-bottom: 10px;">参数设置</el-button>
                        <el-button size="default" @click="openAgcPlanCurvesDialog" style="width: 100%;">计划曲线</el-button>
                      </div>
                    </el-col>
                  </el-row>
                </div>
                
                <!-- 电站负荷 -->
                <div class="control-section">
                  <div class="section-title">电站负荷</div>
                  <div class="control-grid load-grid">
                    <!-- 目标有功 - 改为只读 -->
                    <div class="control-item">
                      <div class="item-label">目标有功(kW)</div>
                      <div class="item-value readonly">{{ formatNumber(targetActivePower) }}</div>
                    </div>
                    
                    <!-- 当前有功 -->
                    <div class="control-item">
                      <div class="item-label">当前有功(kW)</div>
                      <div class="item-value readonly">{{ formatNumber(currentActivePower) }}</div>
                    </div>
                    
                    <!-- 可调上限 -->
                    <div class="control-item">
                      <div class="item-label">可调上限(kW)</div>
                      <div class="item-value readonly">{{ formatNumber(agcUpperLimit) }}</div>
                    </div>
                    
                    <!-- 可调下限 -->
                    <div class="control-item">
                      <div class="item-label">可调下限(kW)</div>
                      <div class="item-value readonly">{{ formatNumber(agcLowerLimit) }}</div>
                    </div>
                    
                    <!-- 系统频率 -->
                    <div class="control-item">
                      <div class="item-label">系统频率(Hz)</div>
                      <div class="item-value readonly">{{ formatNumber(systemFrequency) }}</div>
                    </div>
                  </div>
                </div>
              </div>
            </el-tab-pane>
            
            <!-- AVC控制标签页 -->
            <el-tab-pane label="AVC控制" name="avc">
            <div class="control-content">
              <!-- AVC折线图 - 放在最上面 -->
              <div class="chart-section chart-top">
                <div class="section-title">AVC曲线图</div>
                <!-- VoltageReactiveChart组件渲染 -->
                <div style="margin-top: 20px;">
                  <VoltageReactiveChart :psid="currentPsid" :eqid="mainActiveTab" />
                </div>
              </div>
                
                <!-- AVC控制面板、模式切换、参数设置 - 三列布局 -->
                <div class="control-section">
                  <el-row :gutter="20">
                    <!-- 左侧：AVC控制面板 -->
                    <el-col :span="10">
                      <div class="section-title">AVC控制面板</div>
                      <div class="control-grid">
                        <!-- 电站AVC功能投退 -->
                        <div class="control-item">
                          <div class="item-label">电站AVC功能投退</div>
                          <div class="radio-group">
                            <el-radio-group v-model="avcFunctionState" size="small">
                              <el-radio :label="1">投入</el-radio>
                              <el-radio :label="0">退出</el-radio>
                            </el-radio-group>
                          </div>
                        </div>
                        
                        <!-- 电站AVC调节方式 -->
                        <div class="control-item">
                          <div class="item-label">电站AVC调节方式</div>
                          <div class="radio-group">
                            <el-radio-group v-model="avcControlMode" size="small">
                              <el-radio :label="1">开环指导</el-radio>
                              <el-radio :label="0">闭环调节</el-radio>
                            </el-radio-group>
                          </div>
                        </div>
                        
                        <!-- 电站AVC控制权限 -->
                        <div class="control-item">
                          <div class="item-label">电站AVC控制权限</div>
                          <div class="radio-group">
                            <el-radio-group v-model="avcControlAuthority" size="small">
                              <el-radio :label="1">站内控制</el-radio>
                              <el-radio :label="0">调度控制</el-radio>
                            </el-radio-group>
                          </div>
                        </div>
                      </div>
                    </el-col>
                    
                    <!-- 中间：AVC模式切换 -->
                    <el-col :span="7">
                      <div class="section-title">AVC模式切换</div>
                      <div class="avc-mode-switch" style="margin-top: 20px; text-align: center;">
                        <el-radio-group v-model="avcMode" @change="onAvcModeChange">
                          <el-radio-button label="voltage">电压模式</el-radio-button>
                          <el-radio-button label="reactive">无功模式</el-radio-button>
                        </el-radio-group>
                      </div>
                    </el-col>
                    
                    <!-- 右侧：AVC参数设置 -->
                    <el-col :span="7">
                      <div class="section-title">AVC参数设置</div>
                      <div class="control-row" style="margin-top: 20px;">
                        <el-button size="default" @click="openAvcParametersDialog" style="width: 100%; margin-bottom: 10px;">参数设置</el-button>
                        <el-button size="default" @click="openAvcPlanCurvesDialog" style="width: 100%;">计划曲线</el-button>
                      </div>
                    </el-col>
                  </el-row>
                </div>
                
                <!-- 电站负荷 -->
                <div class="control-section">
                  <div class="section-title">电站负荷</div>
                  <div class="control-grid load-grid">
                    <!-- 目标电压 -->
                    <div class="control-item">
                      <div class="item-label">目标电压(kV)</div>
                      <div class="item-value readonly">{{ formatNumber(targetVoltage) }}</div>
                    </div>
                    
                    <!-- 当前电压 -->
                    <div class="control-item">
                      <div class="item-label">当前电压(kV)</div>
                      <div class="item-value readonly">{{ formatNumber(currentVoltage) }}</div>
                    </div>
                    
                    <!-- 目标无功 -->
                    <div class="control-item">
                      <div class="item-label">目标无功(kVar)</div>
                      <div class="item-value readonly">{{ formatNumber(targetReactive) }}</div>
                    </div>
                    
                    <!-- 当前无功 -->
                    <div class="control-item">
                      <div class="item-label">当前无功(kVar)</div>
                      <div class="item-value readonly">{{ formatNumber(currentReactive) }}</div>
                    </div>
                    
                    <!-- 系统阻抗 -->
                    <div class="control-item">
                      <div class="item-label">系统阻抗</div>
                      <div class="item-value readonly">{{ formatNumber(systemImpedance) }}</div>
                    </div>
                  </div>
                </div>
            </div>
          </el-tab-pane>
        </el-tabs>
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
              value-format="YYYY-MM-DD"
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
    
    <!-- AGC参数设置弹窗 - 改为可编辑 -->
    <el-dialog v-model="agcParametersDialogVisible" title="AGC参数设置" width="600px">
      <el-form :model="agcParametersForm" label-width="140px" label-position="left">
        <el-form-item label="AGC功能退出模式">
          <el-input-number v-model="agcParametersForm.agcFunctionExit" :step="1" style="width: 100%" />
        </el-form-item>
        <el-form-item label="AGC调节步长(kW)">
          <el-input-number v-model="agcParametersForm.agcStepSize" :step="0.1" style="width: 100%" />
        </el-form-item>
        <el-form-item label="AGC步长周期(秒)">
          <el-input-number v-model="agcParametersForm.agcStepPeriod" :step="1" style="width: 100%" />
        </el-form-item>
        <el-form-item label="AGC抖动区间(kW)">
          <el-input-number v-model="agcParametersForm.agcVibrationRange" :step="0.1" style="width: 100%" />
        </el-form-item>
        <el-form-item label="AGC调控周期(秒)">
          <el-input-number v-model="agcParametersForm.agcControlPeriod" :step="1" style="width: 100%" />
        </el-form-item>
        <el-form-item label="AGC微调系数">
          <el-input-number v-model="agcParametersForm.agcMicroAdjustmentCoefficient" :step="0.01" style="width: 100%" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="agcParametersDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="saveAgcParameters">保存</el-button>
        </span>
      </template>
    </el-dialog>
    
    <!-- AVC参数设置弹窗 -->
    <el-dialog v-model="avcParametersDialogVisible" title="AVC参数设置" width="600px">
      <el-form :model="avcParametersForm" label-width="140px" label-position="left">
        <el-form-item label="调节步长(kV)">
          <el-input-number v-model="avcParametersForm.avcStepSize" :step="0.1" style="width: 100%" />
        </el-form-item>
        <el-form-item label="步长周期(秒)">
          <el-input-number v-model="avcParametersForm.avcStepPeriod" :step="1" style="width: 100%" />
        </el-form-item>
        <el-form-item label="抖动区间(kV)">
          <el-input-number v-model="avcParametersForm.avcVibrationRange" :step="0.1" style="width: 100%" />
        </el-form-item>
        <el-form-item label="调控周期(秒)">
          <el-input-number v-model="avcParametersForm.avcControlPeriod" :step="1" style="width: 100%" />
        </el-form-item>
        <el-form-item label="系统阻抗">
          <el-input-number v-model="avcParametersForm.avcSystemImpedance" :step="0.1" style="width: 100%" />
        </el-form-item>
        <el-form-item label="调节范围最小值(kV)">
          <el-input-number v-model="avcParametersForm.avcAdjustmentRangeMin" :step="0.1" style="width: 100%" />
        </el-form-item>
        <el-form-item label="调节范围最大值(kV)">
          <el-input-number v-model="avcParametersForm.avcAdjustmentRangeMax" :step="0.1" style="width: 100%" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="avcParametersDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="saveAvcParameters">保存</el-button>
        </span>
      </template>
    </el-dialog>
    
    <!-- 计划曲线弹窗 -->
    <el-dialog v-model="planCurvesDialogVisible" :title="getPlanCurveTitle()" width="800px">
      <!-- 只保留本地曲线，去掉调度曲线 -->
      <el-table :data="localCurveData" style="width: 100%" max-height="300">
        <el-table-column prop="time" label="时间" width="150">
          <template #default="scope">
            <el-time-picker
              v-model="scope.row.time"
              format="HH:mm:ss"
              value-format="HH:mm:ss"
              placeholder="选择时间"
              style="width: 100%"
            />
          </template>
        </el-table-column>
        <el-table-column prop="targetValue" :label="getPlanCurveValueLabel()" width="150">
          <template #default="scope">
            <el-input-number 
              v-model="scope.row.targetValue" 
              :step="getPlanCurveStep()" 
              style="width: 100%" 
              controls-position="right"
            />
          </template>
        </el-table-column>
        <el-table-column label="操作" width="80">
          <template #default="scope">
            <el-button type="danger" link @click="removeLocalCurvePoint(scope.$index)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div style="margin-top: 10px;">
        <el-button type="primary" @click="addLocalCurvePoint">添加时间点</el-button>
      </div>
      
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="planCurvesDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="savePlanCurves">保存</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import {
  createAgvcBwdHis,
  deleteAgvcBwdHis,
  deleteAgvcBwdHisByIds,
  updateAgvcBwdHis,
  findAgvcBwdHis,
  getAgvcBwdHisList,
  getAgvcBwdHistory,
  generateAgvcBwdTestData,
  // 新增的API
  updateAgcParameters,
  updateAvcParameters,
  updatePlanCurves,
  getAgcParameters,
  getAvcParameters,
  getPlanCurves,
  updateAgcStatus,
  getAgcStatus,
  updateAvcStatus,
  getAvcStatus,
  getBwdRealtimeData
} from '@/api/agvc/agvcBwdHis'

import { getAgvcBwdSettingList, updateAgvcBwdSetting } from '@/api/agvc/agvcBwdSetting'

// 导入图表组件
import PowerChart from '@/view/agvc/agvcChart/powerChart.vue'
import VoltageReactiveChart from '@/view/agvc/agvcChart/voltageReactiveChart.vue'

// 全量引入格式化工具 请按需保留
import { getDictFunc, formatDate, formatBoolean, filterDict ,filterDataSource, returnArrImg, onDownloadFile } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive, nextTick, onMounted, onUnmounted, watch } from 'vue'
import { useAppStore } from "@/pinia"
import * as echarts from 'echarts'
import { ArrowDown } from '@element-plus/icons-vue'





defineOptions({
    name: 'AgvcBwdHis'
})

// 提交按钮loading
const btnLoading = ref(false)
const appStore = useAppStore()

// 控制更多查询条件显示/隐藏状态
const showAllQuery = ref(false)

// 自动化生成的字典（可能为空）以及字段
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
  const table = await getAgvcBwdHisList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
}

getTableData()

// =========== 实时控制面板相关 ===========
// 设备列表
const deviceList = ref([])
// 当前电站编号
const currentPsid = ref(1)

// 获取当前选中的设备
const getCurrentDevice = () => {
  return deviceList.value.find(device => device.number === mainActiveTab.value)
}

// 设备切换事件处理
const onDeviceChange = async (deviceNumber) => {
  console.log('========== onDeviceChange 开始 ==========')
  const device = deviceList.value.find(d => d.number === deviceNumber)
  if (device) {
    console.log('选中设备:', device)
    // 更新电站编号
    if (device.psid) {
      currentPsid.value = device.psid
    }
    
    // 加载新设备的参数和数据
    await loadDeviceParameters(device.number)
  }
}

// el-dropdown 命令处理
const handleDeviceCommand = async (deviceNumber) => {
  console.log('========== handleDeviceCommand 开始 ==========')
  console.log('选中设备编号:', deviceNumber)
  
  // 更新当前选中的设备
  mainActiveTab.value = deviceNumber
  
  // 调用设备切换处理逻辑
  await onDeviceChange(deviceNumber)
}

// 获取并网点设备列表
const getDeviceList = async () => {
  console.log('========== getDeviceList 开始 ==========')
  try {
    const res = await getAgvcBwdSettingList({ page: 1, pageSize: 1000 })
    console.log('1. 获取设备列表响应:', res)
    
    if (res.code === 0) {
      deviceList.value = res.data.list.map(item => ({
        ID: item.ID, // 保存ID用于更新
        number: item.eqid,
        name: item.name,
        psid: item.psid,
        capacity: item.capacity,
        gridConnectionDate: item.gridConnectionDate
      }))
      
      // 设置当前电站编号
      if (deviceList.value.length > 0 && deviceList.value[0].psid) {
        currentPsid.value = deviceList.value[0].psid
      }
      console.log('2. 设备列表长度:', deviceList.value.length)
      
      // 默认选择第一个设备
      if (deviceList.value.length > 0) {
        mainActiveTab.value = deviceList.value[0].number
        console.log('3. 默认选中设备:', deviceList.value[0])
        
        // 等待DOM渲染完成
        console.log('4. 等待DOM渲染...')
        await nextTick()
        console.log('5. DOM渲染完成')
        
        
        // 2. 再加载数据（数据加载后会自动调用updateAgcChart更新图表）
        // console.log('8. 开始加载设备参数...')
        // await loadDeviceParameters(deviceList.value[0].number)
        // console.log('9. 设备参数加载完成')
        
        // console.log('10. 开始加载计划曲线...')
        // await loadDevicePlanCurves(deviceList.value[0].number)
        // console.log('11. 计划曲线加载完成')
      }
    }
    console.log('========== getDeviceList 结束 ==========')
  } catch (error) {
    console.error('❌ 获取并网点设备列表失败:', error)
    console.error('错误堆栈:', error.stack)
  }
}

// 选择设备
const selectDevice = async (device) => {
  console.log('========== selectDevice 开始 ==========')
  // 设置主标签页为当前设备
  mainActiveTab.value = device.number
  
  // 更新电站编号
  if (device.psid) {
    currentPsid.value = device.psid
  }
  
  // 等待DOM更新
  await nextTick()
  
  // 2. 再加载该设备的参数和曲线
  await loadDeviceParameters(device.number)
  //await loadDevicePlanCurves(device.number)
}

// 加载设备参数
const loadDeviceParameters = async (number) => {
  if (!number) {
    console.error('设备编号为空')
    return
  }
  
  console.log('========== loadDeviceParameters 开始 ==========')
  console.log('设备编号:', number)
  try {    
    // 加载AGC状态
    console.log('3. 加载AGC状态...')
    const statusRes = await getAgcStatus({ number })
    console.log('AGC状态响应:', statusRes)
    if (statusRes.code === 0 && statusRes.data) {
      agcFunctionState.value = statusRes.data.agcFunctionState
      agcControlMode.value = statusRes.data.agcControlMode
      agcControlAuthority.value = statusRes.data.agcControlAuthority
      targetActivePower.value = statusRes.data.targetActivePower
      agcUpperLimit.value = statusRes.data.agcUpperLimit || 0
      agcLowerLimit.value = statusRes.data.agcLowerLimit || 0
    }
    
    // 加载AVC状态（从agvc_bwd_his表）
    console.log('3.5. 加载AVC状态...')
    const avcStatusRes = await getAvcStatus({ number })
    console.log('AVC状态响应:', avcStatusRes)
    if (avcStatusRes.code === 0 && avcStatusRes.data) {
      avcFunctionState.value = avcStatusRes.data.avcFunctionState
      avcControlMode.value = avcStatusRes.data.avcControlMode
      avcControlAuthority.value = avcStatusRes.data.avcControlAuthority
      targetReactive.value = avcStatusRes.data.targetReactivePower || 0
      targetVoltage.value = avcStatusRes.data.targetVoltage || 0

    }
    
    // 加载并网点实时数据
    console.log('3.6. 加载并网点实时数据...')
    const realtimeRes = await getBwdRealtimeData({ number })
    console.log('实时数据响应:', realtimeRes)
    if (realtimeRes.code === 0 && realtimeRes.data) {
      // 更新AGC相关实时数据
      currentActivePower.value = realtimeRes.data.currentActivePower || 0
      // targetActivePower.value = realtimeRes.data.targetActivePower || 0
      systemFrequency.value = realtimeRes.data.systemFrequency || 50
      
      // 更新AVC相关实时数据
      currentVoltage.value = realtimeRes.data.currentVoltage || 0
      
      currentReactive.value = realtimeRes.data.currentReactive || 0
      systemImpedance.value = realtimeRes.data.systemImpedance || 0
      
      console.log('实时数据更新完成')
    }
    
    // 无论是否有数据，都更新AGC图表（显示底图）
    // console.log('7. 调用 updateAgcChart...')
    // updateAgcChart(number)
    
    // isLoadingAgcStatus = false
    console.log('========== loadDeviceParameters 完成 ==========')
  } catch (error) {
    console.error('❌ 加载设备参数失败:', error)
    console.error('错误堆栈:', error.stack)
    isLoadingAgcStatus = false
  }
}

// 加载设备计划曲线
const loadDevicePlanCurves = async (number) => {
  try {
    const res = await getPlanCurves({ number })
    if (res.code === 0) {
      // AGC和AVC分开存储，根据当前模式加载
      if (controlActiveTab.value === 'agc' && res.data.agc) {
        localCurveData.value = res.data.agc.localCurve || []
      } else if (controlActiveTab.value === 'avc') {
        // 根据AVC模式加载对应的计划曲线
        if (avcMode.value === 'voltage' && res.data['avc-voltage']) {
          localCurveData.value = res.data['avc-voltage'].localCurve || []
        } else if (avcMode.value === 'reactive' && res.data['avc-reactive']) {
          localCurveData.value = res.data['avc-reactive'].localCurve || []
        }
      }
    }
  } catch (error) {
    console.error('加载计划曲线失败:', error)
  }
}

// 活动标签页

// 当前显示值（模拟数据，TODO:pengchen 后续算法更新后补上）
const currentActivePower = ref(100.5)
const currentVoltage = ref(10.2)
const systemFrequency = ref(50.0)
const systemImpedance = ref(0.8)
const targetReactive = ref(0) // 目标无功
const currentReactive = ref(0) // 当前无功

// AGC状态值
const agcFunctionState = ref(1) // 1: 投入, 0: 退出
const agcControlMode = ref(1) // 1: 闭环指导, 0: 开环指导
const agcControlAuthority = ref(1) // 1: 调度控制, 0: 站内控制

// AVC状态值
const avcFunctionState = ref(1) // 1: 投入, 0: 退出
const avcControlMode = ref(1) // 1: 开环指导, 0: 闭环调节
const avcControlAuthority = ref(1) // 1: 站内控制, 0: 调度控制

// 已禁用loading状态，不再显示loading动画和遮罩
let isLoadingAgcStatus = false

// 监听AGC状态变化并自动保存
watch([agcFunctionState, agcControlMode, agcControlAuthority], async () => {
  // 获取当前选中的设备
  const currentDevice = deviceList.value.find(device => device.number === mainActiveTab.value)
  if (!currentDevice) return
  
  // try {
  //   const res = await updateAgcStatus({
  //     number: currentDevice.number,
  //     agcFunctionState: agcFunctionState.value,
  //     agcControlMode: agcControlMode.value,
  //     agcControlAuthority: agcControlAuthority.value
  //   })
  //   if (res.code !== 0) {
  //     console.error('保存AGC状态失败:', res.msg)
  //   }
  // } catch (error) {
  //   console.error('保存AGC状态失败:', error)
  // }
})

// 监听AVC状态变化并自动保存
watch([avcFunctionState, avcControlMode, avcControlAuthority], async () => {
  // 获取当前选中的设备
  const currentDevice = deviceList.value.find(device => device.number === mainActiveTab.value)
  if (!currentDevice) return
  
  // try {
  //   // 更新AVC状态（与agvc_bwd_his表）
  //   const res = await updateAvcStatus({
  //     number: currentDevice.number,
  //     avcFunctionState: avcFunctionState.value,
  //     avcControlMode: avcControlMode.value,
  //     avcControlAuthority: avcControlAuthority.value
  //   })
  //   if (res.code !== 0) {
  //     console.error('保存AVC状态失败:', res.msg)
  //   }
  // } catch (error) {
  //   console.error('保存AVC状态失败:', error)
  // }
})

// 定时刷新实时数据
let realtimeTimer = null

// 刷新实时数据
const refreshRealtimeData = async () => {
  const currentDevice = deviceList.value.find(device => device.number === mainActiveTab.value)
  if (!currentDevice) {
    return
  }
  
  try {
    // 从历史数据中获取最新数据作为实时数据
    const now = new Date()
    const endTime = now.toISOString().split('T')[0]
    const startTime = new Date(now.getTime() - 60000).toISOString().split('T')[0] // 最近1分钟
    
    const realtimeRes = await getAgvcBwdHistory({ eqid: currentDevice.number, startTime, endTime })
    if (realtimeRes.code === 0 && realtimeRes.data && realtimeRes.data.length > 0) {
      // 提取最新数据更新电站负荷字段
      const sortedData = realtimeRes.data.sort((a, b) => new Date(a.time) - new Date(b.time))
      const latestData = {}
      sortedData.forEach(item => {
        latestData[item.point] = item.value
      })
      
      // 更新AGC相关实时数据
      currentActivePower.value = latestData['100'] || currentActivePower.value
      targetActivePower.value = latestData['101'] || targetActivePower.value
      agcUpperLimit.value = latestData['401'] || agcUpperLimit.value
      agcLowerLimit.value = latestData['402'] || agcLowerLimit.value
      systemFrequency.value = latestData['10'] || systemFrequency.value
      
      // 更新AVC相关实时数据
      targetVoltage.value = latestData['102'] || targetVoltage.value
      currentVoltage.value = latestData['103'] || currentVoltage.value
      targetReactive.value = latestData['104'] || targetReactive.value
      currentReactive.value = latestData['105'] || currentReactive.value
      systemImpedance.value = latestData['106'] || systemImpedance.value
      
      console.log('实时数据刷新完成')
    }
  } catch (error) {
    console.error('刷新实时数据失败:', error)
  }
}

// 启动定时刷新
const startRealtimeTimer = () => {
  // 清除现有定时器
  if (realtimeTimer) {
    clearInterval(realtimeTimer)
  }
  
  // 设置30秒刷新一次
  realtimeTimer = setInterval(refreshRealtimeData, 30000)
  console.log('实时数据定时刷新已启动，间隔30秒')
}

// 停止定时刷新
const stopRealtimeTimer = () => {
  if (realtimeTimer) {
    clearInterval(realtimeTimer)
    realtimeTimer = null
    console.log('实时数据定时刷新已停止')
  }
}

// 监听主标签页变化（并网点切换）
// watch(mainActiveTab, (newNumber, oldNumber) => {
//   if (!newNumber || newNumber === oldNumber) return
  
//   console.log('并网点切换:', oldNumber, '->', newNumber)
  
//   // 停止旧的定时器
//   stopRealtimeTimer()
  
//   // 加载新设备的参数和曲线
//   loadDeviceParameters(newNumber)
//   // loadDevicePlanCurves(newNumber)
  
//   // 延迟初始化图表
//   nextTick(() => {
//     if (controlActiveTab.value === 'agc') {
//       initAgcChart()
//     } else {
//       initAvcChart()
//     }
//   })
  
//   // 启动新的定时器
//   startRealtimeTimer()
// })

// 目标值
const targetActivePower = ref(120)
const targetVoltage = ref(10.5)

// 格式化数字显示
const formatNumber = (value) => {
  if (value === null || value === undefined) return '--'
  return value.toFixed(2)
}

// AGC需要从数据库读取的字段（TODO:pengchen 后续由窗接口提供）
const agcUpperLimit = ref(0) // 可调上限(kW) - TODO:pengchen
const agcLowerLimit = ref(0) // 可调下限(kW) - TODO:pengchen

// 目标值变化处理
const onTargetActivePowerChange = (value) => {
  console.log('目标有功变化:', value)
  // TODO: 发送控制指令
}

const onTargetVoltageChange = (value) => {
  console.log('目标电压变化:', value)
  // TODO: 发送控制指令
}

// AGC参数设置弹窗
const agcParametersDialogVisible = ref(false)
const agcParametersForm = ref({
  ID: null,
  number: '',
  name: '',
  voltageLevel: 0,
  agcFunctionExit: 0, // 改为数字类型
  agcStepSize: 0,
  agcStepPeriod: 0,
  agcVibrationRange: 0,
  agcControlPeriod: 0,
  agcMicroAdjustmentCoefficient: 0
})

// 打开AGC参数设置弹窗
const openAgcParametersDialog = async () => {
  // 获取当前选中的设备
  const currentDevice = deviceList.value.find(device => device.number === mainActiveTab.value)
  if (!currentDevice) {
    ElMessage.warning('请先选择设备')
    return
  }
  
  // 从并网点配置中获取数据
  try {
    const res = await getAgvcBwdSettingList({ number: currentDevice.number })
    if (res.code === 0 && res.data.list.length > 0) {
      const setting = res.data.list[0]
      agcParametersForm.value = {
        ID: setting.ID, // 保存ID用于更新
        number: setting.number || '',
        name: setting.name || '',
        voltageLevel: setting.voltageLevel || '',
        agcFunctionExit: setting.agcFunctionExit || '',
        agcStepSize: setting.agcStepSize || 0,
        agcStepPeriod: setting.agcStepPeriod || 0,
        agcVibrationRange: setting.agcVibrationRange || 0,
        agcControlPeriod: setting.agcControlPeriod || 0,
        agcMicroAdjustmentCoefficient: setting.agcMicroAdjustmentCoefficient || 0
      }
      agcParametersDialogVisible.value = true
    } else {
      ElMessage.error('未找到并网点配置数据')
    }
  } catch (error) {
    ElMessage.error('获取并网点配置数据失败: ' + error.message)
  }
}

// 保存AGC参数（同时更新并网点配置）
const saveAgcParameters = async () => {
  try {
    // 直接更新并网点配置
    const res = await updateAgvcBwdSetting({
      ID: agcParametersForm.value.ID, // 需要ID
      number: agcParametersForm.value.number,
      agcFunctionExit: agcParametersForm.value.agcFunctionExit,
      agcStepSize: agcParametersForm.value.agcStepSize,
      agcStepPeriod: agcParametersForm.value.agcStepPeriod,
      agcVibrationRange: agcParametersForm.value.agcVibrationRange,
      agcControlPeriod: agcParametersForm.value.agcControlPeriod,
      agcMicroAdjustmentCoefficient: agcParametersForm.value.agcMicroAdjustmentCoefficient
    })
    
    if (res.code === 0) {
      ElMessage.success('AGC参数保存成功')
      agcParametersDialogVisible.value = false
    } else {
      ElMessage.error('保存失败: ' + res.msg)
    }
  } catch (error) {
    ElMessage.error('保存失败: ' + error.message)
  }
}

// AVC参数设置弹窗
const avcParametersDialogVisible = ref(false)
const avcParametersForm = ref({
  ID: null,
  number: '',
  name: '',
  voltageLevel: 0,
  avcStepSize: 0,
  avcStepPeriod: 0,
  avcVibrationRange: 0,
  avcControlPeriod: 0,
  avcSystemImpedance: 0,
  avcAdjustmentRangeMin: 0,
  avcAdjustmentRangeMax: 0
})

// 打开AVC参数设置弹窗
const openAvcParametersDialog = async () => {
  // 获取当前选中的设备
  const currentDevice = deviceList.value.find(device => device.number === mainActiveTab.value)
  if (!currentDevice) {
    ElMessage.warning('请先选择设备')
    return
  }
  
  // 从并网点配置中获取数据
  try {
    const res = await getAgvcBwdSettingList({ number: currentDevice.number })
    if (res.code === 0 && res.data.list.length > 0) {
      const setting = res.data.list[0]
      avcParametersForm.value = {
        ID: setting.ID, // 保存ID用于更新
        number: setting.number || '',
        name: setting.name || '',
        voltageLevel: setting.voltageLevel || '',
        avcStepSize: setting.avcStepSize || 0,
        avcStepPeriod: setting.avcStepPeriod || 0,
        avcVibrationRange: setting.avcVibrationRange || 0,
        avcControlPeriod: setting.avcControlPeriod || 0,
        avcSystemImpedance: setting.avcSystemImpedance || 0,
        avcAdjustmentRangeMin: setting.avcAdjustmentRangeMin || 0,
        avcAdjustmentRangeMax: setting.avcAdjustmentRangeMax || 0
      }
      avcParametersDialogVisible.value = true
    } else {
      ElMessage.error('未找到并网点配置数据')
    }
  } catch (error) {
    ElMessage.error('获取并网点配置数据失败: ' + error.message)
  }
}

// 保存AVC参数（同时更新并网点配置）
const saveAvcParameters = async () => {
  try {
    // 直接更新并网点配置
    const res = await updateAgvcBwdSetting({
      ID: avcParametersForm.value.ID, // 需要ID
      number: avcParametersForm.value.number,
      avcStepSize: avcParametersForm.value.avcStepSize,
      avcStepPeriod: avcParametersForm.value.avcStepPeriod,
      avcVibrationRange: avcParametersForm.value.avcVibrationRange,
      avcControlPeriod: avcParametersForm.value.avcControlPeriod,
      avcSystemImpedance: avcParametersForm.value.avcSystemImpedance,
      avcAdjustmentRangeMin: avcParametersForm.value.avcAdjustmentRangeMin,
      avcAdjustmentRangeMax: avcParametersForm.value.avcAdjustmentRangeMax
    })
    
    if (res.code === 0) {
      ElMessage.success('AVC参数保存成功')
      avcParametersDialogVisible.value = false
    } else {
      ElMessage.error('保存失败: ' + res.msg)
    }
  } catch (error) {
    ElMessage.error('保存失败: ' + error.message)
  }
}

// 计划曲线弹窗
const planCurvesDialogVisible = ref(false)
const planCurvesType = ref('agc') // agc, avc-voltage, avc-reactive
const localCurveData = ref([])
// 不再需要调度曲线数据

// 获取计划曲线弹窗标题
const getPlanCurveTitle = () => {
  if (planCurvesType.value === 'agc') return 'AGC计划曲线'
  if (planCurvesType.value === 'avc-voltage') return 'AVC电压计划曲线'
  if (planCurvesType.value === 'avc-reactive') return 'AVC无功计划曲线'
  return '计划曲线'
}

// 获取计划曲线值标签
const getPlanCurveValueLabel = () => {
  if (planCurvesType.value === 'agc') return '目标有功(kW)'
  if (planCurvesType.value === 'avc-voltage') return '目标电压(kV)'
  if (planCurvesType.value === 'avc-reactive') return '目标无功(kVar)'
  return '目标值'
}

// 获取计划曲线步长
const getPlanCurveStep = () => {
  if (planCurvesType.value === 'agc') return 1
  return 0.1 // 电压和无功都用 0.1
}

// 打开AGC计划曲线弹窗
const openAgcPlanCurvesDialog = async () => {
  // 获取当前选中的设备
  const currentDevice = deviceList.value.find(device => device.number === mainActiveTab.value)
  if (!currentDevice) {
    ElMessage.warning('请先选择设备')
    return
  }
  planCurvesType.value = 'agc'
  
  // 加载AGC计划曲线
  try {
    const res = await getPlanCurves({ number: currentDevice.number, curveType: 'agc' })
    if (res.code === 0 && res.data) {
      localCurveData.value = res.data.localCurve || []
    }
  } catch (error) {
    console.error('加载AGC计划曲线失败:', error)
  }
  
  planCurvesDialogVisible.value = true
}

// 打开AVC计划曲线弹窗
const openAvcPlanCurvesDialog = async () => {
  // 获取当前选中的设备
  const currentDevice = deviceList.value.find(device => device.number === mainActiveTab.value)
  if (!currentDevice) {
    ElMessage.warning('请先选择设备')
    return
  }
  
  // 根据AVC模式设置曲线类型
  planCurvesType.value = avcMode.value === 'voltage' ? 'avc-voltage' : 'avc-reactive'
  
  // 加载AVC计划曲线
  try {
    const res = await getPlanCurves({ number: currentDevice.number, curveType: planCurvesType.value })
    if (res.code === 0 && res.data) {
      localCurveData.value = res.data.localCurve || []
    }
  } catch (error) {
    console.error('加载AVC计划曲线失败:', error)
  }
  
  planCurvesDialogVisible.value = true
}

// 添加本地曲线点
const addLocalCurvePoint = () => {
  localCurveData.value.push({
    time: '',
    targetValue: 0
  })
}

// 删除本地曲线点
const removeLocalCurvePoint = (index) => {
  localCurveData.value.splice(index, 1)
}

// 保存计划曲线
const savePlanCurves = async () => {
  try {
    const currentDevice = deviceList.value.find(device => device.number === mainActiveTab.value)
    if (!currentDevice) {
      ElMessage.warning('请先选择设备')
      return
    }
    
    console.log('保存计划曲线:', {
      number: currentDevice.number,
      curveType: planCurvesType.value,
      localCurveData: localCurveData.value
    })
    
    const res = await updatePlanCurves({
      number: currentDevice.number,
      curveType: planCurvesType.value,
      planCurves: {
        localCurve: localCurveData.value
      }
    })
    
    if (res.code === 0) {
      ElMessage.success('计划曲线保存成功')
      planCurvesDialogVisible.value = false
      
      // 保存成功后重新加载计划曲线
      loadDevicePlanCurves(currentDevice.number)
    } else {
      ElMessage.error('保存失败: ' + res.msg)
    }
  } catch (error) {
    ElMessage.error('保存失败: ' + error.message)
  }
}

// 显示历史数据弹窗（从主标签右侧按钮进入）
const showHistoryDataDialog = async () => {
  // 获取当前选中的设备
  const current = deviceList.value.find(device => device.number === mainActiveTab.value)
  if (!current) {
    ElMessage.warning('请先选择设备')
    return
  }
  
  // 设置当前设备
  currentDevice.value = current
  
  // 设置默认日期范围为最近7天
  const end = new Date()
  const start = new Date()
  start.setTime(start.getTime() - 3600 * 1000 * 24 * 7)
  historyDateRange.value = [start, end]
  
  // 显示弹窗
  historyDialogVisible.value = true
  
  // 在DOM更新后初始化图表并加载数据
  await nextTick()
  
  // 先初始化图表
  initChart()
  
  // 再加载历史数据（加载完成后会自动调用updateChart）
  await loadHistoryData()
}

// =========== 表格控制部分结束 ===============

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

// 历史数据图表实例（独立命名，避免与AGC/AVC图表冲突）
let historyChartInstance = null

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
  
  try {
    const params = {
      eqid: currentDevice.value.number,  // 使用设备编号
      startTime: historyDateRange.value[0],
      endTime: historyDateRange.value[1]
    }
    
    const res = await getAgvcBwdHistory(params)
    if (res.code === 0) {
      // 处理从 InfluxDB获取的历史数据
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
  console.log('initChart 被调用')
  console.log('chartContainer.value:', chartContainer.value)
  console.log('historyChartInstance:', historyChartInstance)
  
  if (!chartContainer.value) {
    console.warn('chartContainer.value 为空，等待DOM渲染')
    return
  }
  
  // 使用 echarts 绘制图表
  if (!historyChartInstance) {
    console.log('创建新的echarts实例')
    historyChartInstance = echarts.init(chartContainer.value)
  }
  
  updateChart()
}

// 更新图表
const updateChart = () => {
  console.log('updateChart 被调用')
  console.log('historyChartInstance:', historyChartInstance)
  console.log('historyData.value.length:', historyData.value.length)
  
  if (!historyChartInstance) {
    console.warn('historyChartInstance 为空，无法更新图表')
    return
  }
  
  if (!historyData.value.length) {
    console.warn('历史数据为空，不更新图表')
    return
  }
  
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
  historyChartInstance.setOption(option, true)
}

// 在组件卸载时移除监听
onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  // 清除定时器
  if (refreshTimer) {
    clearInterval(refreshTimer)
    refreshTimer = null
  }
  // 清理历史数据图表实例
  if (historyChartInstance) {
    historyChartInstance.dispose()
    historyChartInstance = null
  }
  // 清理AGC和AVC图表实例
  if (agcChartInstance) {
    agcChartInstance.dispose()
  }
  if (avcChartInstance) {
    avcChartInstance.dispose()
  }
})

// ======== 大屏幕图表相关 =========
// 主标签页
const mainActiveTab = ref('')
// 控制标签页
const controlActiveTab = ref('agc')
// AVC模式切换
const avcMode = ref('voltage') // voltage 或 reactive

// 图表实例
let agcChartInstance = null
let avcChartInstance = null

// AGC图表容器引用
const agcChartContainer = ref(null)
// AVC图表容器引用
const avcChartContainer = ref(null)

// AVC模式切换处理
const onAvcModeChange = (mode) => {
  console.log('AVC模式切换到:', mode)
  // 重新加载AVC图表数据
  loadAvcChartData()
}

// 加载AVC图表数据
const loadAvcChartData = () => {
  updateAvcChart()
}

// 定时器变量
let refreshTimer = null

// 在组件挂载时初始化图表
onMounted(() => {
  window.addEventListener('resize', handleResize)
  // 加载并网点设备列表
  getDeviceList()
  
  // 设置定时器，每秒刷新一次设备参数
  refreshTimer = setInterval(() => {
    if (mainActiveTab.value) {
      loadDeviceParameters(mainActiveTab.value)
    }
  }, 5000)
})

// 监听主标签页变化（并网点切换）
watch(mainActiveTab, (newNumber, oldNumber) => {
  if (!newNumber || newNumber === oldNumber) return
  console.log('并网点切换:', oldNumber, '->', newNumber)
  // 加载新设备的参数和曲线
  loadDeviceParameters(newNumber)
  // loadDevicePlanCurves(newNumber)
  
})

// 监听窗口大小变化，重置图表大小
const handleResize = () => {
  // 历史数据图表
  if (historyChartInstance) {
    historyChartInstance.resize()
  }
  // AGC图表
  if (agcChartInstance) {
    agcChartInstance.resize()
  }
  // AVC图表
  if (avcChartInstance) {
    avcChartInstance.resize()
  }
}


</script>

<style scoped>
/* 实时控制面板 */
.realtime-control-panel {
  background: #f5f7fa;
  border: 1px solid #e4e7ed;
  border-radius: 4px;
  padding: 15px;
  margin-bottom: 20px;
}

/* 设备选择器 */
.device-selector {
  margin-bottom: 15px;
}

.selector-title {
  font-size: 16px;
  font-weight: bold;
  margin-bottom: 10px;
  color: #303133;
}

.device-buttons {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.device-button {
  margin-bottom: 5px;
}

/* 控制标签页 */
.control-tabs {
  margin-top: 15px;
}

/* 未选择设备提示 */
.no-device-selected {
  padding: 20px;
  text-align: center;
}

/* 主标签页 */
.main-tabs {
  background: #ffffff;
  border: 1px solid #e4e7ed;
  border-radius: 4px;
  padding: 15px;
  margin-bottom: 20px;
}

/* 控制内容 */
.control-content {
  padding: 15px 0;
}

/* 控制区域 */
.control-section {
  margin-bottom: 20px;
  padding-bottom: 20px;
  border-bottom: 1px solid #f0f0f0;
}

.control-section:last-child {
  border-bottom: none;
  margin-bottom: 0;
  padding-bottom: 0;
}

.section-title {
  font-size: 14px;
  font-weight: bold;
  margin-bottom: 10px;
  color: #606266;
  display: flex;
  align-items: center;
}

.section-title::before {
  content: '';
  width: 4px;
  height: 16px;
  background: #409eff;
  border-radius: 2px;
  margin-right: 8px;
}

.control-row {
  display: flex;
  gap: 10px;
  margin-bottom: 15px;
}

/* 控制网格 */
.control-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 15px;
}

/* 电站负荷特殊网格（整整增大） */
.control-grid.load-grid {
  grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));
  gap: 20px;
  padding: 15px 0;
}

/* 控制项 */
.control-item {
  display: flex;
  flex-direction: column;
}

.item-label {
  font-size: 12px;
  color: #606266;
  margin-bottom: 5px;
}

.item-value {
  font-size: 14px;
  font-weight: bold;
  color: #303133;
}

.item-value.readonly {
  color: #909399;
  font-weight: normal;
}

.radio-group {
  display: flex;
  align-items: center;
}

.radio-group .el-radio {
  margin-right: 15px;
}

/* 图表区域 */
.chart-section {
  margin-top: 30px;
  padding-top: 20px;
  border-top: 2px solid #f0f0f0;
}

/* 图表区域 - 位于最上面 */
.chart-section.chart-top {
  margin-top: 0;
  padding-top: 0;
  border-top: none;
  margin-bottom: 30px;
  padding-bottom: 20px;
  border-bottom: 2px solid #f0f0f0;
}

/* AVC模式切换 */
.avc-mode-switch {
  margin-bottom: 15px;
  text-align: center;
}

/* 图表容器 */
.chart-container {
  height: 400px;
  margin-top: 15px;
  border: 1px solid #e4e7ed;
  border-radius: 4px;
  width: 100%;
  margin: 15px 0 0 0;
  padding: 0;
}

.chart-wrapper {
  width: 100%;
  height: 100%;
  padding: 0;
  margin: 0;
}

/* 优化控制项样式 */
.control-item :deep(.el-input-number) {
  width: 100%;
}

.control-item :deep(.el-input-number .el-input__wrapper) {
  padding: 0 8px;
  height: 32px;
}

.control-item :deep(.el-input-number input) {
  text-align: center;
}

/* 设备选择卡片样式 */
.device-selector-card {
  margin-bottom: 20px;
}

.device-info {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.device-name-row {
  display: flex;
  align-items: center;
}

.device-name-row :deep(.el-select .el-input__wrapper) {
  font-size: 18px;
  font-weight: bold;
  padding: 10px 15px;
}

.device-details {
  display: flex;
  gap: 30px;
  padding-left: 5px;
}

.detail-item {
  display: flex;
  align-items: center;
  font-size: 14px;
}

.detail-label {
  color: #909399;
  margin-right: 8px;
}

.detail-value {
  color: #606266;
  font-weight: 500;
}

/* 响应式调整 */
@media (max-width: 768px) {
  .control-grid {
    grid-template-columns: 1fr;
  }
  
  .control-row {
    flex-direction: column;
    gap: 10px;
  }
  
  .device-buttons {
    flex-direction: column;
  }
  
  .device-details {
    flex-direction: column;
    gap: 10px;
  }
}
</style>
