package utils

import (
	"bytes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"github.com/tjfoc/gmsm/sm4"
)

// Sm4Encrypt sm4 加密
func Sm4Encrypt(key []byte, plainText string) (string, error) {
	if len(key) != 16 {
		return "", errors.New("key length not 16")
	}
	block, err := sm4.NewCipher(key)
	if err != nil {
		return "", err
	}
	var (
		blockSize = block.BlockSize()
		iv        = make([]byte, sm4.BlockSize)
	)
	origData := Pkcs5Padding([]byte(plainText), blockSize)
	blockMode := cipher.NewCBCEncrypter(block, iv)
	cryted := make([]byte, len(origData))
	blockMode.CryptBlocks(cryted, origData)
	return base64.StdEncoding.EncodeToString(cryted), nil
}

// Sm4Decrypt sm4 解密
func Sm4Decrypt(key []byte, cipherText string) ([]byte, error) {
	block, err := sm4.NewCipher(key)
	if err != nil {
		return nil, err
	}
	// 密文 base64反解码
	cbyt, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return nil, err
	}
	iv := make([]byte, sm4.BlockSize)
	blockMode := cipher.NewCBCDecrypter(block, iv)
	origData := make([]byte, len(cbyt))
	blockMode.CryptBlocks(origData, cbyt)
	origData = Pkcs5UnPadding(origData)
	return origData, nil
}

// Pkcs5Padding pkcs5填充
func Pkcs5Padding(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

func Pkcs5UnPadding(src []byte) []byte {
	length := len(src)
	if length == 0 {
		return nil
	}
	unpadding := int(src[length-1])
	return src[:(length - unpadding)]
}
