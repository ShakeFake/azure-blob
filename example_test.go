/**
    @author: yunkaiwang
    @mail: yunkaiwang.tvunetwork.com
    @data: 2023/2/9
**/

package azure

import (
	"fmt"
	"path"
	"strings"
	"testing"
)

type FileInfo struct {
	FilePath string
	filename string
}

func (f *FileInfo) GetName() string {
	if path.IsAbs(f.FilePath) {
		return path.Dir(f.FilePath)
	} else {
		paths := strings.Split(f.FilePath, "/")
		return paths[len(paths)-1]
	}
}

func getFile() []FileInfo {
	var fileinfos []FileInfo
	//fileinfos = append(fileinfos, FileInfo{FilePath: file10m})
	fileinfos = append(fileinfos, FileInfo{FilePath: filePathOne})
	fileinfos = append(fileinfos, FileInfo{FilePath: filePathTwo})
	return fileinfos
}

var (
	filePathOne   = "./source/one.txt"
	filePathTwo   = "./source/two.txt"
	filePathStage = "./source/stage.txt"
)

// todo: 考虑使用 mock 进行替换。
func TestUploadFiles(t *testing.T) {
	files := getFile()

	azure := GetAzure(accountName, accountKey)
	container := azure.GetContainerURL("yunkaitest")

	for _, file := range files {
		success, err := container.GetBlob(file.GetName()).UploadOneFile(file.FilePath)
		if err != nil {
			panic(err)
		}
		fmt.Sprintf("上传文件 %v 结果为 %v", file.FilePath, success)
	}
}

func TestCreateBlob(t *testing.T) {
	azure := GetAzure(accountName, accountKey)
	container := azure.GetContainerURL("yunkaitest")
	blobURL := container.GetBlob("a/b/c")
	success, err := blobURL.UploadOneFile(filePathOne)
	if err != nil {
		panic(err)
	}
	fmt.Sprintf("结果为 %v", success)
}

func TestListContainer(t *testing.T) {
	azure := GetAzure(accountName, accountKey)
	container := azure.GetContainerURL("yunkaitest")
	container.ListBlob()
}

func TestDeleteContainerBlob(t *testing.T) {
	azure := GetAzure(accountName, accountKey)
	container := azure.GetContainerURL("yunkaitest")
	blobLists := container.ListBlob()

	_ = blobLists
	//for _, blob := range blobLists {
	//	fmt.Println(blob)
	//	resp, err := container.GetBlob(blob).BlobURL.Delete(context.Background(), azblob.DeleteSnapshotsOptionNone, azblob.BlobAccessConditions{})
	//	if err != nil {
	//		panic(err)
	//	}
	//	_ = resp
	//}
	container.ListBlob()
}
