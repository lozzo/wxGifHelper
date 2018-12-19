package tools

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

//upload.go 文件上传oss服务

var (
	httpClient *http.Client //http.Client 是全局对象，注意设置超时时间问题
	ossClinet  *oss.Client
	wg         sync.WaitGroup
	bucketName string
)

const (
	maxIdleConnections = 20
	requestTimeout     = 60
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
	Endpoint        string
	AccessKeyID     string
	AccessKeySecret string
	BucketName      string
}

// Init oss 服务初始化
func Init(o *OssConf) {
	var err error
	ossClinet, err = oss.New(o.Endpoint, o.AccessKeyID, o.AccessKeySecret)
	if err != nil {
		panic("oss服务启动失败")
	}
	if err != nil {
		panic("oss服务启动失败")
	}
	bucketName = o.BucketName
	httpClient = createHTTPClient()

}

// FileWithURL 网络文件
type FileWithURL struct {
	URL  string
	Name string
}

// DowAndUploadToOss 网络文件下载然后传到oss,count 控制并发数量
func DowAndUploadToOss(files []*FileWithURL, count int) {
	goroutineCount := make(chan int, count)
	for i, file := range files {
		goroutineCount <- i
		go dowAndUploadToOss(file, goroutineCount)
	}

}

func dowAndUploadToOss(f *FileWithURL, c chan int) {
	resp, err := httpClient.Get(f.URL)
	if err != nil {
		fmt.Println("请求错误：", err)
		return
	}
	defer resp.Body.Close()

	bucket, err := ossClinet.Bucket(bucketName)
	if err != nil {
		fmt.Println("bucket创建失败“：", err)
		return
	}
	err = bucket.PutObject(f.Name, resp.Body)
	if err != nil {
		fmt.Println("上传错误：", err)
		return
	}
	<-c
}
