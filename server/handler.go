package server

import (
	"encoding/json"
	"log"
	"net/http"
)

// Server はHTTPサーバーの構造体
type Server struct {
	ocrService *OCRService
}

// NewServer は新しいサーバーインスタンスを作成します
func NewServer() (*Server, error) {
	ocrService, err := NewOCRService()
	if err != nil {
		return nil, err
	}

	return &Server{
		ocrService: ocrService,
	}, nil
}

// Close はサーバーのリソースを解放します
func (s *Server) Close() error {
	return s.ocrService.Close()
}

// SetupRoutes はHTTPルートを設定します
func (s *Server) SetupRoutes() {
	// 静的ファイルを提供
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	// API エンドポイント
	http.HandleFunc("/api/preprocess", s.handlePreprocess)
	http.HandleFunc("/api/recognize", s.handleRecognize)
}

// ProcessRequest は画像前処理のリクエスト
type ProcessRequest struct {
	ImageData string `json:"imageData"`
	Threshold int    `json:"threshold"`
}

// ProcessResponse は画像前処理のレスポンス
type ProcessResponse struct {
	ProcessedImage string `json:"processedImage"`
	Error          string `json:"error,omitempty"`
}

// RecognizeRequest は文字認識のリクエスト
type RecognizeRequest struct {
	ImageData string `json:"imageData"`
}

// RecognizeResponse は文字認識のレスポンス
type RecognizeResponse struct {
	Text  string `json:"text"`
	Error string `json:"error,omitempty"`
}

// 画像前処理のハンドラ
func (s *Server) handlePreprocess(w http.ResponseWriter, r *http.Request) {
	// POSTメソッドのみ受け付ける
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// リクエストをデコード
	var req ProcessRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSONResponse(w, ProcessResponse{Error: "Invalid request: " + err.Error()}, http.StatusBadRequest)
		return
	}

	// 画像処理を実行
	processedImage, err := s.ocrService.ProcessImage(req.ImageData, req.Threshold)
	if err != nil {
		log.Printf("画像処理エラー: %v", err)
		sendJSONResponse(w, ProcessResponse{Error: "Processing error: " + err.Error()}, http.StatusInternalServerError)
		return
	}

	// 成功レスポンスを返す
	sendJSONResponse(w, ProcessResponse{ProcessedImage: processedImage}, http.StatusOK)
}

// 文字認識のハンドラ
func (s *Server) handleRecognize(w http.ResponseWriter, r *http.Request) {
	// POSTメソッドのみ受け付ける
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// リクエストをデコード
	var req RecognizeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSONResponse(w, RecognizeResponse{Error: "Invalid request: " + err.Error()}, http.StatusBadRequest)
		return
	}

	// 文字認識を実行
	text, err := s.ocrService.RecognizeText(req.ImageData)
	if err != nil {
		log.Printf("文字認識エラー: %v", err)
		sendJSONResponse(w, RecognizeResponse{Error: "Recognition error: " + err.Error()}, http.StatusInternalServerError)
		return
	}

	// 成功レスポンスを返す
	sendJSONResponse(w, RecognizeResponse{Text: text}, http.StatusOK)
}

// JSONレスポンスを送信するヘルパー関数
func sendJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("レスポンスの書き込みエラー: %v", err)
	}
}
