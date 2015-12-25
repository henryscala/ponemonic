// main.go
package pinyinchinesechar

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"regexp"
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
var gDigitCharDigit map[string]int

var gVowel []string = []string{"a", "e", "i", "o", "u"}
var gConsonant []string = []string{"b", "p", "m", "f",
	"d", "t", "n", "l",
	"g", "k", "h", "j",
	"q", "x", "zh", "ch",
	"sh", "r", "z", "c", "s",
	"y", "w"}
var gConsonantNum map[string]int = map[string]int{
	"b": 8, "p": 4, "m": 5, "f": 8,
	"d": 9, "t": 2, "n": 2, "l": 6,
	"g": 7, "k": 3, "h": 3, "j": 9,
	"q": 7, "x": 0, "zh": 0, "ch": 6,
	"sh": 4, "r": 1, "z": 0, "c": 6, "s": 4,
	"y": 1, "w": 5,
}
var gNumConsonant [][]string

//judge whether a string contains number
var gRegexpNumber *regexp.Regexp = regexp.MustCompile(`\d`)

//judge whether a string contains english char
var gRegexpEnglish *regexp.Regexp = regexp.MustCompile(`[a-zA-Z]`)

const (
	chineseCharFileName = "chinese_character_frequency_6763.csv"
)

const (
	inputTypeInt int = iota
	inputTypePinyin
	inputTypeChinese

	numDigit = 10
)

func InputToOutput(input string) (output string) {
	i := strings.Trim(input, " \t\r\n")
	category := judgeInputType(i)
	switch category {
	case inputTypeInt:
		output = NumStrToChineseStr(i)
	case inputTypePinyin:
		output = PinyinStrToChineseStr(i)
	case inputTypeChinese:
		output = ChineseStrToDigitStr(i)
	}
	return
}

func judgeInputType(input string) int {
	if gRegexpNumber.MatchString(input) {
		return inputTypeInt
	}
	if gRegexpEnglish.MatchString(input) {
		return inputTypePinyin
	}
	return inputTypeChinese
}

// Convert the pinyin of a char to the corresponding digit
// e.g. : b -> 8, ban -> 8
func pinyinToDigit(py string) int {

	for _, c := range gConsonant {
		if strings.HasPrefix(py, c) {
			return gConsonantNum[c]
		}

	}

	return 0
}

// Convert a number(might be several digits) to possible Chinese chars
// The output format is fixed
func NumStrToChineseStr(numStr string) string {
	var b bytes.Buffer
	for _, c := range numStr {
		digitChar := string(c)
		digit, ok := gDigitCharDigit[digitChar]

		if ok {
			b.WriteString(fmt.Sprintf("[%d[", digit))
			pinyinList := DigitToConsonant(digit)

			for i, py := range pinyinList {

				chStr := pinyinToChineseStr(py)
				if i == 0 {
					b.WriteString("\n")
				}
				b.WriteString(fmt.Sprintf("  %s -> %s\n", py, chStr))
			}
			b.WriteString(fmt.Sprintf("]]\n\n"))
		} else {
			b.WriteString(fmt.Sprintf("[[%s]]\n\n", digitChar))
		}

	}
	return b.String()
}

// Convert a Chinese str(may contain multiple chars)
// to a number(multiple digits)
// A char may contributate multiple digits
func ChineseStrToDigitStr(str string) string {
	var b bytes.Buffer
	for _, c := range str {
		digits := ChineseCharToDigit(string(c))
		if len(digits) == 0 {
			b.WriteString("[]")
		} else if len(digits) == 1 {
			b.WriteString(fmt.Sprintf("%d", digits[0]))
		} else {
			b.WriteString("[")
			for _, d := range digits {
				b.WriteString(fmt.Sprintf("%d", d))
			}
			b.WriteString("]")
		}
	}
	return b.String()
}

//a single chinese char to digits, might be multiple pronouncation,
//so might be multiple results
func ChineseCharToDigit(char string) []int {
	result := []int{}
	pinyinList, err := ChinseCharToPinyin(char)
	if err != nil {
		return result
	}

	for _, py := range pinyinList {
		result = append(result, pinyinToDigit(py))
	}
	return result
}

// Get a string with the mappings from digit to pinyin char
func DigitToConsonantTable() string {
	var b bytes.Buffer

	for i := 0; i < numDigit; i++ {
		list := DigitToConsonant(i)
		b.WriteString(fmt.Sprintf("%d = %v\n", i, list))
	}
	return b.String()
}

func DigitToConsonant(digit int) []string {
	digit = digit % numDigit
	if gNumConsonant == nil {
		gNumConsonant = make([][]string, numDigit)
		for k, v := range gConsonantNum {
			if gNumConsonant[v] == nil {
				gNumConsonant[v] = make([]string, 0)
			}
			gNumConsonant[v] = append(gNumConsonant[v], k)
		}
	}
	return gNumConsonant[digit]
}

//A list in string format (separated by comma, space char) to string list
func listInStrToList(listStr string) (result []string) {
	var list []string
	list = strings.Split(listStr, " ") //split by space char
	var listOfList [][]string = make([][]string, len(list))
	for i, str := range list {
		listOfList[i] = strings.Split(str, ",") //split by comma
	}
	for _, list := range listOfList {
		result = append(result, list...)
	}
	return result
}

//A list of pinyin(separated by comma, space char) to string that separated
//by predefined delimiters
func PinyinStrToChineseStr(pyList string) string {
	list := pinyinStrToChineseStrList(pyList)
	var b bytes.Buffer
	const lineSep string = "\n--------\n"
	for i, s := range list {
		if i != 0 {
			b.WriteString(lineSep)
		}
		b.WriteString(s)
		b.WriteString("\n")
	}
	return b.String()
}

//A list of pinyin(separated by comma, space char) to string list
func pinyinStrToChineseStrList(pyList string) (result []string) {
	result = listInStrToList(pyList)

	for i := range result {
		result[i] = pinyinToChineseStr(result[i])
	}
	return result
}

//from a single pinyin to get a string with multiple chinese chars
func pinyinToChineseStr(py string) string {
	var b bytes.Buffer
	list := PinyinToChinseChar(py)
	for _, char := range list {
		b.WriteString(char)
	}
	return b.String()
}

//the py is supposed to be only for a single char, like zhong, not zhong guo.
//py may also be partial prefix like zh
func PinyinToChinseChar(py string) []string {
	resultWholeMatch := []string{}
	resultPrefixMatch := []string{}

	for _, c := range gChineseChars {
		for _, pyInList := range c.pinyinList {
			py = strings.ToLower(py)

			if py == pyInList {
				resultWholeMatch = append(resultWholeMatch, c.theChar)
				continue
			}

			if strings.HasPrefix(pyInList, py) {
				resultPrefixMatch = append(resultPrefixMatch, c.theChar)
				continue
			}
		}
	}

	return append(resultWholeMatch, resultPrefixMatch...)
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

func init() {
	var err error
	gChineseChars, err = readChineseCharCsvBytes([]byte(gChineseCharsCsv))
	if err != nil {
		panic(err)
	}

	gDigitCharDigit = make(map[string]int)
	for i := 0; i < numDigit; i++ {
		gDigitCharDigit[fmt.Sprintf("%d", i)] = i
	}
}
