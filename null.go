package gostuff

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
)

// check if there are any null fields in a table
func CheckNullInTable(table string) {
	problems, _ := os.OpenFile("logs/errors.txt", os.O_APPEND|os.O_WRONLY, 0666)
	defer problems.Close()
	log := log.New(problems, "", log.LstdFlags|log.Lshortfile)

	/*
	   #SELECT * FROM ratinghistory WHERE bullet IS NULL OR blitz IS NULL OR standard IS NULL;
	*/
	var colNames []string
	colNames = getColumnNamesInTable(table)
	if colNames == nil {
		return
	}

	concatColNames := ""

	for _, value := range colNames {
		concatColNames += (value + " IS NULL OR ")
	}
	concatColNames = strings.TrimSuffix(concatColNames, "OR ")

	//fmt.Println(concatColNames)

	rows, err := db.Query("SELECT * FROM " + table + " WHERE " + concatColNames)
	if err != nil {
		log.Println(err)
		return
	}
	defer rows.Close()

	var columns []string
	columns, err = rows.Columns()

	if err != nil {
		fmt.Println(err)
		return
	}

	colNum := len(columns)
	var values = make([]interface{}, colNum)

	for i, _ := range values {
		var ii interface{}
		values[i] = &ii
	}

	i := 0

	for rows.Next() {
		i++

		if err := rows.Scan(values...); err != nil {
			log.Println(err)
		}

		if err != nil {
			log.Println(err)
			return
		}
		for i, colName := range columns {
			var rawValue = *(values[i].(*interface{}))
			var rawType = reflect.TypeOf(rawValue)

			if rawValue != nil {
				//value := fmt.Sprintf("%s", rawValue)
				//fmt.Println(colName, rawType, value)
			} else {
				fmt.Println(colName, rawType, rawValue)
				fmt.Println(colName, "is nil!")
			}
		}
	}
	if i == 0 {
		fmt.Println("No null values in", table, "table")
	}
}

// gets all the column names in table
// returns nil if there was an error
func getColumnNamesInTable(table string) []string {

	problems, _ := os.OpenFile("logs/errors.txt", os.O_APPEND|os.O_WRONLY, 0666)
	defer problems.Close()
	log := log.New(problems, "", log.LstdFlags|log.Lshortfile)

	var results []string
	query := "SELECT column_name FROM INFORMATION_SCHEMA.COLUMNS " +
		"WHERE table_name = '" + table + "' ORDER BY ordinal_position"
	rows, err := db.Query(query)

	if err != nil {
		log.Println(err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var field string
		if err := rows.Scan(&field); err != nil {
			log.Println(err)
		}
		results = append(results, field)
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil
	}
	if err != nil {
		log.Println(err)
		return nil
	}

	return results
}
