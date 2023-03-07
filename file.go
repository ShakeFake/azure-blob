/**
    @author: yunkaiwang
    @mail: yunkaiwang.tvunetwork.com
    @data: 2023/2/6
**/

package azure

import (
	"os"
)

// OpenFile 接受绝对路径，打开文件。
func OpenFile(filePath string) (*os.File, error) {
	_, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}

	fileHandler, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	return fileHandler, nil
}
