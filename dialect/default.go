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
// time   : 2020-01-09 10:12 下午
// version: 1.0.0
// desc   : 

package dialect

import (
	"database/sql"
	"fmt"
	"github.com/yhyzgn/glue/external"
	"github.com/yhyzgn/glue/internal"
	"reflect"
)

type Default struct{}

func (*Default) Name() string {
	return "default"
}

func (*Default) Quote(key string) string {
	return fmt.Sprintf(`"%s"`, key)
}

func (*Default) Placeholder(index int) string {
	return "?"
}

func (*Default) Insert(executor internal.Executor, command *external.Command) (sql.Result, error) {
	return nil, nil
}

func (*Default) Update(executor internal.Executor, command *external.Command) (sql.Result, error) {
	return nil, nil
}

func (*Default) SQLType(field *reflect.StructField) string {
	return ""
}

func (*Default) Database(executor internal.Executor) string {
	var name string
	_ = executor.QueryRow("SELECT DATABASE()").Scan(&name)
	return name
}

func (d *Default) HasTable(executor internal.Executor, name string) bool {
	database := d.Database(executor)
	var count int
	_ = executor.QueryRow("SELECT count(1) FROM INFORMATION_SCHEMA.TABLES WHERE table_schema = ? AND table_name = ?", database, name).Scan(&count)
	return count > 0
}

func (d *Default) HasColumn(executor internal.Executor, table, column string) bool {
	database := d.Database(executor)
	var count int
	_ = executor.QueryRow("SELECT count(1) FROM INFORMATION_SCHEMA.COLUMNS WHERE table_schema = ? AND table_name = ? AND column_name = ?", database, table, column).Scan(&count)
	return count > 0
}

func (*Default) CreateTable() {}
