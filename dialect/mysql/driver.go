// Copyright 2019 yhyzgn glue
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
// time   : 2020-01-17 11:04
// version: 1.0.0
// desc   : 

package mysql

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/yhyzgn/glue/internal"
)

type mysql struct {
}

func (*mysql) Name() string {
	return "mysql"
}

func (*mysql) Driver() string {
	return "mysql"
}

func (*mysql) Quote(key string) string {
	return fmt.Sprintf("`%s`", key)
}

func (*mysql) Placeholder(index int) string {
	return "?"
}

func (*mysql) Database() *internal.Command {
	return internal.NewCommand("SELECT DATABASE()")
}
