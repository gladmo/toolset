package makelevel

import (
	"encoding/json"
	"fmt"
	"go/format"
	"reflect"
	"strings"
	"unsafe"
)

// FlatStruct flat struct
type FlatStruct struct {
	isSlice  bool
	sliceLen int
	tagName  string
	
	structRelation map[string][]string
	StructField    []reflect.StructField
}

// BackFill create new reflect.Value BackFill origin data
func (th FlatStruct) BackFill(origin interface{}) reflect.Value {
	val := reflect.Indirect(reflect.ValueOf(origin))
	
	// 数组或切片的单独处理
	if th.isSlice {
		data := reflect.MakeSlice(reflect.SliceOf(reflect.StructOf(th.StructField)), 0, 0)
		
		for i := 0; i < th.sliceLen; i++ {
			item := reflect.Indirect(reflect.New(reflect.StructOf(th.StructField)))
			for field, loc := range th.structRelation {
				fvi := fieldByNames(val.Index(i), loc...)
				
				item.FieldByName(field).Set(GetUnexportedField(fvi))
			}
			
			data = reflect.Append(data, item)
		}
		return data
	}
	
	data := reflect.New(reflect.StructOf(th.StructField))
	data = reflect.Indirect(data)
	for field, loc := range th.structRelation {
		fvi := fieldByNames(val, loc...)
		
		data.FieldByName(field).Set(GetUnexportedField(fvi))
	}
	
	return data
}

// fieldByNames 根据原结构体层次获取原结构中的数值
func fieldByNames(val reflect.Value, fields ...string) reflect.Value {
	val = reflect.Indirect(val)
	for _, field := range fields {
		val = reflect.Indirect(val.FieldByName(field))
	}
	return val
}

// GetUnexportedField 获取结构中未导出的字段
func GetUnexportedField(field reflect.Value) reflect.Value {
	return reflect.Indirect(reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())))
}

// SetUnexportedField 设置结构中未导出的字段
func SetUnexportedField(field reflect.Value, value reflect.Value) {
	reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).
		Elem().Set(value)
}

// BackFillInterface BackFill interface
func (th FlatStruct) BackFillInterface(origin interface{}) interface{} {
	return th.BackFill(origin).Interface()
}

// BackFillJson BackFill json
func (th FlatStruct) BackFillJson(origin interface{}) ([]byte, error) {
	return json.Marshal(th.BackFillInterface(origin))
}

// BackFillTable BackFill table
func (th FlatStruct) BackFillTable(origin interface{}) (data [][]string) {
	var fields []string
	var headers []string
	
	val := reflect.Indirect(th.BackFill(origin))
	
	// 数组或切片的单独处理
	if th.isSlice {
		elem := val.Type().Elem()
		for i := 0; i < elem.NumField(); i++ {
			tag := elem.Field(i).Tag.Get(th.tagName)
			if tag != "" {
				fields = append(fields, elem.Field(i).Name)
				headers = append(headers, tag)
			}
		}
		data = append(data, headers)
		
		for i := 0; i < val.Len(); i++ {
			var row []string
			for _, item := range fields {
				field := val.Index(i).FieldByName(item)
				
				b, _ := json.Marshal(field.Interface())
				
				row = append(row, strings.Trim(string(b), `"`))
			}
			
			data = append(data, row)
		}
		return data
	}
	
	for i := 0; i < val.NumField(); i++ {
		fields = append(fields, val.Type().Field(i).Name)
		headers = append(headers, val.Type().Field(i).Tag.Get(th.tagName))
	}
	
	data = append(data, headers)
	
	var row []string
	for _, item := range fields {
		field := val.FieldByName(item)
		
		b, _ := json.Marshal(field.Interface())
		
		row = append(row, strings.Trim(string(b), `"`))
	}
	
	data = append(data, row)
	return data
}

// BackFillTableInterface BackFill interface table
func (th FlatStruct) BackFillTableInterface(origin interface{}) (data [][]interface{}) {
	var fields []string
	var headers []interface{}
	
	val := reflect.Indirect(th.BackFill(origin))
	
	// 数组或切片的单独处理
	if th.isSlice {
		elem := val.Type().Elem()
		for i := 0; i < elem.NumField(); i++ {
			tag := elem.Field(i).Tag.Get(th.tagName)
			if tag != "" {
				fields = append(fields, elem.Field(i).Name)
				headers = append(headers, tag)
			}
		}
		data = append(data, headers)
		
		for i := 0; i < val.Len(); i++ {
			var row []interface{}
			for _, item := range fields {
				field := val.Index(i).FieldByName(item)
				
				row = append(row, field.Interface())
			}
			
			data = append(data, row)
		}
		return data
	}
	
	for i := 0; i < val.NumField(); i++ {
		fields = append(fields, val.Type().Field(i).Name)
		headers = append(headers, val.Type().Field(i).Tag.Get(th.tagName))
	}
	
	data = append(data, headers)
	
	var row []interface{}
	for _, item := range fields {
		field := val.FieldByName(item)
		
		row = append(row, field.Interface())
	}
	
	data = append(data, row)
	return data
}

// PrintStruct print struct desc
func (th FlatStruct) PrintStruct() {
	b, _ := format.Source([]byte(reflect.New(reflect.StructOf(th.StructField)).Type().String()))
	fmt.Println(string(b))
	
	b, _ = json.Marshal(th.structRelation)
	fmt.Println(string(b))
	fmt.Println(len(th.structRelation))
}
