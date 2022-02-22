// Pakage utils 유틸리티 패키지 전체 어플리케이션에서 사용
package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

var logFn = log.Panic

func HandleErr(err error) {
	if err != nil {
		logFn(err)
	}
}

func ToBytes(i interface{}) []byte {
	var aBuffer bytes.Buffer
	encoder := gob.NewEncoder(&aBuffer)
	HandleErr(encoder.Encode(i))
	return aBuffer.Bytes()
}

// 데이터를 바이트로 변환해줌
func FromBytes(i interface{}, data []byte) {
	encoder := gob.NewDecoder(bytes.NewReader(data)).Decode(i)
	HandleErr(encoder)
}

// 데이터를 Hash로 변환시켜 string 형태롤 반환해줌
func Hash(i interface{}) string {
	s := fmt.Sprintf("%v", i)
	hash := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", hash)
}

func Splitter(s string, sep string , i int) string {
	r := strings.Split(s, sep)
	if len(r)-1 < i {
		return ""
	}
	return r[i]
}

func ToJSON(i interface{}) []byte {
	r, err := json.Marshal(i)
	HandleErr(err)
	return r
}