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
	Key              string
	KenLen           int
	compressBuf      bytes.Buffer
	compressReader   io.ReadCloser
	uncompressReader *bytes.Reader
	compressWriter   *zlib.Writer
}

func (a *AesCbc) Encrypt(src []byte) ([]byte, error) {
	key := evpBytesToKey(a.Key, a.KenLen)
	b, err := encryptAES(src, key)
	if err != nil {
		return nil, err
	}

	if a.compressWriter == nil {
		a.compressWriter, _ = zlib.NewWriterLevel(&a.compressBuf, zlib.BestCompression)
	} else {
		a.compressWriter.Reset(&a.compressBuf)
	}

	_, err = a.compressWriter.Write(b)
	if err != nil {
		_ = a.compressWriter.Close()
		return nil, err
	}

	_ = a.compressWriter.Close()

	data := a.compressBuf.Bytes()
	a.compressBuf.Reset()

	return data, nil
}

func (a *AesCbc) Decrypt(src []byte) ([]byte, error) {

	var err error
	if a.uncompressReader == nil {
		a.uncompressReader = bytes.NewReader(src)
	} else {
		a.uncompressReader.Reset(src)
	}

	if a.compressReader == nil {
		a.compressReader, err = zlib.NewReader(a.uncompressReader)
		if err != nil {
			_ = a.compressReader.Close()
			return nil, err
		}
	} else {
		a.compressReader.(zlib.Resetter).Reset(a.uncompressReader, nil)
	}

	_ = a.compressReader.Close()
	src, err = io.ReadAll(a.compressReader)
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
