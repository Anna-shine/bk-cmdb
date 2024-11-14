/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package y3_6_201911261109

import (
	"context"
	"fmt"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

func initInnerChart(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	idArr := make([]uint64, 0)
	idArr, err := db.NextSequences(ctx, common.BKTableNameChartConfig, len(metadata.InnerChartsArr))
	if err != nil {
		return fmt.Errorf("get next sequences failed, tableName: %s, err: %+v", common.BKTableNameChartConfig, err)
	}

	for index, chart := range InnerChartsArr {
		innerChart := InnerChartsMap[chart]
		innerChart.ConfigID = idArr[index]
		innerChart.CreateTime.Time = time.Now()
		innerChart.OwnerID = conf.OwnerID
		if err := db.Table(common.BKTableNameChartConfig).Insert(ctx, innerChart); err != nil {
			return fmt.Errorf("insert chart config failed, tableName: %s, chart: %+v, err: %+v",
				common.BKTableNameChartConfig, innerChart, err)
		}
	}

	position := ChartPosition{
		BizID: 0,
		Position: metadata.PositionInfo{
			Host: idArr[2:6],
			Inst: idArr[6:],
		},
		OwnerID: "0",
	}

	if err := db.Table(common.BKTableNameChartPosition).Insert(ctx, position); err != nil {
		return fmt.Errorf("insert cahrt position data failed, table: %s, position: %+v, err: %s",
			common.BKTableNameChartPosition, position, err)
	}

	return nil
}

// ChartConfig TODO
type ChartConfig struct {
	ConfigID   uint64        `json:"config_id" bson:"config_id"`
	ReportType string        `json:"report_type" bson:"report_type"`
	Name       string        `json:"name" bson:"name"`
	CreateTime metadata.Time `json:"create_time" bson:"create_time"`
	OwnerID    string        `json:"bk_supplier_account" bson:"bk_supplier_account"`
	ObjID      string        `json:"bk_obj_id" bson:"bk_obj_id"`
	Width      string        `json:"width" bson:"width"`
	ChartType  string        `json:"chart_type" bson:"chart_type"`
	Field      string        `json:"field" bson:"field"`
	XAxisCount int64         `json:"x_axis_count" bson:"x_axis_count"`
}

var (
	// BizModuleHostChart TODO
	BizModuleHostChart = ChartConfig{
		ReportType: common.BizModuleHostChart,
	}

	// HostOsChart TODO
	HostOsChart = ChartConfig{
		ReportType: common.HostOSChart,
		Name:       "按操作系统类型统计",
		ObjID:      "host",
		Width:      "50",
		ChartType:  "pie",
		Field:      "bk_os_type",
		XAxisCount: 10,
	}

	// HostBizChart TODO
	HostBizChart = ChartConfig{
		ReportType: common.HostBizChart,
		Name:       "按业务统计",
		ObjID:      "host",
		Width:      "50",
		ChartType:  "bar",
		XAxisCount: 10,
	}

	// HostCloudChart TODO
	HostCloudChart = ChartConfig{
		ReportType: common.HostCloudChart,
		Name:       "按管控区域统计",
		Width:      "100",
		ObjID:      "host",
		ChartType:  "bar",
		Field:      common.BKCloudIDField,
		XAxisCount: 20,
	}

	// HostChangeBizChart TODO
	HostChangeBizChart = ChartConfig{
		ReportType: common.HostChangeBizChart,
		Name:       "主机数量变化趋势",
		Width:      "100",
		XAxisCount: 20,
	}

	// ModelAndInstCountChart TODO
	ModelAndInstCountChart = ChartConfig{
		ReportType: common.ModelAndInstCount,
	}

	// ModelInstChart TODO
	ModelInstChart = ChartConfig{
		ReportType: common.ModelInstChart,
		Name:       "实例数量统计",
		Width:      "50",
		ChartType:  "bar",
		XAxisCount: 10,
	}

	// ModelInstChangeChart TODO
	ModelInstChangeChart = ChartConfig{
		ReportType: common.ModelInstChangeChart,
		Name:       "实例变更统计",
		Width:      "50",
		ChartType:  "bar",
		XAxisCount: 10,
	}

	// InnerChartsMap TODO
	InnerChartsMap = map[string]ChartConfig{
		common.BizModuleHostChart:   BizModuleHostChart,
		common.ModelAndInstCount:    ModelAndInstCountChart,
		common.HostOSChart:          HostOsChart,
		common.HostBizChart:         HostBizChart,
		common.HostCloudChart:       HostCloudChart,
		common.HostChangeBizChart:   HostChangeBizChart,
		common.ModelInstChart:       ModelInstChart,
		common.ModelInstChangeChart: ModelInstChangeChart,
	}

	// InnerChartsArr TODO
	InnerChartsArr = []string{
		common.BizModuleHostChart,
		common.ModelAndInstCount,
		common.HostOSChart,
		common.HostBizChart,
		common.HostCloudChart,
		common.HostChangeBizChart,
		common.ModelInstChart,
		common.ModelInstChangeChart,
	}
)

// ChartPosition TODO
type ChartPosition struct {
	BizID    int64                 `json:"bk_biz_id" bson:"bk_biz_id"`
	Position metadata.PositionInfo `json:"position" bson:"position"`
	OwnerID  string                `json:"bk_supplier_account" bson:"bk_supplier_account"`
}
