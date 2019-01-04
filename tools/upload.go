package tools

import (
	"bytes"
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"io/ioutil"
	"net/http"
	"tg_gif/common"
	"time"

	"github.com/chai2010/webp"
	"github.com/golang/glog"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

//upload.go 文件上传oss服务

var (
	//HTTPClient HTTPClient
	HTTPClient *http.Client //http.Client 是全局对象，注意设置超时时间问题
	ossClinet  *oss.Client
	bucketName string
)

const (
	maxIdleConnections = 100
	requestTimeout     = 120
)

// createHTTPClient for connection re-use
func createHTTPClient() *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: maxIdleConnections,
		},
		Timeout: time.Duration(requestTimeout) * time.Second,
	}
	return client
}

// OssConf oss配置文档
type OssConf struct {
	Endpoint        string `yaml:"Endpoint"`
	AccessKeyID     string `yaml:"AccessKeyID"`
	AccessKeySecret string `yaml:"AccessKeySecret"`
	BucketName      string `yaml:"BucketName"`
}

// OssInit oss 服务初始化
func OssInit(o *OssConf) {
	var err error
	ossClinet, err = oss.New(o.Endpoint, o.AccessKeyID, o.AccessKeySecret)
	if err != nil {
		panic("oss服务启动失败")
	}
	if err != nil {
		panic("oss服务启动失败")
	}
	bucketName = o.BucketName
	HTTPClient = createHTTPClient()

}

// DowAndUploadToOss 网络文件下载然后传到oss,count 控制并发数量
func DowAndUploadToOss(files []*common.FileWithURL, count int) {
	goroutineCount := make(chan int, count)
	for i, file := range files {
		glog.V(5).Info(file)
		goroutineCount <- i
		go dowWihtGenGifAndUploadToOss(file, goroutineCount)
	}

}

// dowAndUploadToOss 直接上传webp到oss --废弃
func dowAndUploadToOss(f *common.FileWithURL, c chan int) {
	defer func(i chan int) {
		<-i
	}(c)
	if f.URL == "uploaded" {
		return
	}
	resp, err := HTTPClient.Get(f.URL)
	if err != nil {
		glog.Error("请求错误", err)
		return
	}
	defer resp.Body.Close()

	bucket, err := ossClinet.Bucket(bucketName)
	if err != nil {
		glog.Error("bucket创建失败“：", err)
		return
	}
	err = bucket.PutObject(f.Name, resp.Body)
	if err != nil {
		glog.Error("上传错误：", err)
		return
	}
}

// dowWihtGenGifAndUploadToOss 下载然后生成gif再上传到oss
func dowWihtGenGifAndUploadToOss(f *common.FileWithURL, c chan int) {
	defer func(i chan int) {
		<-i
	}(c)
	if f.URL == "uploaded" {
		return
	}
	resp, err := HTTPClient.Get(f.URL)
	if err != nil {
		glog.Error("请求错误", err)
		return
	}
	defer resp.Body.Close()
	RGBAData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Error("读取错误", err)
		return
	}
	w, h, hasAlpha, err := webp.GetInfo(RGBAData)
	if err != nil {
		glog.Error("获取RGBA信息错误", err)
		return
	}

	img, err := webp.DecodeRGBA(RGBAData)
	if err != nil {
		glog.Error("DecodeRGBA错误", err)
		return
	}

	if hasAlpha {
		for x := 0; x < w; x++ {
			for y := 0; y < h; y++ {
				cr := img.RGBAAt(x, y).R
				cg := img.RGBAAt(x, y).G
				cb := img.RGBAAt(x, y).B
				ca := img.RGBAAt(x, y).A
				if ca > 0x00 {
					img.Set(x, y, color.RGBA{cr, cg, cb, 255})
				} else {
					img.Set(x, y, color.Transparent)
				}
			}
		}
	}

	b := bytes.NewBuffer(make([]byte, 0))
	palette := append(palette.WebSafe, color.Transparent)

	anim := gif.GIF{}
	paletted := image.NewPaletted(img.Bounds(), palette)
	draw.FloydSteinberg.Draw(paletted, img.Bounds(), img, image.ZP)
	anim.Image = append(anim.Image, paletted)
	anim.Delay = append(anim.Delay, 15)
	gif.EncodeAll(b, &anim)

	bucket, err := ossClinet.Bucket(bucketName)
	if err != nil {
		glog.Error("bucket创建失败“：", err)
		return
	}
	err = bucket.PutObject(f.Name, b)

	if err != nil {
		glog.Error("上传错误：", err)
		return
	}
}
