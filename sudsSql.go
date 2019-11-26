package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
)

// InitDb Initializes the DB interface
func InitDb(driverName string, dbPath string) {
	db, err = sql.Open("sqlite3", "./foo.db")
	if err != nil {
		log.Fatal(err)
	}
}

// AddColumnToTable will add a column to the given table.
func AddColumnToTable(tableName string, columnName string) error {

	query := `ALTER TABLE ` + tableName + ` ADD COLUMN ` + columnName + ` TEXT NULL DEFAULT ""; `

	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	return nil

}

// GetTableColumns will get an array of all the columns in the table.
func GetTableColumns(tableName string) ([]string, error) {

	// we only need one record for this.
	query := "SELECT * FROM " + tableName + " LIMIT 1"

	rows, err := db.Query(query)

	if err != nil {
		return nil, err
	}

	columns, err := rows.Columns()

	if err != nil {
		return nil, err
	}

	return columns, nil
}

// GetTables will return all the tables in the db
func GetTables() ([]string, error) {
	query := "SELECT name FROM sqlite_master WHERE type ='table' AND name NOT LIKE 'sqlite_%';"

	rows, err := db.Query(query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	names := make([]string, 0)
	for rows.Next() {
		var name string
		rows.Scan(&name)
		names = append(names, name)
	}
	return names, nil
}

// InsertValues inserts the values in the given JSON values.
func InsertValues(insertValues string) {
	if !json.Valid([]byte(insertValues)) {
		fmt.Println("Error: Invalid JSON data")
		fmt.Println("Data: ", insertValues)
		return
	}

	var inVal interface{}
	// Decode bytes b into interface i
	json.Unmarshal([]byte(insertValues), &inVal)

	// convert the top level to an array of interfaces
	tlVal := inVal.([]interface{})

	// there should be only one entry in the top level, so grab that and convert it grabbing the table key
	tableNameTmp := tlVal[0].(map[string]interface{})["table"]
	tableName := tableNameTmp.(string)

	// now let's grab the values key, which is a map, but this converts it to an interface{}
	colsIntf := tlVal[0].(map[string]interface{})["values"]

	// convert the interface to a map so we can iterate over it.
	values := colsIntf.(map[string]interface{})

	err = CreateTable(tableName, values)
	if err != nil {
		fmt.Println("Error: unable to verify table.")
		fmt.Println(err.Error())
		return
	}
	// should be good here, insert the values to the table.
	var columnsSlice, valuesSlice []string

	for column, value := range values {
		columnsSlice = append(columnsSlice, column)
		valueStr := value.(string)
		valuesSlice = append(valuesSlice, valueStr)
	}

	query := "INSERT INTO " + tableName + " ("

	for i := 0; i < len(columnsSlice); i++ {
		query += columnsSlice[i] + ", "
	}

	query = query[0:len(query)-2] + ") VALUES ("

	for i := 0; i < len(valuesSlice); i++ {
		query += "'" + valuesSlice[i] + "', "
	}

	query = query[0:len(query)-2] + ");"

	_, err = db.Exec(query)
	if err != nil {
		fmt.Println("Error: Unable to insert data.")
		fmt.Println(err.Error())
	}

}

// CreateTable creates the table in the db, if it does not exist,
// using the given string as a structure
func CreateTable(tableName string, tableVals map[string]interface{}) error {

	// see if this table name already exists
	tables, err := GetTables()

	if err != nil {
		fmt.Println("Error: unable to get table names")
		fmt.Println(err.Error())
		return err
	}

	if ContainsString(tables, tableName) {

		// get the column names from the table.
		columns, err := GetTableColumns(tableName)
		if err != nil {
			fmt.Println("Error: unable to get column names from table: " + tableName)
			fmt.Println(err.Error())
			return err
		}

		// this table already exists so let's check the columns
		for column := range tableVals {
			if !ContainsString(columns, column) {
				AddColumnToTable(tableName, column)
			}
		}
	} else {
		// this table doesn't exist.  Create it.
		query := "CREATE TABLE " + tableName + "( "
		query += `id INTEGER PRIMARY KEY AUTOINCREMENT,
		t TIMESTAMP	DEFAULT CURRENT_TIMESTAMP, `

		for column := range tableVals {
			query += column + " TEXT NULL DEFAULT \"\", "
		}

		query = query[0 : len(query)-2]
		query += `);`

		_, err := db.Exec(query)
		if err != nil {
			fmt.Println("Error: Unable to create table: " + tableName)
			fmt.Println(err.Error())
			return nil
		}

	}

	return nil
}

// DumpTable dumps the contents of a table to a sting.
func DumpTable(tableName string) (string, error) {

	rows, err := db.Query("SELECT * FROM " + tableName + ";")
	if err != nil {
		return "", err
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return "", err
	}
	count := len(columns)
	tableData := make([]map[string]interface{}, 0)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}
		tableData = append(tableData, entry)
	}
	jsonData, err := json.Marshal(tableData)
	if err != nil {
		return "", err
	}

	return string(jsonData), nil

}
