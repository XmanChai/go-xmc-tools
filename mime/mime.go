package mime

import (
	"bytes"
	"encoding/hex"
	"strconv"
	"strings"
	"sync"
)

var mimeMap sync.Map

func init() {
	mimeMap.Store("ffd8ffe000104a464946", "jpg")  //JPEG (jpg)
	mimeMap.Store("89504e470d0a1a0a0000", "png")  //PNG (png)
	mimeMap.Store("47494638396126026f01", "gif")  //GIF (gif)
	mimeMap.Store("49492a00227105008037", "tif")  //TIFF (tif)
	mimeMap.Store("424d228c010000000000", "bmp")  //16色位图(bmp)
	mimeMap.Store("424d8240090000000000", "bmp")  //24位位图(bmp)
	mimeMap.Store("424d8e1b030000000000", "bmp")  //256色位图(bmp)
	mimeMap.Store("41433130313500000000", "dwg")  //CAD (dwg)
	mimeMap.Store("3c21444f435459504520", "html") //HTML (html)   3c68746d6c3e0  3c68746d6c3e0
	mimeMap.Store("3c68746d6c3e0", "html")        //HTML (html)   3c68746d6c3e0  3c68746d6c3e0
	mimeMap.Store("3c21646f637479706520", "htm")  //HTM (htm)
	mimeMap.Store("48544d4c207b0d0a0942", "css")  //css
	mimeMap.Store("696b2e71623d696b2e71", "js")   //js
	mimeMap.Store("7b5c727466315c616e73", "rtf")  //Rich Text Format (rtf)
	mimeMap.Store("38425053000100000000", "psd")  //Photoshop (psd)
	mimeMap.Store("46726f6d3a203d3f6762", "eml")  //Email [Outlook Express 6] (eml)
	mimeMap.Store("d0cf11e0a1b11ae10000", "doc")  //MS Excel 注意：word、msi 和 excel的文件头一样
	mimeMap.Store("d0cf11e0a1b11ae10000", "vsd")  //Visio 绘图
	mimeMap.Store("5374616E64617264204A", "mdb")  //MS Access (mdb)
	mimeMap.Store("252150532D41646F6265", "ps")
	mimeMap.Store("255044462d312e350d0a", "pdf")  //Adobe Acrobat (pdf)
	mimeMap.Store("2e524d46000000120001", "rmvb") //rmvb/rm相同
	mimeMap.Store("464c5601050000000900", "flv")  //flv与f4v相同
	mimeMap.Store("00000020667479706d70", "mp4")
	mimeMap.Store("49443303000000002176", "mp3")
	mimeMap.Store("000001ba210001000180", "mpg") //
	mimeMap.Store("3026b2758e66cf11a6d9", "wmv") //wmv与asf相同
	mimeMap.Store("52494646e27807005741", "wav") //Wave (wav)
	mimeMap.Store("52494646d07d60074156", "avi")
	mimeMap.Store("4d546864000000060001", "mid") //MIDI (mid)
	mimeMap.Store("504b0304140000000800", "zip")
	mimeMap.Store("526172211a0700cf9073", "rar")
	mimeMap.Store("235468697320636f6e66", "ini")
	mimeMap.Store("504b03040a0000000000", "jar")
	mimeMap.Store("4d5a9000030000000400", "exe")        //可执行文件
	mimeMap.Store("3c25402070616765206c", "jsp")        //jsp文件
	mimeMap.Store("4d616e69666573742d56", "mf")         //MF文件
	mimeMap.Store("3c3f786d6c2076657273", "xml")        //xml文件
	mimeMap.Store("494e5345525420494e54", "sql")        //xml文件
	mimeMap.Store("7061636b616765207765", "java")       //java文件
	mimeMap.Store("406563686f206f66660d", "bat")        //bat文件
	mimeMap.Store("1f8b0800000000000000", "gz")         //gz文件
	mimeMap.Store("6c6f67346a2e726f6f74", "properties") //bat文件
	mimeMap.Store("cafebabe0000002e0041", "class")      //bat文件
	mimeMap.Store("49545346030000006000", "chm")        //bat文件
	mimeMap.Store("04000000010000001300", "mxp")        //bat文件
	mimeMap.Store("504b0304140006000800", "docx")       //docx文件
	mimeMap.Store("d0cf11e0a1b11ae10000", "wps")        //WPS文字wps、表格et、演示dps都是一样的
	mimeMap.Store("6431303a637265617465", "torrent")
	mimeMap.Store("6D6F6F76", "mov")         //Quicktime (mov)
	mimeMap.Store("FF575043", "wpd")         //WordPerfect (wpd)
	mimeMap.Store("CFAD12FEC5FD746F", "dbx") //Outlook Express (dbx)
	mimeMap.Store("2142444E", "pst")         //Outlook (pst)
	mimeMap.Store("AC9EBD8F", "qdf")         //Quicken (qdf)
	mimeMap.Store("E3828596", "pwl")         //Windows Password (pwl)
	mimeMap.Store("2E7261FD", "ram")         //Real Audio (ram)

}

func bytesToHexString(src []byte) string {
	res := bytes.Buffer{}
	if src == nil || len(src) <= 0 {
		return ""
	}
	temp := make([]byte, 0)
	for _, v := range src {
		sub := v & 0xFF
		hv := hex.EncodeToString(append(temp, sub))
		if len(hv) < 2 {
			res.WriteString(strconv.FormatInt(int64(0), 10))
		}
		res.WriteString(hv)
	}
	return res.String()
}

func GetFileType(fs []byte) string {
	var fileType string
	fileCode := bytesToHexString(fs)

	mimeMap.Range(func(key, value interface{}) bool {
		k := key.(string)
		v := value.(string)
		if strings.HasPrefix(fileCode, strings.ToLower(k)) ||
			strings.HasPrefix(k, strings.ToLower(fileCode)) {
			fileType = v
			return false
		}
		return true
	})
	return fileType
}
