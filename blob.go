package azure

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/Azure/azure-storage-blob-go/azblob"
	"io"
	"os"
	"path"
	"sync"
	"time"
)

type Blob struct {
	Name      string
	BlockSize int
	BlobURL   azblob.BlockBlobURL

	mu sync.Mutex
}

type BlobOperator interface {
	Upload()
	Download(blobName string, downloadPath string) bool
}

func (b *Blob) UploadOneFile(filePath string) (bool, error) {
	fileHandler, err := OpenFile(filePath)
	if err != nil {
		return false, err
	}

	ctx := context.Background()
	filename := path.Base(fileHandler.Name())

	firstTime := time.Now().UnixMilli()
	resp, err := b.BlobURL.Upload(
		ctx,
		fileHandler,
		azblob.BlobHTTPHeaders{ContentType: "text/plain"},
		azblob.Metadata{},
		azblob.BlobAccessConditions{},
		azblob.DefaultAccessTier,
		nil,
		azblob.ClientProvidedKeyOptions{},
		azblob.ImmutabilityPolicyOptions{},
	)
	secondTime := time.Now().UnixMilli()

	if err != nil {
		return false, err
	}
	fmt.Printf("上传文件 %v 时间为 %v ms", filename, secondTime-firstTime)
	// 此处暂时搁置
	ShowResp(resp, err)
	return true, nil
}

func (b *Blob) StageUpload(filePath string) (bool, error) {
	fileH, err := OpenFile(filePath)
	if err != nil {
		return false, err
	}
	reader := bufio.NewReader(fileH)

	// todo: 此处需要转换为公共的 public 参数
	// 信号量控制上传
	signal := make(chan int, 10)
	defer close(signal)

	ctx := context.Background()
	blockIds := make([]string, 0)
	blockBuffer := make([]byte, b.BlockSize)

	partBegin := time.Now().UnixMilli()
	index := 0
	for true {
		number, err := reader.Read(blockBuffer)
		if err != nil {
			if err == io.EOF {
				break
			}
		}
		index++
		blockid := BlockIDIntToBase64(index)
		blockIds = append(blockIds, blockid)
		newBlock := make([]byte, len(blockBuffer))
		copy(newBlock, blockBuffer)
		signal <- 1
		go func() {
			// todo: 还要处理是否成功业务。
			b.BlobURL.StageBlock(
				ctx,
				blockid,
				bytes.NewReader(newBlock[:number]),
				azblob.LeaseAccessConditions{},
				nil,
				azblob.ClientProvidedKeyOptions{},
			)
		}()
	}
	partEnd := time.Now().UnixMilli()
	fmt.Printf("part time is %v\n", partEnd-partBegin)

	commitBegin := time.Now().UnixMilli()
	b.BlobURL.CommitBlockList(
		ctx,
		blockIds,
		azblob.BlobHTTPHeaders{},
		azblob.Metadata{},
		azblob.BlobAccessConditions{},
		azblob.DefaultAccessTier,
		nil,
		azblob.ClientProvidedKeyOptions{},
		azblob.ImmutabilityPolicyOptions{},
	)
	commitEnd := time.Now().UnixMilli()
	fmt.Printf("commit time is %v\n", commitEnd-commitBegin)

	return true, nil
}

// Download 下载块
// 如果获取了对应快的 releaseId，在下载的时候，需要将 releaseId 塞入。
func (b *Blob) Download(downloadPath string) bool {

	ctx := context.Background()
	get, err := b.BlobURL.Download(
		ctx,
		0,
		0,
		// 此处可设置 releaseid
		azblob.BlobAccessConditions{},
		false,
		azblob.ClientProvidedKeyOptions{},
	)
	if err != nil {
		panic(err)
	}

	// reader 只能读一次
	reader := get.Body(azblob.RetryReaderOptions{})
	defer reader.Close()

	//downloadData := bytes.Buffer{}
	//downloadData.ReadFrom(reader)
	//fmt.Println(downloadData)

	fileH, err := os.Create(downloadPath)
	defer fileH.Close()
	if err != nil {
		panic(err)
	}
	block := make([]byte, 1024)

	next := true
	for next {
		numberR, err := reader.Read(block)
		if err != nil {
			if err == io.EOF {
				next = false
			} else {
				break
			}
		}
		numberW, err := fileH.Write(block[:numberR])
		if err != nil {
			break
		}
		if numberR != numberW {
			panic(errors.New("读写不一致"))
		}
	}

	return false
}

func (b *Blob) TestBlobMethods() {
	// 从 url 中进行获取部分 url

}
