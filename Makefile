.PHONY: all clean build serve setup deps

# デフォルトターゲット：すべての処理を実行
all: deps setup build serve

# 依存関係のインストール
deps:
	@echo "依存関係をインストールしています..."
	@go mod tidy
	@echo "依存関係のインストールが完了しました。"

# セットアップ：必要なファイルを確認・作成
setup:
	@echo "セットアップを開始します..."
	@mkdir -p static
	@which go > /dev/null || (echo "Goがインストールされていません。インストールしてください。" && exit 1)
	@which tesseract > /dev/null || (echo "Tesseractがインストールされていません。インストールしてください。" && exit 1)
	@echo "セットアップが完了しました。"

# ビルド：Goサーバーをビルド
build:
	@echo "サーバーをビルドしています..."
	go build -o bin/ocr-server cmd/server/main.go
	@echo "ビルドが完了しました。"

# サーブ：サーバーを起動
serve:
	@echo "サーバーを起動しています..."
	@echo "http://localhost:8080 にアクセスしてください"
	@./bin/ocr-server

# クリーン：生成されたファイルを削除
clean:
	@echo "生成されたファイルを削除しています..."
	rm -rf bin
	@echo "クリーンアップが完了しました。"
