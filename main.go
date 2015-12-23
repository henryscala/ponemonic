// main.go
package main

import (
	
	
	"github.com/gopherjs/gopherjs/js"
	"ponemonic/pinyinchinesechar"
)

func main() {
	

	js.Global.Set("pinyinchinesechar", map[string]interface{}{
		"PinyinToChineseStr": pinyinchinesechar.PinyinToChineseStr,

	})

}
