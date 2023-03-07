/**
    @author: yunkaiwang
    @mail: yunkaiwang.tvunetwork.com
    @data: 2023/2/9
**/

package azure

import (
	"context"
	"fmt"
	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/siddontang/go/log"
	"net/url"
	"sync"
)

type Azure struct {
	Name       string `json:"name"`
	Key        string `json:"key"`
	Credential azblob.Credential
	ServiceURL azblob.ServiceURL
	IsInit     bool `json:"isInit"`

	Containers map[string]*Container
	mu         sync.Mutex
}

type Operator interface {
	GetContainerURL(containerName string) *Container
	DeleteContainerURL(containerName string) bool
}

func GetAzure(name string, key string) *Azure {
	cre, err := azblob.NewSharedKeyCredential(name, key)
	if err != nil {
		log.Infof("")
		return nil
	}
	azure := &Azure{
		Name:       name,
		Key:        key,
		Credential: cre,
	}
	u, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net", azure.Name))
	azure.ServiceURL = azblob.NewServiceURL(*u, azblob.NewPipeline(azure.Credential, azblob.PipelineOptions{}))
	azure.Containers = make(map[string]*Container, 0)
	azure.IsInit = true
	return azure
}

// GetContainerURL 创建一个 container。会创建一个新的。
func (a *Azure) GetContainerURL(containerName string) *Container {
	a.mu.Lock()
	defer a.mu.Unlock()
	if c, ok := a.Containers[containerName]; ok {
		return c
	} else {
		containerURL := a.ServiceURL.NewContainerURL(containerName)
		c = &Container{
			Name:         containerName,
			ContainerURL: containerURL,
			UploadFiles:  []UploadList{},
			Blobs:        make(map[string]*Blob, 0),
		}
		ctx := context.Background()
		resp, err := containerURL.GetAccessPolicy(ctx, azblob.LeaseAccessConditions{})
		if err != nil {
			panic(err)
		}
		resp.BlobPublicAccess()
		ShowResp(resp, err)
		a.Containers[containerName] = c
		return c
	}
}

// DeleteContainerURL 删除一个 url
func (a *Azure) DeleteContainerURL(containerName string) bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	if _, ok := a.Containers[containerName]; ok {
		delete(a.Containers, containerName)
		return true
	}
	return true
}
