package monitor

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/monitor"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SysGridConnectionPointHistoryApi struct {}

// GetSysGridConnectionPointHistoryList 获取并网点监控列表
// @Tags SysGridConnectionPointHistory
// @Summary 获取并网点监控列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=[]monitor.SysGridConnectionPointHistory,msg=string} "获取成功"
// @Router /sysGridConnectionPointHistory/getSysGridConnectionPointHistoryList [get]
func (sysGridConnectionPointHistoryApi *SysGridConnectionPointHistoryApi) GetSysGridConnectionPointHistoryList(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	list, err := sysGridConnectionPointHistoryService.GetSysGridConnectionPointHistoryList(ctx)
	if err != nil {
	    global.GVA_LOG.Error("获取失败!", zap.Error(err))
        response.FailWithMessage("获取失败:" + err.Error(), c)
        return
    }
    response.OkWithDetailed(list, "获取成功", c)
}

// GetGridPointHistoryData 获取并网点历史数据
// @Tags SysGridConnectionPointHistory
// @Summary 获取并网点历史数据
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body monitor.GridPointHistoryQuery true "获取并网点历史数据"
// @Success 200 {object} response.Response{data=[]monitor.GridPointHistoryResponse,msg=string} "获取成功"
// @Router /sysGridConnectionPointHistory/getGridPointHistoryData [post]
func (sysGridConnectionPointHistoryApi *SysGridConnectionPointHistoryApi) GetGridPointHistoryData(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	var query monitor.GridPointHistoryQuery
	err := c.ShouldBindJSON(&query)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 设置默认时间范围（最近24小时）
	if query.StartTime == "" {
		query.StartTime = time.Now().Add(-24 * time.Hour).Format(time.RFC3339)
	}
	if query.EndTime == "" {
		query.EndTime = time.Now().Format(time.RFC3339)
	}

	data, err := sysGridConnectionPointHistoryService.GetGridPointHistoryData(ctx, query)
	if err != nil {
	    global.GVA_LOG.Error("获取历史数据失败!", zap.Error(err))
        response.FailWithMessage("获取历史数据失败:" + err.Error(), c)
        return
    }
    
    response.OkWithDetailed(data, "获取成功", c)
}

// GenerateTestData 生成并网点测试数据
// @Tags SysGridConnectionPointHistory
// @Summary 生成并网点测试数据
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param GridPointID query int false "并网点ID（可选，如果不提供则随机选择）"
// @Param pointID query int false "点标识ID（可选，如果不提供则随机选择）"
// @Param count query int false "生成数据点的数量，默认为100"
// @Success 200 {object} response.Response{msg=string} "生成成功"
// @Router /sysGridConnectionPointHistory/generateTestData [post]
func (sysGridConnectionPointHistoryApi *SysGridConnectionPointHistoryApi) GenerateTestData(c *gin.Context) {
	// 检查InfluxDB客户端是否已初始化
	if global.GVA_INFLUXDB == nil {
		response.FailWithMessage("InfluxDB客户端未初始化", c)
		return
	}

	gridPointID := c.Query("GridPointID")
	pointID := c.Query("pointID")
	
	count := 100
	if c.Query("count") != "" {
		fmt.Sscanf(c.Query("count"), "%d", &count)
	}

	// 获取写API
	writeAPI := global.GVA_INFLUXDB.WriteAPIBlocking(global.GVA_CONFIG.InfluxDB.Org, global.GVA_CONFIG.InfluxDB.Bucket)

	// 生成测试数据
	err := sysGridConnectionPointHistoryApi.generateAndWriteTestData(writeAPI, gridPointID, pointID, count)
	if err != nil {
		global.GVA_LOG.Error("生成测试数据失败", zap.Error(err))
		response.FailWithMessage("生成测试数据失败: " + err.Error(), c)
		return
	}

	response.OkWithMessage(fmt.Sprintf("成功生成%d个测试数据点", count), c)
}

// generateAndWriteTestData 生成并写入测试数据
func (sysGridConnectionPointHistoryApi *SysGridConnectionPointHistoryApi) generateAndWriteTestData(writeAPI api.WriteAPIBlocking, gridPointID, pointID string, count int) error {
	// 根据并网点ID获取对应的点标识列表
	var targetPointIDs []int
	
	if gridPointID != "" {
		// 如果指定了并网点ID，则获取该并网点的所有点标识
		gridPointIDInt := 0
		fmt.Sscanf(gridPointID, "%d", &gridPointIDInt)
		targetPointIDs = monitor.GetPointIdentifiersByGridPoint(gridPointIDInt)
	} else if pointID != "" {
		// 如果指定了点标识ID，则只使用该点标识
		pointIDInt := 0
		fmt.Sscanf(pointID, "%d", &pointIDInt)
		targetPointIDs = []int{pointIDInt}
	} else {
		// 如果都没有指定，则使用所有点标识
		targetPointIDs = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43}
	}
	
	if len(targetPointIDs) == 0 {
		return fmt.Errorf("没有找到有效的点标识")
	}
	
	// 生成指定数量的数据点
	now := time.Now()
	for i := 0; i < count; i++ {
		timestamp := now.Add(-time.Duration(count-i) * time.Minute)
		
		// 为每个目标点标识生成数据
		for _, currentPointID := range targetPointIDs {
			// 生成模拟数值，使用正弦函数创建波动数据
			value := 100 + 50*math.Sin(float64(i)/10) + rand.Float64()*10 - 5
			
			// 创建数据点
			p := influxdb2.NewPoint(
				"mqtt",
				map[string]string{
					"point_id": fmt.Sprintf("%d", currentPointID),
				},
				map[string]interface{}{
					"value": value,
				},
				timestamp,
			)
			
			// 写入数据点
			ctx := context.Background()
			err := writeAPI.WritePoint(ctx, p)
			if err != nil {
				return fmt.Errorf("写入InfluxDB失败: %v", err)
			}
		}
	}

	return nil
}