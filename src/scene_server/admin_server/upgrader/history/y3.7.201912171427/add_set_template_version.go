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

package y3_7_201912171427

import (
	"context"
	"fmt"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

func addSetTemplateDefaultVersion(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	filter := map[string]interface{}{
		"version": map[string]interface{}{
			common.BKDBExists: false,
		},
	}
	doc := map[string]interface{}{
		"version": 0,
	}
	if err := db.Table(common.BKTableNameSetTemplate).Update(ctx, filter, doc); err != nil {
		return fmt.Errorf("addSetTemplateDefaultVersion failed, err: %+v", err)
	}
	return nil
}

func addSetDefaultVersion(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	filter := map[string]interface{}{
		"set_template_version": map[string]interface{}{
			common.BKDBExists: false,
		},
	}
	doc := map[string]interface{}{
		"set_template_version": 0,
	}
	if err := db.Table(common.BKTableNameBaseSet).Update(ctx, filter, doc); err != nil {
		return fmt.Errorf("addSetDefaultVersion failed, err: %+v", err)
	}
	return nil
}

func addSetVersionField(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	filter := map[string]interface{}{
		common.BKObjIDField:      common.BKInnerObjIDSet,
		common.BKPropertyIDField: "set_template_version",
	}
	count, err := db.Table(common.BKTableNameObjAttDes).Find(filter).Count(ctx)
	if err != nil {
		return fmt.Errorf("check whether set_template_version attribute exist failed, err: %+v", err)
	}
	if count != 0 {
		return nil
	}

	id, err := db.NextSequence(ctx, common.BKTableNameObjAttDes)
	if err != nil {
		return fmt.Errorf("generate attribute id failed, err: %+v", err)
	}

	now := metadata.Now()
	attribute := map[string]interface{}{
		"id":                     int64(id),
		"bk_supplier_account":    conf.OwnerID,
		"bk_obj_id":              common.BKInnerObjIDSet,
		"bk_property_id":         "set_template_version",
		"bk_property_name":       "集群模板",
		"bk_property_group":      "default",
		"bk_property_group_name": "default",
		"bk_property_index":      0,
		"unit":                   "",
		"placeholder":            "",
		"editable":               true,
		"ispre":                  true,
		"isrequired":             false,
		"isreadonly":             true,
		"isonly":                 false,
		// IsSystem = true 时，字段标记系统内部使用的字段，不会返回到前端
		"bk_issystem": true,
		// IsAPI = true 时，字段对页面不可见
		"bk_isapi":         true,
		"bk_property_type": "int",
		"option":           "",
		"description":      "集群版本，从通集群模板同步",
		"creator":          conf.User,
		"create_time":      &now,
		"last_time":        &now,
	}
	uniqueFields := []string{common.BKObjIDField, common.BKPropertyIDField, "bk_supplier_account"}
	if _, _, err := upgrader.Upsert(ctx, db, common.BKTableNameObjAttDes, attribute, "id", uniqueFields,
		[]string{}); err != nil {
		blog.Errorf("addSetVersionField failed, add set_template_version attribute failed, err: %+v", err)
		return fmt.Errorf("add set_template_version attribute failed, err: %+v", err)
	}
	return nil
}
