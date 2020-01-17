// Copyright 2020 yhyzgn glue
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// author : 颜洪毅
// e-mail : yhyzgn@gmail.com
// time   : 2020-01-11 3:51 下午
// version: 1.0.0
// desc   : 

package dialect

import (
	"fmt"
	"github.com/yhyzgn/glue/internal"
	"github.com/yhyzgn/glue/primary"
	"testing"
)

func TestDefault_CreateTable(t *testing.T) {
	dfs := &internal.Definition{
		TableName: "user",
		Model:     nil,
		Strategy:  &primary.AutoIncrement{},
		Fields: []*internal.Field{
			{
				Name:      "ID",
				Type:      nil,
				ElmType:   nil,
				Column:    "id",
				SQLType:   "BIGINT",
				IsPrimary: true,
				NotNull:   true,
				Default:   nil,
				Comment:   "主键ID",
			},
			{
				Name:      "Code",
				Type:      nil,
				ElmType:   nil,
				Column:    "code",
				SQLType:   "VARCHAR(100)",
				IsPrimary: true,
				NotNull:   true,
				Default:   nil,
				Comment:   "用户编号",
			},
			{
				Name:      "Name",
				Type:      nil,
				ElmType:   nil,
				Column:    "name",
				SQLType:   "VARCHAR(255)",
				IsPrimary: false,
				NotNull:   false,
				Default:   nil,
				Comment:   "姓名",
			},
			{
				Name:      "Age",
				Type:      nil,
				ElmType:   nil,
				Column:    "age",
				SQLType:   "INT",
				IsPrimary: false,
				NotNull:   true,
				Default:   "0",
				Comment:   "年龄",
			},
		},
		PrimaryKeys: []*internal.Field{
			{
				Column: "id",
			},
			{
				Column: "code",
			},
		},
		Indexes: map[string][]*internal.Index{
			"index_normal": {
				{
					Name:   "index_normal",
					Column: "age",
					Type:   internal.IndexNormal,
				},
				{
					Name:   "index_normal",
					Column: "name",
					Type:   internal.IndexNormal,
				},
			},
			"index_unique": {
				{
					Name:   "index_unique",
					Column: "code",
					Type:   internal.IndexUnique,
				},
			},
			"index_fulltext": {
				{
					Name:   "index_fulltext",
					Column: "name",
					Type:   internal.IndexFullText,
				},
			},
		},
		ForeignKeys: map[string][]*internal.ForeignKey{},
	}

	dft := Creator{}

	cmd := dft.CreateTable(dfs)

	fmt.Println(cmd[0].SQL())
}
