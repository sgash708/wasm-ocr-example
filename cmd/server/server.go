package main

import (
	"log"
	"net/http"
	"strconv"
)

// startServer は指定されたポートでHTTPサーバーを起動します
func startServer(port int) {
	// 静的ファイルを提供するハンドラーを設定
	fs := http.FileServer(http.Dir("./build"))
	http.Handle("/", fs)

	// サーバーを起動
	portStr := strconv.Itoa(port)
	log.Printf("HTTPサーバーを起動しています: http://localhost:%s", portStr)
	if err := http.ListenAndServe(":"+portStr, nil); err != nil {
		log.Fatal("サーバーエラー: ", err)
	}
}
