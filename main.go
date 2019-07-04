package main

import (
	"flag"
	"github.com/simplexity-ckcclc/gochannel/api"
	"github.com/simplexity-ckcclc/gochannel/common"
	"github.com/simplexity-ckcclc/gochannel/match"
)

func main() {
	// sign click info
	//sig, e := apicommon.SignPSSAndBase64Encoding("appKey=foo&deviceId=bar&clickTime=1562238598132", "MIICdgIBADANBgkqhkiG9w0BAQEFAASCAmAwggJcAgEAAoGBALCHxsWlkhdXWznXyO8QLK+QR0ZexvD7vPv8jKwmNybYYr2aM6rxwJNq/iDMMTyL547leNUakkTI/mtLZeFru8xjnmxg4teJBtXwuQKYmw5NQ8prhCWCTDBh85gQw6pl+o983leUUXdCxT7zRc1r9SXAg1Oy3qOEg0Ncm7YAzoQvAgMBAAECgYAzaPdcbCGlpo0sxUkBRkadQnlfZw6s88NP53bYU7DQIUhwS04sxIb+57PmvVDBf0UKeo28Eiby3U4q1SRwh72DBm78E5roN6y9ukVKAdpx7tq7EYrTrkvje+z4O3A+63B3JRQyU/qnbtrkJs87TalAPe0f22Dx6QZpqyvBxrozAQJBAOZpPD0aHukG18Xg+xYnSXehJmQJjkh4iO0Iat3k7CXd7w0ukdc1XDcj78hLR+5tOI/5nDfiKRZaK1H7eynX3O8CQQDEIqvaZuOsPx9PhzLzySii7OVfoV1Y5aFlg7TuucyNii0QYgMPeqXuEEZRwBfQXK80tFiXHavl06TbQkysJkzBAkA5kuIygmxm3gbcszMKfhalhecJ6DldcoEEea36dFFtxN8O9CwNEpBQVvJ7ohP/R9tyXnTioeiSZUWd3rEP65iRAkBgE3fJUVM/YeBNjbW405XzUUX+tUXLsRiBaKXtttfrkX8HomtLXtH/LruzefxwVaaBk8I9rAwzVZxQx0ZVoaFBAkEAwkTRUrlzpVDrkHo5tQTLToB6XJ5b+vUNu+L8oN2tgyNjkMY9hTVWv4U3Zp9qSiQ7nE60ndUCwuONhS1deO5bTQ==")
	//if e != nil {
	//	fmt.Println(e)
	//}
	//fmt.Println(sig)

	// sign nonce
	//sig, e := apicommon.SignPSSAndBase64Encoding("", "MIICdwIBADANBgkqhkiG9w0BAQEFAASCAmEwggJdAgEAAoGBAJaxeQZ0W3/DJUJQy5qoV52UQ3Urxrv4yL0dhqBuBhPItFYa4zQxKdapvMnbaiHpLkIwPN1E5rlQDAAPwO2oNPrb9n0v8zoX4DyOFS9uXjjUBstEqgjA8MfWz5ChvlcuoEqDFcMYsAbBMTBOdjVw6XBPuely9ZCu6Zfj7LrSuGqZAgMBAAECgYEAhBDtnDaFqicAlGnyxowanUO/CwVemoaihvtFbXx/Xv9a7MuLq8YagMMpbU8aaLXPkLpt3Q3xlx8MJVGpJ59vLcMrcLDNWPXa9ywjDKKucyrsrYI56Oe2A4HsG1rL1bii8COSpVNuF9vxoYS52z9LHVyG/TS+0IssuNjEjfKG0sECQQDEfUSk7ovUsXd9eUy0zBhCaiJce+V7wreE7G+SIBjj93epTUQeXIb0tnShJDgAqaz8ZzKDBQkE3gribkmHy/6dAkEAxFVr+jAe+l8LfdoxZP8NJbtuD1EDpyuf7dXON3fyAE+/sRF21rmkX7Nr3bEfihJ1L0EgwDZk7OWBLcqlv5J9LQJBAJEef9tcf5PoOntGYlvJvUUYBCbQLs44Iries0x2PkvoUs2MznmqFtaYBw2YpW//4U5NnaXcyyt4HwvbLp2IEZUCQEgtUyXF3Q2UNWhN94y2iwHNFtgAo4QocIIB8O7JZKkiqEkTL4oe80PPdR8qB3s97+CwY7bmCFJiyQupjSeRVf0CQCWQKaJwNreLti5DSDhWNLZwPIkdl5rEcY8Iq7kvCORQPKZUhRqFMZmZRcvPMSnYYNTerwUMVNUYQRwuom/xvo0=")
	//if e != nil {
	//	fmt.Println(e)
	//}
	//fmt.Println(sig)

	var confPath string
	flag.StringVar(&confPath, "conf", "", "gochannel config file")

	flag.Parse()

	// load config
	var err error
	common.Conf, err = common.LoadConf(confPath)
	if err != nil {
		panic(err)
	}

	db, err := common.OpenDB(common.Conf.Core.DSN)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// start match-Server
	go match.Serve()

	// start api-server
	api.Serve()

}
