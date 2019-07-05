package main

import (
	"flag"
	"github.com/simplexity-ckcclc/gochannel/api"
	"github.com/simplexity-ckcclc/gochannel/common"
	"github.com/simplexity-ckcclc/gochannel/match"
)

func main() {
	//sign click info
	//sig, e := apicommon.SignPSSAndBase64Encoding("appKey=appkeyA&channelId=channelA&deviceId=bar&clickTime=1562238598132", "MIICeAIBADANBgkqhkiG9w0BAQEFAASCAmIwggJeAgEAAoGBANfJ0gKwqWWePNLuKtlJS1VKzvMGPKXgwWaZNWP/l7R9L/QPyyoqJoxg3QddL/m/oVgRe9esewZ0xZfyxNXLfGZkIfgSBOLdz80T3zbCizFeu0lYG6m8lb0j4Zr4wFVE4p6b4RtRgJ8kpILhs46Y0VKljtE0MS4mSG6N4nHhAzqZAgMBAAECgYEAjVAHjedvJ7L2lhOOT/llshdpa1E8SkzjmnLeufvZt0L8MlJdc+FimS+dz4LBNka+PFRGy7iSYGn8NEOxj2jQr3Az3gyqvWmLsCnC0Td0sY98EhCQVd9jIs97QUlhhtoHLiEYnX8Wrr5SaabmAYTCVPwcWpBeO/+b7kcL9k5EvVkCQQDY8PPu68SGedRpcGA6ByQP5G8RS4WCulzLk/AcQtSMlVE8QhzWuNJr2AoGDERmB/6gB0KoWvg+tbainVltlwmbAkEA/qO7HK/o5/03Rrfs9aM34lQkLW+muhhxhSKJcTy2CVoyhJ9laTsT4G8uijkj4zEzrW60A2E5uRhnU/Cc8Oe52wJBAI/eS4cY5/3ecZVzJv2Umr/HWDj6ApKNkNiZRVUYpOiOZY82sPVdIH7QiOU14W5gwuXRqs0HdzXvQC1beGELFx0CQQCruHbira5/ZD/2rOpb7KovM1cCXR0uunUzt0rA1pRcUjtnPKcDBBgvbksQY+BTwkZ7WwCClvp6XH6yGL19qIepAkA+1xy7d9GPQ8u0UofVaL68ja3+wFEryHfl1aQCGV2nIrkDmirc7Y3YKo4sbSLAk7zmpDFu9O4IJh5BJXibyg/+")
	//if e != nil {
	//	fmt.Println(e)
	//}
	//fmt.Println(sig)
	//valid, _ := apicommon.VerifyBase64WithRSAPubKey("appKey=appkeyA&channelId=channelA&deviceId=bar&clickTime=1562238598132", "MIGJAoGBANfJ0gKwqWWePNLuKtlJS1VKzvMGPKXgwWaZNWP/l7R9L/QPyyoqJoxg3QddL/m/oVgRe9esewZ0xZfyxNXLfGZkIfgSBOLdz80T3zbCizFeu0lYG6m8lb0j4Zr4wFVE4p6b4RtRgJ8kpILhs46Y0VKljtE0MS4mSG6N4nHhAzqZAgMBAAE=", "iFRMwk83T0YM92tBNRqpIj0h269psZQd+BgMPVkkohzMlaSLylN2AcEhAy7ZV2cWQxUT3A6lrfmFxkBC919u+WsbsP6UgZ5HGOpu9s5+zaPU58hX2yte3zEiHfyV7cnjJAu2R1scshc0NGxNiTP2cQ3HI0u08/5+5MrHXSdjMU8=")
	//fmt.Println(valid)

	// sign nonce
	//sig, e := apicommon.SignPSSAndBase64Encoding("abc", "MIICdwIBADANBgkqhkiG9w0BAQEFAASCAmEwggJdAgEAAoGBAJaxeQZ0W3/DJUJQy5qoV52UQ3Urxrv4yL0dhqBuBhPItFYa4zQxKdapvMnbaiHpLkIwPN1E5rlQDAAPwO2oNPrb9n0v8zoX4DyOFS9uXjjUBstEqgjA8MfWz5ChvlcuoEqDFcMYsAbBMTBOdjVw6XBPuely9ZCu6Zfj7LrSuGqZAgMBAAECgYEAhBDtnDaFqicAlGnyxowanUO/CwVemoaihvtFbXx/Xv9a7MuLq8YagMMpbU8aaLXPkLpt3Q3xlx8MJVGpJ59vLcMrcLDNWPXa9ywjDKKucyrsrYI56Oe2A4HsG1rL1bii8COSpVNuF9vxoYS52z9LHVyG/TS+0IssuNjEjfKG0sECQQDEfUSk7ovUsXd9eUy0zBhCaiJce+V7wreE7G+SIBjj93epTUQeXIb0tnShJDgAqaz8ZzKDBQkE3gribkmHy/6dAkEAxFVr+jAe+l8LfdoxZP8NJbtuD1EDpyuf7dXON3fyAE+/sRF21rmkX7Nr3bEfihJ1L0EgwDZk7OWBLcqlv5J9LQJBAJEef9tcf5PoOntGYlvJvUUYBCbQLs44Iries0x2PkvoUs2MznmqFtaYBw2YpW//4U5NnaXcyyt4HwvbLp2IEZUCQEgtUyXF3Q2UNWhN94y2iwHNFtgAo4QocIIB8O7JZKkiqEkTL4oe80PPdR8qB3s97+CwY7bmCFJiyQupjSeRVf0CQCWQKaJwNreLti5DSDhWNLZwPIkdl5rEcY8Iq7kvCORQPKZUhRqFMZmZRcvPMSnYYNTerwUMVNUYQRwuom/xvo0=")
	//if e != nil {
	//	fmt.Println(e)
	//}
	//fmt.Println(sig)

	var confPath string
	flag.StringVar(&confPath, "conf", "", "gochannel config file")
	flag.Parse()

	// load config
	var err error
	if err = common.LoadConf(confPath); err != nil {
		panic(err)
	}

	// start match-Server
	go match.Serve()

	// start api-server
	api.Serve()

}
