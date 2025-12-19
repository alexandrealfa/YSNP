package internal

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type WalletFile struct {
	Version          int    `json:"version"`
	WalletName       string `json:"wallet_name"`
	Algorithm        string `json:"algorithm"`
	CreatedAt        string `json:"created_at"`
	PublicKeyDERB64  string `json:"public_key_der_b64"`
	PrivateKeyDERB64 string `json:"private_key_der_b64"`
}

func SaveWalletJSON(path, walletName string, priv *rsa.PrivateKey) error {
	// Private -> PKCS#8 DER
	privDER, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return err
	}

	// Public -> PKIX (SubjectPublicKeyInfo) DER
	pubDER, err := x509.MarshalPKIXPublicKey(&priv.PublicKey)
	if err != nil {
		return err
	}

	w := WalletFile{
		Version:          1,
		WalletName:       walletName,
		Algorithm:        "RSA-2048",
		CreatedAt:        time.Now().UTC().Format(time.RFC3339),
		PublicKeyDERB64:  base64.StdEncoding.EncodeToString(pubDER),
		PrivateKeyDERB64: base64.StdEncoding.EncodeToString(privDER),
	}

	data, err := json.MarshalIndent(w, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}

func LoadWalletJSON(path string) (*rsa.PrivateKey, *rsa.PublicKey, *WalletFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, nil, err
	}

	var w WalletFile
	if err := json.Unmarshal(data, &w); err != nil {
		return nil, nil, nil, err
	}

	privDER, err := base64.StdEncoding.DecodeString(w.PrivateKeyDERB64)
	if err != nil {
		return nil, nil, nil, err
	}
	pubDER, err := base64.StdEncoding.DecodeString(w.PublicKeyDERB64)
	if err != nil {
		return nil, nil, nil, err
	}

	// Parse private PKCS#8
	privAny, err := x509.ParsePKCS8PrivateKey(privDER)
	if err != nil {
		return nil, nil, nil, err
	}
	priv, ok := privAny.(*rsa.PrivateKey)
	if !ok {
		return nil, nil, nil, fmt.Errorf("private key não é RSA")
	}

	// Parse public PKIX
	pubAny, err := x509.ParsePKIXPublicKey(pubDER)
	if err != nil {
		return nil, nil, nil, err
	}
	pub, ok := pubAny.(*rsa.PublicKey)
	if !ok {
		return nil, nil, nil, fmt.Errorf("public key não é RSA")
	}

	return priv, pub, &w, nil
}
