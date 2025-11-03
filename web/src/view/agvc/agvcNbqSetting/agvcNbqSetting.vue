
<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" class="demo-form-inline" @keyup.enter="onSubmit">
            <el-form-item label="逆变器编号" prop="inverterNo">
  <el-input v-model.number="searchInfo.inverterNo" placeholder="搜索条件" />
</el-form-item>
            
            <el-form-item label="逆变器名称" prop="name">
  <el-input v-model="searchInfo.name" placeholder="搜索条件" />
</el-form-item>
            
            <el-form-item label="参与调节" prop="isParticipateAdjust">
  <el-select v-model="searchInfo.isParticipateAdjust" clearable placeholder="请选择">
    <el-option key="true" label="是" value="true"></el-option>
    <el-option key="false" label="否" value="false"></el-option>
  </el-select>
</el-form-item>
            
            <el-form-item label="标杆逆变器" prop="isBenchmarkInverter">
  <el-select v-model="searchInfo.isBenchmarkInverter" clearable placeholder="请选择">
    <el-option key="true" label="是" value="true"></el-option>
    <el-option key="false" label="否" value="false"></el-option>
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
            <ExportTemplate  template-id="agvc_AgvcNbqSetting" />
            <ExportExcel  template-id="agvc_AgvcNbqSetting" filterDeleted/>
            <ImportExcel  template-id="agvc_AgvcNbqSetting" @on-success="getTableData" />
        </div>
        <el-table
        ref="multipleTable"
        style="width: 100%"
        tooltip-effect="dark"
        :data="tableData"
        row-key="ID"
        @selection-change="handleSelectionChange"
        @sort-change="sortChange"
        >
        <el-table-column type="selection" width="55" />
        
            <el-table-column sortable align="left" label="编号" prop="inverterNo" width="100" />

            <el-table-column align="left" label="并网点编号" prop="bwdNo" width="120" />

            <el-table-column align="left" label="名称" prop="name" width="120" />

            <el-table-column align="left" label="有功额定(kW)" prop="ratedActivePower" width="120" />

            <el-table-column align="left" label="无功额定(kVar)" prop="ratedReactivePower" width="120" />

            <el-table-column align="left" label="抖动区间" prop="jitterRange" width="100" />

            <el-table-column align="left" label="死区区间" prop="deadbandRange" width="100" />

            <el-table-column align="left" label="升额优先级" prop="upgradePriority" width="100" />

            <el-table-column align="left" label="降额优先级" prop="downgradePriority" width="100" />

            <el-table-column align="left" label="参与调节" prop="isParticipateAdjust" width="100">
    <template #default="scope">{{ formatBoolean(scope.row.isParticipateAdjust) }}</template>
</el-table-column>
            <el-table-column align="left" label="标杆设备" prop="isBenchmarkInverter" width="100">
    <template #default="scope">{{ formatBoolean(scope.row.isBenchmarkInverter) }}</template>
</el-table-column>
        <el-table-column align="left" label="操作" fixed="right" :min-width="appStore.operateMinWith">
            <template #default="scope">
            <el-button  type="primary" link class="table-button" @click="getDetails(scope.row)"><el-icon style="margin-right: 5px"><InfoFilled /></el-icon>查看</el-button>
            <el-button  type="primary" link icon="edit" class="table-button" @click="updateAgvcNbqSettingFunc(scope.row)">编辑</el-button>
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
            <el-form-item label="逆变器编号:" prop="inverterNo">
    <el-input v-model="formData.inverterNo" :clearable="true" placeholder="请输入逆变器编号" />
</el-form-item>
            <el-form-item label="并网点编号:" prop="bwdNo">
    <el-input v-model="formData.bwdNo" :clearable="true" placeholder="请输入并网点编号" />
</el-form-item>
            <el-form-item label="逆变器名称:" prop="name">
    <el-input v-model="formData.name" :clearable="true" placeholder="请输入逆变器名称" />
</el-form-item>
            <el-form-item label="额定有功功率:" prop="ratedActivePower">
    <el-input-number v-model="formData.ratedActivePower" style="width:100%" :precision="2" :clearable="true" />
</el-form-item>
            <el-form-item label="额定无功功率:" prop="ratedReactivePower">
    <el-input-number v-model="formData.ratedReactivePower" style="width:100%" :precision="2" :clearable="true" />
</el-form-item>
            <el-form-item label="功率抖动区间:" prop="jitterRange">
    <el-input-number v-model="formData.jitterRange" style="width:100%" :precision="2" :clearable="true" />
</el-form-item>
            <el-form-item label="功率死区区间:" prop="deadbandRange">
    <el-input-number v-model="formData.deadbandRange" style="width:100%" :precision="2" :clearable="true" />
</el-form-item>
            <el-form-item label="升额优先级:" prop="upgradePriority">
    <el-input-number v-model="formData.upgradePriority" style="width:100%" :precision="2" :clearable="true" />
</el-form-item>
            <el-form-item label="降级优先级:" prop="downgradePriority">
    <el-input-number v-model="formData.downgradePriority" style="width:100%" :precision="2" :clearable="true" />
</el-form-item>
            <el-form-item label="参与调节:" prop="isParticipateAdjust">
    <el-switch v-model="formData.isParticipateAdjust" active-color="#13ce66" inactive-color="#ff4949" active-text="是" inactive-text="否" clearable ></el-switch>
</el-form-item>
            <el-form-item label="标杆逆变器:" prop="isBenchmarkInverter">
    <el-switch v-model="formData.isBenchmarkInverter" active-color="#13ce66" inactive-color="#ff4949" active-text="是" inactive-text="否" clearable ></el-switch>
</el-form-item>
          </el-form>
    </el-drawer>

    <el-drawer destroy-on-close :size="appStore.drawerSize" v-model="detailShow" :show-close="true" :before-close="closeDetailShow" title="查看">
            <el-descriptions :column="1" border>
                    <el-descriptions-item label="逆变器编号">
    {{ detailForm.inverterNo }}
</el-descriptions-item>
                    <el-descriptions-item label="并网点编号">
    {{ detailForm.bwdNo }}
</el-descriptions-item>
                    <el-descriptions-item label="逆变器名称">
    {{ detailForm.name }}
</el-descriptions-item>
                    <el-descriptions-item label="额定有功功率">
    {{ detailForm.ratedActivePower }}
</el-descriptions-item>
                    <el-descriptions-item label="额定无功功率">
    {{ detailForm.ratedReactivePower }}
</el-descriptions-item>
                    <el-descriptions-item label="功率抖动区间">
    {{ detailForm.jitterRange }}
</el-descriptions-item>
                    <el-descriptions-item label="功率死区区间">
    {{ detailForm.deadbandRange }}
</el-descriptions-item>
                    <el-descriptions-item label="升额优先级">
    {{ detailForm.upgradePriority }}
</el-descriptions-item>
                    <el-descriptions-item label="降级优先级">
    {{ detailForm.downgradePriority }}
</el-descriptions-item>
                    <el-descriptions-item label="参与调节">
    {{ detailForm.isParticipateAdjust }}
</el-descriptions-item>
                    <el-descriptions-item label="标杆逆变器">
    {{ detailForm.isBenchmarkInverter }}
</el-descriptions-item>
            </el-descriptions>
        </el-drawer>

  </div>
</template>

<script setup>
import {
  createAgvcNbqSetting,
  deleteAgvcNbqSetting,
  deleteAgvcNbqSettingByIds,
  updateAgvcNbqSetting,
  findAgvcNbqSetting,
  getAgvcNbqSettingList
} from '@/api/agvc/agvcNbqSetting'

// 全量引入格式化工具 请按需保留
import { getDictFunc, formatDate, formatBoolean, filterDict ,filterDataSource, returnArrImg, onDownloadFile } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive } from 'vue'
import { useAppStore } from "@/pinia"

// 导出组件
import ExportExcel from '@/components/exportExcel/exportExcel.vue'
// 导入组件
import ImportExcel from '@/components/exportExcel/importExcel.vue'
// 导出模板组件
import ExportTemplate from '@/components/exportExcel/exportTemplate.vue'


defineOptions({
    name: 'AgvcNbqSetting'
})

// 提交按钮loading
const btnLoading = ref(false)
const appStore = useAppStore()

// 控制更多查询条件显示/隐藏状态
const showAllQuery = ref(false)

// 自动化生成的字典（可能为空）以及字段
const formData = ref({
            inverterNo: '',
            bwdNo: '',
            name: '',
            ratedActivePower: 0,
            ratedReactivePower: 0,
            jitterRange: 0,
            deadbandRange: 0,
            upgradePriority: 0,
            downgradePriority: 0,
            isParticipateAdjust: false,
            isBenchmarkInverter: false,
        })



// 验证规则
const rule = reactive({
  inverterNo: [
    { required: true, message: '请输入逆变器编号', trigger: 'blur' },
    { pattern: /^\d+$/, message: '逆变器编号只能输入数字', trigger: 'blur' }
  ],
  bwdNo: [
    { required: true, message: '请输入并网点编号', trigger: 'blur' },
    { pattern: /^\d+$/, message: '并网点编号只能输入数字', trigger: 'blur' }
  ]
})

const elFormRef = ref()
const elSearchFormRef = ref()

// =========== 表格控制部分 ===========
const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({})
// 排序
const sortChange = ({ prop, order }) => {
  const sortMap = {
    CreatedAt:"created_at",
    ID:"id",
            inverterNo: 'inverter_no',
  }

  let sort = sortMap[prop]
  if(!sort){
   sort = prop.replace(/[A-Z]/g, match => `_${match.toLowerCase()}`)
  }

  searchInfo.value.sort = sort
  searchInfo.value.order = order
  getTableData()
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
    if (searchInfo.value.isParticipateAdjust === ""){
        searchInfo.value.isParticipateAdjust=null
    }
    if (searchInfo.value.isBenchmarkInverter === ""){
        searchInfo.value.isBenchmarkInverter=null
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
  const table = await getAgvcNbqSettingList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
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
            deleteAgvcNbqSettingFunc(row)
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
      const res = await deleteAgvcNbqSettingByIds({ ids: IDs })
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
const updateAgvcNbqSettingFunc = async(row) => {
    const res = await findAgvcNbqSetting({ id: row.ID })
    type.value = 'update'
    if (res.code === 0) {
        formData.value = res.data
        dialogFormVisible.value = true
    }
}


// 删除行
const deleteAgvcNbqSettingFunc = async (row) => {
    const res = await deleteAgvcNbqSetting({ id: row.ID })
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
        inverterNo: '',
        bwdNo: '',
        name: '',
        ratedActivePower: 0,
        ratedReactivePower: 0,
        jitterRange: 0,
        deadbandRange: 0,
        upgradePriority: 0,
        downgradePriority: 0,
        isParticipateAdjust: false,
        isBenchmarkInverter: false,
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
                  res = await createAgvcNbqSetting(formData.value)
                  break
                case 'update':
                  res = await updateAgvcNbqSetting(formData.value)
                  break
                default:
                  res = await createAgvcNbqSetting(formData.value)
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
  const res = await findAgvcNbqSetting({ id: row.ID })
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
