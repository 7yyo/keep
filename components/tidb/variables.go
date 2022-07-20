package tidb

import (
	"fmt"
	"github.com/pingcap/tidb/sessionctx/variable"
	"reflect"
	"strings"
)

func vars() {
	reflectConfig(newVar())
}

func newVar() reflect.Value {
	return reflect.ValueOf(variable.NewSessionVars())
}

func reflectConfig(v reflect.Value) {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	n := v.NumField()
	for i := 0; i < n; i++ {
		f := v.Field(i)
		switch f.Kind() {
		case reflect.String:
			fmt.Printf("%s: %s\n", v.Type().Field(i).Name, f.String())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			fmt.Printf("%s: %d\n", v.Type().Field(i).Name, f.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			fmt.Printf("%s: %d\n", v.Type().Field(i).Name, f.Uint())
		case reflect.Bool:
			fmt.Printf("%s: %t\n", v.Type().Field(i).Name, f.Bool())
		case reflect.Struct:
			tl := strings.Split(f.String(), " ")[0]
			tl = strings.ReplaceAll(tl, "<", "")
			tl = strings.ReplaceAll(tl, "variables.", "")
			fmt.Println(fmt.Sprintf("[%s]", tl))
			reflectConfig(f)
		case reflect.Map:
			break
		}
	}
}
