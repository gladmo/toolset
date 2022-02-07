package makelevel

import (
	"fmt"
	"reflect"
	"strings"
)

// MakeLevel 结构体找平参数
type MakeLevel struct {
	tagName      string
	tagSeparator string
}

// NewMakeLevel new make level tool
//
// 	ml := NewMakeLevel("header")
//	lf := ml.MakeLevelStruct(slice)
//  fmt.Println(lf.BackFillJson(slice))
func NewMakeLevel(tagName string) *MakeLevel {
	return &MakeLevel{
		tagName:      tagName,
		tagSeparator: ".",
	}
}

// MakeLevelStruct 找平结构体，只支持[]struct{} & struct{}
func (th *MakeLevel) MakeLevelStruct(s interface{}) (lf FlatStruct) {
	var isSlice bool
	var sliceLen int
	
	v := reflect.Indirect(reflect.ValueOf(s))
	if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
		isSlice = true
		sliceLen = v.Len()
		s = v.Index(0).Interface()
	} else if v.Kind() != reflect.Struct {
		panic("MakeLevelStruct must []struct or struct")
	}
	
	sf, structRelation := th.ReadStruct(s)
	return FlatStruct{
		isSlice:        isSlice,
		sliceLen:       sliceLen,
		tagName:        th.tagName,
		structRelation: structRelation,
		StructField:    sf,
	}
}

func (th *MakeLevel) ReadStruct(st interface{}) (sf []reflect.StructField, structRelation map[string][]string) {
	sf, structRelation = th.parseStructField(reflect.ValueOf(st))
	
	for i, field := range sf {
		lowercase := field.Name
		sf[i].Name = strings.Title(field.Name)
		
		structRelation[sf[i].Name] = structRelation[lowercase]
		delete(structRelation, lowercase)
	}
	return
}

// parseStructField 递归解析原结构体
func (th *MakeLevel) parseStructField(val reflect.Value) (sf []reflect.StructField, structRelation map[string][]string) {
	structRelation = make(map[string][]string)
	
	val = reflect.Indirect(val)
	
	for i := 0; i < val.NumField(); i++ {
		f := val.Field(i)
		
		switch f.Kind() {
		case reflect.Struct:
			filedName := val.Type().Field(i).Name
			
			sfs, sr := th.parseStructField(f)
			for idx, structField := range sfs {
				sfs[idx].Name = fmt.Sprintf(`%s_%s`, filedName, structField.Name)
				sfs[idx].Tag = reflect.StructTag(
					fmt.Sprintf(`%s:"%s%s%s"`,
						th.tagName, filedName, th.tagSeparator, structField.Tag.Get(th.tagName),
					))
				structRelation[sfs[idx].Name] = append([]string{filedName}, sr[structField.Name]...)
			}
			sf = append(sf, sfs...)
		default:
			filed := val.Type().Field(i)
			filedName := filed.Name
			
			sf = append(sf, reflect.StructField{
				Name:      filedName,
				Type:      filed.Type,
				Tag:       filed.Tag,
				Anonymous: false,
			})
			
			structRelation[filedName] = []string{filedName}
		}
	}
	
	return
}
