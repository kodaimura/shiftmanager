package main

import (
	"io"
	"os"
	"fmt"
	"strings"
	
	"github.com/kodaimura/ddlparse"

	"shiftmanager/config"
)

var cf *config.Config

func main() {
	cf = config.GetConfig()
	args := os.Args

	if len(args) < 2 {
		fmt.Println("Error: Please specify the path to the DDL as the 1st argument")
		return
	}
	generate(args[1:])
}

func generate(args []string) {
	generateModel(args)
	generateRepository(args)
}

func readFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(fmt.Sprintf("Error: Failed to retrieve the file %s. \n%s", path, err.Error()))
		return "", err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		fmt.Println(fmt.Sprintf("Error: Failed to read the file %s. \n%s", path, err.Error()))
		return "", err
	}
	return string(content), nil	
}

func parseDDL(ddl string) ([]ddlparse.Table, error) {
	var tables []ddlparse.Table
	var err error
	if (cf.DBDriver == "postgres") {
		tables, err = ddlparse.ParsePostgreSQL(ddl)
	} else if (cf.DBDriver == "mysql") {
		tables, err = ddlparse.ParseMySQL(ddl)
	} else if (cf.DBDriver == "sqlite3") {
		tables, err = ddlparse.ParseSQLite(ddl)
	} else {
		fmt.Println(fmt.Sprintf("Error: DB_DRIVER=%s is not supported.", cf.DBDriver))
	}
	
	return tables, err
}

func writeFile(path, content string) error {
	file, err := os.Create(path)
	if err != nil {
		fmt.Println(fmt.Sprintf("Error: Failed to create the file %s. \n%s", path, err.Error()))
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		fmt.Println(fmt.Sprintf("Error: Failed to write to the file %s. \n%s", path, err.Error()))
		return err
	}
	fmt.Println(fmt.Sprintf("%s ✅", path))
	return nil
}

func filterTable(tables []ddlparse.Table, names []string) []ddlparse.Table {
	if len(names) == 0 {
		return tables
	}
	var ret []ddlparse.Table
	for _, table := range tables {
		for _, name := range names {
			if name == table.Name {
				ret = append(ret, table)
			}
		}
	}
	return ret
}

func getParsedTables(args []string) ([]ddlparse.Table, error) {
	ddl, err := readFile(args[0])
	if err != nil {
		return []ddlparse.Table{}, err
	}

	tables, err := parseDDL(ddl)
	if err != nil {
		return []ddlparse.Table{}, err
	}

	return filterTable(tables, args[1:]), nil

}

func snakeToPascal(snake string) string {
	ls := strings.Split(strings.ToLower(snake), "_")
	for i, s := range ls {
		ls[i] = strings.ToUpper(s[0:1]) + s[1:]
	}
	return strings.Join(ls, "")
}

func snakeToCamel(snake string) string {
	ls := strings.Split(strings.ToLower(snake), "_")
	for i, s := range ls {
		if i != 0 {
			ls[i] = strings.ToUpper(s[0:1]) + s[1:]
		}
	}
	return strings.Join(ls, "")
}

func getSnakeInitial(snake string) string {
	ls := strings.Split(strings.ToLower(snake), "_")
	ret := ""
	for _, s := range ls {
		ret += s[0:1]
	}
	return ret
}

func getFieldName(columnName, tableName string) string {
	cn := strings.ToLower(columnName)
	tn := strings.ToLower(tableName)
	pf := tn + "_"
	if (strings.HasPrefix(cn, pf)) {
		cn = cn[len(pf):]
	}
	return snakeToPascal(cn)
}

func generateModel(args []string) error {
	tables, err := getParsedTables(args)
	if err != nil {
		fmt.Println(fmt.Sprintf("Error: Failed to parse the DDL. \n%s", err.Error()))
		return err
	}

	for _, table := range tables {
		code := generateModelCode(table)
		path := fmt.Sprintf("internal/model/%s.go", strings.ToLower(table.Name))
		err = writeFile(path, code)
		if err != nil {
			return err
		}
	}
	return nil
}

func isNullColumn(column ddlparse.Column, constraints ddlparse.TableConstraint) bool {
	if (column.Constraint.IsNotNull) {
		return false
	}
	if (column.Constraint.IsPrimaryKey) {
		return false
	}
	if (column.Constraint.IsAutoincrement) {
		return false
	}

	for _, pk := range constraints.PrimaryKey {
		for _, name := range pk.ColumnNames {
			if (column.Name == name) {
				return false
			}
		}
	}
	return true
}

func generateModelCode(table ddlparse.Table) string {
	code := "package model\n\n\n"
	tn := strings.ToLower(table.Name)
	code += "type " + snakeToPascal(tn) + " struct {\n"
	for _, column := range table.Columns {
		cn := strings.ToLower(column.Name)
		code += "\t" + getFieldName(cn ,tn) + " ";
		if isNullColumn(column, table.Constraints) {
			code += "*"
		}
		code += dataTypeToGoType(column.DataType.Name) + " "
		code += "`db:\"" + cn + "\" json:\"" + cn + "\"`\n"
	}
	code += "}"
	return code
}

func generateRepository(args []string) error {
	tables, err := getParsedTables(args)
	if err != nil {
		return err
	}

	for _, table := range tables {
		code := generateRepositoryCode(table)
		path := fmt.Sprintf("internal/repository/%s.go", strings.ToLower(table.Name))
		err = writeFile(path, code)
		if err != nil {
			return err
		}
	}
	return nil
}

func dataTypeToGoType(dataType string) string {
	dataType = strings.ToUpper(dataType)

	if (strings.Contains(dataType, "INT") || strings.Contains(dataType, "BIT") || strings.Contains(dataType, "SERIAL")) {
		return "int"
	} else if strings.Contains(dataType, "NUMERIC") ||
		strings.Contains(dataType, "DECIMAL") ||
		strings.Contains(dataType, "FLOAT") ||
		strings.Contains(dataType, "REAL") ||
		strings.Contains(dataType, "DOUBLE") {
		return "float64"
	} else {
		return "string"
	}
}

const TEMPLATE = 
`package repository

import (
	"database/sql"

	"%s/internal/core/db"
	"%s/internal/model"
)


type %sRepository interface {
	Get(%s *model.%s) ([]model.%s, error)
	GetOne(%s *model.%s) (model.%s, error)
	Insert(%s *model.%s, tx *sql.Tx) error
	Update(%s *model.%s, tx *sql.Tx) error
	Delete(%s *model.%s, tx *sql.Tx) error
}


type %sRepository struct {
	db *sql.DB
}

func New%sRepository() %sRepository {
	db := db.GetDB()
	return &%sRepository{db}
}


%s


%s


%s


%s


%s`

const TEMPLATE_GET =
`func (rep *%sRepository) Get(%s *model.%s) ([]model.%s, error) {
	where, binds := db.BuildWhereClause(%s)
	query := %s + where
	rows, err := rep.db.Query(query, binds...)
	defer rows.Close()

	if err != nil {
		return []model.%s{}, err
	}

	ret := []model.%s{}
	for rows.Next() {
		%s := model.%s{}
		err = rows.Scan(%s)
		if err != nil {
			return []model.%s{}, err
		}
		ret = append(ret, %s)
	}

	return ret, nil
}`

const TEMPLATE_GETONE =
`func (rep *%sRepository) GetOne(%s *model.%s) (model.%s, error) {
	var ret model.%s
	where, binds := db.BuildWhereClause(%s)
	query := %s + where

	err := rep.db.QueryRow(query, binds...).Scan(%s)

	return ret, err
}`

const TEMPLATE_INSERT =
`func (rep *%sRepository) Insert(%s *model.%s, tx *sql.Tx) error {
	cmd := %s
	binds := []interface{}{%s}

	var err error
	if tx != nil {
        _, err = tx.Exec(cmd, binds...)
    } else {
        _, err = rep.db.Exec(cmd, binds...)
    }
	
	return err
}`

const TEMPLATE_UPDATE =
`func (rep *%sRepository) Update(%s *model.%s, tx *sql.Tx) error {
	cmd := %s
	binds := []interface{}{%s}
	
	var err error
	if tx != nil {
        _, err = tx.Exec(cmd, binds...)
    } else {
        _, err = rep.db.Exec(cmd, binds...)
    }
	
	return err
}`

const TEMPLATE_DELETE =
`func (rep *%sRepository) Delete(%s *model.%s, tx *sql.Tx) error {
	where, binds := db.BuildWhereClause(%s)
	cmd := "DELETE FROM %s " + where

	var err error
	if tx != nil {
        _, err = tx.Exec(cmd, binds...)
    } else {
        _, err = rep.db.Exec(cmd, binds...)
    }
	
	return err
}`

func generateRepositoryCode(table ddlparse.Table) string {
	tn := strings.ToLower(table.Name)
	tnc := snakeToCamel(tn)
	tnp := snakeToPascal(tn)
	tni := getSnakeInitial(tn)

	return fmt.Sprintf(
		TEMPLATE, cf.AppName, cf.AppName,
		tnp, tni, tnp, tnp, tni, tnp, tnp, tni, tnp, tni, tnp, tni, tnp,
		tnc, tnp, tnp, tnc,
		generateRepositoryGetCode(table),
		generateRepositoryGetOneCode(table),
		generateRepositoryInsertCode(table),
		generateRepositoryUpdateCode(table),
		generateRepositoryDeleteCode(table),
	)
}

func generateRepositoryGetCode(table ddlparse.Table) string {
	tn := strings.ToLower(table.Name)
	tnc := snakeToCamel(tn)
	tnp := snakeToPascal(tn)
	tni := getSnakeInitial(tn)

	query := "\n\t`SELECT\n"
	for i, c := range table.Columns {
		if i == 0 {
			query += fmt.Sprintf("\t\t%s", c.Name)
		} else {
			query += fmt.Sprintf("\n\t\t,%s", c.Name)
		}
	}
	query += fmt.Sprintf("\n\t FROM %s `", tn)

	scan := "\n"
	for _, c := range table.Columns {
		scan += fmt.Sprintf("\t\t\t&%s.%s,\n", tni, getFieldName(c.Name ,tn))
	}
	scan += "\t\t"

	return fmt.Sprintf(
		TEMPLATE_GET,
		tnc, tni, tnp, tnp, tni, 
		query,
		tnp, tnp, tni, tnp,
		scan,
		tnp, tni,
	) 
}


func generateRepositoryGetOneCode(table ddlparse.Table) string {
	tn := strings.ToLower(table.Name)
	tnc := snakeToCamel(tn)
	tnp := snakeToPascal(tn)
	tni := getSnakeInitial(tn)

	query := "\n\t`SELECT\n"
	for i, c := range table.Columns {
		if i == 0 {
			query += fmt.Sprintf("\t\t%s", c.Name)
		} else {
			query += fmt.Sprintf("\n\t\t,%s", c.Name)
		}
	}
	query += fmt.Sprintf("\n\t FROM %s `", tn)

	scan := "\n"
	for _, c := range table.Columns {
		scan += fmt.Sprintf("\t\t&ret.%s,\n", getFieldName(c.Name ,tn))
	}
	scan += "\t"

	return fmt.Sprintf(
		TEMPLATE_GETONE,
		tnc, tni, tnp, tnp, tnp, tni, 
		query,
		scan,
	) 
}


func getBindVar(dbDriver string, n int) string {
	if dbDriver == "postgres" {
		return fmt.Sprintf("$%d", n)
	} else {
		return "?"
	}
}


func concatBindVariableWithCommas(dbDriver string, bindCount int) string {
	var ls []string
	for i := 1; i <= bindCount; i++ {
		ls = append(ls, getBindVar(dbDriver, i))
	}
	return strings.Join(ls, ",")
}


func isInsertColumn(c ddlparse.Column) bool {
	if c.Constraint.IsAutoincrement {
		return false
	}
	if strings.Contains(strings.ToUpper(c.DataType.Name), "SERIAL") {
		return false
	}
	if strings.Contains(c.Name, "_at") || strings.Contains(c.Name, "_AT") {
		return false
	}

	return true
}


func generateRepositoryInsertCode(table ddlparse.Table) string {
	
	tn := strings.ToLower(table.Name)
	tnc := snakeToCamel(tn)
	tnp := snakeToPascal(tn)
	tni := getSnakeInitial(tn)

	query := fmt.Sprintf("\n\t`INSERT INTO %s (\n", tn)
	bindCount := 0
	for _, c := range table.Columns {
		if isInsertColumn(c) {
			bindCount += 1
			if bindCount == 1 {
				query += fmt.Sprintf("\t\t%s", c.Name)
			} else {
				query += fmt.Sprintf("\n\t\t,%s", c.Name)
			}
		}	
	}
	query += fmt.Sprintf("\n\t ) VALUES(%s)`\n", concatBindVariableWithCommas(cf.DBDriver, bindCount))

	binds := "\n"
	for _, c := range table.Columns {
		if isInsertColumn(c) {
			binds += fmt.Sprintf("\t\t%s.%s,\n", tni, getFieldName(c.Name ,tn))
		}
	}
	binds += "\t"

	return fmt.Sprintf(
		TEMPLATE_INSERT,
		tnc, tni, tnp,
		query,
		binds,
	) 
}


func isUpdateColumn(c ddlparse.Column) bool {
	if c.Constraint.IsAutoincrement {
		return false
	}
	if strings.Contains(strings.ToUpper(c.DataType.Name), "SERIAL") {
		return false
	}
	if c.Constraint.IsPrimaryKey {
		return false
	}
	if strings.Contains(c.Name, "_at") || strings.Contains(c.Name, "_AT") {
		return false
	}

	return true
}


func generateRepositoryUpdateCode(table ddlparse.Table) string {
	tn := strings.ToLower(table.Name)
	tnc := snakeToCamel(tn)
	tnp := snakeToPascal(tn)
	tni := getSnakeInitial(tn)

	query := fmt.Sprintf("\n\t`UPDATE %s\n\t SET ", tn) 
	bindCount := 0
	for _, c := range table.Columns {
		if isUpdateColumn(c) {
			bindCount += 1
			if bindCount == 1 {
				query += fmt.Sprintf("%s = %s\n", c.Name, getBindVar(cf.DBDriver, bindCount))
			} else {
				query += fmt.Sprintf("\t\t,%s = %s\n", c.Name, getBindVar(cf.DBDriver, bindCount))
			}
		}	
	}
	query += "\t WHERE "
	isFirst := true
	for _, c := range table.Columns {
		if c.Constraint.IsPrimaryKey {
			bindCount += 1
			if isFirst {
				query += fmt.Sprintf("%s = %s", c.Name, getBindVar(cf.DBDriver, bindCount))
				isFirst = false
			} else {
				query += fmt.Sprintf("\n\t   AND %s = %s", c.Name, getBindVar(cf.DBDriver, bindCount))
			}
		}
	}
	query += "`"

	binds := "\n"
	for _, c := range table.Columns {
		if isUpdateColumn(c) {
			binds += fmt.Sprintf("\t\t%s.%s,\n", tni, getFieldName(c.Name ,tn))
		}
	}
	for _, c := range table.Columns {
		if c.Constraint.IsPrimaryKey {
			binds += fmt.Sprintf("\t\t%s.%s,\n", tni, getFieldName(c.Name ,tn))
		}
	}
	binds += "\t"

	return fmt.Sprintf(
		TEMPLATE_UPDATE,
		tnc, tni, tnp,
		query,
		binds,
	) 
}


func generateRepositoryDeleteCode(table ddlparse.Table) string {
	tn := strings.ToLower(table.Name)
	tnc := snakeToCamel(tn)
	tnp := snakeToPascal(tn)
	tni := getSnakeInitial(tn)

	return fmt.Sprintf(TEMPLATE_DELETE, tnc, tni, tnp, tni, tn) 
}