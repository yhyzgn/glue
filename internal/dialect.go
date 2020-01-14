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
	"reflect"
)

type Dialect interface {
	Name() string

	Driver() string

	Quote(key string) string

	Quotes(keys ...string) []string

	Placeholder(index int) string

	InsertExecutor(executor Executor, command *Command) (sql.Result, error)

	UpdateExecutor(executor Executor, command *Command) (sql.Result, error)

	SQLType(field *reflect.StructField) string

	Database() *Command

	HasTable(name string) *Command

	CreateTable(definition *Definition) []*Command

	Columns(table string) *Command

	HasColumn(table, column string) *Command

	ModifyColumn(table, column, rename, tpy, comment string, notNull bool, defValue interface{}) *Command

	AddColumn(table, column, tpy, comment string, notNull bool, defValue interface{}) *Command

	DropColumn(table, column string) *Command

	HasIndex(table, name string) *Command

	RemoveIndex(table, name string) *Command

	HasForeignKey(table, name string) *Command

	AddForeignKey(table string, key *ForeignKey) *Command

	RemoveForeignKey(table, name string) *Command

	DefaultValue() string

	BuildKeyName(kind, table string, fields ...string) string

	Insert() *Command

	Delete() *Command

	Remove() *Command

	Update() *Command

	Select() *Command

	Count(cmd *Command) *Command

	Page(cmd *Command, page, size int) *Command
}
