package main

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/color"
	"image/png"
	"syscall/js"
)

func main() {
	// グローバルな関数を登録する
	js.Global().Set("preprocessImage", js.FuncOf(preprocessImage))
	js.Global().Set("binarizeImage", js.FuncOf(binarizeImage))

	// このプログラムが終了しないようにチャネルで待機する
	<-make(chan bool)
}

// 画像の前処理を行う関数
func preprocessImage(this js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return "No image data provided"
	}

	// Base64エンコードされた画像データを取得
	base64Data := args[0].String()

	// Base64デコード
	reader := base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(base64Data))

	// 画像を読み込む
	img, _, err := image.Decode(reader)
	if err != nil {
		return "Failed to decode image: " + err.Error()
	}

	// グレースケール変換
	grayImg := toGrayscale(img)

	// 結果をBase64にエンコード
	var buf bytes.Buffer
	err = png.Encode(&buf, grayImg)
	if err != nil {
		return "Failed to encode processed image: " + err.Error()
	}

	// Base64に変換して返す
	return base64.StdEncoding.EncodeToString(buf.Bytes())
}

// 二値化処理を行う関数
func binarizeImage(this js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return "No image data provided"
	}

	threshold := 128 // デフォルトの閾値
	if len(args) > 1 {
		threshold = args[1].Int()
	}

	// Base64エンコードされた画像データを取得
	base64Data := args[0].String()

	// Base64デコード
	reader := base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(base64Data))

	// 画像を読み込む
	img, _, err := image.Decode(reader)
	if err != nil {
		return "Failed to decode image: " + err.Error()
	}

	// グレースケール変換
	grayImg := toGrayscale(img)

	// 二値化処理
	binaryImg := binarize(grayImg, uint8(threshold))

	// 結果をBase64にエンコード
	var buf bytes.Buffer
	err = png.Encode(&buf, binaryImg)
	if err != nil {
		return "Failed to encode processed image: " + err.Error()
	}

	// Base64に変換して返す
	return base64.StdEncoding.EncodeToString(buf.Bytes())
}

// カラー画像をグレースケールに変換する関数
func toGrayscale(img image.Image) *image.Gray {
	bounds := img.Bounds()
	grayImg := image.NewGray(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			grayImg.Set(x, y, img.At(x, y))
		}
	}

	return grayImg
}

// グレースケール画像を二値化する関数
func binarize(img *image.Gray, threshold uint8) *image.Gray {
	bounds := img.Bounds()
	result := image.NewGray(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if img.GrayAt(x, y).Y > threshold {
				result.SetGray(x, y, color.Gray{Y: 255}) // 白
			} else {
				result.SetGray(x, y, color.Gray{Y: 0}) // 黒
			}
		}
	}

	return result
}
