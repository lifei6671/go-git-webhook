gojson
===
[![Build Status](https://drone.io/github.com/widuu/gojson/status.png)](https://drone.io/github.com/widuu/gojson/latest)

---
**[The official website](http://www.widuu.com)**

##中文文档

###简介

>gojson是快速解析json数据的一个golang包，你使用它可以快速的查找json内的数据

###安装

 	 原作者 go get github.com/widuu/gojson
 	 修改过的包 go get github.com/zyx4843/gojson

###修改说明
	
	增加函数IsValid()判断数据是否有效
	Getindex防止出现越界
	Tostring增加bool类型

###使用简介

结构

	type Js struct {
		data interface{}
	}

(1) `func Json(data) *Js` data为string类型,初始化Js结构，解析json并且return `Js.data`

	json := `{"from":"en","to":"zh"}`
	c1 := gojson.Json(json) //&{map[from:en to:zh]}

(2) `func (*Js) Get() *js` 获取简单json中的某个值，递归查找，return `Js.data`

	json := `{"from":"en","to":"zh","trans_result":{"src":"today","dst":"\u4eca\u5929"},"result":["src","today","dst","\u4eca\u5929"]}`
	
 	c2 := gojson.Json(json).Get("trans_result").Get("dst")
	fmt.Println(c2) //&{今天}

	c2 := gojson.Json(json).Get("from")
	fmt.Println(c2) //&{en}

(3) `func (*Js)Tostring()string` 将单个数据转化成string类型,因为string类型转其它类型都比较好转就让数据返回string

	c2 := gojson.Json(json).Get("from").Tostring()
	fmt.Println(c2) //en

(4) `func (j *Js) Getpath(args ...string) *Js ` 通过输入string的多个参数来获取某个值，json数据一定要是递归的

 	c4 := gojson.Json(json).Getpath("trans_result", "src").Tostring()
	fmt.Println(c4)  //today

(5) `func (j *Js) Arrayindex(i int) string` 获取Json数据中数组结构的值，根据输入的num来返回对应的值，仅限于处理{"result":["src","today","dst","\u4eca\u5929"]}中[]内的值

	json := `{"from":"en","to":"zh","trans_result":{"src":"today","dst":"\u4eca\u5929"},"result":["src","today","dst","\u4eca\u5929"]}`
	c7 := gojson.Json(json).Get("result").Arrayindex(1)
	fmt.Println(c7) //src

(6) `func (j *Js) Getkey(key string, i int) *Js` 这个函数是针对数据中有重复数据，取值，使用js.data必须是[]interface{}类型，这个是百度翻译时候返回的js可能会用到

 	json1 := `{"from":"en","to":"zh","trans_result":[{"src":"today","dst":"\u4eca\u5929"},{"src":"tomorrow","dst":"\u660e\u5929"}]}`
	c8 := gojson.Json(json1).Get("trans_result").Getkey("src", 1).Tostring()
	fmt.Println(c8) //则返回trans_result第一组中的src today

(7) `func (j *Js) ToArray() (k, d []string) `将json数据转换成key []string{} value []string{} 一一对应的数组，只能使用到二级 不能到多级

	
 	c9k, c9v := gojson.Json(json1).Get("trans_result").ToArray()
	fmt.Println(c9k, c9v) //[src dst src dst] [today 今天 tomorrow 明天]

	c3k, c3v := gojson.Json(json).Getindex(1).ToArray()
	fmt.Println(c3k, c3v) //	[from] [en]

(8) `func (j *Js) Getindex(i int) *Js ` 根据i返回json内的数据，可以逐级查找

	
	json1 := `{"from":"en","to":"zh","trans_result":[{"src":"today","dst":"\u4eca\u5929"},{"src":"tomorrow","dst":"\u660e\u5929"}]}`

	c10 := gojson.Json(json1).Getindex(3).Getindex(1).Getindex(1).Get("src").Tostring()
	fmt.Println(c10) //today

(9) `func (j *Js) StringtoArray() []string` 将{"result":["src","today","dst","\u4eca\u5929"]}数据json中的result对应的数据，返回成[]string的slice

	c11 := gojson.Json(json).Get("result").StringtoArray()
	fmt.Println(c11) //[src today dst 今天]

(10) `func (j *Js) Type()` 打印测试用，打印数据类型

	gojson.Json(json).Get("result").Type()  //[]interface {}

---


##English Document

###Introduction

>gojson是快速解析json数据的一个golang包，你使用它可以快速的查找json内的数据

###Install

 	 go get github.com/widuu/gojson

###Use

structure

	type Js struct {
		data interface{}
	}

(1) `func Json(data) *Js` data is string type,initialize Js structure, parse the json and return `Js.data`

	json := `{"from":"en","to":"zh"}`
	c1 := gojson.Json(json) //&{map[from:en to:zh]}

(2) `func (*Js) Get() *js` Retrieve a value from a simple json, recursive lookup，return `Js.data`

	json := `{"from":"en","to":"zh","trans_result":{"src":"today","dst":"\u4eca\u5929"},"result":["src","today","dst","\u4eca\u5929"]}`
	
 	c2 := gojson.Json(json).Get("trans_result").Get("dst")
	fmt.Println(c2) //&{今天}

	c2 := gojson.Json(json).Get("from")
	fmt.Println(c2) //&{en}

(3) `func (*Js)Tostring()string` To a single data into type string, for type string turned other types are better let return to string data

	c2 := gojson.Json(json).Get("from").Tostring()
	fmt.Println(c2) //en

(4) `func (j *Js) Getpath(args ...string) *Js ` Get a value by multiple parameters of the input string, the json data must be recursive

 	c4 := gojson.Json(json).Getpath("trans_result", "src").Tostring()
	fmt.Println(c4)  //today

(5) `func (j *Js) Arrayindex(i int) string` Get the value of the array structure in the Json data, based on the input of num to return the value of the corresponding is limited to processing {" result ": [" SRC", "today", "DST", "\ u4eca \ u5929"]} [] in value

	json := `{"from":"en","to":"zh","trans_result":{"src":"today","dst":"\u4eca\u5929"},"result":["src","today","dst","\u4eca\u5929"]}`
	c7 := gojson.Json(json).Get("result").Arrayindex(1)
	fmt.Println(c7) //src

(6) `func (j *Js) Getkey(key string, i int) *Js` There is duplicate data in this function is for data, the values, use js. Data must be [] interface {} type, this is baidu translation when the returned js might need

 	json1 := `{"from":"en","to":"zh","trans_result":[{"src":"today","dst":"\u4eca\u5929"},{"src":"tomorrow","dst":"\u660e\u5929"}]}`
	c8 := gojson.Json(json1).Get("trans_result").Getkey("src", 1).Tostring()
	fmt.Println(c8) //则返回trans_result第一组中的src today

(7) `func (j *Js) ToArray() (k, d []string) ` Converting the json data into key string [] {} the value string [] {} one-to-one array, can only be used to level 2 not to multistage

	
 	c9k, c9v := gojson.Json(json1).Get("trans_result").ToArray()
	fmt.Println(c9k, c9v) //[src dst src dst] [today 今天 tomorrow 明天]

	c3k, c3v := gojson.Json(json).Getindex(1).ToArray()
	fmt.Println(c3k, c3v) //	[from] [en]

(8) `func (j *Js) Getindex(i int) *Js ` According to the I returned in the json data, can find step by step

	
	json1 := `{"from":"en","to":"zh","trans_result":[{"src":"today","dst":"\u4eca\u5929"},{"src":"tomorrow","dst":"\u660e\u5929"}]}`

	c10 := gojson.Json(json1).Getindex(3).Getindex(1).Getindex(1).Get("src").Tostring()
	fmt.Println(c10) //today

(9) `func (j *Js) StringtoArray() []string` Will {" result ":" SRC ", "today", "DST", "\ u4eca \ u5929"]} corresponding data, the result of the json data returned into [] string slice

	c11 := gojson.Json(json).Get("result").StringtoArray()
	fmt.Println(c11) //[src today dst 今天]

(10) `func (j *Js) Type()` Print test, the data type

	gojson.Json(json).Get("result").Type()  //[]interface {}








