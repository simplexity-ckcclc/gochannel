# GoChannel

这是穷人版的渠道匹配，简化了后台管理接口及第三方平台对接接口，仅保留逻辑。

### gochannel-api
api server，接收点击上报信息，校验，入库
1. 点击上报接口 
```
http://${host}:8480/ad/click?appKey=${appkey}&channelId=${channel}&deviceId=${deviceId}&clickTime=${clickTime}&sig=${sig}
```
|字段|含义|
|--|--|
|appKey | 应用APP|
|channel | 渠道标识|
|deviceId | imei或者idfa|
|clickTime | 点击时间，13位|
|sig | 对“appKey=${appkey}&channelId=${channel}&deviceId=${deviceId}&clickTime=${clickTime}”使用该channel对应对RSA私钥签名，然后进行base64编码|

注：请求需要url encode，返回结果格式为(非0为错误)
```json
{
    "code": 0,
    "message": "Success"
}
``` 

2. 注册channel（内部接口，管理员权限）
```
 http://localhost:8480/internal/channel/:channel/register?nonce=${nonce}&sig=${sig}
```
|字段|含义|
|--|--|
|nonce | 调用方随机生成，防止重放攻击|
|channel | 渠道标识|
|sig | 对nonce字段使用管理员RSA私钥签名，然后进行base64编码|

注：请求需要url encode，返回结果格式为(非0为错误)。其中，公钥PublicKey和私钥PrivateKey均经过base64编码。
```json
{
    "code": 0,
    "message": "Success",
    "data": {
        "AppKey": "appkeyA",
        "ChannelId": "channelB",
        "PublicKey": "MIGJAoGBAK5CmRtblI8gV9QPJ/8zloQYvaiDgDhwU288OfP4ysxDJEHFqw1JzLIlVuy/iHVzMYHERM5iOAyQwYkllY7eJMBCbs6XvtaasLtwOKStv6h5kFPIaTfEk6UIK8SPKD9VzFwTh09T5BDhZeMfhwgiBJhvz8OT/H2Qp8beYLbHgp89AgMBAAE=",
        "PrivateKey": "MIICdwIBADANBgkqhkiG9w0BAQEFAASCAmEwggJdAgEAAoGBAK5CmRtblI8gV9QPJ/8zloQYvaiDgDhwU288OfP4ysxDJEHFqw1JzLIlVuy/iHVzMYHERM5iOAyQwYkllY7eJMBCbs6XvtaasLtwOKStv6h5kFPIaTfEk6UIK8SPKD9VzFwTh09T5BDhZeMfhwgiBJhvz8OT/H2Qp8beYLbHgp89AgMBAAECgYEAj8EA7UCvXSMhUR7vr+eu02pVix5wOB7xtWHJrSogokEBOAEJCv1Gj++dtdCLkdhljteKq3b7JeKExc7rgeBgD506bXxWoczVUOVCflo5MTEu+CrwKTelRu38d8u0J3fz6XCiF2Z3d1T3icfKJOa8hmEvmTZG7RiZLQXl8ubzjAECQQDkIYmfI8RbhQbeIlpAh8zJk3erShxFIxFGdEUZLswKUHRwJNEDicfr0fh+/okLPGCl1eNWwN8PoydTs0zL6sCxAkEAw4xUF2vhXSOIp4L9/44FdX+oGBTvhNN0Ep1VJZwsGNIw9FHdDTucnhP119InIASwkIUnM06x+iHD1gvsYkvKTQJAOVBXs/yXa2rLY+l7hTTY9VewO/99hL0frPSvG3mPV5QI/NezD1GBQbTZ2oX2RjVgDqni2LvSDqqtybCFPcH3sQJADSM5ZxVtX6eKf86SKAAvp7Q649tMOD1ImBOP6+XxJH3Cojd5xXDS1d/7bVOGI2WNQzhe6NiJpodsM847RGNZJQJBAIltELnDCC/5rIMUvJtvKOrGx+6XsWTs4JMnrsbE2PEGr+OEqkHTY5qUpX/DgJ7T87l1Mf8jzcjCum28BbW5Vwc="
    }
}
```

3. 撤销channel（内部接口，管理员权限）s
```
 http://localhost:8480/internal/channel/:channel/evict?nonce=${nonce}&sig=${sig}
```
字段含义及要求同注册channel。