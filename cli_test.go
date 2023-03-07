package azure

import (
	"context"
	"fmt"
	"github.com/Azure/azure-storage-blob-go/azblob"
	"testing"
)

// 测试 cli， 并且打印对应的 container
func TestInit(t *testing.T) {
	err := Init(accountName, accountKey)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	resp, err := CliManager.ServiceURL.ListContainersSegment(ctx, azblob.Marker{}, azblob.ListContainersSegmentOptions{})
	if err != nil {
		panic(err)
	}
	for k, c := range resp.ContainerItems {
		fmt.Println(k, c)
	}
}
