package main

import (
	"filestore_server/handler"
	"fmt"
	"net/http"
)

func main() {
	//静态资源处理
	//http.Handle("/static/",http.StripPrefix("/static/",http.FileServer(assets.AssetFS())))
	http.Handle("/static/",
		http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// 文件存取接口
	http.HandleFunc("/file/upload/", handler.UploadHandler)
	http.HandleFunc("/file/upload/suc", handler.UploadSucHandler)
	http.HandleFunc("/file/meta", handler.GetFileMetaHandler)
	http.HandleFunc("/file/query", handler.FileQueryHandler)
	http.HandleFunc("/file/download", handler.DownloadHandler)
	http.HandleFunc("/file/update", handler.FileMetaUpdateHandler)
	http.HandleFunc("/file/delete", handler.FileDeleteHandler)

	//秒传接口
	http.HandleFunc("/file/fastupload", handler.HTTPInterceptor(handler.TryFastUploadHandler))

	http.HandleFunc("/user/signup", handler.SignupHandler)
	http.HandleFunc("/user/signin", handler.SignInHandler)
	http.HandleFunc("/user/info", handler.HTTPInterceptor(handler.UserInfoHandler))

	fmt.Println("FileStore server is starting……")

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Printf("Failed to start server,err:%s", err.Error())
	}
	fmt.Println("FileStore server is started and running……")
}
