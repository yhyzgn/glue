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
// time   : 2020-01-09 9:38 下午
// version: 1.0.0
// desc   : 

package internal

import (
	"database/sql"
	"github.com/yhyzgn/glue/external"
	"reflect"
)

type Dialect interface {
	Name() string

	Quote(key string) string

	Placeholder(index int) string

	Insert(executor Executor, command *external.Command) (sql.Result, error)

	Update(executor Executor, command *external.Command) (sql.Result, error)

	SQLType(field *reflect.StructField) string

	Database(executor Executor) string

	HasTable(executor Executor, name string) bool

	HasColumn(executor Executor, table, column string) bool

	CreateTable()
}
