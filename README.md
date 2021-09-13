# GO-XMC-TOOLS 轮子开会

![](logo.png)
### 1.起因

为了更好的学习和使用golang进行服务端开发，东平西凑地整理了很多DEMO和框架，在使用的过程中
慢慢体会到制造轮子真的是学习的好方法。所以，为了记录，更是为了加深学习的印象，特此开了此仓库
并将很多我用到的，摸索的轮子集中于此。其中很多代码还没有经过优化，仅供记录和分享。


### 2.轮子列表
+ qrcode
+ mime

### 3.轮子解析

#### github.com/xmanchai/go-xmc-tools/coding/qrcode

##### Mod 依赖

+ github.com/skip2/go-qrcode
+ github.com/nfnt/resize
+ github.com/ajstarks/svgo

接口定义:
```go
type Qrcode interface {
	New() ([]byte, error) //流模式输出
	String(is []byte) string //BASE64 字符串输出
	Write(is []byte, f string) error // 文件输出
}
```
```go
/*
默认结构及初始化函数 Default
*/
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

使用方法
nd := qrcode.NewDefaultStruct()
nd.Content = "要生成的内容"
s,_ := nd.New()
nd.String(s)
nd.Write(s,"文件名称") // 后缀会根据生成类型自动添加
```



WithColor
WithLogo
WithFrame
