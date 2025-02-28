package server

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"strings"

	"github.com/otiai10/gosseract/v2"
)

// OCRService はOCR処理のためのサービス構造体
type OCRService struct {
	client *gosseract.Client
}

// NewOCRService は新しいOCRServiceを作成します
func NewOCRService() (*OCRService, error) {
	client := gosseract.NewClient()
	err := client.SetLanguage("jpn", "eng")
	if err != nil {
		return nil, fmt.Errorf("言語設定エラー: %w", err)
	}

	return &OCRService{
		client: client,
	}, nil
}

// Close はOCRServiceのリソースを解放します
func (s *OCRService) Close() error {
	return s.client.Close()
}

// ProcessImage は画像を前処理して二値化します
func (s *OCRService) ProcessImage(base64Image string, threshold int) (string, error) {
	// Base64から画像データを取得
	imgData, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(base64Image, "data:image/png;base64,"))
	if err != nil {
		imgData, err = base64.StdEncoding.DecodeString(strings.TrimPrefix(base64Image, "data:image/jpeg;base64,"))
		if err != nil {
			return "", fmt.Errorf("Base64デコードエラー: %w", err)
		}
	}

	// 画像をデコード
	img, format, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		return "", fmt.Errorf("画像デコードエラー: %w", err)
	}

	// グレースケール変換
	grayImg := toGrayscale(img)

	// 二値化処理
	binaryImg := binarize(grayImg, uint8(threshold))

	// 処理後の画像をエンコード
	var buf bytes.Buffer
	if format == "png" {
		err = png.Encode(&buf, binaryImg)
	} else {
		err = jpeg.Encode(&buf, binaryImg, &jpeg.Options{Quality: 90})
	}
	if err != nil {
		return "", fmt.Errorf("画像エンコードエラー: %w", err)
	}

	// Base64に変換して返す
	return fmt.Sprintf("data:image/%s;base64,%s", format, base64.StdEncoding.EncodeToString(buf.Bytes())), nil
}

// RecognizeText は画像からテキストを認識します
func (s *OCRService) RecognizeText(base64Image string) (string, error) {
	// Base64から画像データを取得
	var imgData []byte
	var err error

	if strings.HasPrefix(base64Image, "data:image/") {
		// data:image/xxx;base64, プレフィックスを削除
		parts := strings.SplitN(base64Image, ",", 2)
		if len(parts) != 2 {
			return "", fmt.Errorf("不正なBase64形式")
		}
		imgData, err = base64.StdEncoding.DecodeString(parts[1])
	} else {
		// プレフィックスのない通常のBase64
		imgData, err = base64.StdEncoding.DecodeString(base64Image)
	}

	if err != nil {
		return "", fmt.Errorf("Base64デコードエラー: %w", err)
	}

	// 一時ファイルに画像を保存せずに直接バイトデータを設定
	err = s.client.SetImageFromBytes(imgData)
	if err != nil {
		return "", fmt.Errorf("画像設定エラー: %w", err)
	}

	// テキスト認識を実行
	text, err := s.client.Text()
	if err != nil {
		return "", fmt.Errorf("テキスト認識エラー: %w", err)
	}

	return text, nil
}

// 画像をグレースケールに変換
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

// グレースケール画像を二値化
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
