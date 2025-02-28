package main

import (
	"flag"
)

func main() {
	port := flag.Int("port", 8080, "サーバーのポート番号")
	flag.Parse()

	// サーバーを起動
	startServer(*port)
}
