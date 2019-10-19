package common

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"errors"
)

const _rsaKeySize = 1024

func VerifyBase64WithRSAPubKey(plainText string, publicKey string, sig string) (bool, error) {
	decodedPubKey, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return false, err
	}

	pubKey, err := x509.ParsePKCS1PublicKey([]byte(decodedPubKey))
	if err != nil {
		return false, err
	}

	digest := sha256.Sum256([]byte(plainText))
	decodedSig, _ := base64.StdEncoding.DecodeString(sig)
	if err = rsa.VerifyPSS(pubKey, crypto.SHA256, digest[:], decodedSig, nil); err != nil {
		return false, nil
	}
	return true, nil
}

func GenerateRSAKeyPair() (string, string, error) {
	priKeyI, err := rsa.GenerateKey(rand.Reader, _rsaKeySize)
	if err != nil {
		return "", "", err
	}

	priKeyBytes, err := x509.MarshalPKCS8PrivateKey(priKeyI)
	if err != nil {
		return "", "", err
	}

	// Extract Public Key from RSA Private Key
	pubKeyI := priKeyI.PublicKey
	pubKeyBytes := x509.MarshalPKCS1PublicKey(&pubKeyI)
	return base64.StdEncoding.EncodeToString(pubKeyBytes[:]), base64.StdEncoding.EncodeToString(priKeyBytes[:]), nil
}

func SignPSSAndBase64Encoding(plaintext string, privateKey string) (string, error) {
	decodedPriKey, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		return "", err
	}

	priKeyI, err := x509.ParsePKCS8PrivateKey([]byte(decodedPriKey))
	if err != nil {
		return "", err
	}
	priKey, ok := priKeyI.(*rsa.PrivateKey)
	if !ok {
		return "", errors.New("cannot cast private Key string to rsa private key")
	}

	hashed := sha256.Sum256([]byte(plaintext))
	signature, err := rsa.SignPSS(rand.Reader, priKey, crypto.SHA256, hashed[:], nil)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(signature), nil
}
