# Go + WebAssembly OCRデモ

ブラウザ上でGoとWebAssemblyを使って画像処理を行い、Tesseract.jsで文字認識するデモアプリケーションです。

## 機能

- 画像のアップロード
- Go/WebAssemblyによる画像の前処理（グレースケール化、二値化）
- 二値化閾値の動的調整
- Tesseract.jsによる日本語/英語OCR
- すべての処理がブラウザ上で完結（サーバーへのアップロード不要）

## インストールと実行

### 前提条件

- Goがインストールされていること

### セットアップ

1. Makefileを使って環境をセットアップ

```bash
make
```

このコマンド一つで以下の処理が自動的に行われます：
- 依存関係のインストール
- 必要なディレクトリの作成
- Go WebAssemblyのコンパイル
- ローカルHTTPサーバーの起動

2. ブラウザでアクセス

`http://localhost:8080` にアクセスするとアプリケーションが表示されます。

## 使い方

1. 「ファイルを選択」ボタンで画像ファイルをアップロード
2. 「前処理」ボタンをクリックして画像処理を実行
3. 必要に応じてスライダーで二値化閾値を調整
4. 「文字認識」ボタンをクリックしてOCR処理を実行
5. 認識結果が下部に表示されます

## プロジェクト構成

```
.
├── Makefile
├── README.md
├── build
│   ├── index.html
│   ├── main.wasm
│   └── wasm_exec.js
├── cmd
│   ├── server        # サーバー起動用のエントリーポイント
│   │   ├── main.go
│   │   └── server.go
│   └── wasm          # WASMにコンパイルされるGoコード（画像処理機能）
│       └── main.go
├── go.mod
└── index.html
```

## 技術解説

### Go WebAssembly

GoコードをWebAssemblyにコンパイルして、ブラウザ上で実行しています。これにより、クライアントサイドでありながらGoの高速な画像処理を利用できます。

コンパイルは `GOOS=js GOARCH=wasm go build` コマンドで行われ、`wasm_exec.js` ファイルがGoとJavaScriptの橋渡しを担当します。

### 画像処理

画像処理は以下の手順で行われます：

1. Base64エンコードされた画像データをGoコードに渡す
2. 画像をグレースケールに変換
3. 指定された閾値で二値化処理
4. 処理後の画像をBase64で返す

### OCR処理

OCR処理はTesseract.jsライブラリを使用しています：

1. 前処理された画像をTesseract.jsに渡す
2. 日本語と英語の言語モデルを使用して文字認識
3. 認識結果をブラウザに表示

## カスタマイズ

- `wasm/main.go` の画像処理アルゴリズムを変更することで、異なる前処理を実装できます
- `index.html` のUIや処理フローをカスタマイズできます
- 言語設定を変更することで、他の言語のOCRも可能です

## 参考資料

- [Go WebAssembly公式ドキュメント](https://golang.org/doc/go1.11#wasm)
- [Tesseract.jsドキュメント](https://github.com/naptha/tesseract.js)
