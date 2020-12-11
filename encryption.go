package go_utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
)

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

//AES加密
func AesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = PKCS7Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

//AES解密
func AesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS7UnPadding(origData)
	return origData, nil
}

// 加解密算法
type gaes struct {
	key   string                 `json:"-"`
	Value map[string]interface{} `json:"value"`
}

// 初始化
func NewGAes(key string) *gaes {
	return &gaes{
		key:   key,
		Value: make(map[string]interface{}),
	}
}

// 解码
func (a *gaes) Decode(data string) error {

	decodeString, _ := base64.StdEncoding.DecodeString(data)

	// 解密请求头
	decrypt, err := AesDecrypt(decodeString, []byte(a.key))
	if err != nil {
		return err
	}

	// 反序列化
	err = json.Unmarshal(decrypt, &a)
	if err != nil {
		return err
	}
	return nil

}

// 编码
func (a *gaes) Encode() (string, error) {

	// 序列化
	data, err := json.Marshal(a)
	if err != nil {
		return "", err
	}

	// Aes加密
	decrypt, err := AesEncrypt(data, []byte(a.key))
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(decrypt), nil

}

