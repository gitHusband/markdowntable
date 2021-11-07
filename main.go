package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
)

// 目标：将 JSON 数据转换成 markdown 的表格
// 使用方法：markdowntable [-in jsonFilePath] [-out markdownFilePath]

// 声明这个类型主要是想给 []string 添加方法，方便调用而已
type stringSlice []string

// 如果JSON中的某个元素仅仅包含这些字段，这个元素将不再遍历，称之为 最终元素
// 并且将这些字段的值组合成元素的详情，也就是 html 表格的最后一列 - 详情
// 1. required - 必须包含这些字段的元素才能认为是最终元素
// 2. optional - 这些字段是可选的
// 3. 如果包含额外的字段，则不能将其认作为 最终元素
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

		if _, ok := m["options"]; ok {
			e.options = getJsonSlice(m["options"])
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
		s += fmt.Sprintf("<strong style=\"font-size: 15px\">%v</strong>", e.header)
	}

	if e.desc != "" {
		s += fmt.Sprintf("<br/> <em style=\"color: #888888\">%v</em>", e.desc)
	}

	if e.defaultValue != "" {
		s += fmt.Sprintf("<br/> <b>默认：<ins>%v</ins></b>", e.defaultValue)
	}

	if e.options != nil {
		s += "<br/> <b>可能的值：</b> <ul>"
		for _, v := range e.options {
			s += fmt.Sprintf("<li>%v</li>", v)
		}
		s += "</ul>"
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

// 保存每一行的所有列
var trs trSlice

// 默认输入的用户JSON 文件
const defaultJsonFile = "./testeasy.json"

// 默认输出的 Markdown 文件
const defaultMarkdownFile = ""

func main() {
	var jsonData map[string]interface{}

	var inputPath = flag.String("in", defaultJsonFile, "输入JSON文件的路径")
	var outputPath = flag.String("out", defaultMarkdownFile, "输出Markdown文件的路径")
	flag.Parse()

	jsonFile := *inputPath
	markdownFile := *outputPath
	if markdownFile == "" {
		markdownFile = getFullFilePathWithoutSuffix(jsonFile) + ".md"
	}

	// 从文件中读取JSON
	jsonData = readJsonFile(jsonFile)
	// fmt.Printf("Json 数据：%#v\n", jsonData)

	maxColspan = getMaxColspan(jsonData)
	// fmt.Printf("最大 maxColspan: %v\n", maxColspan)

	// 遍历生成markdown 表格 字符串
	markdownTable := createMarkdownTable(jsonData)

	// fmt.Printf("\n%v\n", markdownTable)

	// 输出 markdown 文件
	writeMarkdownFile(markdownTable, markdownFile)
}

// 从文件中读取JSON
func readJsonFile(jsonFilePath string) map[string]interface{} {
	filePtr, err := os.Open(jsonFilePath)
	if err != nil {
		fmt.Printf("文件打开失败 [Err: %s]\n", err.Error())
		os.Exit(0)
	}

	defer filePtr.Close()

	var anyJson map[string]interface{}
	decoder := json.NewDecoder(filePtr)

	err = decoder.Decode(&anyJson)

	if err != nil {
		fmt.Printf("解码失败 [Err: %s]\n", err.Error())
		os.Exit(0)
	}

	return anyJson
}

// 将 markdown 表格字符串写入文件
// 文件名是 json文件名 + “.md”
func writeMarkdownFile(tableHtml string, markdownFilePath string) {
	file, err := os.OpenFile(markdownFilePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Printf("Markdown 数据：\n%v", tableHtml)
		fmt.Printf("打开文件错误\n\t %v\n", err)
		os.Exit(0)
	}

	defer file.Close()

	writer := bufio.NewWriter(file)
	writer.WriteString(tableHtml)
	writer.Flush()
	fmt.Printf("成功创建 Markdown 表格文件：%v\n", markdownFilePath)
}

func getFullFilePathWithoutSuffix(filePath string) string {
	var fileNameWithSuffix string
	fileNameWithSuffix = path.Base(filePath)
	var fileSuffix string
	fileSuffix = path.Ext(fileNameWithSuffix)

	var filePathOnly string
	filePathOnly = strings.TrimSuffix(filePath, fileSuffix)

	return filePathOnly
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

// JSON 数组会被转换成 []interface 类型
// 所以需要将其转换成 []string 类型
func getJsonSlice(js interface{}) stringSlice {
	var options stringSlice

	switch js.(type) {
	case string, int, bool:
		options = append(options, fmt.Sprintf("%v", js))
	case []interface{}:
		for _, v := range js.([]interface{}) {
			options = append(options, fmt.Sprintf("%v", v))
		}
	default:
		panic("options must be slice, string, int, bool")
	}

	return options
}

// 获取map 所有的key
func getMapKeys(m map[string]interface{}) stringSlice {
	keys := make(stringSlice, 0)
	for k, _ := range m {
		keys = append(keys, k)
	}
	return keys
}

// 如果空接口 是实际类型 是 map, 获取它所有的key
func getEmptyInterfaceKeys(x interface{}) stringSlice {
	m, ok := x.(map[string]interface{})
	if !ok {
		panic("不是 map[string]interface{}，无法获取key")
	}

	return getMapKeys(m)
}

// 获取map 所有的key，并按字母排序
func getSortMapKeys(m map[string]interface{}) stringSlice {
	keys := getMapKeys(m)
	// 对切片进行排序
	sort.Strings(keys)

	return keys
}

// 创建 Markdown 表格
func createMarkdownTable(jsonData map[string]interface{}) string {
	headerHtml := createTableHeader()
	bodyHtml := createTableBody(jsonData)

	style := "style=\"width:100%\""
	tableHtml := fmt.Sprintf("<table %v>\n%v%v</table>\n", style, headerHtml, bodyHtml)

	return tableHtml
}

// 创建 表格头
func createTableHeader() string {
	var headerHtml string
	for i := 0; i < maxColspan; i++ {
		var thHtml string
		var thContent string

		if i == 0 {
			thContent = "参数"
		} else if i == maxColspan-1 {
			thContent = "释义"
		} else {
			thContent = "子参"
		}

		thHtml += fmt.Sprintf("\t<th>%v</th>\n", thContent)

		headerHtml += thHtml
	}
	headerHtml = fmt.Sprintf("<thead>\n%v</thead>\n", headerHtml)

	return headerHtml
}

// 创建 表格内容
func createTableBody(jsonData map[string]interface{}) string {
	var bodyHtml string
	setTableTr(jsonData, maxColspan, 0, "")
	trsHtml := createTableTrHtml()
	bodyHtml += fmt.Sprintf("<tbody>\n%v</tbody>\n", trsHtml)

	return bodyHtml
}

// 创建 表格行
func createTableTrHtml() string {
	var trsHtml string

	for _, tr := range trs {
		var trHtml string
		trHtml += fmt.Sprintf("\t<tr>\n")

		for _, td := range tr {
			trHtml += fmt.Sprintf("\t\t%v\n", td)
		}

		trHtml += fmt.Sprintf("\t</tr>\n")

		trsHtml += trHtml
	}

	return trsHtml
}

// 递归 JSON 数据，生成每一行的html 并保存到全局变量 trs
// 遍历当前 jsonData 的所有key 生成 <tr>s
// 1. 如果是最终元素，则取key生成参数名<td>, key值生成详情<td>，合并tds 生成一个 <tr>
// 2. 如果不是最终元素
//   2.1 一次性遍历获取 key的值(jsonData2) 内的所有的第一个子元素生成参数名<td>(s), 最后一个key值生成详情<td>，合并<td>s 生成一个 <tr>
//   2.2 递归遍历 key的值(jsonData2) 的所有key(排除第一个key) 生成 <tr>s
func setTableTr(jsonData map[string]interface{}, colspan int, depth int, keyPath string) {
	// fmt.Printf("当前JOSN data: %#v\n", jsonData)
	keys := getSortMapKeys(jsonData)

	keysLen := len(keys)
	keyOffset := maxColspan - colspan
	currentColspan := colspan - depth - 1
	// fmt.Printf("@@@ 最大列数是：%v， 当前列数是: %v, keyOffset 是：%v, keysLen 是：%v\n", maxColspan, colspan, keyOffset, keysLen)

	var currentKeyPath string

	if keyOffset > 0 && !isEndElement(jsonData[keys[keyOffset-1]]) {
		if keyPath == "" {
			currentKeyPath = keys[keyOffset-1]
		} else {
			currentKeyPath = keyPath + "." + keys[keyOffset-1]
		}
		setTableTr(jsonData[keys[keyOffset-1]].(map[string]interface{}), colspan, depth+1, currentKeyPath)
	}

	// for _, key := range keys {
	for i := keyOffset; i < keysLen; i++ {
		var tr tdSlice
		var td string

		currentKey := keys[i]

		if keyPath == "" {
			currentKeyPath = currentKey
		} else {
			currentKeyPath = keyPath + "." + currentKey
		}

		jsonDataValue := jsonData[currentKey]

		if isEndElement(jsonDataValue) {
			// fmt.Printf("=== 这是最终元素 - %v\n", currentKeyPath)

			td = createTableTdHtml(currentKey, 1, 1)
			tr = append(tr, td)

			elmDetails := new(elementDetails)
			elmDetails.setup(jsonDataValue)
			// fmt.Printf("### elmDetails: %#v\n", elmDetails)
			details := elmDetails.generateDetailsHtml()
			// fmt.Printf("*** elmDetails HTML: %#v\n\n", details)
			td = createTableTdHtml(details, currentColspan, 1)
			// fmt.Printf("*** %v - td HTML: %#v\n\n", currentKeyPath, td)
			tr = append(tr, td)
			appendTrs(tr)

			// fmt.Printf("$$$ tds - %v：%#v\n\n", currentKeyPath, tr)
		} else {
			// fmt.Printf("### 这不是最终元素 - %v\n", currentKeyPath)
			// fmt.Printf("### 所以需要先插入当前字段 %v，再一次性递归获取所有子字段的第一个字段\n", currentKey)

			// 递归获取key 值内的所有的第一个子元素生成参数名<td>, 最后一个key值生成详情<td>
			td = createTableTdHtml(currentKey, 1, getRowspan(jsonDataValue.(map[string]interface{})))
			tr = append(tr, td)

			firstColumns := getFirstColumns(jsonDataValue.(map[string]interface{}), currentColspan)
			tr = append(tr, firstColumns...)
			appendTrs(tr)

			// fmt.Printf("$$$ tds - %v：%#v\n\n", currentKeyPath, tr)

			setTableTr(jsonDataValue.(map[string]interface{}), currentColspan, depth, currentKeyPath)
		}

	}
}

// 创建 表格列
func createTableTdHtml(content string, colspan int, rowspan int) string {
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

// 插入全局变量 trs
// trs - 保存每一行的所有列
func appendTrs(tr tdSlice) {
	trs = append(trs, tr)
}

// 获取整个表格最大有多少列
func getMaxColspan(jsonData map[string]interface{}) int {
	colspan := 1
	for _, v := range jsonData {
		curColspan := 1
		if isEndElement(v) {
			curColspan = 2
		} else {
			curColspan += getMaxColspan(v.(map[string]interface{}))
		}

		if curColspan > colspan {
			colspan = curColspan
		}
	}
	return colspan
}

// 获取当前字段有多少个子字段，即当前字段占几行
func getRowspan(jsonData map[string]interface{}) int {
	rowspan := 0
	for _, v := range jsonData {
		if isEndElement(v) {
			rowspan += 1
		} else {
			rowspan += getRowspan(v.(map[string]interface{}))
		}
	}

	return rowspan
}

// 递归获取当前字段内第一个字段的td 的集合
func getFirstColumns(jsonData map[string]interface{}, colspan int) tdSlice {
	var tds tdSlice
	var td string

	keys := getSortMapKeys(jsonData)
	firstKey := keys[0]

	jsonDataValue := jsonData[firstKey]

	if isEndElement(jsonDataValue) {
		// fmt.Printf("====== 这是最终元素 - %v\n", firstKey)

		td = createTableTdHtml(firstKey, 1, 1)
		tds = append(tds, td)

		elmDetails := new(elementDetails)
		elmDetails.setup(jsonDataValue)
		// fmt.Printf("****** elmDetails: %#v\n", elmDetails)
		details := elmDetails.generateDetailsHtml()
		// fmt.Printf("****** elmDetails HTML: %#v\n\n", details)
		td := createTableTdHtml(details, colspan-1, 1)
		// fmt.Printf("****** %v - td HTML: %#v\n\n", firstKey, td)
		tds = append(tds, td)
	} else {
		// fmt.Printf("###### 这不是最终元素 - %v\n", firstKey)
		// fmt.Printf("###### 所以需要先插入当前字段 %v，再一次性递归获取所有子字段的第一个字段\n", firstKey)

		td = createTableTdHtml(firstKey, 1, getRowspan(jsonDataValue.(map[string]interface{})))
		tds = append(tds, td)

		nextTd := getFirstColumns(jsonDataValue.(map[string]interface{}), colspan-1)
		tds = append(tds, nextTd...)
	}

	return tds
}
