package main

import (
	"flag"
	"github.com/simplexity-ckcclc/gochannel/api"
	"github.com/simplexity-ckcclc/gochannel/common"
	"github.com/simplexity-ckcclc/gochannel/match"
)

func main() {
	//sig, e := apicommon.SignPSSAndBase64Encoding("appKey=ab&deviceId=bar", "MIICeAIBADANBgkqhkiG9w0BAQEFAASCAmIwggJeAgEAAoGBAKPpAoCJ1+EUa8yehglcixFC5k5QbcA/SM1ilLoCvy9BZOKF3A6H2BArgyc1RHWeov3axjW7YNufCNPpTjckuypoITpv5UDMxV53atA1ezHp+13muX0qINX6j1XYMb0bSqcCpCYNCdefO+HLtnyitvCsltooGQWSf8byrJbYwB+RAgMBAAECgYEAg41Ew9NeHzjbmBt26laLCuyNmTc2DsD79lNzmKMRvKSYirHyvvrKL5gsqDA5ZMlQebu3r3JXN405cZLjgqCJUiOqAMhaE6f6klILNEusFKonk+jI5fAMslhkgdP41yU9T1FEmTk+KzNZYCeMgaxhJS+u7s1oX+JG2d878dGUM60CQQDETENRvvj5BiCzhzxxgqz3rwaARRAsx9V67x9tfPAWQG6uq1vseuv5nIcUDK9T5l07EuIz9gBZRnGbSAjOAvkjAkEA1cMK9DiNIhSbcoPMMHveXfDU3rq2LNd9Iy649R7va1ofujGxOTPy8ej8JZnchw8PjtV0YI7x6vQ2GX/4cSkBuwJAHc+7NZH8O82Lb9hs/Iws+pyxLw/OCg77Q+VG75jW2XpFlO9fUYXFiq4T8Z6Pjf1hUVRn2B5XJTfGjx+cfrUC4wJBAKgFFAIcxUppejoLwJ7HbmTWnOupRPKAOrtByV3agAQYpeGbl5rH64kcQb1ob/+05dy2iTAwi5TLeg6XGPgRbGsCQQDD1aqB8B5HtxN1f+iLuW2KfiIfpw1QjMNZLabW0a7y3Eg+dHL1JmSjtvYaN+2FPXiKh8iaZ0l4Kzso5Jvg9ZPK")
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
