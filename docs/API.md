
### API Server，接收点击上报信息，校验，入库
1. 点击上报接口 
```
POST http://${host}:8480/ad/click?appKey=${appkey}&channelId=${channel}&deviceId=${deviceId}&clickTime=${clickTime}&sig=${sig}
```
|字段|含义|
|--|--|
|appKey | 应用APP|
|channel | 渠道标识|
|deviceId | imei或者idfa|
|clickTime | 点击时间，13位|
|sig | 对“${appkey}-${channel}-${deviceId}-${clickTime}”使用该（appKey， channel）对应的RSA私钥签名，然后进行base64编码|

注：请求需要url encode，返回结果格式为(非0为错误)
```json
{
    "code": 0,
    "message": "Success"
}
``` 

2. 注册channel（内部接口，管理员权限）
```
 POST http://localhost:8480/internal/channel/register?nonce=${nonce}&sig=${sig}
```
|字段|含义|
|--|--|
|nonce | 调用方随机生成，防止重放攻击|
|sig | 对nonce字段使用管理员RSA私钥签名，然后进行base64编码|

请求示例
```json
{
    "app_key": "appkeyA",   // APP标志
    "channel": "channelA",  // 渠道标志
    "channel_type": "ios",  // 渠道类型，ios， android，etc.
    "callback_url": "http://localhost:7810/channel-callback"    // 匹配结果回调url
}
```

注：请求需要url encode，返回结果格式为(非0为错误)。其中，公钥PublicKey和私钥PrivateKey均经过base64编码。
```json
{
    "code": 0,
    "message": "Success",
    "data": {
        "app_key": "appkeyA",
        "channel": "channelB",
        "pub_key": "MIGJAoGBAK5CmRtblI8gV9QPJ/8zloQYvaiDgDhwU288OfP4ysxDJEHFqw1JzLIlVuy/iHVzMYHERM5iOAyQwYkllY7eJMBCbs6XvtaasLtwOKStv6h5kFPIaTfEk6UIK8SPKD9VzFwTh09T5BDhZeMfhwgiBJhvz8OT/H2Qp8beYLbHgp89AgMBAAE=",
        "pri_key": "MIICdwIBADANBgkqhkiG9w0BAQEFAASCAmEwggJdAgEAAoGBAK5CmRtblI8gV9QPJ/8zloQYvaiDgDhwU288OfP4ysxDJEHFqw1JzLIlVuy/iHVzMYHERM5iOAyQwYkllY7eJMBCbs6XvtaasLtwOKStv6h5kFPIaTfEk6UIK8SPKD9VzFwTh09T5BDhZeMfhwgiBJhvz8OT/H2Qp8beYLbHgp89AgMBAAECgYEAj8EA7UCvXSMhUR7vr+eu02pVix5wOB7xtWHJrSogokEBOAEJCv1Gj++dtdCLkdhljteKq3b7JeKExc7rgeBgD506bXxWoczVUOVCflo5MTEu+CrwKTelRu38d8u0J3fz6XCiF2Z3d1T3icfKJOa8hmEvmTZG7RiZLQXl8ubzjAECQQDkIYmfI8RbhQbeIlpAh8zJk3erShxFIxFGdEUZLswKUHRwJNEDicfr0fh+/okLPGCl1eNWwN8PoydTs0zL6sCxAkEAw4xUF2vhXSOIp4L9/44FdX+oGBTvhNN0Ep1VJZwsGNIw9FHdDTucnhP119InIASwkIUnM06x+iHD1gvsYkvKTQJAOVBXs/yXa2rLY+l7hTTY9VewO/99hL0frPSvG3mPV5QI/NezD1GBQbTZ2oX2RjVgDqni2LvSDqqtybCFPcH3sQJADSM5ZxVtX6eKf86SKAAvp7Q649tMOD1ImBOP6+XxJH3Cojd5xXDS1d/7bVOGI2WNQzhe6NiJpodsM847RGNZJQJBAIltELnDCC/5rIMUvJtvKOrGx+6XsWTs4JMnrsbE2PEGr+OEqkHTY5qUpX/DgJ7T87l1Mf8jzcjCum28BbW5Vwc="
    }
}
```

3. 撤销channel（内部接口，管理员权限）
```
 POST http://localhost:8480/internal/channell/evict?appKey=${appKey}&channel=${channel}&nonce=${nonce}&sig=${sig}
```
字段含义及要求同注册channel。