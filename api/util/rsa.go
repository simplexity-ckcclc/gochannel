package util

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

    ds, _ := base64.StdEncoding.DecodeString(sig)
    if err = rsa.VerifyPKCS1v15(pub, crypto.SHA1, digest, ds); err != nil {
        fmt.Println("Verify signature not match.")
        return false, nil
    }
    return true, nil
}
