package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"wasm-ocr-example/server"
)

func main() {
	// コマンドライン引数からポートを取得
	port := flag.Int("port", 8080, "サーバーのポート番号")
	flag.Parse()

	// サーバーインスタンスを作成
	srv, err := server.NewServer()
	if err != nil {
		log.Fatalf("サーバーの初期化エラー: %v", err)
	}
	defer srv.Close()

	// ルートを設定
	srv.SetupRoutes()

	// 終了シグナルを受け取るためのチャネル
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// サーバーの起動
	portStr := strconv.Itoa(*port)
	log.Printf("サーバーを起動しています: http://localhost:%s", portStr)

	// サーバーをゴルーチンで起動
	go func() {
		if err := http.ListenAndServe(":"+portStr, nil); err != nil && err != http.ErrServerClosed {
			log.Fatalf("サーバーエラー: %v", err)
		}
	}()

	// 終了シグナルを待機
	<-quit
	log.Println("サーバーを終了しています...")
}
