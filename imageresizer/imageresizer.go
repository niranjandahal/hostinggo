

package imageresizer


import (


	"bytes"
	"encoding/base64"
	"html/template"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"strconv"


	"github.com/nfnt/resize"

	
)

type ImageData struct {
	Base64Image string
}

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

	img, format, err := image.Decode(file)
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

	var buf bytes.Buffer
	switch format {
	case "jpeg":
		err = jpeg.Encode(&buf, resizedImg, nil)
	case "png":
		err = png.Encode(&buf, resizedImg)
	default:
		http.Error(w, "Unsupported image format", http.StatusUnsupportedMediaType)
		return
	}

	if err != nil {
		http.Error(w, "Failed to encode image", http.StatusInternalServerError)
		return
	}

	base64Image := base64.StdEncoding.EncodeToString(buf.Bytes())
	data := ImageData{
		Base64Image: base64Image,
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