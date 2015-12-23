// main.go
package pinyinchinesechar

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io/ioutil"
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

func PinyinToChineseStr(py string) string {
	var b bytes.Buffer
	list:=PinyinToChinseChar(py)
	for _,char := range list {
		b.WriteString(char)
	}
	return b.String() 
}

//the py is supposed to be only for a single char, like zhong, not zhong guo.
//py may also be partial prefix like zh
func PinyinToChinseChar(py string) []string {
	resultWholeMatch :=[]string{}
	resultPrefixMatch := []string{}

	for _, c := range gChineseChars {
		for _, pyInList := range c.pinyinList {
			py = strings.ToLower(py)

			if py == pyInList {
				resultWholeMatch = append(resultWholeMatch,c.theChar)
				continue
			}			
			
			if strings.HasPrefix(pyInList, py) {
				resultPrefixMatch = append(resultPrefixMatch, c.theChar)
				continue
			}
		}
	}
	
	return append(resultWholeMatch,resultPrefixMatch...)
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
func ChinseCharToPinyin(theChar string) ([]string, error) {
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

func readChineseCharCsvBytes(content []byte) ([]tChineseChar, error) {
	result := []tChineseChar{}

	reader := bytes.NewReader(content)
	csvReader := csv.NewReader(reader)
	csvReader.Comment = '#'

	var metas [][]string
	metas, err := csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("read csv content failed")
	}

	for _, metaFields := range metas {
		theChar := strings.TrimSpace(metaFields[0])
		frequency, err := strconv.Atoi(strings.TrimSpace(metaFields[1]))
		if err != nil {
			return nil, fmt.Errorf("cannot convert to number %s", metaFields[1])
		}
		pinyinStr, err := ChinseCharToPinyin(theChar)
		if err != nil {
			return nil, err
		}
		c := tChineseChar{theChar: theChar, relativeFrequency: frequency, pinyinList: pinyinStr}
		result = append(result, c)
	}
	return result, nil
}

func readChineseCharCsvFile(filePath string) ([]tChineseChar, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return readChineseCharCsvBytes(content)
}

func init(){
	var err error 
	gChineseChars, err = readChineseCharCsvBytes([]byte(gChineseCharsCsv))
	if err!=nil {
		panic(err)
	}
}


