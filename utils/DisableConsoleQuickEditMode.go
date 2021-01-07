package utils

// https://blog.csdn.net/qq_34062754/article/details/107084616
func IntPtr(n int) uintptr {
	return uintptr(n)
}
func UIntPtr(n uint) uintptr {
	return uintptr(n)
}

//
//func DisableConsoleQuickEditMode() (result bool) {
//	var kernel32DLL, err = syscall.LoadLibrary("kernel32.dll")
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	var funcStdHandle uintptr
//	funcStdHandle, err = syscall.GetProcAddress(kernel32DLL, "GetStdHandle")
//	if funcStdHandle == uintptr(0) {
//		fmt.Println(err, funcStdHandle)
//		return
//	}
//
//	var funcGetConsoleMode uintptr
//	funcGetConsoleMode, err = syscall.GetProcAddress(kernel32DLL, "GetConsoleMode")
//	if funcGetConsoleMode == uintptr(0) {
//		fmt.Println(err, funcGetConsoleMode)
//		return
//	}
//
//	var funcSetConsoleMode uintptr
//	funcSetConsoleMode, err = syscall.GetProcAddress(kernel32DLL, "SetConsoleMode")
//	fmt.Println(err, funcSetConsoleMode)
//	if funcSetConsoleMode == uintptr(0) {
//		return
//	}
//
//	var stdHandle uintptr
//	stdHandle, _, err = syscall.Syscall9(funcStdHandle, 1, IntPtr(syscall.STD_INPUT_HANDLE),
//		0, 0, 0, 0, 0, 0, 0, 0)
//	if stdHandle == uintptr(0) {
//		fmt.Println(stdHandle, err)
//		return
//	}
//	var mode uint
//	var bret uintptr
//	bret, _, err = syscall.Syscall9(funcGetConsoleMode, 2, stdHandle, uintptr(unsafe.Pointer(&mode)),
//		0, 0, 0, 0, 0, 0, 0)
//	if bret != uintptr(1) {
//		return
//	}
//	// ENABLE_QUICK_EDIT_MODE 0x0040
//	// ENABLE_INSERT_MODE 0x0020
//	// ENABLE_MOUSE_INPUT 0x0010
//	mode &= ^uint(0x0040) //移除快速编辑模式
//	mode &= ^uint(0x0020) //移除插入模式
//	mode &= ^uint(0x0010)
//
//	var ret uintptr
//	ret, _, err = syscall.Syscall9(funcSetConsoleMode, 2, stdHandle, UIntPtr(mode),
//		0, 0, 0, 0, 0, 0, 0)
//	if ret == uintptr(1) {
//		result = true
//	}
//	return
//}
