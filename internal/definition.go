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
// time   : 2020-01-09 10:55 下午
// version: 1.0.0
// desc   : 

package internal

import (
	"reflect"
)

type Definition struct {
	TableName   string
	Model       *Model
	Strategy    Strategy
	Fields      []*Field
	PrimaryKeys []*Field
	Indexes     map[string][]*Index
	ForeignKeys map[string][]*ForeignKey
}

type Model struct {
	Type    reflect.Type
	ElmType reflect.Type
}

type Field struct {
	Name        string
	Type        reflect.Type
	ElmType     reflect.Type
	Column      string
	SQLType     string
	IsPrimary   bool
	NotNull     bool
	Default     interface{}
	Comment     string
}

type Index struct {
	Name   string
	Column string
	Type   IndexType
}

type ForeignKey struct {
	Name      string
	Column    string
	Table     string
	Reference string
}
