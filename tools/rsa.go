package tools

import (
	"fmt"
	"github.com/simplexity-ckcclc/gochannel/api/common"
)

func SignAndPrintClickInfo() {
	clickInfo := "appKey=appkeyA&channelId=channelA&deviceId=bar&clickTime=1562238598132"
	priKey := "MIICdwIBADANBgkqhkiG9w0BAQEFAASCAmEwggJdAgEAAoGBAM9OBRM2aL/ZDEdWKyGaaxzekBfk6p1z6YtA6QdgNE1Hg0WyMEMgIWNMtdH0rvAhuK6OBhOeUKRgeSGOtXed6huwv3LJEP22JwLITbvUl4C9Ybdb1y9Bxmwqa9g8Gne+Xk3dBvvgI40st/CsXbBIF6zHObWDj6PFRCQa63t5+BRZAgMBAAECgYB1pqfGsZhdWQdY7RRpa8PijIVmqipk1cXznBEkeHr2aOGdinVNg0yvmHeQArfN3LV9i2jzdWP7Bi142A8xJdQYgXsWrFGNTyLHYLa3qTqEZOKPOyQsaBs+jT1mzMflNt6plXIAYhpTpuPAw4yAqoQetnLNu9PBajbwvex7z0vDAQJBAOHDdmrapbWpCPniUo6Y0H/6MpAL4Hk/a7KvJ8QLX2VE43h262DLmwOKlZO4Pym16P9mHKIsTFLf/BREo9hWkSkCQQDrEa3vc8TJhnOw5D5dwNB7r/r58NYHbNCkec086rnsklgD0g7mX+d4ZxPgpxCStcgq++NVpH/rB/KG1EHOs9+xAkEA1wB77MUvnPJO7xL/hne30LkooA//hdjFKxUt7MDb56iUbOvru3ILvXKkglqsJH/uVhQb3sILKb3P5kl8NBI1CQJAGEWgDxyGEkT2xyoaInYZUNwv7wTmJKggtwr4nTSjdAD8Y5CaB0GZ1f3WuJinhm6Mt5uAssQkjTEp4rAiB2TdUQJBANC5seVNZAR2A2wsOhc1+sRaDoNj+awyZHDE+pIAYYwJpbSX7N69khEpreI2EFxE7rbHd0ce6Sd2c3r6rdbSHbo="

	sig, e := common.SignPSSAndBase64Encoding(clickInfo, priKey)
	if e != nil {
		fmt.Println(e)
	}
	fmt.Println(sig)
}

func SignAndPrintNonce() {
	nonce := "abc"
	priKey := "MIICdwIBADANBgkqhkiG9w0BAQEFAASCAmEwggJdAgEAAoGBAJaxeQZ0W3/DJUJQy5qoV52UQ3Urxrv4yL0dhqBuBhPItFYa4zQxKdapvMnbaiHpLkIwPN1E5rlQDAAPwO2oNPrb9n0v8zoX4DyOFS9uXjjUBstEqgjA8MfWz5ChvlcuoEqDFcMYsAbBMTBOdjVw6XBPuely9ZCu6Zfj7LrSuGqZAgMBAAECgYEAhBDtnDaFqicAlGnyxowanUO/CwVemoaihvtFbXx/Xv9a7MuLq8YagMMpbU8aaLXPkLpt3Q3xlx8MJVGpJ59vLcMrcLDNWPXa9ywjDKKucyrsrYI56Oe2A4HsG1rL1bii8COSpVNuF9vxoYS52z9LHVyG/TS+0IssuNjEjfKG0sECQQDEfUSk7ovUsXd9eUy0zBhCaiJce+V7wreE7G+SIBjj93epTUQeXIb0tnShJDgAqaz8ZzKDBQkE3gribkmHy/6dAkEAxFVr+jAe+l8LfdoxZP8NJbtuD1EDpyuf7dXON3fyAE+/sRF21rmkX7Nr3bEfihJ1L0EgwDZk7OWBLcqlv5J9LQJBAJEef9tcf5PoOntGYlvJvUUYBCbQLs44Iries0x2PkvoUs2MznmqFtaYBw2YpW//4U5NnaXcyyt4HwvbLp2IEZUCQEgtUyXF3Q2UNWhN94y2iwHNFtgAo4QocIIB8O7JZKkiqEkTL4oe80PPdR8qB3s97+CwY7bmCFJiyQupjSeRVf0CQCWQKaJwNreLti5DSDhWNLZwPIkdl5rEcY8Iq7kvCORQPKZUhRqFMZmZRcvPMSnYYNTerwUMVNUYQRwuom/xvo0="

	sig, e := common.SignPSSAndBase64Encoding(nonce, priKey)
	if e != nil {
		fmt.Println(e)
	}
	fmt.Println(sig)
}
