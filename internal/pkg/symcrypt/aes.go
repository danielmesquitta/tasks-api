package symcrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"

	"github.com/danielmesquitta/tasks-api/internal/config"
	"github.com/danielmesquitta/tasks-api/internal/domain/entity"
)

type AESCrypto struct {
	env *config.Env
}

func NewAESCrypto(env *config.Env) *AESCrypto {
	return &AESCrypto{
		env: env,
	}
}

// Decrypt decrypts given text in AES 256 CBC
func (c *AESCrypto) Decrypt(
	encrypted string,
) (plaintext string, err error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", entity.NewErr(err)
	}

	block, err := aes.NewCipher([]byte(c.env.CipherSecretKey))
	if err != nil {
		return "", entity.NewErr(err)
	}

	if len(ciphertext)%aes.BlockSize != 0 {
		return "", fmt.Errorf("block size cant be zero")
	}

	mode := cipher.NewCBCDecrypter(block, []byte(c.env.InitializationVector))
	mode.CryptBlocks(ciphertext, ciphertext)
	ciphertext = pkcs5UnPadding(ciphertext)

	return string(ciphertext), nil
}

// Encrypt encrypts given text in AES 256 CBC
func (c *AESCrypto) Encrypt(
	plaintext string,
) (encrypted string, err error) {
	var plainTextBlock []byte
	length := len(plaintext)

	if length%16 != 0 {
		extendBlock := 16 - (length % 16)
		plainTextBlock = make([]byte, length+extendBlock)
		copy(
			plainTextBlock[length:],
			bytes.Repeat([]byte{uint8(extendBlock)}, extendBlock),
		)
	} else {
		plainTextBlock = make([]byte, length)
	}

	copy(plainTextBlock, plaintext)
	block, err := aes.NewCipher([]byte(c.env.CipherSecretKey))
	if err != nil {
		return "", entity.NewErr(err)
	}

	ciphertext := make([]byte, len(plainTextBlock))
	mode := cipher.NewCBCEncrypter(block, []byte(c.env.InitializationVector))
	mode.CryptBlocks(ciphertext, plainTextBlock)

	str := base64.StdEncoding.EncodeToString(ciphertext)

	return str, nil
}

// pkcs5UnPadding  pads a certain blob of data with
// necessary data to be used in AES block cipher
func pkcs5UnPadding(src []byte) []byte {
	length := len(src)
	unPadding := int(src[length-1])

	return src[:(length - unPadding)]
}
