package goscs

import (
	"bytes"
	"crypto/aes"
)

const BLOCKSIZE int = aes.BlockSize

func addPadding(slice []byte) []byte {
	padding := BLOCKSIZE - len(slice)%BLOCKSIZE
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(slice, padText...)
}

func removePadding(slice []byte) []byte {
	length := len(slice)
	unpadding := int(slice[length-1])
	return slice[:(length - unpadding)]
}
