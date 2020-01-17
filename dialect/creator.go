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

type Creator struct {
	driver internal.Driver
}

func New(driver internal.Driver) *Creator {
	return &Creator{driver: driver}
}

func (c *Creator) GetDriver() internal.Driver {
	return c.driver
}

func (c *Creator) Name() string {
	return c.driver.Name()
}

func (c *Creator) Driver() string {
	return c.driver.Driver()
}

func (c *Creator) Quote(key string) string {
	return c.driver.Quote(key)
}

func (c *Creator) Placeholder(index int) string {
	return c.driver.Placeholder(index)
}

func (c *Creator) Database() string {
	return c.driver.Database()
}

func (*Creator) InsertExecutor(executor internal.Executor, command *internal.Command) (sql.Result, error) {
	return nil, nil
}

func (*Creator) UpdateExecutor(executor internal.Executor, command *internal.Command) (sql.Result, error) {
	return nil, nil
}

func (*Creator) SQLType(field *reflect.StructField) string {
	return ""
}

func (c *Creator) HasTable(name string) *internal.Command {
	return internal.NewCommand("SELECT").
		TabLine("COUNT(*)").
		Line("FROM").
		TabLine("INFORMATION_SCHEMA.TABLES").
		Line("WHERE").
		TabLine("table_schema = ?").
		TabLine("AND table_name = ?").
		Arguments(c.driver.Database(), name)
}

func (c *Creator) CreateTable(definition *internal.Definition) []*internal.Command {
	if definition == nil || definition.Fields == nil || len(definition.Fields) == 0 {
		return nil
	}
	ln := len(definition.Fields)
	// 一张表只允许一个 AUTO_INCREMENT 主键，需要标识判断
	hasAutoIncrement := false

	cmd := internal.NewCommand("CREATE TABLE").Space(c.driver.Quote(definition.TableName)).Space("(")
	for idx, field := range definition.Fields {
		cmd.TabLine(c.driver.Quote(field.Column)).Space(field.SQLType)
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
			keys = append(keys, c.driver.Quote(field.Column))
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
					keys = append(keys, c.driver.Quote(index.Column))
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
				cmd.Append(fmt.Sprintf("INDEX %s (%s)", c.driver.Quote(name), strings.Join(keys, ", ")))
			}
		}
	}
	if len(definition.ForeignKeys) > 0 {
		// 有外键
		for name, fks := range definition.ForeignKeys {
			if fks != nil && len(fks) > 0 {
				refs := make([]string, 0)
				for _, index := range fks {
					refs = append(refs, c.driver.Quote(index.Reference))
				}
				cmd.Append(",").TabLine(fmt.Sprintf("CONSTRAINT %s FOREIGN KEY (%s) REFERENCES %s (%s)", name, c.driver.Quote(fks[0].Column), c.driver.Quote(fks[0].Table), strings.Join(refs, ", ")))
			}
		}
	}
	cmd.Line(")")

	return []*internal.Command{cmd}
}

func (c *Creator) Columns(table string) *internal.Command {
	return internal.NewCommand("SELECT").
		TabLine("*").
		Line("FROM").
		TabLine("INFORMATION_SCHEMA.COLUMNS").
		Line("WHERE").
		TabLine("table_schema = ?").
		TabLine("AND table_name = ?").
		Line("ORDER BY").
		TabLine("ORDINAL_POSITION ASC").
		Arguments(c.driver.Database(), table)
}

func (c *Creator) HasColumn(table, column string) *internal.Command {
	return internal.NewCommand("SELECT").
		TabLine("COUNT(*)").
		Line("FROM").
		TabLine("INFORMATION_SCHEMA.COLUMNS").
		Line("WHERE").
		TabLine("table_schema = ?").
		TabLine("AND table_name = ?").
		TabLine("AND column_name = ?").
		Arguments(c.driver.Database(), table, column)
}

func (c *Creator) ModifyColumn(table, column, rename, tpy, comment string, notNull bool, defValue interface{}) *internal.Command {
	return internal.NewCommand(fmt.Sprintf("ALTER TABLE %v CHANGE COLUMN %v %v %v %v COMMENT '%v'", c.driver.Quote(table), c.driver.Quote(column), c.driver.Quote(rename), tpy, extraOfColumn(notNull, defValue), comment))
}

func (c *Creator) AddColumn(table, column, tpy, comment string, notNull bool, defValue interface{}) *internal.Command {
	return internal.NewCommand(fmt.Sprintf("ALTER TABLE %v ADD COLUMN %v %v %v COMMENT '%v'", c.driver.Quote(table), c.driver.Quote(column), tpy, extraOfColumn(notNull, defValue), comment))
}

func (c *Creator) DropColumn(table, column string) *internal.Command {
	return internal.NewCommand(fmt.Sprintf("ALTER TABLE %v DROP COLUMN %v", c.driver.Quote(table), c.driver.Quote(column)))
}

func (c *Creator) HasIndex(table, name string) *internal.Command {
	return internal.NewCommand("SELECT").
		TabLine("COUNT(*)").
		Line("FROM").
		TabLine("INFORMATION_SCHEMA.STATISTICS").
		Line("WHERE").
		TabLine("table_schema = ?").
		TabLine("AND table_schema = ?").
		TabLine("AND index_name = ?").
		Arguments(c.driver.Database(), c.driver.Quote(table), c.driver.Quote(name))
}

func (c *Creator) RemoveIndex(table, name string) *internal.Command {
	return internal.NewCommand(fmt.Sprintf("DROP INDEX %v", c.driver.Quote(name)))
}

func (c *Creator) HasForeignKey(table, name string) *internal.Command {
	return internal.NewCommand("SELECT").
		TabLine("COUNT(*)").
		Line("FROM").
		TabLine("INFORMATION_SCHEMA.TABLE_CONSTRAINTS").
		Line("WHERE").
		TabLine("CONSTRAINT_SCHEMA = ?").
		TabLine("AND TABLE_NAME = ?").
		TabLine("AND CONSTRAINT_NAME = ?").
		TabLine("AND CONSTRAINT_TYPE = 'FOREIGN KEY'").
		Arguments(c.driver.Database(), table, name)
}

func (c *Creator) AddForeignKey(table string, key *internal.ForeignKey) *internal.Command {
	return internal.NewCommand(fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT %s FOREIGN KEY (%s) REFERENCES %s(%s)", c.driver.Quote(table), key.Name, key.Column, c.driver.Quote(key.Table), key.Reference))
}

func (c *Creator) RemoveForeignKey(table, name string) *internal.Command {
	return internal.NewCommand(fmt.Sprintf("ALTER TABLE %s DROP FOREIGN KEY (%s)", c.driver.Quote(table), name))
}

func (*Creator) DefaultValue() string {
	return "DEFAULT VALUES"
}

func (*Creator) BuildKeyName(kind, table string, fields ...string) string {
	return fmt.Sprintf("%s_%s,%s", kind, table, strings.Join(fields, "_"))
}

func (c *Creator) Insert(value *internal.ExecValue) *internal.Command {
	if value == nil {
		return nil
	}
	return internal.NewCommand(fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", c.driver.Quote(value.Table), )).Arguments(value.Values...)
}

func (*Creator) Delete(value *internal.ExecValue) *internal.Command {
	return nil
}

func (*Creator) Remove(value *internal.ExecValue) *internal.Command {
	return nil
}

func (*Creator) Update(value *internal.ExecValue) *internal.Command {
	return nil
}

func (*Creator) Select(value *internal.ExecValue) *internal.Command {
	return nil
}

func (*Creator) Count(cmd *internal.Command) *internal.Command {
	return internal.NewCommand(fmt.Sprintf("SELECT COUNT(*) FROM (%s) AS T", cmd.SQL())).Arguments(cmd.Args()...)
}

func (c *Creator) Page(cmd *internal.Command, page, size int) *internal.Command {
	if page <= 0 {
		page = 0
	}
	offset := (page - 1) * size
	args := append(cmd.Args(), offset, size)
	return internal.NewCommand(fmt.Sprintf("SELECT * FROM (%s) AS T LIMIT ?, ?", cmd.SQL())).Arguments(args...)
}

func extraOfColumn(notNull bool, defValue interface{}) (result string) {
	if notNull {
		result += "NOT "
	}
	result += "NULL"

	if defValue != nil {
		result += " DEFAULT "
		tp := reflect.ValueOf(defValue)
		value := fmt.Sprintf("%v", defValue)
		if tp.Kind() == reflect.String {
			result += fmt.Sprintf("'%v'", value)
		} else {
			result += value
		}
	}

	return
}
