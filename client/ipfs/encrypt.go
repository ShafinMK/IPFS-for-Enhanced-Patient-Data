package ipfs

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

type EncryptedChunk struct {
	Data []byte
	IV   []byte
	Tag  []byte
}

// EncryptChunk applies AES-256-GCM with the per-chunk key (Algorithm 5, line 7).
func EncryptChunk(key, plaintext []byte) (*EncryptedChunk, error) {
	block, err := aes.NewCipher(key)
	if err != nil { return nil, err }
	gcm, err := cipher.NewGCM(block)
	if err != nil { return nil, err }
	iv := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, iv); err != nil { return nil, err }
	sealed := gcm.Seal(nil, iv, plaintext, nil)
	return &EncryptedChunk{
		Data: sealed[:len(sealed)-gcm.Overhead()],
		IV:   iv,
		Tag:  sealed[len(sealed)-gcm.Overhead():],
	}, nil
}

func DecryptChunk(key []byte, c *EncryptedChunk) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil { return nil, err }
	gcm, err := cipher.NewGCM(block)
	if err != nil { return nil, err }
	return gcm.Open(nil, c.IV, append(c.Data, c.Tag...), nil)
}