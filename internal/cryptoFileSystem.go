package internal

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
)

type Envelope struct {
	Alg       string `json:"alg"`        // "RSA-OAEP-SHA256"
	Enc       string `json:"enc"`        // "A256GCM"
	EncKeyB64 string `json:"enc_key"`    // RSA-encrypted AES key
	NonceB64  string `json:"nonce"`      // GCM nonce
	CipherB64 string `json:"ciphertext"` // AES-GCM ciphertext (includes auth tag)
}

func GenerateRsaKeys() (privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey, err error) {
	privateKey, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	publicKey = &privateKey.PublicKey

	return
}

func DecryptJSONWithPrivateKey(priv *rsa.PrivateKey, envelopeJSON []byte) ([]byte, error) {
	var env Envelope
	if err := json.Unmarshal(envelopeJSON, &env); err != nil {
		return nil, err
	}
	if env.Alg != "RSA-OAEP-SHA256" || env.Enc != "A256GCM" {
		return nil, errors.New("unsupported envelope algorithms")
	}

	encKey, err := base64.StdEncoding.DecodeString(env.EncKeyB64)
	if err != nil {
		return nil, err
	}
	nonce, err := base64.StdEncoding.DecodeString(env.NonceB64)
	if err != nil {
		return nil, err
	}
	ciphertext, err := base64.StdEncoding.DecodeString(env.CipherB64)
	if err != nil {
		return nil, err
	}

	// 1) RSA-OAEP decrypt AES key with PRIVATE key
	label := []byte("json-envelope-v1")
	aesKey, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, priv, encKey, label)
	if err != nil {
		return nil, err
	}

	// 2) AES-GCM decrypt
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	if len(nonce) != gcm.NonceSize() {
		return nil, errors.New("invalid nonce size")
	}

	var aad []byte = nil
	plain, err := gcm.Open(nil, nonce, ciphertext, aad)
	if err != nil {
		// Se falhar aqui, pode ser: chave errada, arquivo alterado, nonce errado, etc.
		return nil, err
	}

	return plain, nil
}

func EncryptJSONWithPublicKey(pub *rsa.PublicKey, jsonPlain []byte) ([]byte, error) {
	// 1) random AES-256 key
	aesKey := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, aesKey); err != nil {
		return nil, err
	}

	// 2) AES-GCM encrypt
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// AAD opcional: metadata autenticada (não criptografada). Pode deixar nil.
	var aad []byte = nil

	ciphertext := gcm.Seal(nil, nonce, jsonPlain, aad)

	// 3) RSA-OAEP encrypt AES key with PUBLIC key
	label := []byte("json-envelope-v1") // label precisa ser igual na decrypt
	encKey, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, pub, aesKey, label)
	if err != nil {
		return nil, err
	}

	env := Envelope{
		Alg:       "RSA-OAEP-SHA256",
		Enc:       "A256GCM",
		EncKeyB64: base64.StdEncoding.EncodeToString(encKey),
		NonceB64:  base64.StdEncoding.EncodeToString(nonce),
		CipherB64: base64.StdEncoding.EncodeToString(ciphertext),
	}

	return json.MarshalIndent(env, "", "  ")
}
