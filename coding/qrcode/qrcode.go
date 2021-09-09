package xmc_qrcode

import (
	"bytes"
	"fmt"
	"github.com/nfnt/resize"
	qr "github.com/skip2/go-qrcode"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"log"
	"reflect"
)

type Qrcode interface {
	New() ([]byte, error)
}

type Default struct {
	Content      string
	RecoverLevel qr.RecoveryLevel
	Size         int
}

func (d Default) New() ([]byte, error) {
	return qr.Encode(d.Content, d.RecoverLevel, d.Size)
}

func NewDefaultStruct() Default {
	return Default{
		Content:      "Default struct",
		RecoverLevel: qr.Medium,
		Size:         256,
	}
}

type WithColor struct {
	Content      string
	RecoverLevel qr.RecoveryLevel
	Size         int
	Foreground   color.Color
	Background   color.Color
}

func NewWithColorStruct() WithColor {
	return WithColor{
		Content:    "With color",
		Size:       128,
		Foreground: color.Black,
		Background: color.White,
	}
}

func (wc WithColor) New() ([]byte, error) {
	c, err := qr.New(wc.Content, wc.RecoverLevel)
	c.ForegroundColor = wc.Foreground
	c.BackgroundColor = wc.Background
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	return c.PNG(wc.Size)
}

type WithLogo struct {
	Content      string
	RecoverLevel qr.RecoveryLevel
	Size         int
	LogoStream   interface{}
	Foreground   color.Color
	Background   color.Color
}

func NewWithLogoStruct() WithLogo {
	return WithLogo{
		Content:      "With logo",
		Size:         256,
		RecoverLevel: qr.Medium,
		Foreground:   color.Black,
		Background:   color.White,
		LogoStream:   image.Image(image.Rect(0, 0, 10, 10)),
	}
}

func (wl WithLogo) New() ([]byte, error) {
	var src image.Image
	qc, err := qr.New(wl.Content, wl.RecoverLevel)
	qc.BackgroundColor = wl.Background
	qc.ForegroundColor = wl.Foreground
	//psize := len(wl.Content)
	if err != nil {
		return nil, err
	}
	qcs, err := qc.PNG(wl.Size)
	if err != nil {
		return nil, err
	}
	dst, err := _convertByteToDrawImage(qcs)
	if err != nil {
		return nil, err
	}
	fmt.Println(reflect.TypeOf(wl.LogoStream).String())
	switch reflect.TypeOf(wl.LogoStream).String() {
	case "*image.NRGBA":
		fmt.Println("this type is image.Image")
		src = wl.LogoStream.(image.Image)
	case "image.Image":
		fmt.Println("this type is image.Image")
		src = wl.LogoStream.(image.Image)
	case "[]byte":
		fmt.Println("this type is Byte[]")
		src, _ = _convertByteToImage(wl.LogoStream.([]byte))
	case "multipart.sectionReadCloser":
		b := bytes.NewBuffer(nil)
		_, _ = io.Copy(b, wl.LogoStream.(io.Reader))
		ori, _ := _convertByteToImage(b.Bytes())
		src, _ = _resizeImage(ori, uint(wl.Size), len(wl.Content))
	}
	//src = wl.LogoStream.(image.Image)
	//fmt.Println(src)
	dstBound := dst.Bounds()
	srcBound := src.Bounds()
	dstRect := image.Rect((dstBound.Max.X/2)-(srcBound.Max.X/2), (dstBound.Max.X/2)-(srcBound.Max.X/2), (dstBound.Max.X/2)+(srcBound.Max.X/2), (dstBound.Max.X/2)+(srcBound.Max.X/2))
	draw.Draw(dst, dstRect, src, srcBound.Min, draw.Over)
	buf := new(bytes.Buffer)
	err = png.Encode(buf, dst)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

type WithFrame struct {
	Content      string
	RecoverLevel qr.RecoveryLevel
	Size         int
	QrcodeStream []byte
	FrameStream  []byte
	Foreground   color.Color
	Background   color.Color
}

func (wf WithFrame) New() ([]byte, error) {
	return nil, nil
}

type WithLogoFrame struct {
	Content      string
	RecoverLevel qr.RecoveryLevel
	Size         int
	LogoStream   []byte
	QrcodeStream []byte
	FrameStream  []byte
	Foreground   color.Color
	Background   color.Color
}

func (wlf WithLogoFrame) New() ([]byte, error) {
	return nil, nil
}

func _convertByteToImage(is []byte) (image.Image, error) {
	i, _, err := image.Decode(bytes.NewReader(is))
	if err != nil {
		return nil, err
	}
	return i, nil
}

func _convertByteToDrawImage(is []byte) (draw.Image, error) {
	i, err := _convertByteToImage(is)
	if err == nil {
		di := image.NewRGBA(i.Bounds())
		draw.Draw(di, i.Bounds(), i, i.Bounds().Min, draw.Src)
		return di, nil
	}
	return nil, err
}

func _resizeImage(si image.Image, size uint, l int) (image.Image, error) {
	var ns uint
	if l >= ContentThreshold {
		ns = size / MaxLogoRatio
	} else {
		ns = size / MinLogoRatio
	}
	m := resize.Resize(ns, 0, si, resize.Bilinear)
	return m, nil

}

var (
	MaxLogoRatio     uint = 4
	MinLogoRatio     uint = 8
	ContentThreshold      = 40
)
