<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" class="demo-form-inline" @keyup.enter="onSubmit">
      <el-form-item label="已读" prop="isRead">
        <el-select v-model="searchInfo.isRead" clearable placeholder="请选择">
          <el-option label="是" value="是"></el-option>
          <el-option label="否" value="否"></el-option>
        </el-select>
      </el-form-item>
      
      <el-form-item label="数据类型" prop="dataType">
        <el-select v-model="searchInfo.dataType" clearable placeholder="请选择数据类型">
          <el-option
            v-for="item in DataTypeOptions"
            :key="item.value"
            :label="item.label"
            :value="item.value"
          />
        </el-select>
      </el-form-item>
      
      <el-form-item label="发生时间" prop="happenTimeRange">
        <el-date-picker
          v-model="searchInfo.happenTimeRange"
          type="datetimerange"
          range-separator="至"
          start-placeholder="开始时间"
          end-placeholder="结束时间"
          value-format="YYYY-MM-DD HH:mm:ss"
        />
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
            <!-- 移除新增和删除按钮，因为历史事件记录不可编辑 -->
        </div>
        <el-table
        ref="multipleTable"
        style="width: 100%"
        tooltip-effect="dark"
        :data="tableData"
        row-key="ID"
        @selection-change="handleSelectionChange"
        >
        <!-- 移除选择列 -->
        
<!--        <el-table-column sortable align="left" label="日期" prop="CreatedAt" width="180">-->
<!--            <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>-->
<!--        </el-table-column>-->
        
            <el-table-column align="left" label="电站名称" prop="psid" width="120" />

            <el-table-column align="left" label="设备编号" prop="eqid" width="120" />

            <el-table-column align="left" label="设备类型" prop="eqType" width="120" />

            <el-table-column align="left" label="数据类型" prop="dataType" width="120">
    <template #default="scope">
    {{ filterDict(scope.row.dataType,DataTypeOptions) }}
    </template>
</el-table-column>
            <el-table-column align="left" label="点号" prop="datapoint" width="120" />

            <el-table-column align="left" label="发生时间" prop="happenTime" width="180">
   <template #default="scope">{{ formatDate(scope.row.happenTime) }}</template>
</el-table-column>
            <el-table-column align="left" label="数值" prop="dataValue" width="120" />

            <el-table-column align="left" label="已读" prop="isRead" width="120" />

            <el-table-column align="left" label="名称" prop="name" width="120" />

            <el-table-column align="left" label="恢复时间" prop="recoverTime" width="180">
   <template #default="scope">{{ formatDate(scope.row.recoverTime) }}</template>
</el-table-column>
        <el-table-column align="left" label="操作" fixed="right" :min-width="appStore.operateMinWith">
            <template #default="scope">
            <el-button  type="primary" link class="table-button" @click="getDetails(scope.row)"><el-icon style="margin-right: 5px"><InfoFilled /></el-icon>查看</el-button>
            <!-- 移除编辑按钮，因为记录不可编辑 -->
            <!-- 移除删除按钮，因为记录不可编辑 -->
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
            <el-form-item label="电站名称:" prop="psid">
    <el-input v-model.number="formData.psid" :clearable="true" placeholder="请输入电站名称" />
</el-form-item>
            <el-form-item label="设备编号:" prop="eqid">
    <el-input v-model.number="formData.eqid" :clearable="true" placeholder="请输入设备编号" />
</el-form-item>
            <el-form-item label="设备类型:" prop="eqType">
    <el-input v-model.number="formData.eqType" :clearable="true" placeholder="请输入设备类型" />
</el-form-item>
            <el-form-item label="数据类型:" prop="dataType">
    <el-select v-model="formData.dataType" placeholder="请选择数据类型" style="width:100%" filterable :clearable="true">
        <el-option v-for="(item,key) in DataTypeOptions" :key="key" :label="item.label" :value="item.value" />
    </el-select>
</el-form-item>
            <el-form-item label="点号:" prop="datapoint">
    <el-input v-model="formData.datapoint" :clearable="true" placeholder="请输入点号" />
</el-form-item>
            <el-form-item label="发生时间:" prop="happenTime">
    <el-date-picker v-model="formData.happenTime" type="date" style="width:100%" placeholder="选择日期" :clearable="true" />
</el-form-item>
            <el-form-item label="数值:" prop="dataValue">
    <el-input v-model="formData.dataValue" :clearable="true" placeholder="请输入数值" />
</el-form-item>
            <el-form-item label="已读:" prop="isRead">
    <el-input v-model="formData.isRead" :clearable="true" placeholder="请输入已读" />
</el-form-item>
            <el-form-item label="名称:" prop="name">
    <el-input v-model="formData.name" :clearable="true" placeholder="请输入名称" />
</el-form-item>
            <el-form-item label="恢复时间:" prop="recoverTime">
    <el-date-picker v-model="formData.recoverTime" type="date" style="width:100%" placeholder="选择日期" :clearable="true" />
</el-form-item>
          </el-form>
    </el-drawer>

    <el-drawer destroy-on-close :size="appStore.drawerSize" v-model="detailShow" :show-close="true" :before-close="closeDetailShow" title="查看">
            <el-descriptions :column="1" border>
                    <el-descriptions-item label="电站名称">
    {{ detailForm.psid }}
</el-descriptions-item>
                    <el-descriptions-item label="设备编号">
    {{ detailForm.eqid }}
</el-descriptions-item>
                    <el-descriptions-item label="设备类型">
    {{ detailForm.eqType }}
</el-descriptions-item>
                    <el-descriptions-item label="数据类型">
    {{ detailForm.dataType }}
</el-descriptions-item>
                    <el-descriptions-item label="点号">
    {{ detailForm.datapoint }}
</el-descriptions-item>
                    <el-descriptions-item label="发生时间">
    {{ detailForm.happenTime }}
</el-descriptions-item>
                    <el-descriptions-item label="数值">
    {{ detailForm.dataValue }}
</el-descriptions-item>
                    <el-descriptions-item label="已读">
    {{ detailForm.isRead }}
</el-descriptions-item>
                    <el-descriptions-item label="名称">
    {{ detailForm.name }}
</el-descriptions-item>
                    <el-descriptions-item label="恢复时间">
    {{ detailForm.recoverTime }}
</el-descriptions-item>
            </el-descriptions>
        </el-drawer>

  </div>
</template>

<script setup>
import {
  createSysHisEvent,
  deleteSysHisEvent,
  deleteSysHisEventByIds,
  updateSysHisEvent,
  findSysHisEvent,
  getSysHisEventList
} from '@/api/system/sysHisEvent'

// 全量引入格式化工具 请按需保留
import { getDictFunc, formatDate, formatBoolean, filterDict ,filterDataSource, returnArrImg, onDownloadFile } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive } from 'vue'
import { useAppStore } from "@/pinia"




defineOptions({
    name: 'SysHisEvent'
})

// 提交按钮loading
const btnLoading = ref(false)
const appStore = useAppStore()

// 控制更多查询条件显示/隐藏状态
const showAllQuery = ref(false)

// 自动化生成的字典（可能为空）以及字段
const DataTypeOptions = ref([])
const formData = ref({
            psid: undefined,
            eqid: undefined,
            eqType: undefined,
            dataType: '',
            datapoint: '',
            happenTime: new Date(),
            dataValue: '',
            isRead: '',
            name: '',
            recoverTime: new Date(),
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
  const table = await getSysHisEventList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
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
    DataTypeOptions.value = await getDictFunc('DataType')
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
            deleteSysHisEventFunc(row)
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
      const res = await deleteSysHisEventByIds({ IDs })
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
const updateSysHisEventFunc = async(row) => {
    const res = await findSysHisEvent({ ID: row.ID })
    type.value = 'update'
    if (res.code === 0) {
        formData.value = res.data
        dialogFormVisible.value = true
    }
}


// 删除行
const deleteSysHisEventFunc = async (row) => {
    const res = await deleteSysHisEvent({ ID: row.ID })
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
        psid: undefined,
        eqid: undefined,
        eqType: undefined,
        dataType: '',
        datapoint: '',
        happenTime: new Date(),
        dataValue: '',
        isRead: '',
        name: '',
        recoverTime: new Date(),
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
                  res = await createSysHisEvent(formData.value)
                  break
                case 'update':
                  res = await updateSysHisEvent(formData.value)
                  break
                default:
                  res = await createSysHisEvent(formData.value)
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
  const res = await findSysHisEvent({ ID: row.ID })
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