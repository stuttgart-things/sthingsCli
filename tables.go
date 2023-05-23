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
