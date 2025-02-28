.PHONY: all clean build serve setup deps

all: deps setup build serve

# 依存関係のインストール
deps:
	@echo "依存関係をインストールしています..."
	@go mod init wasm-ocr-example || true
	@echo "依存関係のインストールが完了しました。"

# セットアップ：必要なファイルを確認・作成
setup:
	@echo "セットアップを開始します..."
	@mkdir -p build
	@mkdir -p cmd/server
	@mkdir -p cmd/wasm
	@which go > /dev/null || (echo "Goがインストールされていません。インストールしてください。" && exit 1)
	@cp "$(shell go env GOROOT)/lib/wasm/wasm_exec.js" build/
	@cp index.html build/
	@cp index.js build/
	@echo "セットアップが完了しました。"

# ビルド：GoコードをWASMにコンパイル
build:
	@echo "WASMファイルをコンパイルしています..."
	GOOS=js GOARCH=wasm go build -o build/main.wasm ./cmd/wasm/main.go
	@echo "コンパイルが完了しました。"

# サーブ：ローカルサーバーを起動
serve:
	@echo "ローカルサーバーを起動しています..."
	@go run cmd/server/main.go cmd/server/server.go
	@echo "http://localhost:8080 にアクセスしてください"

# クリーン：生成されたファイルを削除
clean:
	@echo "生成されたファイルを削除しています..."
	rm -rf build
	@echo "クリーンアップが完了しました。"
