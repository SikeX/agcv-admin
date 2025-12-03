import service from '@/utils/request'

/**
 * 获取电站出力图表数据
 * @param {Object} params 查询参数
 * @param {number} params.psid 电站编号
 * @param {number} params.eqid 设备编号
 * @param {string} params.startTime 开始时间
 * @param {string} params.endTime 结束时间
 * @returns {Promise} 图表数据
 */
export const getPowerChartData = (params) => {
  return service({
    url: '/agvcChart/getPowerChartData',
    method: 'get',
    donNotShowLoading: true,
    params: params
  })
}

/**
 * 获取电压和无功图表数据
 * @param {Object} params 查询参数
 * @param {number} params.psid 电站编号
 * @param {number} params.eqid 设备编号
 * @param {string} params.startTime 开始时间
 * @param {string} params.endTime 结束时间
 * @returns {Promise} 图表数据
 */
export const getVoltageReactiveChartData = (params) => {
  return service({
    url: '/agvcChart/getVoltageReactiveChartData',
    method: 'get',
    donNotShowLoading: true,
    params: params
  })
}
