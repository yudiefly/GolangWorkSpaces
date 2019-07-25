package main

import (
	//"fmt"
	"log"

	"github.com/fsnotify/fsnotify"
)

func main() {
	//创建一个监控对象
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()
	//添加要监控的对象，文件或文件夹
	err = watcher.Add("./tmp")
	if err != nil {
		log.Fatal(err)
	}

	//我们另启一个goroutine来处理监控对象的事件
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				//判断事件发生的类型，如下5种
				// Create 创建
				// Write 写入
				// Remove 删除
				// Rename 重命名
				// Chmod 修改权限
				log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("写入文件 : ", event.Name)
				}
				if event.Op&fsnotify.Create == fsnotify.Create {
					log.Println("创建文件 : ", event.Name)
				}
				if event.Op&fsnotify.Remove == fsnotify.Remove {
					log.Println("删除文件 : ", event.Name)
				}
				if event.Op&fsnotify.Rename == fsnotify.Rename {
					log.Println("重命名文件 : ", event.Name)
				}
				if event.Op&fsnotify.Chmod == fsnotify.Chmod {
					log.Println("修改权限 : ", event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	//循环（阻塞主进程）
	select {}

}
