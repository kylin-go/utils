package pubkey

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

//var pubKey = `-----BEGIN PUBLIC KEY-----
//MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCAtw74sSN6eLcpnCyBbBN2mu29
//7/2uEZUqCYS2uYCYqIV3b/RvhrRDqlzxvxpXuXYpwLrU/SElVvbao/WnX8/g5WE5
//alION1NNtoQgdZVt/AcWiJowXIN2T6BVYx3JebPPSFC/Hhr5TX/EPKL6X7YArrbQ
//5j5t0EkUIfE0kuWLBQIDAQAB
//-----END PUBLIC KEY-----`
//
//var privateKey = `-----BEGIN RSA PRIVATE KEY-----
//MIICWwIBAAKBgQCAtw74sSN6eLcpnCyBbBN2mu297/2uEZUqCYS2uYCYqIV3b/Rv
//hrRDqlzxvxpXuXYpwLrU/SElVvbao/WnX8/g5WE5alION1NNtoQgdZVt/AcWiJow
//XIN2T6BVYx3JebPPSFC/Hhr5TX/EPKL6X7YArrbQ5j5t0EkUIfE0kuWLBQIDAQAB
//AoGAOxUME76nxOdVWA2+Zhf4ZShfeaCIJtceS6H7364Np8UvInBq2KiR5T91k2f/
//jQXuBeNYPz0D8nJVNG4vbAkwT2FCv4zWnji+37tYsJizdKN1itKM09pstlUy6vD1
//wdfmT3c4uV5oozO3yw6DQ6jdDZqxwj3VZSyaNHLGyGz0588CQQCDcN9BWNun+sgM
//DVLUIoe20/Vv4S12tYln/+W4mGb435bC0NNFwznI+FxsYlt8uOwqNBHu2uq0i76A
//OkbSz1snAkEA+rDocxodB+5CQcACqlzPH8APUqLBIXrPnDAzp7Txty0lfjWWFLQi
//nw8FjiLTbzWMTjkIkn1TsxXiV/PovpTz8wJAZzEqVadpbAvbGnsrWBhz6/mka12h
//z9zeL6Qbuj0MOr9vISvJcq++oiU6imz93oFgCBIxMhD0yyIbQZh/GeppaQJAWILL
//n5ARvfIWfKZxinr4OkqSXmfObqaw1IGES2ssiLMs8LG0ypyLOMOR/4w2QKaUYi3q
//4+XR/oH0h387psZqlwJAJGxtIDp8ufRgQIkrIzG+Xeyl3IoeIoZATNmzzdot5o4W
//f5NQX4zU6v4DzUwksPNT7KsWT990HRFmAgmtqaLveA==
//-----END RSA PRIVATE KEY-----`

// 支持超出117字节数据分段加密
func RsaPkcs1PubEncrypt(data, publicKey []byte) ([][]byte, error) {
	if len(data) == 0 {
		return nil, errors.New("加密数据不能为空")
	}
	var cipherText [][]byte
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	for i := 0; i*117 < len(data); i++ {
		endIndex := (i + 1) * 117
		if endIndex > len(data) {
			endIndex = len(data)
		}
		d, e := rsa.EncryptPKCS1v15(rand.Reader, pub, data[i*117:endIndex])
		if e != nil {
			return nil, e
		}
		cipherText = append(cipherText, d)
	}
	return cipherText, nil
}

// 支持超出117字节数据分段解密
func RsaPkcs1PrivateDecrypt(cipherTexts [][]byte, privateKey []byte) ([]byte, error) {
	if len(cipherTexts) == 0 {
		return nil, errors.New("RSA私钥解密操作密文不能为空")
	}
	var rData []byte
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key error!")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	for _, data := range cipherTexts {
		d, e := rsa.DecryptPKCS1v15(rand.Reader, priv, data)
		if e != nil {
			return nil, e
		}
		rData = append(rData, d...)
	}
	return rData, nil
}
