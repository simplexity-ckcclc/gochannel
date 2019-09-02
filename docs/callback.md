
### Callback，回传匹配结果
1. 回传接口 
```
POST ${callback_url}
```
请求Body示例
```json
{
    "app_key": "appkeyA",
    "channel": "channelA",
    "device_id": "1234567890", 
    "click_time": 1567353600000,
    "activate_time": 1567353610000
}
```
|字段|含义|
|--|--|
|app_key | 应用APP|
|channel | 渠道标识|
|device_id | imei或者idfa|
|click_time | 点击时间，13位|
|activate_time | 激活时间，13位|

回复Body示例
```json
{
    "code": 0,
    "message": "Success"
}
``` 

|code|含义|
|--|--|
|0 | 成功|
|10000 | 参数错误|