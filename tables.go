/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package cli

import (
	"reflect"

	"github.com/jedib0t/go-pretty/v6/table"
)

func CreateTableHeader(s interface{}) table.Row {

	gtype := reflect.TypeOf(s)
	var header table.Row

	for i := 0; i < gtype.NumField(); i++ {
		header = append(header, gtype.Field(i).Name)
	}

	return header
}

func CreateTableRows(s interface{}) table.Row {

	gtype := reflect.TypeOf(s)
	numFields := gtype.NumField()
	rg := GetUnderlyingAsValue(s)
	var row table.Row
	for i := 0; i < numFields; i++ {

		var value string
		switch reflect.TypeOf(rg.Field(i).Interface()).String() {
		case "string":
			value = rg.Field(i).Interface().(string)
			break
		case "int":
			value = strconv.Itoa(rg.Field(i).Interface().(int))
			break

		}
		row = append(row, value)

	}
	return row
}
