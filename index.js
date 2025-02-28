// WebAssemblyのロード
const go = new Go();
let wasmLoaded = false;

// ページの読み込みが完了したらWASMを初期化
document.addEventListener('DOMContentLoaded', async function() {
  try {
    // WASMモジュールを読み込む
    const result = await WebAssembly.instantiateStreaming(
      fetch("main.wasm"),
      go.importObject
    );
    go.run(result.instance);
    wasmLoaded = true;
    console.log("WebAssembly module loaded");
  } catch (err) {
    console.error("Failed to load WebAssembly module:", err);
  }

  // Tesseractが正しく読み込まれたか確認
  if (typeof Tesseract !== 'undefined') {
    console.log("Tesseract.js loaded successfully");
  } else {
    console.error("Tesseract.js not loaded, trying to load again");
    // スクリプトタグで直接読み込む試行
    const script = document.createElement('script');
    script.src = "https://cdn.jsdelivr.net/npm/tesseract.js@v5.0.3/dist/tesseract.min.js";
    document.head.appendChild(script);
  }
});

// DOM要素の取得
const imageInput = document.getElementById('imageInput');
const originalCanvas = document.getElementById('originalCanvas');
const processedCanvas = document.getElementById('processedCanvas');
const preprocessBtn = document.getElementById('preprocessBtn');
const recognizeBtn = document.getElementById('recognizeBtn');
const thresholdSlider = document.getElementById('thresholdSlider');
const thresholdValue = document.getElementById('thresholdValue');
const recognizedText = document.getElementById('recognizedText');
const loadingMessage = document.getElementById('loadingMessage');

// カメラ機能のための要素
const cameraBtn = document.getElementById('cameraBtn');
const cameraContainer = document.getElementById('cameraContainer');
const video = document.getElementById('video');
const captureBtn = document.getElementById('captureBtn');
const closeBtn = document.getElementById('closeBtn');

// カメラストリーム
let mediaStream = null;

// キャンバスコンテキスト
const originalCtx = originalCanvas.getContext('2d');
const processedCtx = processedCanvas.getContext('2d');

// 画像データ
let originalImage = null;
let processedImageData = null;

// ファイル選択時の処理
imageInput.addEventListener('change', function(e) {
  if (e.target.files.length === 0) return;

  const file = e.target.files[0];
  const reader = new FileReader();

  reader.onload = function(event) {
    const img = new Image();
    img.onload = function() {
      // キャンバスサイズを設定
      originalCanvas.width = img.width;
      originalCanvas.height = img.height;
      processedCanvas.width = img.width;
      processedCanvas.height = img.height;

      // 元画像を描画
      originalCtx.drawImage(img, 0, 0);
      originalImage = img;

      // 処理後キャンバスをクリア
      processedCtx.clearRect(0, 0, processedCanvas.width, processedCanvas.height);

      // 認識結果をクリア
      recognizedText.textContent = '';
    };
    img.src = event.target.result;
  };

  reader.readAsDataURL(file);
});

// 前処理ボタンのクリックイベント
preprocessBtn.addEventListener('click', function() {
  if (!wasmLoaded) {
    alert("WebAssemblyモジュールがまだロードされていません");
    return;
  }

  if (!originalImage) {
    alert("画像を選択してください");
    return;
  }

  try {
    // 元画像をBase64エンコード
    const dataURL = originalCanvas.toDataURL('image/png');
    const base64Data = dataURL.replace(/^data:image\/(png|jpg|jpeg);base64,/, '');

    // WASMで二値化処理
    const threshold = parseInt(thresholdSlider.value);
    const processedData = binarizeImage(base64Data, threshold);

    // 処理後画像を表示
    const processedImg = new Image();
    processedImg.onload = function() {
      processedCtx.drawImage(processedImg, 0, 0);
      processedImageData = processedImg;
    };
    processedImg.src = 'data:image/png;base64,' + processedData;
  } catch (err) {
    console.error("Image processing error:", err);
    alert("画像処理中にエラーが発生しました: " + err.message);
  }
});

// 文字認識ボタンのクリックイベント
recognizeBtn.addEventListener('click', async function() {
  if (!processedImageData) {
    alert("まず画像の前処理を行ってください");
    return;
  }

  // ローディングメッセージを表示
  loadingMessage.style.display = 'block';
  recognizedText.textContent = '';

  try {
    // Tesseractがグローバルに存在するか確認
    if (typeof Tesseract === 'undefined') {
      throw new Error('Tesseractライブラリが読み込まれていません。');
    }

    // 簡易版のOCR実行（Worker APIを避ける）
    console.log("Starting OCR with simplified approach...");

    // 画像データをDataURLとして取得
    const imageData = processedCanvas.toDataURL('image/png');

    try {
      // Tesseract v5の新しいAPI
      const result = await Tesseract.recognize(
        imageData,
        'jpn+eng',
        {
          logger: progress => {
            console.log('OCR進行状況:', progress);
            if (progress.status === 'recognizing text') {
              loadingMessage.textContent = `文字認識処理中... ${(progress.progress * 100).toFixed(1)}%`;
            }
          }
        }
      );

      // 認識結果を表示
      console.log("Recognition completed");
      recognizedText.textContent = result.data.text || "テキストが認識できませんでした";
    } catch (err) {
      console.error("OCR simplified approach failed:", err);

      // フォールバック：画像を直接使用
      try {
        console.log("Trying fallback method...");
        const img = new Image();
        img.src = imageData;

        const result = await Tesseract.recognize(
          img,
          'jpn+eng',
          { logger: m => console.log(m) }
        );

        recognizedText.textContent = result.data.text || "テキストが認識できませんでした";
      } catch (fallbackErr) {
        console.error("Fallback OCR failed:", fallbackErr);
        recognizedText.textContent = "OCR処理に失敗しました: " + fallbackErr.message;
      }
    }
  } catch (error) {
    console.error("OCR処理中にエラーが発生しました:", error);
    recognizedText.textContent = "エラー: " + error.message;
  } finally {
    // ローディングメッセージを非表示
    loadingMessage.style.display = 'none';
  }
});

// 閾値スライダーの変更イベント
thresholdSlider.addEventListener('input', function() {
  thresholdValue.textContent = this.value;
});

// カメラボタンのクリックイベント
cameraBtn.addEventListener('click', async function() {
  try {
    // カメラストリームを取得
    mediaStream = await navigator.mediaDevices.getUserMedia({
      video: { facingMode: 'environment' }, // バックカメラを優先（モバイル端末の場合）
      audio: false
    });

    // ビデオ要素にストリームを設定
    video.srcObject = mediaStream;
    video.play();

    // カメラコンテナを表示
    cameraContainer.style.display = 'block';
  } catch (err) {
    console.error('カメラへのアクセスに失敗しました:', err);
    alert('カメラへのアクセスに失敗しました。カメラの使用を許可するか、別の方法で画像をアップロードしてください。');
  }
});

// キャプチャボタンのクリックイベント
captureBtn.addEventListener('click', function() {
  if (!video.srcObject) {
    alert('カメラが起動していません');
    return;
  }

  // キャンバスサイズをビデオサイズに合わせる
  originalCanvas.width = video.videoWidth;
  originalCanvas.height = video.videoHeight;
  processedCanvas.width = video.videoWidth;
  processedCanvas.height = video.videoHeight;

  // ビデオフレームをキャンバスに描画
  originalCtx.drawImage(video, 0, 0, originalCanvas.width, originalCanvas.height);

  // 画像をキャプチャしたのでoriginalImageを設定
  const img = new Image();
  img.src = originalCanvas.toDataURL('image/png');
  originalImage = img;

  // 処理後キャンバスをクリア
  processedCtx.clearRect(0, 0, processedCanvas.width, processedCanvas.height);

  // 認識結果をクリア
  recognizedText.textContent = '';

  // カメラコンテナを非表示
  cameraContainer.style.display = 'none';

  // カメラストリームを停止
  stopCameraStream();
});

// カメラを閉じるボタンのクリックイベント
closeBtn.addEventListener('click', function() {
  // カメラコンテナを非表示
  cameraContainer.style.display = 'none';

  // カメラストリームを停止
  stopCameraStream();
});

// カメラストリームを停止する関数
function stopCameraStream() {
  if (mediaStream) {
    mediaStream.getTracks().forEach(track => track.stop());
    mediaStream = null;
    video.srcObject = null;
  }
}
