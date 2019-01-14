package rsa

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"fmt"
)

func VerifySig(plainText string, publicKey string, sig string) (bool, error) {
	key, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		fmt.Println(err)
		return false, err
	}

	pubKey, err := x509.ParsePKIXPublicKey(key)
	if err != nil {
		fmt.Println(err)
		return false, err
	}

	pub := pubKey.(*rsa.PublicKey)
	hash := sha1.New()
	hash.Write([]byte(plainText))
	digest := hash.Sum(nil)

	decodedSig, _ := base64.StdEncoding.DecodeString(sig)
	if err = rsa.VerifyPSS(pub, crypto.SHA1, digest, decodedSig, nil); err != nil {
		fmt.Println("Verify signature not match.")
		return false, nil
	}
	return true, nil
}
