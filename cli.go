/**
    @author: yunkaiwang
    @mail: yunkaiwang.tvunetwork.com
    @data: 2023/2/6
**/

package azure

import (
	"fmt"
	"github.com/Azure/azure-pipeline-go/pipeline"
	"github.com/Azure/azure-storage-blob-go/azblob"
	"log"
	"net/url"
)

var (
	CliManager  Azure
	accountName = ""
	accountKey  = ""
)

// Init 初始化操作用来
func Init(name string, key string) error {
	credential, err := azblob.NewSharedKeyCredential(name, key)
	if err != nil {
		return err
	}

	// 目前处理，共用一个 pipeline
	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})
	u, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net", name))

	CliManager.mu.Lock()
	defer CliManager.mu.Unlock()
	// init global azure
	CliManager.Name = name
	CliManager.Key = key
	CliManager.ServiceURL = azblob.NewServiceURL(*u, p)
	CliManager.IsInit = true

	return nil
}

// ShowResp 直接打印 azure 的返回信息
func ShowResp(response pipeline.Response, err error) {
	if err != nil {
		if stgErr, ok := err.(azblob.StorageError); !ok {
			log.Fatal(err)
		} else {
			fmt.Print("Failure: " + stgErr.Response().Status + "\n")
		}
	} else {
		if get, ok := response.(*azblob.DownloadResponse); ok {
			get.Body(azblob.RetryReaderOptions{}).Close()
		}
		fmt.Print("Success: " + response.Response().Status + "\n")
	}
}
