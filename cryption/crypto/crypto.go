package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"errors"
)

//补码
//AES加密数据块分组长度必须为128bit(byte[16])，
//密钥长度可以是128bit(byte[16])、192bit(byte[24])、256bit(byte[32])中的任意一个。
func pkcs7Padding(ciphertext *[]byte, blocksize int) {
	padding := blocksize - len(*ciphertext)%blocksize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	*ciphertext = append(*ciphertext, padtext...)
}

//去码
func pkcs7UnPadding(origData *[]byte) {
	length := len(*origData)
	unpadding := int((*origData)[length-1])
	*origData = (*origData)[:(length - unpadding)]
}

func AesEncrypt(data *[]byte, key string) (*[]byte, error) {
	// 转成字节数组
	k := []byte(key)
	// 分组秘钥
	// NewCipher该函数限制了输入k的长度必须为16, 24或者32
	block, err := aes.NewCipher(k)
	if err != nil {
		return nil, err
	}
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 补全码
	pkcs7Padding(data, blockSize)
	// 加密模式
	blockMode := cipher.NewCBCEncrypter(block, k[:blockSize])
	// 创建数组
	cryted := make([]byte, len(*data))
	// 加密
	blockMode.CryptBlocks(cryted, *data)
	//return base64.StdEncoding.EncodeToString(cryted), nil
	return &cryted, nil
}

func AesDecrypt(cipherByte *[]byte, key string) (*[]byte, error) {
	var err error
	defer func() {
		if p := recover(); p != nil {
			if v, ok := p.(string); ok == true {
				err = errors.New(v)
			} else {
				err = errors.New("解密错误")
			}
		}
	}()
	// 转成字节数组
	// cipherByte, _ := base64.StdEncoding.DecodeString(cipherText)
	k := []byte(key)
	// 分组秘钥
	block, err := aes.NewCipher(k)
	if err != nil {
		return nil, err
	}
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 加密模式
	blockMode := cipher.NewCBCDecrypter(block, k[:blockSize])
	// 创建数组
	orig := make([]byte, len(*cipherByte))
	// 解密
	blockMode.CryptBlocks(orig, *cipherByte)
	// 去补全码
	pkcs7UnPadding(&orig)
	return &orig, nil
}

func DesEncrypt(data *[]byte, key string) (*[]byte, error) {
	// 转成字节数组
	k := []byte(key)
	// 分组秘钥
	// NewCipher该函数限制了输入k的长度必须为16, 24或者32
	block, err := des.NewCipher(k)
	if err != nil {
		return nil, err
	}
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 补全码
	pkcs7Padding(data, blockSize)
	// 加密模式
	blockMode := cipher.NewCBCEncrypter(block, k[:blockSize])
	// 创建数组
	cryted := make([]byte, len(*data))
	// 加密
	blockMode.CryptBlocks(cryted, *data)
	// return base64.StdEncoding.EncodeToString(cryted), nil
	return &cryted, nil
}

func DesDecrypt(cipherByte *[]byte, key string) (*[]byte, error) {
	var err error
	defer func() {
		if p := recover(); p != nil {
			if v, ok := p.(string); ok == true {
				err = errors.New(v)
			} else {
				err = errors.New("解密错误")
			}
		}
	}()
	// 转成字节数组
	// cipherByte, _ := base64.StdEncoding.DecodeString(cipherText)
	k := []byte(key)
	// 分组秘钥
	block, err := des.NewCipher(k)
	if err != nil {
		return nil, err
	}
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 加密模式
	blockMode := cipher.NewCBCDecrypter(block, k[:blockSize])
	// 创建数组
	orig := make([]byte, len(*cipherByte))
	// 解密
	blockMode.CryptBlocks(orig, *cipherByte)
	// 去补全码
	pkcs7UnPadding(&orig)
	return &orig, nil
}
