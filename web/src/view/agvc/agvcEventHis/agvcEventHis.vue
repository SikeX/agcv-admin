
<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" class="demo-form-inline" @keyup.enter="onSubmit">
      <el-form-item label="电站名称" prop="psid">
        <el-select v-model="searchInfo.psid" clearable placeholder="请选择电站名称" filterable>
          <el-option
            v-for="item in psidOptions"
            :key="item.value"
            :label="item.label"
            :value="item.value"
          />
        </el-select>
      </el-form-item>
      
      <el-form-item label="设备类型" prop="eqType">
        <el-select v-model="searchInfo.eqType" clearable placeholder="请选择设备类型" filterable>
          <el-option
            v-for="item in DeviceTypeOptions"
            :key="item.value"
            :label="item.label"
            :value="item.value"
          />
        </el-select>
      </el-form-item>
      
      <el-form-item label="设备编号" prop="eqid">
        <el-input v-model="searchInfo.eqid" placeholder="请输入设备编号" clearable />
      </el-form-item>
      
      <el-form-item label="发生时间" prop="recordTimeRange">
        <el-date-picker
          v-model="searchInfo.recordTimeRange"
          type="datetimerange"
          range-separator="至"
          start-placeholder="开始时间"
          end-placeholder="结束时间"
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
            <!-- 历史事件记录不可新增、编辑、删除 -->
        </div>
        <el-table
        ref="multipleTable"
        style="width: 100%"
        tooltip-effect="dark"
        :data="tableData"
        row-key="ID"
        >
        <!-- 移除选择列 -->
        
            <el-table-column align="left" label="电站名称" prop="psid" min-width="100">
    <template #default="scope">
    {{ filterDict(scope.row.psid,psidOptions) }}
    </template>
</el-table-column>
            <el-table-column align="left" label="设备编号" prop="eqid" min-width="100" />

            <el-table-column align="left" label="设备类型" prop="eqType" min-width="100">
    <template #default="scope">
    {{ filterDict(scope.row.eqType,DeviceTypeOptions) }}
    </template>
</el-table-column>
            <el-table-column align="left" label="数据类型" prop="dataType" min-width="100">
    <template #default="scope">
    {{ filterDict(scope.row.dataType,DataTypeOptions) }}
    </template>
</el-table-column>
            <el-table-column align="left" label="点号" prop="datapoint" min-width="80" />

            <el-table-column align="left" label="发生时间" prop="RecordTime" min-width="150">
   <template #default="scope">{{ formatDate(scope.row.RecordTime) }}</template>
</el-table-column>
            <el-table-column align="left" label="数值" prop="dataValue" min-width="100">
    <template #default="scope">
    {{ filterDict(scope.row.dataValue,DataValueOptions) }}
    </template>
</el-table-column>

            <el-table-column align="left" label="已读" prop="isRead" min-width="80">
    <template #default="scope">
    {{ filterDict(scope.row.isRead,IsReadOptions) }}
    </template>
</el-table-column>

        <el-table-column align="left" label="操作" fixed="right" width="100">
            <template #default="scope">
            <el-button  type="primary" link class="table-button" @click="getDetails(scope.row)"><el-icon style="margin-right: 5px"><InfoFilled /></el-icon>查看</el-button>
            <!-- 历史事件记录不可编辑、删除 -->
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

    <el-drawer destroy-on-close :size="appStore.drawerSize" v-model="detailShow" :show-close="true" :before-close="closeDetailShow" title="查看">
            <el-descriptions :column="1" border>
                    <el-descriptions-item label="电站名称">
    {{ filterDict(detailForm.psid, psidOptions) }}
</el-descriptions-item>
                    <el-descriptions-item label="设备编号">
    {{ detailForm.eqid }}
</el-descriptions-item>
                    <el-descriptions-item label="设备类型">
    {{ filterDict(detailForm.eqType, DeviceTypeOptions) }}
</el-descriptions-item>
                    <el-descriptions-item label="数据类型">
    {{ filterDict(detailForm.dataType, DataTypeOptions) }}
</el-descriptions-item>
                    <el-descriptions-item label="点号">
    {{ detailForm.datapoint }}
</el-descriptions-item>
                    <el-descriptions-item label="发生时间/恢复时间">
    {{ formatDate(detailForm.RecordTime) }}
</el-descriptions-item>
                    <el-descriptions-item label="数值">
    {{ filterDict(detailForm.dataValue, DataValueOptions) }}
</el-descriptions-item>
                    <el-descriptions-item label="已读">
    {{ filterDict(detailForm.isRead, IsReadOptions) }}
</el-descriptions-item>
            </el-descriptions>
        </el-drawer>

  </div>
</template>

<script setup>
import {
  getAgvcEventHisList,
  findAgvcEventHis
} from '@/api/agvc/agvcEventHis'

// 全量引入格式化工具 请按需保留
import { getDictFunc, formatDate, filterDict } from '@/utils/format'
import { ref, reactive } from 'vue'
import { useAppStore } from "@/pinia"


defineOptions({
    name: 'AgvcEventHis'
})

const appStore = useAppStore()

// 控制更多查询条件显示/隐藏状态
const showAllQuery = ref(false)

// 自动化生成的字典（可能为空）以及字段
const DataTypeOptions = ref([])
const psidOptions = ref([])
const DeviceTypeOptions = ref([])
const DataValueOptions = ref([])  // 数值字典 (发生或恢复)
const IsReadOptions = ref([])     // 已读字典
const elSearchFormRef = ref([])

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
  const table = await getAgvcEventHisList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
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
    psidOptions.value = await getDictFunc('psid')
    DeviceTypeOptions.value = await getDictFunc('DeviceType')
    DataValueOptions.value = await getDictFunc('HappenOrRecover')
    IsReadOptions.value = await getDictFunc('YesOrNo') 
}

// 获取需要的字典 可能为空 按需保留
setOptions()

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
  const res = await findAgvcEventHis({ ID: row.ID })
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
