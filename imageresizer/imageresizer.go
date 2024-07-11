package imageresizer

import (
	"fmt"
	"html/template"
	"image"
	"image/jpeg"
	_ "image/png"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/nfnt/resize"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.ServeFile(w, r, "imageresizer/static/imageresizer.html")
		return
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Failed to upload image", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		http.Error(w, "Failed to decode image", http.StatusInternalServerError)
		return
	}

	width, err := strconv.Atoi(r.FormValue("width"))
	if err != nil {
		http.Error(w, "Invalid width", http.StatusBadRequest)
		return
	}

	height, err := strconv.Atoi(r.FormValue("height"))
	if err != nil {
		http.Error(w, "Invalid height", http.StatusBadRequest)
		return
	}

	resizedImg := resize.Resize(uint(width), uint(height), img, resize.Lanczos3)

	fileName := fmt.Sprintf("resized_%d.jpg", time.Now().UnixNano())
	filePath := filepath.Join("imageresizer/static/uploads", fileName)
	os.MkdirAll(filepath.Dir(filePath), os.ModePerm)

	outFile, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Failed to create file", http.StatusInternalServerError)
		return
	}
	defer outFile.Close()

	err = jpeg.Encode(outFile, resizedImg, nil)
	if err != nil {
		http.Error(w, "Failed to encode image", http.StatusInternalServerError)
		return
	}

	data := map[string]string{
		"FilePath": "/imageresizer/static/uploads/" + fileName,
	}

	tmpl, err := template.ParseFiles("imageresizer/static/download.html")
	if err != nil {
		http.Error(w, "Failed to load template", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Path[len("/imageresizer/download/"):]
	filePath := filepath.Join("imageresizer/static/uploads", fileName)

	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "Failed to open file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Content-Type", "image/jpeg")

	http.ServeContent(w, r, fileName, time.Now(), file)
}
