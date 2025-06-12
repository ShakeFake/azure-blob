package azure

import (
	"context"
	"fmt"
	"github.com/Azure/azure-storage-blob-go/azblob"
	"log"
	"path"
	"sync"
	"time"
)

type Container struct {
	Name         string
	ContainerURL azblob.ContainerURL
	UploadFiles  []UploadList
	Blobs        map[string]*Blob

	mu sync.Mutex
	wg sync.WaitGroup
}

type ContainerOperator interface {
}

type UploadList struct {
	FilePath string
}

// Create 直接创建一个 container
func (c *Container) Create(containerName string) {
	ctx := context.Background()
	// 在创建时，指定使用 container 的权限。
	resp, err := c.ContainerURL.Create(ctx, azblob.Metadata{}, azblob.PublicAccessContainer)
	if err != nil {
		panic(err)
	}
	fmt.Println(&resp)
}

// Delete 删除一个container
func (c *Container) Delete(containerName string) {
	ctx := context.Background()
	resp, err := c.ContainerURL.Delete(ctx, azblob.ContainerAccessConditions{})
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}

func (c *Container) ListBlob() []string {
	var blobList []string
	ctx := context.Background()
	for marker := (azblob.Marker{}); marker.NotDone(); {
		fmt.Println(c.ContainerURL)
		listBlob, err := c.ContainerURL.ListBlobsFlatSegment(ctx, marker, azblob.ListBlobsSegmentOptions{})
		if err != nil {
			log.Fatal(err)
		}
		marker = listBlob.NextMarker
		for _, blobInfo := range listBlob.Segment.BlobItems {
			// todo: 这块考虑能不能对接导出文件列表等。Excel 等。
			blobList = append(blobList, blobInfo.Name)
		}
		time.Sleep(time.Second)
	}
	return blobList
}

func (c *Container) Upload() {
	// 做一个定时上传
	c.mu.Lock()
	files := c.UploadFiles
	c.mu.Unlock()

	for _, file := range files {
		filename := path.Base(file.FilePath)
		blob := c.GetBlob(filename)
		go blob.UploadOneFile(file.FilePath)
	}
	c.wg.Wait()
}

func (c *Container) GetBlob(blobName string) *Blob {
	c.mu.Lock()
	defer c.mu.Unlock()
	if b, ok := c.Blobs[blobName]; ok {
		return b
	} else {
		blobURL := c.ContainerURL.NewBlockBlobURL(blobName)
		b := &Blob{
			Name:    blobName,
			BlobURL: blobURL,
		}
		c.Blobs[blobName] = b
		return b
	}
}

func (c *Container) DeleteBlob(blobName string) {
	blobURL := c.ContainerURL.NewBlockBlobURL(blobName)
	ctx := context.Background()
	resp, err := blobURL.Delete(ctx, azblob.DeleteSnapshotsOptionNone, azblob.BlobAccessConditions{})
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}
