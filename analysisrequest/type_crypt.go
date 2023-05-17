package analysisrequest

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"io"
	"strings"
)

const keyphrase = "leodido"

func hashMD5(input string) string {
	b := []byte(input)
	h := md5.Sum(b)

	return hex.EncodeToString(h[:]) // by referring to it as a string
}

func gcm(key string) (cipher.AEAD, error) {
	block, err := aes.NewCipher([]byte(hashMD5(key)))
	if err != nil {
		return nil, err
	}

	return cipher.NewGCM(block)
}

func enc(value []byte, key string) ([]byte, error) {
	// Create a new cipher with a nonce
	// Using the Galois Counter Mode
	gcm, err := gcm(key)
	if err != nil {
		return nil, err
	}

	// Create the nonce
	nonce := make([]byte, gcm.NonceSize())
	_, _ = io.ReadFull(rand.Reader, nonce)

	// Encrypt using the nonce
	return gcm.Seal(nonce, nonce, value, nil), nil
}

func dec(value []byte, key string) ([]byte, error) {
	gcm, err := gcm(key)
	if err != nil {
		return nil, err
	}

	size := gcm.NonceSize()
	nonce, text := value[:size], value[size:]

	orig, err := gcm.Open(nil, nonce, text, nil)
	if err != nil {
		return nil, err
	}

	return orig, nil
}

func (t Type) Encrypt(also ...string) ([]byte, error) {
	// FIXME: not sure we wanna do this or not
	// c := t.Components()
	// if c.Parent != nil {
	// 	return nil, fmt.Errorf("enrichers are not meant to create verdicts but only to enrich them, thus there is no need to encrypt them")
	// }

	text := t.String()
	for _, more := range also {
		if more != "" && !strings.HasPrefix(more, "@") {
			text += "@" + more
		}
	}

	return enc([]byte(text), keyphrase)
}

func Decrypt(value []byte) ([]byte, error) {
	return dec(value, keyphrase)
}
