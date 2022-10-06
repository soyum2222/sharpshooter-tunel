package crypto

import (
	"bytes"
	"compress/zlib"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"io"
)

type AesCbc struct {
	Key    string
	KenLen int
}

func (a *AesCbc) Encrypt(src []byte) ([]byte, error) {
	key := evpBytesToKey(a.Key, a.KenLen)
	b, err := encryptAES(src, key)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	w := zlib.NewWriter(&buf)
	_, err = w.Write(b)
	if err != nil {
		return nil, err
	}
	w.Close()

	return buf.Bytes(), nil
}

func (a *AesCbc) Decrypt(src []byte) ([]byte, error) {

	b := bytes.NewReader(src)

	r, err := zlib.NewReader(b)
	if err != nil {
		return nil, err
	}

	r.Close()
	src, err = io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	key := evpBytesToKey(a.Key, a.KenLen)
	return decryptAES(src, key)
}

func padding(src []byte, blocksize int) []byte {
	padnum := blocksize - len(src)%blocksize
	pad := bytes.Repeat([]byte{byte(padnum)}, padnum)
	return append(src, pad...)
}

func unpadding(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])
	if length-unpadding >= length || length-unpadding < 0 {
		return []byte{}
	}
	return src[:(length - unpadding)]
}

func encryptAES(src []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	src = padding(src, block.BlockSize())
	blockmode := cipher.NewCBCEncrypter(block, key)
	blockmode.CryptBlocks(src, src)
	return src, nil
}

func decryptAES(src []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockmode := cipher.NewCBCDecrypter(block, key)
	blockmode.CryptBlocks(src, src)
	src = unpadding(src)
	return src, nil
}

func evpBytesToKey(password string, keyLen int) (key []byte) {
	const md5Len = 16

	cnt := (keyLen-1)/md5Len + 1
	m := make([]byte, cnt*md5Len)
	copy(m, md5sum([]byte(password)))

	// Repeatedly call md5 until bytes generated is enough.
	// Each call to md5 uses data: prev md5 sum + password.
	d := make([]byte, md5Len+len(password))
	start := 0
	for i := 1; i < cnt; i++ {
		start += md5Len
		copy(d, m[start-md5Len:start])
		copy(d[md5Len:], password)
		copy(m[start:], md5sum(d))
	}
	return m[:keyLen]
}

func md5sum(d []byte) []byte {
	h := md5.New()
	h.Write(d)
	return h.Sum(nil)
}
