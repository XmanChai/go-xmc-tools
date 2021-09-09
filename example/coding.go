package main

import (
	"fmt"
	xmcQrcode "go-xmc-tools/coding/qrcode"
	"image"
	"image/color"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

func main() {
	testDefault()
	testWithColor()
	testWithLogo()

	server := &http.Server{
		Addr: "127.0.0.1:8810",
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/postdata", UploadFile)
	server.Handler = mux
	log.Fatalln(server.ListenAndServe())

}

func UploadFile(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	logo, _, err := r.FormFile("logo")
	fmt.Printf("type is %T\nvalue is %v", logo, logo)
	if err != nil {
		log.Fatalln(err)
	}
	defer func(logo multipart.File) {
		_ = logo.Close()
	}(logo)
	start := time.Now()
	nwl := xmcQrcode.NewWithLogoStruct()
	// 更换颜色，在这里操作，默认是白底黑字
	nwl.Background = color.RGBA{R: 255, G: 255, A: 255}
	nwl.Foreground = color.RGBA{R: 0, G: 100, B: 255, A: 255}
	nwl.LogoStream = logo
	nwl.Content = "https://github.com/xmanchai/go-xmc-tools/coding/qrcode"
	s, err := nwl.New()
	if err != nil {
		log.Fatalf("the error is %v", err)
	}
	w.Header().Set("Content-Type", "image/png")
	end := time.Since(start)
	fs, err := os.Create("post.png")
	_, err = fs.Write(s)
	if err != nil {
		log.Fatalf("the error is %v", err)
	}
	fmt.Printf("Total %d bytes writed in %v\nUse WithLogo struct only set Content\n", s, end)
	_, _ = w.Write(s)
}

func testWithLogo() {
	start := time.Now()
	nwl := xmcQrcode.NewWithLogoStruct()
	// 更换颜色，在这里操作，默认是白底黑字
	nwl.Background = color.RGBA{R: 255, G: 255, A: 255}
	nwl.Foreground = color.RGBA{R: 0, G: 100, B: 255, A: 255}
	lf, err := os.Open("logo.png")
	lfs, _, err := image.Decode(lf)
	nwl.LogoStream = lfs
	s, err := nwl.New()
	if err != nil {
		log.Fatalf("the error is %v", err)
	}
	fs, err := os.Create("withlogo.png")
	w, err := fs.Write(s)
	if err != nil {
		log.Fatalf("the error is %v", err)
	}

	defer func(fs *os.File) {
		err := fs.Close()
		if err != nil {
			return
		}
	}(fs)
	end := time.Since(start)
	fmt.Printf("Total %d bytes writed in %v\nUse WithLogo struct only set Content\n", w, end)
}

func testWithColor() {
	start := time.Now()
	nwc := xmcQrcode.NewWithColorStruct()
	// 更换颜色，在这里操作，默认是白底黑字
	nwc.Background = color.RGBA{R: 255, G: 255, A: 255}
	nwc.Foreground = color.RGBA{R: 0, G: 100, B: 255, A: 255}
	s, err := nwc.New()
	if err != nil {
		log.Fatalf("the error is %v", err)
	}
	fs, err := os.Create("color.png")
	w, err := fs.Write(s)
	if err != nil {
		log.Fatalf("the error is %v", err)
	}

	defer func(fs *os.File) {
		err := fs.Close()
		if err != nil {
			return
		}
	}(fs)
	end := time.Since(start)
	fmt.Printf("Total %d bytes writed in %v\nUse WithColor struct only set Content\n", w, end)
}

func testDefault() {
	start := time.Now()
	nd := xmcQrcode.NewDefaultStruct()
	// 如果想更改默认值，请直接修改nd
	s, err := nd.New()
	if err != nil {
		log.Fatalf("the error is %v", err)
	}
	fs, err := os.Create("default.png")
	w, err := fs.Write(s)
	if err != nil {
		log.Fatalf("the error is %v", err)
	}

	defer func(fs *os.File) {
		err := fs.Close()
		if err != nil {
			return
		}
	}(fs)
	end := time.Since(start)
	fmt.Printf("Total %d bytes writed in %v\nUse Default struct only set Content\n", w, end)
}
