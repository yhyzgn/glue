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
	"github.com/yhyzgn/glue/internal"
	"github.com/yhyzgn/glue/primary"
	"reflect"
	"strings"
)

type Default struct{}

func (*Default) Name() string {
	return "default"
}

func (*Default) Quote(key string) string {
	return fmt.Sprintf(`"%s"`, key)
}

func (d *Default) Quotes(keys ...string) []string {
	if keys == nil || len(keys) == 0 {
		return nil
	}
	for idx := range keys {
		keys[idx] = d.Quote(keys[idx])
	}
	return keys
}

func (*Default) Placeholder(int) string {
	return "?"
}

func (*Default) Insert(executor internal.Executor, command *internal.Command) (sql.Result, error) {
	return nil, nil
}

func (*Default) Update(executor internal.Executor, command *internal.Command) (sql.Result, error) {
	return nil, nil
}

func (*Default) SQLType(field *reflect.StructField) string {
	return ""
}

func (*Default) Database() *internal.Command {
	return internal.NewCommand("SELECT DATABASE()")
}

func (d *Default) HasTable(name string) *internal.Command {
	return internal.NewCommand("SELECT").
		Tab("COUNT(1)").
		Line("FROM").
		Tab("INFORMATION_SCHEMA.TABLES").
		Line("WHERE").
		Tab("table_schema = ?").
		Tab("AND table_name = ?").
		Arguments(d.Database(), name)
}

func (d *Default) HasColumn(table, column string) *internal.Command {
	return internal.NewCommand("SELECT").
		Tab("COUNT(1)").
		Line("FROM").
		Tab("INFORMATION_SCHEMA.COLUMNS").
		Line("WHERE").
		Tab("table_schema = ?").
		Tab("AND table_name = ?").
		Tab("AND column_name = ?").
		Arguments(d.Database(), table, column)
}

func (d *Default) CreateTable(definition *internal.Definition) []*internal.Command {
	if definition == nil || definition.Fields == nil || len(definition.Fields) == 0 {
		return nil
	}
	ln := len(definition.Fields)
	// 一张表只允许一个 AUTO_INCREMENT 主键，需要标识判断
	hasAutoIncrement := false

	cmd := internal.NewCommand("CREATE TABLE").Space(d.Quote(definition.TableName)).Space("(")
	for idx, field := range definition.Fields {
		cmd.TabLine(d.Quote(field.Column)).Space(field.SQLType)
		if field.NotNull {
			cmd.Space("NOT")
		}
		cmd.Space("NULL")
		if field.IsPrimary {
			stg := definition.Strategy
			if !hasAutoIncrement && stg != nil && reflect.TypeOf(stg).Elem() == reflect.TypeOf(primary.AutoIncrement{}) {
				// AutoIncrement
				cmd.Space("AUTO_INCREMENT")
				hasAutoIncrement = true
			}
		}
		if field.Default != nil {
			cmd.Space("DEFAULT").Space(field.Default.(string))
		}
		if field.Comment != "" {
			cmd.Space("COMMENT").Space(fmt.Sprintf("'%s'", field.Comment))
		}
		if idx < ln-1 {
			cmd.Append(",")
		}
	}
	if len(definition.PrimaryKeys) > 0 {
		// 有主键
		keys := make([]string, 0)
		for _, field := range definition.PrimaryKeys {
			keys = append(keys, d.Quote(field.Column))
		}
		cmd.Append(",").TabLine(fmt.Sprintf("PRIMARY KEY(%s)", strings.Join(keys, ", ")))
	}
	if len(definition.Indexes) > 0 {
		// 有索引
		for name, indexes := range definition.Indexes {
			if indexes != nil && len(indexes) > 0 {
				keys := make([]string, 0)
				indexType := indexes[0].Type
				for _, index := range indexes {
					keys = append(keys, d.Quote(index.Column))
				}
				cmd.Append(",")
				switch indexType {
				case internal.IndexUnique:
					cmd.TabLine("UNIQUE ")
					break
				case internal.IndexFullText:
					cmd.TabLine("FULLTEXT ")
					break
				case internal.IndexSpatial:
					cmd.TabLine("SPATIAL ")
					break
				default:
					cmd.TabLine("")
					break
				}
				cmd.Append(fmt.Sprintf("INDEX %s (%s)", d.Quote(name), strings.Join(keys, ", ")))
			}
		}
	}
	if len(definition.ForeignKeys) > 0 {
		// 有外键
		for name, fks := range definition.ForeignKeys {
			if fks != nil && len(fks) > 0 {
				refs := make([]string, 0)
				for _, index := range fks {
					refs = append(refs, d.Quote(index.Reference))
				}
				cmd.Append(",").TabLine(fmt.Sprintf("CONSTRAINT %s FOREIGN KEY (%s) REFERENCES %s (%s)", d.Quote(name), d.Quote(fks[0].Column), d.Quote(fks[0].Table), strings.Join(refs, ", ")))
			}
		}
	}
	cmd.Line(")")

	return []*internal.Command{cmd}
}
