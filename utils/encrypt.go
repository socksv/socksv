package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
)

var Key = "S12@49I789AOqdef"

func AesEncrypt(orig []byte, key string) []byte {
	k := []byte(key)

	block, _ := aes.NewCipher(k)
	blockSize := block.BlockSize()
	orig = padding(orig, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, k[:blockSize])
	cryted := make([]byte, len(orig))
	blockMode.CryptBlocks(cryted, orig)

	return cryted
}
func AesDecrypt(cryted []byte, key string) []byte {
	k := []byte(key)

	block, _ := aes.NewCipher(k)
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, k[:blockSize])
	orig := make([]byte, len(cryted))
	blockMode.CryptBlocks(orig, cryted)
	orig = unPadding(orig)
	return orig
}

func padding(ciphertext []byte, blocksize int) []byte {
	padding := blocksize - len(ciphertext)%blocksize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func unPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
