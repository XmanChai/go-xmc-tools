package main

import (
	"fmt"
	xmcQrcode "github.com/xmanchai/go-xmc-tools/coding/qrcode"
	"image/color"
	"io/ioutil"
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

	//fmt.Println(readyDrawSVG())
	//readyDrawSVG()
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
	nwl.Size = 1024
	nwl.MimeType = xmcQrcode.MimeSVG
	s, err := nwl.New()
	if err != nil {
		log.Fatalf("the error is %v", err)
	}
	w.Header().Set("Content-Type", "image/svg+xml")
	end := time.Since(start)
	fs, err := os.Create("post.svg")
	_, err = fs.Write(s)
	if err != nil {
		log.Fatalf("the error is %v", err)
	}
	fmt.Printf("Total %d bytes writed in %v\nUse WithLogo struct only set Content\n", len(s), end)
	_, _ = w.Write(s)
}

func testWithLogo() {
	//var lfs image.Image
	start := time.Now()
	nwl := xmcQrcode.NewWithLogoStruct()
	// 更换颜色，在这里操作，默认是白底黑字
	nwl.Background = color.RGBA{R: 255, G: 255, A: 255}
	nwl.Foreground = color.RGBA{R: 0, G: 100, B: 255, A: 255}
	nwl.MimeType = xmcQrcode.MimePNG
	nwl.Size = 512
	nwl.Content = "https://tse4-mm.cn.bing.net/th/id/OIP-C.-lbqdC5gjVz0pgzk-mzvpQHaD8?w=311&h=180&c=7&r=0&o=5&pid=1.7"
	lf, err := os.Open("logo.png")
	if err != nil {
		log.Fatalln(err)
	}
	defer func(lf *os.File) {
		err := lf.Close()
		if err != nil {

		}
	}(lf)
	all, _ := ioutil.ReadAll(lf)

	nwl.LogoStream = all
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
	//fmt.Println(nwl.String(s))
	end := time.Since(start)
	fmt.Printf("Total %d bytes writed in %v\nUse WithLogo struct only set Content\n", w, end)
}

func testWithColor() {
	start := time.Now()
	nwc := xmcQrcode.NewWithColorStruct()
	// 更换颜色，在这里操作，默认是白底黑字
	nwc.Background = color.RGBA{R: 255, G: 0, A: 150}
	nwc.Foreground = color.RGBA{R: 25, G: 100, B: 255, A: 0}
	nwc.MimeType = xmcQrcode.MimeSVG
	//nwc.SVGPixel = xmcQrcode.SVGCircle
	s, err := nwc.New()
	if err != nil {
		log.Fatalf("the error is %v", err)
	}
	fs, err := os.Create("color.svg")
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
	nd.Content = "https://tse4-mm.cn.bing.net/th/id/OIP-C.-lbqdC5gjVz0pgzk-mzvpQHaD8?w=311&h=180&c=7&r=0&o=5&pid=1.7"
	nd.MimeType = xmcQrcode.MimeSVG
	// 如果想更改默认值，请直接修改nd
	//nd.Size = 512
	s, err := nd.New()
	if err != nil {
		log.Fatalf("the error is %v", err)
	}
	nd.Write(s, "default")
	end := time.Since(start)
	fmt.Printf("Total %d bytes writed in %v\nUse Default struct only set Content\n", len(s), end)
}
