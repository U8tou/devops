package datacv

import (
	"strconv"
	"strings"
)

func IntSliceToStrSlice(args []int) []string {
	var strSlice []string
	if len(args) == 0 {
		return strSlice
	}
	for _, number := range args {
		strSlice = append(strSlice, strconv.Itoa(number))
	}
	return strSlice
}

func StrSliceToIntSlice(args []string) []int {
	var intSlice []int
	if len(args) == 0 {
		return intSlice
	}
	for _, str := range args {
		i, err := strconv.Atoi(str)
		if err != nil {
			return make([]int, 0)
		}
		intSlice = append(intSlice, i)
	}
	return intSlice
}

func StrToInt(str string) int64 {
	num, _ := strconv.ParseInt(str, 10, 64)
	return num
}

func IntToStr(it int64) string {
	return strconv.FormatInt(it, 10)
}

func StrToInts(it string) []int {
	strs := strings.Split(it, ",")
	ls := make([]int, len(strs))
	for idx, str := range strs {
		num, _ := strconv.ParseInt(str, 10, 64)
		ls[idx] = int(num)
	}
	return ls
}

func IntsToStr(it []int) string {
	str := ""
	for idx, i := range it {
		if idx == 0 {
			str = strconv.Itoa(i)
		} else {
			str += "," + strconv.Itoa(i)
		}
	}
	return str
}

// Int64SliceToStrSlice []int64 转 []string
func Int64SliceToStrSlice(args []int64) []string {
	strSlice := make([]string, len(args))
	if len(args) == 0 {
		return strSlice
	}
	for idx, number := range args {
		strSlice[idx] = strconv.FormatInt(number, 10)
	}
	return strSlice
}

// StrSliceToInt64Slice []string 转 []int64
func StrSliceToInt64Slice(args []string) []int64 {
	int64Slice := make([]int64, len(args))
	if len(args) == 0 {
		return make([]int64, 0)
	}
	for idx, str := range args {
		i, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return make([]int64, 0)
		}
		int64Slice[idx] = i
	}
	return int64Slice
}
