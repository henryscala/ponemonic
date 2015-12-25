package pinyinchinesechar

import (
	"testing"
)

func TestJudgeInputType(t *testing.T) {

	if judgeInputType("a1,2") != inputTypeInt {
		t.Fail()
	}
	if judgeInputType("cd中国e") != inputTypePinyin {
		t.Fail()
	}
	if judgeInputType(" 天安门 ") != inputTypeChinese {
		t.Fail()
	}
}

func TestNumStrToChineseStr(t *testing.T) {
	input := "1"
	output := NumStrToChineseStr(input)
	t.Log("input", input)
	t.Log("output", output)
}
