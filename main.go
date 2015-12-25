// main.go
package main

import (
	"github.com/gopherjs/gopherjs/js"
	"ponemonic/pinyinchinesechar"
)

func main() {

	js.Global.Set("pinyinchinesechar", map[string]interface{}{
		"InputToOutput":         pinyinchinesechar.InputToOutput,
		"PinyinStrToChineseStr": pinyinchinesechar.PinyinStrToChineseStr,
		"NumStrToChineseStr":    pinyinchinesechar.NumStrToChineseStr,
		"ChineseStrToDigitStr":  pinyinchinesechar.ChineseStrToDigitStr,
		"DigitToConsonantTable": pinyinchinesechar.DigitToConsonantTable,
	})

}
