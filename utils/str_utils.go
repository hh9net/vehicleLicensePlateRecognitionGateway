//字符串工具
package utils

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"strings"
	"time"
	//uuid "github.com/satori/go.uuid"
)

var randomStrSource = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

//获取随机字符串
func GetRandomStr(length int) string {
	result := make([]byte, length)
	r := rand.New(rand.NewSource(time.Now().UnixNano() + rand.Int63())) //增大随机性
	for i := 0; i < length; i++ {
		result[i] = randomStrSource[r.Intn(len(randomStrSource))]
	}
	return string(result)
}

//生成纯数字的随机字符串
func GetRandomNumStr(length int) string {
	result := make([]byte, length)
	r := rand.New(rand.NewSource(time.Now().UnixNano() + rand.Int63())) //增大随机性
	for i := 0; i < length; i++ {
		result[i] = byte('0' + r.Intn(10)) //0 - 9
	}
	return string(result)
}

//func UUID() (string, error) {
//	id, err := uuid.NewV1()
//	if err != nil {
//		return "", err
//	}
//	return id.String(), nil
//}

// bool to int64(1,0)
func Uint8ToBool(data uint8) bool {
	if data > 0 {
		return true
	}
	return false
}

func ToString(i interface{}) string {
	if b, err := json.Marshal(i); err == nil {
		return string(b)
	}
	return "转换失败..."
}

func MakeUidWithTime(prefix string, random_len int) string {
	const time_base_format = "060102030405"
	uid := prefix + time.Now().Format(time_base_format) + GetRandomNumStr(random_len)
	return uid
}

func GetSha1(data string) string {
	t := sha1.New()
	io.WriteString(t, data)
	return fmt.Sprintf("%x", t.Sum(nil))
}

func Encode(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func StringExist(srv string, dis string) bool {
	return strings.Contains(srv, dis)
}

//返回一个32位md5加密后的字符串
func GetMD5Encode(data string) string {
	h := md5.New()
	h.Write([]byte(data))

	return hex.EncodeToString(h.Sum(nil))
}

//返回一个16位md5加密后的字符串
func Get16MD5Encode(data string) string {
	return GetMD5Encode(data)[8:24]
}

//去除重复字符串
func RemoveRepeatedElement(arr []string) (newArr []string) {
	newArr = make([]string, 0)
	for i := 0; i < len(arr); i++ {
		repeat := false
		for j := i + 1; j < len(arr); j++ {
			if arr[i] == arr[j] {
				repeat = true
				break
			}
		}
		if !repeat {
			newArr = append(newArr, arr[i])
		}
	}
	return newArr
}

// 已知list中元素"amber","jack"

func slice() {
	arr1 := [...]string{"a1", "b1", "c1", "d1", "m222"}
	arr2 := [...]string{"a1", "c1"}
	// 初始化map

	set := make(map[string]struct{})
	set2 := make(map[string]struct{})
	// 上面2部可替换为set := make(map[string]struct{})

	// 将list内容传递进map,只根据key判断，所以不需要关心value的值，用struct{}{}表示
	for _, value := range arr1 {
		set[value] = struct{}{}
	}

	for _, value := range arr2 {
		set2[value] = struct{}{}
	}

	for _, v := range arr1 {
		// 检查元素是否在map
		if _, ok := set2[v]; ok {
			fmt.Println(v, " is in the list")
		} else {
			fmt.Println(v, " is not in the list")

		}
	}

}
