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
// time   : 2020-01-09 9:45 下午
// version: 1.0.0
// desc   : 

package internal

type Command struct {
	SQL  string
	Args []interface{}
}

func NewCommand(sql string) *Command {
	return &Command{SQL: sql}
}

func (c *Command) Append(sql string) *Command {
	c.SQL += sql
	return c
}

func (c *Command) Space(sql string) *Command {
	c.SQL += " " + sql
	return c
}

func (c *Command) Line(sql string) *Command {
	c.SQL += "\n" + sql
	return c
}

func (c *Command) Tab(sql string) *Command {
	return c.Tabs(sql, 1)
}

func (c *Command) Tabs(sql string, tabs int) *Command {
	c.SQL += c.tabs(sql, tabs)
	return c
}

func (c *Command) TabLine(sql string) *Command {
	return c.TabsLine(sql, 1)
}

func (c *Command) TabsLine(sql string, tabs int) *Command {
	return c.Line(c.tabs(sql, tabs))
}

func (c *Command) Arguments(args ...interface{}) *Command {
	if args != nil && len(args) > 0 {
		if c.Args == nil {
			c.Args = make([]interface{}, 0)
		}
		c.Args = append(c.Args, args...)
	}
	return c
}

func (c *Command) tabs(sql string, tabs int) string {
	for i := 0; i < tabs; i++ {
		sql = "\t" + sql
	}
	return sql
}
