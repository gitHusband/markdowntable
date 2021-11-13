# markdowntable - JSON 转 Markdown 表格的工具

> 目录

[TOC]

## 1. 使用方法

```
./markdowntable [-in jsonFilePath] [-out markdownFilePath] [-sort default]

选项：
  in: 输入JSON 文件，默认 "./info.json"
  out: 输出 markdown 文件，默认 "./info.md"
  sort: 字段排序方式
      - default：默认按JSON原本字段顺序。
      - asc：递增排序
      - desc：递减排序
```
> 支持 MacOS, Linux 和 Windows 平台

**注：**

Markdown 按JSON原本字段顺序排序的实现？

[jsonkeys - 获取 JSON key 的先后顺序](https://github.com/gitHusband/goutils/tree/master/jsonkeys)
```
当使用GO 标准库 `encoding/json` 解析动态JSON 的时候，我们将结果解析为 `map[string]interface{}`。

而 GO `map` 类型的key 是无序的，也就是说你不能确定JSON key 的先后顺序。

如果你需要确定 JSON key 的顺序，可以使用 `jsonkeys` 包。
```

## 2. 如何设置字段的内容？
### 2.1 以字符串结尾
如果一个字段的值不再包含子元素，值是 字符串，那么字符串就是字段的内容，插入到表格最后一列

### 2.2 以特定结构结尾
特定结构指：

```
{
    "header": "标题 - 这是一个标题",
    "desc": "描述 - 这是一个描述内容",
    "defaultValue": "默认值：Hello world",
    "options": [
        "字段接受的值1 - Hello world",
        "字段接受的值2 - 你好世界"
    ]
}
```
如果一个字段的值是以上的特定结构，那么该特定结构就是字段的内容，插入到表格最后一列。

说明：

```
"header" - 必须包含该字段
"desc" - 必须包含该字段
"defaultValue" - 可选字段
"options" - 可选字段, 支持数组

# 注：不支持自定义字段
```


## 3. 例子
以以下JSON 为例
```
{
    "1_1": {
        "1_2": {
            "1_3": {
                "header": "标题 - 第 1 行 第 4 列",
                "desc": "描述 - 这是一个描述内容",
                "defaultValue": "Hello world",
                "options": [
                    "Hello world",
                    "你好世界"
                ]
            },
            "2_3": {
                "2_4": "第 2 行 第 5 列",
                "3_4": "第 3 行 第 5 列"
            }
        },
        "4_2": {
            "4_3": {
                "4_4": "第 4 行 第 5 列",
                "5_4": "第 5 行 第 5 列"
            },
            "6_3": "第 6 行 第 4 列",
            "7_3": {
                "7_4": "第 7 行 第 5 列",
                "8_4": {
                    "header": "标题 - 第 8 行 第 5 列",
                    "desc": "描述 - 这是一个描述内容",
                    "defaultValue": "Hello world",
                    "options": [
                        "Hello world",
                        "你好世界"
                    ]
                }
            }
        }
    },
    "9_1": "第 9 行 第 2 列",
    "a0_1": {
        "a0_2": {
            "a0_3": {
                "a0_4": "第 10 行 第 5 列"
            }
        }
    },
    "a1_1": {
        "a1_2": "第 11 行 第 3 列"
    }
}
```
它将转换成 表格

<table style="width:100%">
<thead>
	<th>参数</th>
	<th>子参</th>
	<th>子参</th>
	<th>子参</th>
	<th>释义</th>
</thead>
<tbody>
	<tr>
		<td rowspan="8">1_1</td>
		<td rowspan="3">1_2</td>
		<td>1_3</td>
		<td colspan="2"><strong style="font-size: 15px">标题 - 第 1 行 第 4 列</strong><br/> <em style="color: #888888">描述 - 这是一个描述内容</em><br/> <b>默认：<ins>Hello world</ins></b><br/> <b>可能的值：</b> <ul><li>Hello world</li><li>你好世界</li></ul></td>
	</tr>
	<tr>
		<td rowspan="2">2_3</td>
		<td>2_4</td>
		<td><strong style="font-size: 15px">第 2 行 第 5 列</strong></td>
	</tr>
	<tr>
		<td>3_4</td>
		<td><strong style="font-size: 15px">第 3 行 第 5 列</strong></td>
	</tr>
	<tr>
		<td rowspan="5">4_2</td>
		<td rowspan="2">4_3</td>
		<td>4_4</td>
		<td><strong style="font-size: 15px">第 4 行 第 5 列</strong></td>
	</tr>
	<tr>
		<td>5_4</td>
		<td><strong style="font-size: 15px">第 5 行 第 5 列</strong></td>
	</tr>
	<tr>
		<td>6_3</td>
		<td colspan="2"><strong style="font-size: 15px">第 6 行 第 4 列</strong></td>
	</tr>
	<tr>
		<td rowspan="2">7_3</td>
		<td>7_4</td>
		<td><strong style="font-size: 15px">第 7 行 第 5 列</strong></td>
	</tr>
	<tr>
		<td>8_4</td>
		<td><strong style="font-size: 15px">标题 - 第 8 行 第 5 列</strong><br/> <em style="color: #888888">描述 - 这是一个描述内容</em><br/> <b>默认：<ins>Hello world</ins></b><br/> <b>可能的值：</b> <ul><li>Hello world</li><li>你好世界</li></ul></td>
	</tr>
	<tr>
		<td>9_1</td>
		<td colspan="4"><strong style="font-size: 15px">第 9 行 第 2 列</strong></td>
	</tr>
	<tr>
		<td>a0_1</td>
		<td>a0_2</td>
		<td>a0_3</td>
		<td>a0_4</td>
		<td><strong style="font-size: 15px">第 10 行 第 5 列</strong></td>
	</tr>
	<tr>
		<td>a1_1</td>
		<td>a1_2</td>
		<td colspan="3"><strong style="font-size: 15px">第 11 行 第 3 列</strong></td>
	</tr>
</tbody>
</table>