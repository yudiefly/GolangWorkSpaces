package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"
)

//加载ini文件
func loadIni(fileName string, data interface{}) (err error) {
	//0.参数的校验
	//0.1传进的data参数必须是指针类型（因为需要在函数中给其赋值）
	t := reflect.TypeOf(data)
	fmt.Println(t, t.Kind())
	if t.Kind() != reflect.Ptr {
		//err = fmt.Errorf("data shoud be a pointer...") //格式化输出后返回一个error类型
		err = errors.New("data shoud be a pointer") //创建一个错误
		return
	}
	//0.2传进来的data参数必须是结构体类型指针（因为配置文件中有各种键值对需要赋值给结构体的字段）
	if t.Elem().Kind() != reflect.Struct {
		err = errors.New("data param should be a struct pointer") //创建一个错误
		return
	}
	//1.读取文件到字节类型数据
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return
	}
	// string(b) //将文件内容转换成字符串
	lineSlice := strings.Split(string(b), "\r\n")
	fmt.Printf("%#v\n", lineSlice)
	//2.逐行(一行一行）读取数据
	var structName string
	for idx, line := range lineSlice {
		//去掉字符串首尾的空格
		line = strings.TrimSpace(line)
		//如果空行则跳过
		if len(line) == 0 {
			continue
		}
		//2.1如果是注释就忽略（跳过）
		if strings.HasPrefix(line, ";") || strings.HasPrefix(line, "#") {
			continue
		}
		//2.2如果是[开头的表示是节（section)
		if strings.HasPrefix(line, "[") {
			if line[0] != '[' || line[len(line)-1] != ']' {
				err = fmt.Errorf("line:%d syntax error", idx+1)
				return
			}
			//把这一行首尾的[]去掉，取中间的内容把首尾的空格去掉拿到内容
			sectionName := strings.TrimSpace(line[1 : len(line)-1])
			if len(sectionName) == 0 {
				err = fmt.Errorf("line:%d syntax error", idx+1)
				return
			}
			//根据sectionName去data里面根据反射找到对应的结构体

			for i := 0; i < t.Elem().NumField(); i++ {
				feild := t.Elem().Field(i)
				if sectionName == feild.Tag.Get("ini") {
					//说明找到了对应的嵌套结构体，把字段名记下来
					structName = feild.Name
					fmt.Printf("找到%s对应的嵌套结构体%s", sectionName, structName)
				}

			}

		} else {
			//2.3如果不是[开头就是=分割的键值对
			//1.以等号分割这一行，等号左边是key，等号右边是value
			if strings.Index(line, "=") == -1 || strings.HasPrefix(line, "=") {
				err = fmt.Errorf("line:%d syntax error", idx+1)
				return
			}
			index := strings.Index(line, "=")
			key := strings.TrimSpace(line[:index])
			value := strings.TrimSpace(line[index+1:])
			//2.根据structName去data里面把对应的嵌套结构体取出来
			v := reflect.ValueOf(data)
			sValue := v.Elem().FieldByName(structName) //拿到嵌套结构体的值信息
			sType := sValue.Type()                     //拿到嵌套结构体的类型信息

			if sType.Kind() != reflect.Struct {
				err = fmt.Errorf("data中的%s字段应该是一个结构体", structName)
				return
			}
			var fieldName string
			var fieldType reflect.StructField
			//3.遍历嵌套结构体的每一个字段，判断tag是不是等于key
			for i := 0; i < sValue.NumField(); i++ {
				field := sType.Field(i) //tag信息是存储在类型信息中的
				fieldType = field
				if field.Tag.Get("ini") == key {
					//找到了对应的字段
					fieldName = field.Name
					break
				}
			}
			//4.如果key=tag，给这个字段赋值
			//4.1根据fieldName去取出这个字段
			if len(fieldName) == 0 {
				//在结构体中找不到对应的字段的字符
				continue
			}
			fieldObj := sValue.FieldByName(fieldName)
			//4.2对其进行赋值
			fmt.Println(fieldName, fieldType.Type.Kind()) //fieldObj.Type().Kind()
			switch fieldType.Type.Kind() {
			case reflect.String:
				fieldObj.SetString(value)
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				var valueInt int64
				valueInt, err = strconv.ParseInt(value, 10, 64)
				if err != nil {
					err = fmt.Errorf("line:%d value type error", idx+1)
					return
				}
				fieldObj.SetInt(valueInt)
			case reflect.Bool:
				var valueBool bool
				valueBool, err = strconv.ParseBool(value)
				if err != nil {
					err = fmt.Errorf("line:%d value type error", idx+1)
					return
				}
				fieldObj.SetBool(valueBool)
			case reflect.Float32, reflect.Float64:
				var valueFloat float64
				valueFloat, err = strconv.ParseFloat(value, 64)
				if err != nil {
					err = fmt.Errorf("line:%d value type error", idx+1)
					return
				}
				fieldObj.SetFloat(valueFloat)
			}
		}
	}
	return
}
