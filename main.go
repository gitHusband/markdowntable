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
const jsonFile = "./testeasy.json"

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

type tdSlice []string
type trSlice []tdSlice

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
func (e *elementDetails) generateDetailsHtml() string {
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

// 最多有多少列
var maxColspan int
var trs trSlice

func main() {
	var jsonData map[string]interface{}

	// 从文件中读取JSON
	jsonData = readJsonFile(jsonFile)
	fmt.Printf("%#v\n", jsonData)

	maxColspan = getMaxColspan(jsonData) + 1
	fmt.Printf("最大 maxColspan: %v\n", maxColspan)

	// 遍历生成markdown 表格 字符串
	createMarkdownTable(jsonData)
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

// 判断某个JSON 元素是否是 最终元素
func isEndElement(x interface{}) bool {
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

// 获取map 所有的key
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

func getSortMapKeys(m map[string]interface{}) stringSlice {
	keys := getMapKeys(m)
	// 对切片进行排序
	sort.Strings(keys)

	return keys
}

func createMarkdownTable(jsonData map[string]interface{}) string {
	createTableBody(jsonData)

	return ""
}

func createTableBody(jsonData map[string]interface{}) {
	createTableTr(jsonData, maxColspan, "")
}

func createTableTr(jsonData map[string]interface{}, colspan int, keyPath string) {
	keys := getSortMapKeys(jsonData)

	keysLen := len(keys)
	keyOffset := maxColspan - colspan
	fmt.Printf("@@@ 最大列数是：%v， 当前列数是: %v, keyOffset 是：%v, keysLen 是：%v\n", maxColspan, colspan, keyOffset, keysLen)

	// for _, key := range keys {
	for i := keyOffset; i < keysLen; i++ {
		var tds tdSlice
		var td string

		currentKey := keys[i]
		currentColspan := colspan - 1

		var currentKeyPath string
		if keyPath == "" {
			currentKeyPath = currentKey
		} else {
			currentKeyPath = keyPath + "." + currentKey
		}

		jsonDataValue := jsonData[currentKey]

		if isEndElement(jsonDataValue) {
			// fmt.Printf("=== 这是最终元素 - %v：%#v\n", currentKeyPath, jsonDataValue)
			fmt.Printf("=== 这是最终元素 - %v\n", currentKeyPath)

			td = createTableTd(currentKey, 1, 1)
			tds = append(tds, td)

			elmDetails := new(elementDetails)
			elmDetails.setup(jsonDataValue)
			// fmt.Printf("### elmDetails: %#v\n", elmDetails)
			details := elmDetails.generateDetailsHtml()
			// fmt.Printf("*** elmDetails HTML: %#v\n\n", details)
			td = createTableTd(details, currentColspan, 1)
			// fmt.Printf("*** %v - td HTML: %#v\n\n", currentKeyPath, td)
			tds = append(tds, td)
			fmt.Printf("$$$ tds - %v：%#v\n", currentKeyPath, tds)

		} else {
			// fmt.Printf("### 这不是最终元素 - %v：%#v\n", currentKeyPath, jsonDataValue)
			fmt.Printf("### 这不是最终元素 - %v\n", currentKeyPath)

			td = createTableTd(currentKey, 1, getRowspan(jsonDataValue.(map[string]interface{})))
			tds = append(tds, td)

			fmt.Printf("### 所以需要先插入当前字段 %v，再一次性遍历获取所有子字段的第一个字段\n", currentKey)
			firstColumns := getFirstColumns(jsonDataValue.(map[string]interface{}), currentColspan)
			tds = append(tds, firstColumns...)

			fmt.Printf("$$$ tds - %v：%#v\n", currentKeyPath, tds)

			createTableTr(jsonDataValue.(map[string]interface{}), currentColspan, currentKeyPath)
		}

	}
}

func createTableTd(content string, colspan int, rowspan int) string {
	var td string
	var colspanDetail string
	var rowspanDetail string

	if colspan > 1 {
		colspanDetail = fmt.Sprintf(" colspan=\"%d\"", colspan)
	}

	if rowspan > 1 {
		rowspanDetail = fmt.Sprintf(" rowspan=\"%d\"", rowspan)
	}

	td = fmt.Sprintf("<td%v%v>%s</td>", colspanDetail, rowspanDetail, content)

	return td
}

func createTableHeader() {

}

// 获取整个表格最大有多少列
func getMaxColspan(jsonData map[string]interface{}) int {
	maxColspan := 1
	for _, v := range jsonData {
		curColspan := 1
		if isEndElement(v) {
			curColspan = 1
		} else {
			curColspan += getMaxColspan(v.(map[string]interface{}))
		}

		if curColspan > maxColspan {
			maxColspan = curColspan
		}
	}

	return maxColspan
}

func getColspan() {

}

// 获取当前字段有多少个子字段，即当前字段占几行
func getRowspan(jsonData map[string]interface{}) int {
	rowspan := 0
	for _, v := range jsonData {
		curRowspan := 1
		if isEndElement(v) {
			curRowspan = 1
		} else {
			curRowspan = getRowspan(v.(map[string]interface{}))
		}

		rowspan += curRowspan
	}

	return rowspan
}

func getFirstColumns(jsonData map[string]interface{}, colspan int) tdSlice {
	var tds tdSlice
	var td string

	keys := getSortMapKeys(jsonData)
	firstKey := keys[0]

	jsonDataValue := jsonData[firstKey]

	if isEndElement(jsonDataValue) {
		// fmt.Printf("====== 这是最终元素 - %v：%#v\n", firstKey, jsonDataValue)
		fmt.Printf("====== 这是最终元素 - %v\n", firstKey)

		td = createTableTd(firstKey, 1, 1)
		tds = append(tds, td)

		elmDetails := new(elementDetails)
		elmDetails.setup(jsonDataValue)
		// fmt.Printf("****** elmDetails: %#v\n", elmDetails)
		details := elmDetails.generateDetailsHtml()
		// fmt.Printf("****** elmDetails HTML: %#v\n\n", details)
		td := createTableTd(details, colspan-1, 1)
		// fmt.Printf("****** %v - td HTML: %#v\n\n", firstKey, td)
		tds = append(tds, td)
	} else {
		// fmt.Printf("###### 这不是最终元素 - %v：%#v\n", firstKey, jsonDataValue)
		fmt.Printf("###### 这不是最终元素 - %v\n", firstKey)
		fmt.Printf("###### 所以需要先插入当前字段 %v，再一次性遍历获取所有子字段的第一个字段\n", firstKey)

		td = createTableTd(firstKey, 1, getRowspan(jsonDataValue.(map[string]interface{})))
		tds = append(tds, td)

		nextTd := getFirstColumns(jsonDataValue.(map[string]interface{}), colspan-1)
		tds = append(tds, nextTd...)
	}

	return tds
}
