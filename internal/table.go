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
// time   : 2020-01-09 10:35 下午
// version: 1.0.0
// desc   : 

package internal

type Table interface {
	TableName() string

	PrimaryStrategy() Strategy

	BeforeInsert()

	BeforeUpdate()
}

type TableModel struct{}

func (*TableModel) TableName() string {
	return ""
}

func (*TableModel) PrimaryStrategy() Strategy {
	return nil
}

func (*TableModel) BeforeInsert() {
}

func (*TableModel) BeforeUpdate() {
}
