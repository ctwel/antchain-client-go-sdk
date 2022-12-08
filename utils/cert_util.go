package utils

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
)

func Sign(plain, priKey string) (string, error) {
	h := sha256.New()
	h.Write([]byte(plain))
	d := h.Sum(nil)

	privateKey, err := getPrivateKey(priKey)
	if err != nil {
		return "", err
	}
	sig, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, d)
	if err != nil {
		return "", fmt.Errorf("fail to SignPKCS1v15 plain:%+v priKey:%+v err:%+v", plain, priKey, err)
	}

	return fmt.Sprintf("%x", sig), nil
}

func getPrivateKey(priKey string) (*rsa.PrivateKey, error) {
	priv, err := ioutil.ReadFile(priKey)
	if err != nil {
		return nil, fmt.Errorf("fail to read priKey priKey:%+v err:%+v", priKey, err)
	}

	privPem, _ := pem.Decode(priv)
	if privPem == nil {
		return nil, fmt.Errorf("RSA private key is illegal,please check key")
	}
	var privPemBytes []byte
	if privPem.Type != "PRIVATE KEY" {
		return nil, fmt.Errorf("RSA private key is of the wrong type,actual Pem Type:%+v expect Pem Type:%+v", privPem.Type, "PRIVATE KEY")
	}
	privPemBytes = privPem.Bytes
	var parsedKey interface{}
	if parsedKey, err = x509.ParsePKCS1PrivateKey(privPemBytes); err != nil {
		if parsedKey, err = x509.ParsePKCS8PrivateKey(privPemBytes); err != nil { // note this returns type `interface{}`
			return nil, err
		}
	}
	var privateKey *rsa.PrivateKey
	var ok bool
	privateKey, ok = parsedKey.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("unable to convert to  RSA private key")
	}
	return privateKey, nil
}
