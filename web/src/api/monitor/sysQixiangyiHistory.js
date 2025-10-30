import service from '@/utils/request'

/**
 * 获取气象仪配置列表
 * @returns {Promise} 气象仪配置列表
 */
export const getSysQixiangyiHistoryList = () => {
  return service({
    url: '/sysQixiangyiHistory/getSysQixiangyiHistoryList',
    method: 'get'
  })
}

/**
 * 获取气象仪历史数据
 * @param {Object} data 查询参数
 * @param {number} data.qixiangyiId 气象仪ID
 * @param {string} data.startTime 开始时间
 * @param {string} data.endTime 结束时间
 * @returns {Promise} 历史数据
 */
export const getQixiangyiHistoryData = (data) => {
  return service({
    url: '/sysQixiangyiHistory/getQixiangyiHistoryData',
    method: 'get',
    params: data
  })
}

/**
 * 生成测试数据
 * @param {Object} data 参数
 * @param {number} data.qixiangyiId 气象仪ID
 * @returns {Promise} 生成结果
 */
export const generateTestData = (data) => {
  return service({
    url: '/sysQixiangyiHistory/generateTestData',
    method: 'post',
    params: data  // 改为params，通过查询参数传递
  })
}
