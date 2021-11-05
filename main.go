package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
)

// 目标：将 JSON 数据转换成 markdown 的表格

// 用户自定义JSON 文件
const jsonFile = "./test.json"

// 声明这个类型主要是想给 []string 添加方法，方便调用而已
type stringSlice []string

// 如果JSON中的某个元素仅仅包含这些字段，这个元素将不再遍历，称之为 最终元素
// 并且将这些字段的值组合成元素的详情，也就是 html 表格的最后一列 - 详情
// 1. required - 必须包含这些字段的元素才能认为是最终元素
// 2. optional - 这些字段是可选的
// 2. 如果包含额外的字段，则不能将其认作为 最终元素
type endElementKeysStruct struct {
	required stringSlice
	optional stringSlice
}

type elementDetails struct {
	header       string
	desc         string
	defaultValue string
	options      []string
}

// 判断一个字符串切片中是否包含某个字符串
func (ss stringSlice) include(s string) bool {
	for _, val := range ss {
		if val == s {
			return true
		}
	}

	return false
}

// 判断一个字符串切片中是否包含另一个字符串切片
func (ss stringSlice) includes(sa []string) bool {
	for _, val := range sa {
		if !ss.include(val) {
			return false
		}
	}

	return true
}

// 根据最终元素 设置对应详情字段的内容
func (e *elementDetails) setup(x interface{}) {
	// 如果不是 x.(map[string]interface{} 类型，比如是字符串，则认为是最终元素
	if m, ok := x.(map[string]interface{}); ok {
		if v, ok := m["header"]; ok {
			if s, ok := v.(string); ok {
				e.header = s
			}
		}

		if v, ok := m["desc"]; ok {
			if s, ok := v.(string); ok {
				e.desc = s
			}
		}

		if v, ok := m["defaultValue"]; ok {
			if s, ok := v.(string); ok {
				e.defaultValue = s
			}
		}

		if v, ok := m["options"]; ok {
			if s, ok := v.([]string); ok {
				e.options = s
			}
		}
	} else if v, ok := x.(string); ok {
		e.header = v
	} else if v, ok := x.(int); ok {
		e.header = strconv.Itoa(v)
	} else if v, ok := x.(bool); ok {
		e.header = strconv.FormatBool(v)
	} else if v, ok := x.([]string); ok {
		e.options = v
	}
}

// 根据详情内容生成 html 最后一列的内容
func (e *elementDetails) generateHtml() string {
	s := ""
	if e.header != "" {
		s += fmt.Sprintf("<strong>%v</strong> <br/> ", e.header)
	}

	if e.desc != "" {
		s += fmt.Sprintf("<p>%v</p> <br/>", e.desc)
	}

	if e.defaultValue != "" {
		s += fmt.Sprintf("<p>%v</p> <br/>", e.defaultValue)
	}

	for _, v := range e.options {
		s += fmt.Sprintf("<p>%v</p> <br/>", v)
	}

	return s
}

// 以此为例，判断一个字段是否是最终元素，那么他必须有且仅有以下4个字段：
// header, desc, defaultValue, options
var endElementKeys = endElementKeysStruct{
	required: stringSlice{"header", "desc"},
	optional: stringSlice{"defaultValue", "options"},
}

func main() {
	var jsonMap map[string]interface{}

	// 从文件中读取JSON
	jsonMap = readJsonFile(jsonFile)
	fmt.Printf("%#v\n", jsonMap)

	// 遍历生成markdown 表格 字符串
	rangeJson(jsonMap, "")
}

// 从文件中读取JSON
func readJsonFile(file string) map[string]interface{} {
	filePtr, err := os.Open(jsonFile)
	if err != nil {
		fmt.Printf("文件打开失败 [Err: %s]\n", err.Error())
		return nil
	}

	defer filePtr.Close()

	var anyJson map[string]interface{}
	decoder := json.NewDecoder(filePtr)

	err = decoder.Decode(&anyJson)

	if err != nil {
		fmt.Printf("解码失败 [Err: %s]\n", err.Error())
		return nil
	} else {
		fmt.Println("解码成功")
	}

	return anyJson
}

// 遍历 JSON
func rangeJson(anyJson map[string]interface{}, keyPath string) {
	var keyList []string
	// 将map数据遍历复制到切片中
	for k := range anyJson {
		keyList = append(keyList, k)
	}
	// 对切片进行排序
	sort.Strings(keyList)

	for _, key := range keyList {
		var currentKeyPath string
		if keyPath == "" {
			currentKeyPath = key
		} else {
			currentKeyPath = keyPath + "." + key
		}

		anyJsonValue := anyJson[key]

		if isEndElement(anyJsonValue, key) {
			fmt.Printf("=== %v 是最终元素，%v\n", key, anyJsonValue)
		} else {
			fmt.Printf("+++ %v: %v\n", currentKeyPath, anyJsonValue)
			rangeJson(anyJsonValue.(map[string]interface{}), currentKeyPath)
		}

		elmDetails := new(elementDetails)
		elmDetails.setup(anyJsonValue)
		fmt.Printf("### elmDetails: %#v\n", elmDetails)
		s := elmDetails.generateHtml()
		fmt.Printf("*** elmDetails HTML: %#v\n\n", s)

		// switch anyJsonValue.(type) {
		// case map[string]interface{}:
		// 	elmDetails.setup(anyJsonValue)
		// 	fmt.Printf("+++ %v: %v\n", currentKeyPath, anyJsonValue)
		// case string, int, bool:

		// default:
		// 	fmt.Printf("%v: %v\n", currentKeyPath, anyJson[key])
		// }
	}
}

// 判断某个JSON 元素是否是 最终元素
func isEndElement(x interface{}, key string) bool {
	m, ok := x.(map[string]interface{})

	// 如果不是 x.(map[string]interface{} 类型，比如是字符串，则认为是最终元素
	if !ok {
		return true
	}

	mKeys := getMapKeys(m)

	// 如果不包含必要的字段，则不能认为是最终字段
	if !mKeys.includes(endElementKeys.required) {
		return false
	}

	// 必要字段和可选字段是最终元素支持的所有字段
	// 如果存在其他字段，则不能认为是最终字段
	allAllowedEndKeys := append(endElementKeys.required, endElementKeys.optional...)

	for _, v := range mKeys {
		if !allAllowedEndKeys.include(v) {
			return false
		}
	}

	return true
}

func getMapKeys(m map[string]interface{}) stringSlice {
	keys := make(stringSlice, 0)
	for k, _ := range m {
		keys = append(keys, k)
	}
	return keys
}

func getEmptyInterfaceKeys(x interface{}) stringSlice {
	m, ok := x.(map[string]interface{})
	if !ok {
		panic("不是 map[string]interface{}，无法获取key")
	}

	return getMapKeys(m)
}
