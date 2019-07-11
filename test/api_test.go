package main

import (
	"fmt"
	"gcs/module/system"
	"github.com/gogf/gf/g/util/gconv"
	"reflect"
	"testing"
)

func TestRun(t *testing.T) {
	foreach("1", "2")
	var paramsStr []string
	paramsStr = append(paramsStr, "3")
	paramsStr = append(paramsStr, "4")
	foreach(paramsStr)

	var params []interface{}
	params = append(params, "5")
	params = append(params, "6")
	foreach(params)
	foreach(params...)

}

func TestReflect(t *testing.T) {
	model := system.SysConfig{}
	model.UpdateId = 1
	re := reflect.ValueOf(model).FieldByName("BaseModel")
	updateId := gconv.Int(re.FieldByName("UpdateId").Interface())
	fmt.Println(updateId)
}

func foreach(args ...interface{}) {
	for _, a := range args {
		fmt.Println("#" + gconv.String(a))
	}
}
