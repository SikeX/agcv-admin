
<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" class="demo-form-inline" @keyup.enter="onSubmit">
      <el-form-item label="并网点名称" prop="name">
        <el-input v-model="searchInfo.name" placeholder="请输入并网点名称" />
      </el-form-item>
      <el-form-item label="设备编号" prop="number">
        <el-input v-model="searchInfo.number" placeholder="请输入设备编号" />
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
        
            <el-table-column align="left" label="名称" prop="name" width="120" />

            <el-table-column align="left" label="编号" prop="number" width="100" />

            <el-table-column align="left" label="电压(kV)" prop="voltageLevel" width="100" />

            <el-table-column align="left" label="AGC退出模式" prop="agcFunctionExit" width="130" />

            <el-table-column align="left" label="AGC步长(kW)" prop="agcStepSize" width="130" />

            <el-table-column align="left" label="AGC周期(秒)" prop="agcStepPeriod" width="130" />

            <el-table-column align="left" label="AGC抖动(kW)" prop="agcVibrationRange" width="130" />

            <el-table-column align="left" label="AGC调控周期(秒)" prop="agcControlPeriod" width="140" />

            <el-table-column align="left" label="AGC微调系数" prop="agcMicroAdjustmentCoefficient" width="130" />

            <el-table-column align="left" label="AVC步长(kV)" prop="avcStepSize" width="130" />

            <el-table-column align="left" label="AVC周期(秒)" prop="avcStepPeriod" width="130" />

            <el-table-column align="left" label="AVC抖动(kV)" prop="avcVibrationRange" width="130" />

            <el-table-column align="left" label="AVC调控周期(秒)" prop="avcControlPeriod" width="140" />

            <el-table-column align="left" label="AVC系统阻抗" prop="avcSystemImpedance" width="130" />

            <el-table-column align="left" label="AVC最小值(kV)" prop="avcAdjustmentRangeMin" width="130" />

            <el-table-column align="left" label="AVC最大值(kV)" prop="avcAdjustmentRangeMax" width="130" />

        <el-table-column align="left" label="操作" fixed="right" :min-width="appStore.operateMinWith">
            <template #default="scope">
            <el-button  type="primary" link class="table-button" @click="getDetails(scope.row)"><el-icon style="margin-right: 5px"><InfoFilled /></el-icon>查看</el-button>
            <el-button  type="primary" link icon="edit" class="table-button" @click="updateAgvcBwdSettingFunc(scope.row)">编辑</el-button>
            <el-button   type="primary" link icon="delete" @click="deleteRow(scope.row)">删除</el-button>
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
            <el-form-item label="并网点名称:" prop="name">
    <el-input v-model="formData.name" :clearable="true" placeholder="请输入并网点名称" />
</el-form-item>
            <el-form-item label="设备编号:" prop="number">
    <el-input v-model="formData.number" :clearable="true" placeholder="请输入设备编号" />
</el-form-item>
            <el-form-item label="电压等级(kV):" prop="voltageLevel">
    <el-input-number v-model="formData.voltageLevel" style="width:100%" :precision="2" :clearable="true" />
</el-form-item>
            <el-form-item label="AGC功能退出模式:" prop="agcFunctionExit">
    <el-input v-model.number="formData.agcFunctionExit" :clearable="true" placeholder="请输入AGC功能退出模式" />
</el-form-item>
            <el-form-item label="AGC调节步长(kW):" prop="agcStepSize">
    <el-input-number v-model="formData.agcStepSize" style="width:100%" :precision="2" :clearable="true" />
</el-form-item>
            <el-form-item label="AGC步长周期(秒):" prop="agcStepPeriod">
    <el-input v-model.number="formData.agcStepPeriod" :clearable="true" placeholder="请输入AGC步长周期(秒)" />
</el-form-item>
            <el-form-item label="AGC抖动区间(kW):" prop="agcVibrationRange">
    <el-input-number v-model="formData.agcVibrationRange" style="width:100%" :precision="2" :clearable="true" />
</el-form-item>
            <el-form-item label="AGC调控周期(秒):" prop="agcControlPeriod">
    <el-input v-model.number="formData.agcControlPeriod" :clearable="true" placeholder="请输入AGC调控周期(秒)" />
</el-form-item>
            <el-form-item label="AGC微调系数:" prop="agcMicroAdjustmentCoefficient">
    <el-input-number v-model="formData.agcMicroAdjustmentCoefficient" style="width:100%" :precision="2" :clearable="true" />
</el-form-item>
            <el-form-item label="AVC调节步长(kV):" prop="avcStepSize">
    <el-input-number v-model="formData.avcStepSize" style="width:100%" :precision="2" :clearable="true" />
</el-form-item>
            <el-form-item label="AVC步长周期(秒):" prop="avcStepPeriod">
    <el-input v-model.number="formData.avcStepPeriod" :clearable="true" placeholder="请输入AVC步长周期(秒)" />
</el-form-item>
            <el-form-item label="AVC抖动区间(kV):" prop="avcVibrationRange">
    <el-input-number v-model="formData.avcVibrationRange" style="width:100%" :precision="2" :clearable="true" />
</el-form-item>
            <el-form-item label="AVC调控周期(秒):" prop="avcControlPeriod">
    <el-input v-model.number="formData.avcControlPeriod" :clearable="true" placeholder="请输入AVC调控周期(秒)" />
</el-form-item>
            <el-form-item label="AVC系统阻抗:" prop="avcSystemImpedance">
    <el-input-number v-model="formData.avcSystemImpedance" style="width:100%" :precision="2" :clearable="true" />
</el-form-item>
            <el-form-item label="AVC调节范围最小值(kV):" prop="avcAdjustmentRangeMin">
    <el-input-number v-model="formData.avcAdjustmentRangeMin" style="width:100%" :precision="2" :clearable="true" />
</el-form-item>
            <el-form-item label="AVC调节范围最大值(kV):" prop="avcAdjustmentRangeMax">
    <el-input-number v-model="formData.avcAdjustmentRangeMax" style="width:100%" :precision="2" :clearable="true" />
</el-form-item>
          </el-form>
    </el-drawer>

    <el-drawer destroy-on-close :size="appStore.drawerSize" v-model="detailShow" :show-close="true" :before-close="closeDetailShow" title="查看">
            <el-descriptions :column="1" border>
                    <el-descriptions-item label="并网点名称">
    {{ detailForm.name }}
</el-descriptions-item>
                    <el-descriptions-item label="设备编号">
    {{ detailForm.number }}
</el-descriptions-item>
                    <el-descriptions-item label="电压等级(kV)">
    {{ detailForm.voltageLevel }}
</el-descriptions-item>
                    <el-descriptions-item label="AGC功能退出模式">
    {{ detailForm.agcFunctionExit }}
</el-descriptions-item>
                    <el-descriptions-item label="AGC调节步长(kW)">
    {{ detailForm.agcStepSize }}
</el-descriptions-item>
                    <el-descriptions-item label="AGC步长周期(秒)">
    {{ detailForm.agcStepPeriod }}
</el-descriptions-item>
                    <el-descriptions-item label="AGC抖动区间(kW)">
    {{ detailForm.agcVibrationRange }}
</el-descriptions-item>
                    <el-descriptions-item label="AGC调控周期(秒)">
    {{ detailForm.agcControlPeriod }}
</el-descriptions-item>
                    <el-descriptions-item label="AGC微调系数">
    {{ detailForm.agcMicroAdjustmentCoefficient }}
</el-descriptions-item>
                    <el-descriptions-item label="AVC调节步长(kV)">
    {{ detailForm.avcStepSize }}
</el-descriptions-item>
                    <el-descriptions-item label="AVC步长周期(秒)">
    {{ detailForm.avcStepPeriod }}
</el-descriptions-item>
                    <el-descriptions-item label="AVC抖动区间(kV)">
    {{ detailForm.avcVibrationRange }}
</el-descriptions-item>
                    <el-descriptions-item label="AVC调控周期(秒)">
    {{ detailForm.avcControlPeriod }}
</el-descriptions-item>
                    <el-descriptions-item label="AVC系统阻抗">
    {{ detailForm.avcSystemImpedance }}
</el-descriptions-item>
                    <el-descriptions-item label="AVC调节范围最小值(kV)">
    {{ detailForm.avcAdjustmentRangeMin }}
</el-descriptions-item>
                    <el-descriptions-item label="AVC调节范围最大值(kV)">
    {{ detailForm.avcAdjustmentRangeMax }}
</el-descriptions-item>
            </el-descriptions>
        </el-drawer>

  </div>
</template>

<script setup>
import {
  createAgvcBwdSetting,
  deleteAgvcBwdSetting,
  deleteAgvcBwdSettingByIds,
  updateAgvcBwdSetting,
  findAgvcBwdSetting,
  getAgvcBwdSettingList
} from '@/api/agvc/agvcBwdSetting'

// 全量引入格式化工具 请按需保留
import { getDictFunc, formatDate, formatBoolean, filterDict ,filterDataSource, returnArrImg, onDownloadFile } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive } from 'vue'
import { useAppStore } from "@/pinia"




defineOptions({
    name: 'AgvcBwdSetting'
})

// 提交按钮loading
const btnLoading = ref(false)
const appStore = useAppStore()

// 控制更多查询条件显示/隐藏状态
const showAllQuery = ref(false)

// 自动化生成的字典（可能为空）以及字段
const formData = ref({
            name: '',
            number: '',
            voltageLevel: 0,
            agcFunctionExit: undefined,
            agcStepSize: 0,
            agcStepPeriod: undefined,
            agcVibrationRange: 0,
            agcControlPeriod: undefined,
            agcMicroAdjustmentCoefficient: 0,
            avcStepSize: 0,
            avcStepPeriod: undefined,
            avcVibrationRange: 0,
            avcControlPeriod: undefined,
            avcSystemImpedance: 0,
            avcAdjustmentRangeMin: 0,
            avcAdjustmentRangeMax: 0,
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
const searchInfo = ref({
  name: '',
  number: ''
})
// 重置
const onReset = () => {
  searchInfo.value = {
    name: '',
    number: ''
  }
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
  const table = await getAgvcBwdSettingList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
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
            deleteAgvcBwdSettingFunc(row)
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
      const res = await deleteAgvcBwdSettingByIds({ IDs })
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
const updateAgvcBwdSettingFunc = async(row) => {
    const res = await findAgvcBwdSetting({ ID: row.ID })
    type.value = 'update'
    if (res.code === 0) {
        formData.value = res.data
        dialogFormVisible.value = true
    }
}


// 删除行
const deleteAgvcBwdSettingFunc = async (row) => {
    const res = await deleteAgvcBwdSetting({ ID: row.ID })
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
        name: '',
        number: '',
        voltageLevel: 0,
        agcFunctionExit: undefined,
        agcStepSize: 0,
        agcStepPeriod: undefined,
        agcVibrationRange: 0,
        agcControlPeriod: undefined,
        agcMicroAdjustmentCoefficient: 0,
        avcStepSize: 0,
        avcStepPeriod: undefined,
        avcVibrationRange: 0,
        avcControlPeriod: undefined,
        avcSystemImpedance: 0,
        avcAdjustmentRangeMin: 0,
        avcAdjustmentRangeMax: 0,
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
                  res = await createAgvcBwdSetting(formData.value)
                  break
                case 'update':
                  res = await updateAgvcBwdSetting(formData.value)
                  break
                default:
                  res = await createAgvcBwdSetting(formData.value)
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
  const res = await findAgvcBwdSetting({ ID: row.ID })
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

<style>

</style>
