package main

import (
	"fmt"
	"github.com/rohithputha/jab/memtbl"
	"time"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	memtbl.InitTables("latest.log")
	memtbl.CreateTable("testTable", "id", "string", []string{"name"}, []string{"string"})
	table, err := memtbl.TablesObj.NameTblMap["testTable"]
	if err {
		fmt.Println(err)
	}
	e := table.CreateAndAddKV("hello", "world")
	if e != nil {
		return
	}
	fmt.Println(string(table.CreateAndGetKV("hello").GetData()))
	time.Sleep(1 * time.Second)
}
