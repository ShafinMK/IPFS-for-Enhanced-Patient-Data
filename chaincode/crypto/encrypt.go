// Package main implements the encryption and key-management procedures
// described in Section 3.4 and Algorithm 5.
package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"

	"golang.org/x/crypto/hkdf"
)

// EncryptRecord applies AES-256-GCM with a fresh 256-bit data key (Algorithm 5, lines 14-16).
func EncryptRecord(plaintext []byte) (ciphertext, iv, tag, dataKey []byte, err error) {
	dataKey = make([]byte, 32)
	if _, err = io.ReadFull(rand.Reader, dataKey); err != nil { return }
	block, err := aes.NewCipher(dataKey)
	if err != nil { return }
	gcm, err := cipher.NewGCM(block)
	if err != nil { return }
	iv = make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, iv); err != nil { return }
	sealed := gcm.Seal(nil, iv, plaintext, nil)
	tag = sealed[len(sealed)-gcm.Overhead():]
	ciphertext = sealed[:len(sealed)-gcm.Overhead()]
	return
}

func DecryptRecord(ciphertext, iv, tag, dataKey []byte) ([]byte, error) {
	block, err := aes.NewCipher(dataKey)
	if err != nil { return nil, err }
	gcm, err := cipher.NewGCM(block)
	if err != nil { return nil, err }
	return gcm.Open(nil, iv, append(ciphertext, tag...), nil)
}

// WrapKey performs RSA-OAEP key wrapping for a single recipient (Algorithm 5, line 18).
func WrapKey(pub *rsa.PublicKey, dataKey []byte) ([]byte, error) {
	return rsa.EncryptOAEP(sha256.New(), rand.Reader, pub, dataKey, nil)
}

func UnwrapKey(priv *rsa.PrivateKey, wrapped []byte) ([]byte, error) {
	return rsa.DecryptOAEP(sha256.New(), rand.Reader, priv, wrapped, nil)
}

// DeriveChunkKey derives the per-chunk key CK_i = HKDF(FK, i) (Algorithm 5, line 6).
func DeriveChunkKey(fileKey []byte, index int) ([]byte, error) {
	if len(fileKey) < 16 { return nil, errors.New("file key too short") }
	r := hkdf.New(sha256.New, fileKey, nil, []byte{byte(index)})
	out := make([]byte, 32)
	if _, err := io.ReadFull(r, out); err != nil { return nil, err }
	return out, nil
}

// DeriveFileKey derives FK = HKDF(MK_P, salt || fileID) (Algorithm 5, line 3).
func DeriveFileKey(masterKey []byte, salt, fileID []byte) ([]byte, error) {
	r := hkdf.New(sha256.New, masterKey, salt, fileID)
	out := make([]byte, 32)
	if _, err := io.ReadFull(r, out); err != nil { return nil, err }
	return out, nil
}

// helpers used by the contract
func b64(b []byte) string { return base64.StdEncoding.EncodeToString(b) }
func unb64(s string) ([]byte, error) { return base64.StdEncoding.DecodeString(s) }