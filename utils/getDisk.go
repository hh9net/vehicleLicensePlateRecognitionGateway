package utils

/*
import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)

func GetDisk() string{
	dir, _ := os.Getwd()
	log.Println("当前路径的磁盘：", dir)
	drive:=	strings.Split(dir,"\\")
	kernel32, err := syscall.LoadLibrary("Kernel32.dll")
	if err != nil {
		log.Panic(err)
	}
	defer syscall.FreeLibrary(kernel32)
	GetDiskFreeSpaceEx, err := syscall.GetProcAddress(syscall.Handle(kernel32), "GetDiskFreeSpaceExW")

	if err != nil {
		log.Panic(err)
	}

	lpFreeBytesAvailable := int64(0)
	lpTotalNumberOfBytes := int64(0)
	lpTotalNumberOfFreeBytes := int64(0)
	r, a, b := syscall.Syscall6(uintptr(GetDiskFreeSpaceEx), 4,

		//uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("D:"))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(drive[0]))),
		uintptr(unsafe.Pointer(&lpFreeBytesAvailable)),
		uintptr(unsafe.Pointer(&lpTotalNumberOfBytes)),
		uintptr(unsafe.Pointer(&lpTotalNumberOfFreeBytes)), 0, 0)

	log.Printf("Available  %dmb", lpFreeBytesAvailable/1024/1024.0)
	log.Printf("Total      %dmb", lpTotalNumberOfBytes/1024/1024.0)
	log.Printf("Free       %dmb", lpTotalNumberOfFreeBytes/1024/1024.0)

	v := float64(lpTotalNumberOfFreeBytes/1024/1024.0) / float64((lpTotalNumberOfBytes / 1024 / 1024.0))

	log.Print("还可以使用比:", v)

	v, _ = strconv.ParseFloat(fmt.Sprintf("%.4f", v), 64)
	//s1 := strconv.FormatFloat(v, 'E', -2, 64)//float32s2 := strconv.FormatFloat(v, 'E', -1, 64)//float64

	log.Println("还可以使用比:", v*100,"%")
	log.Println(r, a, b)

	return  "还可以使用磁盘比:"+fmt.Sprintf("%v", v*100)+"%"
}

*/
