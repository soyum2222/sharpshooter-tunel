package crypto

import (
	"fmt"
	"testing"
)

func TestAesCbc_Decrypt(t *testing.T) {

	cbc := AesCbc{
		Key:    "123",
		KenLen: 16,
	}

	b, err := cbc.Encrypt([]byte("123"))
	fmt.Println(err)
	fmt.Println(b)

	b, err = cbc.Decrypt(b)
	fmt.Println(err)
	fmt.Println(string(b))

}
