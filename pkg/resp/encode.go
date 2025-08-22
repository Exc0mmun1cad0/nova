package resp

import (
	"bytes"
	"fmt"
	"strconv"
)

var (
	NullString = []byte("$-1\r\n")
	NullArray = []byte("*0\r\n")
)

func EncodeSimpleString(str string) []byte {
	res := fmt.Sprintf("+%s\r\n", str)
	return []byte(res)
}

func EncodeError(errMsg string) []byte {
	res := fmt.Sprintf("-%s\r\n", errMsg)
	return []byte(res)
}

func EncodeString(str string) []byte {
	res := fmt.Sprintf("$%d\r\n%s\r\n", len(str), str)
	return []byte(res)
}

func EncodeArray(strs []string) []byte {
	var b bytes.Buffer

	b.WriteByte('*')
	b.WriteString(strconv.Itoa(len(strs)))
	b.WriteString("\r\n")

	for _, str := range strs {
		b.Write(EncodeString(str))
	}

	return b.Bytes()
}

func EncodeInt(num int) []byte {
	res := fmt.Sprintf(":%d\r\n", num)
	return []byte(res)
}
