/**
    @author: yunkaiwang
    @mail: yunkaiwang.tvunetwork.com
    @data: 2023/2/9
**/

package azure

import (
	"encoding/base64"
	"encoding/binary"
)

func BlockIDBinaryToBase64(blockID []byte) string {
	return base64.StdEncoding.EncodeToString(blockID)
}

func BlockIDBase64ToBinary(blockID string) []byte {
	binary, _ := base64.StdEncoding.DecodeString(blockID)
	return binary
}

// BlockIDIntToBase64 将整数 id 转化为 base64 标识
func BlockIDIntToBase64(blockID int) string {
	binaryBlockID := (&[4]byte{})[:]
	return BlockIDBinaryToBase64(binaryBlockID)
}

// BlockIDBase64ToInt 将 base64 标识转化为 int id
func BlockIDBase64ToInt(blockID string) int {
	return int(binary.LittleEndian.Uint32(BlockIDBase64ToBinary(blockID)))
}
