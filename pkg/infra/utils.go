package infra

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

func str2Bytes(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bh))
}

func uint16ToBytes(num uint16) []byte {
	data := make([]byte, 2)
	binary.LittleEndian.PutUint16(data, num)
	return data
}

func uint32ToBytes(num uint32) []byte {
	data := make([]byte, 4)
	binary.LittleEndian.PutUint32(data, num)
	return data
}

func uint64ToBytes(num uint64) []byte {
	data := make([]byte, 8)
	binary.LittleEndian.PutUint64(data, num)
	return data
}

func getFieldIndex(tag, key string, s interface{}) string {
	rt := reflect.TypeOf(s)
	if rt.Kind() != reflect.Struct {
		panic(fmt.Errorf("this is not a qualified type, use struct please"))
	}

	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		v := strings.Split(f.Tag.Get(key), ",")[0]
		if strings.ToUpper(v) == strings.ToUpper(tag) {
			return f.Name
		}
	}

	panic("the structure has no specified label")
}

func FormatFloat(num float64, decimal int) float64 {
	// default multiplication by 1
	d := float64(1)
	if decimal > 0 {
		// 10 to the N power
		d = math.Pow10(decimal)
	}

	res := strconv.FormatFloat(math.Trunc(num*d)/d, 'f', -1, 64)
	reDecimal, err := strconv.ParseFloat(res, 64)
	if err != nil {
		log.Println(err)
		return num
	}
	return reDecimal
}

func HexToUnicode(hexStr string) string {
	s, err := hex.DecodeString(hexStr)
	if err != nil {
		log.Printf("this is not a hexadecimal string : %s", err)
		return ""
	}

	return string(s)
}

func TimeNowString() string {
	now := time.Now()
	return fmt.Sprintf("%d-%d-%d", now.Year(), now.Month(), now.Day())
}

func FormattingTime(date string) string {
	if date == "" {
		log.Println("DATE is null")
		return ""
	}

	// failed to parse date field [2021-11-7] with format [yyyy-MM-dd]
	reFormatFullTime := func(spDateInt []int) string {
		var spDate []string

		// for months and days we treat them as two digits
		for i, v := range spDateInt {
			s := strconv.Itoa(v)

			// eliminate the year and determine the month and day
			if i != 0 {
				if len(s) != 2 {
					s = fmt.Sprintf("0%s", s)
				}
			}

			spDate = append(spDate, s)
		}

		return strings.Join(spDate, "-")
	}

	// yyyy-MM-dd
	parse, err := time.Parse("2006-01-02", date)
	if err != nil {
		log.Println(err)
		return ""
	}

	return reFormatFullTime([]int{parse.Year(), int(parse.Month()), parse.Day()})
}
