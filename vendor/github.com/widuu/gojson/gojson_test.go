package gojson

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	json := `{"from":"en","to":"zh","trans_result":{"src":"today","dst":"\u4eca\u5929"},"result":["src","today","dst","\u4eca\u5929"]}`
	c1 := Json(json).Get("from").Tostring()
	fmt.Println(c1) //stdout en
	c2 := Json(json).Get("trans_result").Get("dst")
	fmt.Println(c2) //stdout 今天
	c3k, c3v := Json(json).Getindex(1).ToArray()
	fmt.Println(c3k, c3v) //	[from] [en]

	c4 := Json(json).Getpath("trans_result", "src").Tostring()
	fmt.Println(c4) //today

	c5 := Json(json).Get("trans_result").Getindex(2)
	fmt.Println(c5) //&{map[dst:今天]}
	json2 := `{"from":"en","to":"zh"}`
	c6 := Json(json2)
	fmt.Println(c6)

	c7 := Json(json).Get("result").Arrayindex(1)
	fmt.Println(c7) //src

	json1 := `{"from":"en","to":"zh","trans_result":[{"src":"today","dst":"\u4eca\u5929"},{"src":"tomorrow","dst":"\u660e\u5929"}]}`
	c8 := Json(json1).Get("trans_result").Getkey("src", 1).Tostring()
	fmt.Println(c8) //today

	c9k, c9v := Json(json1).Get("trans_result").ToArray()
	fmt.Println(c9k, c9v) //[src dst src dst] [today 今天 tomorrow 明天]

	c10 := Json(json1).Getindex(3).Getindex(1).Getindex(1).Get("src").Tostring()
	fmt.Println(c10) //today

	c11 := Json(json).Get("result").StringtoArray()
	fmt.Println(c11) //[src today dst 今天]

	Json(json).Get("result").Type() //[]interface {}

}
