package hash

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"errors"
)

func interface2bytes(context interface{}) ([]byte, error) {
	var c []byte
	switch context.(type) {
	case string:
		data := context.(string)
		c = []byte(data)
	case []byte:
		c = context.([]byte)
	default:
		return nil, errors.New("接收参数值支持string, []byte两种类型")
	}
	return c, nil
}

func Md5Encode(context interface{}) (string, error) {
	c, err := interface2bytes(context)
	if err != nil {
		return "", err
	}
	h := md5.New()
	if _, err := h.Write(c); err != nil {
		return "", err
	}
	cipherStr := h.Sum(nil)
	return hex.EncodeToString(cipherStr), nil
}

func Sha256Encode(context interface{}) (string, error) {
	c, err := interface2bytes(context)
	if err != nil {
		return "", err
	}
	h := sha256.New()
	if _, err := h.Write(c); err != nil {
		return "", err
	}
	cipherStr := h.Sum(nil)
	return hex.EncodeToString(cipherStr), nil
}

func Sha512Encode(context interface{}) (string, error) {
	c, err := interface2bytes(context)
	if err != nil {
		return "", err
	}
	h := sha512.New()
	if _, err := h.Write(c); err != nil {
		return "", err
	}
	cipherStr := h.Sum(nil)
	return hex.EncodeToString(cipherStr), nil
}

func Sha1Encode(context interface{}) (string, error) {
	c, err := interface2bytes(context)
	if err != nil {
		return "", err
	}
	h := sha1.New()
	if _, err := h.Write(c); err != nil {
		return "", err
	}
	cipherStr := h.Sum(nil)
	return hex.EncodeToString(cipherStr), nil
}
