import service from '@/utils/request'

/**
 * 获取并网点监控列表
 * @returns {Promise} 并网点监控列表数据
 */
export const getSysGridConnectionPointHistoryList = () => {
  return service({
    url: '/sysGridConnectionPointHistory/getSysGridConnectionPointHistoryList',
    method: 'get'
  })
}

/**
 * 获取并网点历史数据
 * @param {object} data 查询参数
 * @param {number} data.GridPointID 并网点ID
 * @param {string} data.StartTime 开始时间
 * @param {string} data.EndTime 结束时间
 * @returns {Promise} 并网点历史数据
 */
export const getGridPointHistoryData = (data) => {
  return service({
    url: '/sysGridConnectionPointHistory/getGridPointHistoryData',
    method: 'post',
    data: data
  })
}

/**
 * 创建并网点配置
 * @param {object} data 并网点配置数据
 * @returns {Promise} 创建结果
 */
export const createSysGridConnectionPointHistory = (data) => {
  return service({
    url: '/sysGridConnectionPointHistory/createSysGridConnectionPointHistory',
    method: 'post',
    data: data
  })
}

/**
 * 更新并网点配置
 * @param {object} data 并网点配置数据
 * @returns {Promise} 更新结果
 */
export const updateSysGridConnectionPointHistory = (data) => {
  return service({
    url: '/sysGridConnectionPointHistory/updateSysGridConnectionPointHistory',
    method: 'put',
    data: data
  })
}

/**
 * 查找并网点配置
 * @param {object} data 查询参数
 * @returns {Promise} 并网点配置数据
 */
export const findSysGridConnectionPointHistory = (data) => {
  return service({
    url: '/sysGridConnectionPointHistory/findSysGridConnectionPointHistory',
    method: 'get',
    params: data
  })
}

/**
 * 生成测试数据
 * @param {object} params 参数
 * @param {number} params.count 生成数据点的数量
 * @returns {Promise} 生成结果
 */
export const generateTestData = (params) => {
  return service({
    url: '/sysGridConnectionPointHistory/generateTestData',
    method: 'post',
    params: params
  })
}