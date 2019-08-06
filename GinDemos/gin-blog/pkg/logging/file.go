package logging

import (
	"fmt"
	"log"
	"os"
	"time"
)

var (
	LogSavePath = "runtime/logs/"
	LogSaveName = "log"
	LogFileExt  = "log"
	TimeFormat  = "20190806"
)

func getLogFilePath() string {
	return fmt.Sprintf("%s", LogSavePath)
}

func getLogFileFullPath() string {
	prefixPath := getLogFilePath()
	suffixPath := fmt.Sprintf("%s%s.%s", LogSaveName, time.Now().Format(TimeFormat), LogFileExt)

	return fmt.Sprintf("%s%s", prefixPath, suffixPath)
}

func openLogFile(filePath string) *os.File {
	_, err := os.Stat(filePath)
	switch {
	case os.IsNotExist(err):
		mkDir(getLogFilePath())
	case os.IsPermission(err):
		log.Fatal("Permission:%v", err)
	}
	handle, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Fail to OpenFile:%v", err)
	}
	return handle
}

func mkDir(filePath string) {
	dir, _ := os.Getwd()
	err := os.Mkdir(dir+"/"+getLogFilePath(), os.ModePerm)
	if err != nil {
		panic(err)
	}
}

/*
相关知识点备注:
【os.Stat】：返回文件信息结构描述文件。如果出现错误，会返回*PathError

	type PathError struct {
	    Op   string
	    Path string
	    Err  error
	}

【os.IsNotExist】：能够接受ErrNotExist、syscall的一些错误，它会返回一个布尔值，能够得知文件不存在或目录不存在
【os.IsPermission】：能够接受ErrPermission、syscall的一些错误，它会返回一个布尔值，能够得知权限是否满足
【os.OpenFile】：调用文件，支持传入文件名称、指定的模式调用文件、文件权限，返回的文件的方法可以用于I/O。如果出现错误，则为*PathError。

	const (
	    // Exactly one of O_RDONLY, O_WRONLY, or O_RDWR must be specified.
	    O_RDONLY int = syscall.O_RDONLY // 以只读模式打开文件
	    O_WRONLY int = syscall.O_WRONLY // 以只写模式打开文件
	    O_RDWR   int = syscall.O_RDWR   // 以读写模式打开文件
	    // The remaining values may be or'ed in to control behavior.
	    O_APPEND int = syscall.O_APPEND // 在写入时将数据追加到文件中
	    O_CREATE int = syscall.O_CREAT  // 如果不存在，则创建一个新文件
	    O_EXCL   int = syscall.O_EXCL   // 使用O_CREATE时，文件必须不存在
	    O_SYNC   int = syscall.O_SYNC   // 同步IO
	    O_TRUNC  int = syscall.O_TRUNC  // 如果可以，打开时
	)

【os.Getwd】：返回与当前目录对应的根路径名
【os.MkdirAll】：创建对应的目录以及所需的子目录，若成功则返回nil，否则返回error
【os.ModePerm】：const定义ModePerm FileMode = 0777
*/
