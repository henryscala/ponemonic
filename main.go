// main.go
package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"

	pinyin "github.com/mozillazg/go-pinyin"
)

type tChineseChar struct {
	theChar           string
	relativeFrequency int
	pinyinList        []string
}

var gChineseChars []tChineseChar
var gPinyinArgs *pinyin.Args
var gVowel []string = []string{"a", "e", "i", "o", "u"}
var gConsonant []string = []string{"b", "p", "m", "f",
	"d", "t", "n", "l",
	"g", "k", "h", "j",
	"q", "x", "zh", "ch",
	"sh", "r", "z", "c", "s",
	"y", "w"}

const (
	chinese_char_file_name = "chinese_character_frequency_6763.csv"
)

//the py is supposed to be only for a single char, like zhong, not zhong guo.
//py may also be partial prefix like zh
func pinyinToChinseChar(py string) []string {
	result := make([]string, 0)
	
	for _, c := range gChineseChars {
		for _, pyInList := range c.pinyinList {
			py = strings.ToLower(py)
			
			if strings.HasPrefix(pyInList, py) {
				result = append(result, c.theChar)
			}
		}
	}
	return result
}

func removeDuplicates(list []string) []string {
	m := make(map[string]string)
	result := make([]string, 0)
	for _, s := range list {
		m[s] = s
	}
	for key, _ := range m {
		result = append(result, key)
	}
	return result
}

//return list of pinyin ( multiple pinyin for a char may exist)
//theChar shall be a single chinese char
func chinseCharToPinyin(theChar string) ([]string, error) {
	if gPinyinArgs == nil {
		gPinyinArgs = new(pinyin.Args)
		*gPinyinArgs = pinyin.NewArgs()
		gPinyinArgs.Heteronym = true
	}
	listOfPinyinList := pinyin.Pinyin(theChar, *gPinyinArgs)
	if len(listOfPinyinList) <= 0 {
		return nil, fmt.Errorf("didnot get any pinyin from the chinese char %s", theChar)
	}
	pinyinList := listOfPinyinList[0]
	if len(pinyinList) <= 0 {
		return nil, fmt.Errorf("get empty pinyin from the chinese char %s", theChar)
	}

	//pinyinList may contain multiple same pinyin based on my test, make it unique
	return removeDuplicates(pinyinList), nil
}

func readChineseCharCsv(filePath string) ([]tChineseChar, error) {
	result := []tChineseChar{}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("cannot open file %s", filePath)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	csvReader := csv.NewReader(reader)
	csvReader.Comment = '#'

	var metas [][]string
	metas, err = csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("read csv file failed %s", filePath)
	}

	for _, metaFields := range metas {
		theChar := strings.TrimSpace(metaFields[0])
		frequency, err := strconv.Atoi(strings.TrimSpace(metaFields[1]))
		if err != nil {
			return nil, fmt.Errorf("cannot convert to number %s", metaFields[1])
		}
		pinyinStr, err := chinseCharToPinyin(theChar)
		if err != nil {
			return nil, err
		}
		c := tChineseChar{theChar: theChar, relativeFrequency: frequency, pinyinList: pinyinStr}
		result = append(result, c)
	}
	return result, nil
}

func main() {
	var err error
	gChineseChars, err = readChineseCharCsv(chinese_char_file_name)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}

	fmt.Printf("there are %d chars totally\n", len(gChineseChars))

	fmt.Printf("%s %v", os.Args[1], pinyinToChinseChar(os.Args[1]))

}
