// Package xmc_qrcode
// @Description: 利用skip2的QRCODE库生成矩阵，并输出PNG,JPEG,SVG
//TODO 建立命令行应用 补充带附加外框的函数
package xmc_qrcode

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	svg "github.com/ajstarks/svgo"
	"github.com/nfnt/resize"
	qr "github.com/skip2/go-qrcode"
	"github.com/xmanchai/go-xmc-tools/mime"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"os"
	"reflect"
)

// Qrcode
//  @Description:
//
type Qrcode interface {
	New() ([]byte, error)
	String(is []byte) string
	Write(is []byte, f string) error
}

const (
	MimePNG = iota //默认输出格式
	MimeJPEG
	MimeSVG
)

const (
	SVGRect = iota //默认SVG绘图像素
	SVGCircle
)

type Default struct {
	Content      string
	RecoverLevel qr.RecoveryLevel
	Size         int
	MimeType     int //PNG | JPEG | SVG
	SVGPixel     int //SVGRect | SVGCircle
}

func NewDefaultStruct() Default {
	return Default{
		Content:      "Default struct",
		RecoverLevel: qr.Medium,
		MimeType:     MimePNG,
		SVGPixel:     SVGRect,
	}
}

// New 函数名称
//  @描述: 创建一个默认的二维码 ，纠错率为中，输出类型为PNG
//  @结构 d
//  @返回值 []byte
//  @返回值 error
//
func (d Default) New() ([]byte, error) {
	q, err := qr.New(d.Content, d.RecoverLevel)
	size := len(q.Bitmap()) * DefaultImagePixel //默认不设置大小，根据二维码矩阵的维度+1*默认像素步长得到大小
	if d.Size == 0 {
		d.Size = size
	}
	if err != nil {
		return nil, err
	}
	dd, err := _convertByteToTarget(q, d.MimeType, d.SVGPixel, d.Size)
	if err != nil {
		return nil, err
	}
	return dd, nil
}

func (d Default) String(is []byte) string {
	return base64.StdEncoding.EncodeToString(is) //返回BASE64编码字符串
}

func (d Default) Write(is []byte, f string) error {
	var fn string
	switch d.MimeType {
	case MimeSVG:
		fn = ".svg"
	case MimeJPEG:
		fn = ".jpg"
	default:
		fn = ".png"
	}
	fs, err := os.Create(f + fn)
	if err != nil {
		return err
	}
	_, err = fs.Write(is)
	if err != nil {
		return err
	}
	_ = fs.Close()
	return nil
}

type WithColor struct {
	Content      string
	RecoverLevel qr.RecoveryLevel
	Size         int
	Foreground   color.Color
	Background   color.Color
	MimeType     int //PNG | JPEG | SVG
	SVGPixel     int
}

func NewWithColorStruct() WithColor {
	return WithColor{
		Content:    "With color",
		Foreground: color.Black,
		Background: color.White,
		MimeType:   MimePNG,
		SVGPixel:   SVGRect,
	}
}

func (wc WithColor) New() ([]byte, error) {
	c, err := qr.New(wc.Content, wc.RecoverLevel)
	c.ForegroundColor = wc.Foreground
	c.BackgroundColor = wc.Background
	size := len(c.Bitmap())
	if wc.Size == 0 {
		wc.Size = size
	}
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	wcd, err := _convertByteToTarget(c, wc.MimeType, wc.SVGPixel, wc.Size)
	if err != nil {
		return nil, err
	}
	return wcd, nil
}

func (wc WithColor) String(is []byte) string {
	return base64.StdEncoding.EncodeToString(is)
}

func (wc WithColor) Write(is []byte, f string) error {
	fs, err := os.Create(f)
	if err != nil {
		return err
	}
	_, err = fs.Write(is)
	if err != nil {
		return err
	}
	_ = fs.Close()
	return nil
}

type WithLogo struct {
	Content      string
	RecoverLevel qr.RecoveryLevel
	Size         int
	LogoStream   interface{}
	Foreground   color.Color
	Background   color.Color
	MimeType     int //PNG | JPEG | SVG
	SVGPixel     int
}

func NewWithLogoStruct() WithLogo {
	return WithLogo{
		Content:      "With logo",
		RecoverLevel: qr.Medium,
		Foreground:   color.Black,
		Size:         256,
		Background:   color.White,
		LogoStream:   image.Image(image.Rect(0, 0, 10, 10)),
		MimeType:     MimePNG,
		SVGPixel:     SVGRect,
	}
}

// New 函数名称
//  @描述: 生成带有LOGO的二维码 SVG的输出格式无法放置LOGO，请选择PNG或者JPEG
//  @结构 wl
//  @返回值 []byte
//  @返回值 error
//
func (wl WithLogo) New() ([]byte, error) {
	var src image.Image
	qc, err := qr.New(wl.Content, wl.RecoverLevel)
	qc.BackgroundColor = wl.Background
	qc.ForegroundColor = wl.Foreground

	if wl.Size < MinWithLogoSize || wl.Size > MaxWithLogoSize {
		wl.Size = MinWithLogoSize
	}
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

	//fmt.Println(reflect.TypeOf(wl.LogoStream).String())

	switch reflect.TypeOf(wl.LogoStream).String() {
	case "*image.NRGBA":
		src = wl.LogoStream.(image.Image)
		src, _ = _resizeImage(wl.LogoStream.(image.Image), uint(wl.Size))
	case "image.Image":
		src, _ = _resizeImage(wl.LogoStream.(image.Image), uint(wl.Size))
	case "[]uint8":
		ori, _ := _convertByteToImage(wl.LogoStream.([]byte))
		src, err = _resizeImage(ori, uint(wl.Size))
	case "multipart.sectionReadCloser":
		b := bytes.NewBuffer(nil)
		_, _ = io.Copy(b, wl.LogoStream.(io.Reader))
		ori, err := _convertByteToImage(b.Bytes())
		if err != nil {
			return nil, err
		}
		src, err = _resizeImage(ori, uint(wl.Size))
		if err != nil {
			return nil, err
		}
	default:
		src = wl.LogoStream.(image.Image)
	}
	dstBound := dst.Bounds()
	srcBound := src.Bounds()
	dstRect := image.Rect((dstBound.Max.X/2)-(srcBound.Max.X/2), (dstBound.Max.X/2)-(srcBound.Max.X/2), (dstBound.Max.X/2)+(srcBound.Max.X/2), (dstBound.Max.X/2)+(srcBound.Max.X/2))
	draw.Draw(dst, dstRect, src, srcBound.Min, draw.Over)

	if wl.MimeType == MimeSVG {
		r, err := _convertByteToTarget(qc, wl.MimeType, wl.SVGPixel, wl.Size)
		if err != nil {
			return nil, err
		}
		return r, nil
	}
	var r []byte
	var buf = new(bytes.Buffer)
	switch wl.MimeType {
	case MimeSVG:
		r, err = _convertByteToTarget(qc, wl.MimeType, wl.SVGPixel, wl.Size)

	case MimeJPEG:
		err = jpeg.Encode(buf, dst, &jpeg.Options{Quality: DefaultJpegQuality})
		r = buf.Bytes()
	default:
		err = png.Encode(buf, dst)
		r = buf.Bytes()
	}
	return r, err
}

func (wl WithLogo) String(is []byte) string {
	return base64.StdEncoding.EncodeToString(is)
}

func (wl WithLogo) Write(is []byte, f string) error {
	fs, err := os.Create(f)
	if err != nil {
		return err
	}
	_, err = fs.Write(is)
	if err != nil {
		return err
	}
	_ = fs.Close()
	return nil
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

func (wf WithFrame) Write(is []byte, f string) error {
	fs, err := os.Create(f)
	if err != nil {
		return err
	}
	_, err = fs.Write(is)
	if err != nil {
		return err
	}
	_ = fs.Close()
	return nil
}

// _convertByteToImage 函数名称
//  @描述: 转换字节流为image.Image
//  @参数 is
//  @返回值 image.Image
//  @返回值 error
//
func _convertByteToImage(is []byte) (image.Image, error) {
	var i image.Image
	var err error
	switch mime.GetFileType(is[:10]) {
	case "png":
		i, _, err = image.Decode(bytes.NewReader(is))
	case "jpg":
		i, err = jpeg.Decode(bytes.NewReader(is))
	default:
		return nil, errors.New("不支持的图片格式")
	}
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

func _resizeImage(si image.Image, size uint) (image.Image, error) {
	var ns uint
	ns = size / 5
	m := resize.Resize(ns, 0, si, resize.Bilinear)
	return m, nil

}

// _convertByteToTarget 函数名称
//  @描述:
//  @参数 q *qr.Qrcode 实例
//  @参数 t type of target 0 png 1 jpg 2 svg
//  @参数 mt SVG矩阵填充类型 0 rect 1 circle
//  @参数 size q.Image 方法需要非SVG图像的参数
//  @返回值 []byte
//  @返回值 error
//
func _convertByteToTarget(q *qr.QRCode, t, mt, size int) ([]byte, error) {
	var width, height int
	buf := new(bytes.Buffer)
	switch t {
	case 2:
		iw := bytes.NewBuffer(nil)
		s := svg.New(iw)
		bits := q.Bitmap()
		switch mt {
		case 1:
			width, height = (len(bits)+1)*MinSVGCircleRatio, (len(bits)+1)*MinSVGCircleRatio
		default:
			width, height = (len(bits)+1)*DefaultSVGPixel, (len(bits)+1)*DefaultSVGPixel
		}
		s.Start(width, height)
		_drawSVGMatrix(bits, s, q, mt, nil)
		s.End()
		return iw.Bytes(), nil
	case 1:
		err := jpeg.Encode(buf, q.Image(size), &jpeg.Options{Quality: DefaultJpegQuality})
		if err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	default:
		err := png.Encode(buf, q.Image(size))
		if err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	}
}

type svgRect struct {
	x, y, w, h int
	options    string
}
type svgCircle struct {
	x, y, r int
	options string
}

func _drawSVGRect(s *svg.SVG, r *svgRect) {
	s.Rect(r.x, r.y, r.w, r.h, r.options)
}

func _drawSVGCircle(s *svg.SVG, c *svgCircle) {
	s.Circle(c.x, c.y, c.r, c.options)
}

func _drawSVGMatrix(bits [][]bool, s *svg.SVG, q *qr.QRCode, t int, cc color.Color) {
	for y := range bits {
		for x := range bits[y] {
			if bits[y][x] != false {
				rf := &svgRect{y * MinSVGRect, x * MinSVGRect, MinSVGRect, MinSVGRect, fmt.Sprintf("fill:%s", _convertRGBAToString(s, q.ForegroundColor))}
				cf := &svgCircle{y*MinSVGCircleRatio + MinSVGCircleRadius, x*MinSVGCircleRatio + MinSVGCircleRadius, MinSVGCircleRadius, fmt.Sprintf("fill:%s", _convertRGBAToString(s, q.ForegroundColor))}
				switch t {
				case 1:
					_drawSVGCircle(s, cf)
				default:
					_drawSVGRect(s, rf)
				}

			} else {
				rb := &svgRect{y * MinSVGRect, x * MinSVGRect, MinSVGRect, MinSVGRect, fmt.Sprintf("fill:%s", _convertRGBAToString(s, q.BackgroundColor))}
				cb := &svgCircle{y*MinSVGCircleRatio + MinSVGCircleRadius, x*MinSVGCircleRatio + MinSVGCircleRadius, MinSVGCircleRadius, fmt.Sprintf("fill:%s", _convertRGBAToString(s, q.BackgroundColor))}
				switch t {
				case 1:
					_drawSVGCircle(s, cb)
				default:
					_drawSVGRect(s, rb)
				}
			}
		}
	}
}

func _convertRGBAToString(s *svg.SVG, c color.Color) string {
	r, g, b, a := c.RGBA()
	return s.RGBA(int(r), int(g), int(b), float64(a))
}

var (
	DefaultImagePixel  = 5
	DefaultSVGPixel    = 5
	MinSVGRect         = 5 // 最小SVG矩阵
	MinSVGCircleRadius = 3 // 最小SVG圆点半径
	MinSVGCircleRatio  = 6
	DefaultJpegQuality = 100
	MinWithLogoSize    = 256
	MaxWithLogoSize    = 1024
)
