# TaskScheduler定时任务系统

## api接口

### 1.创建 /xtimer/creatTimer POST

创建请求

| 字段                   | 类型               | 是否必填 |
| ---------------------- | ------------------ | -------- |
| app                    | String             | 是       |
| name                   | String             | 是       |
| cron                   | String             | 是       |
| notifyHTTPParam        | NotifyHTTPParam    | 是       |
| NotifyHTTPParam.url    | String             | 是       |
| NotifyHTTPParam.method | String             | 是       |
| NotifyHTTPParam.header | Map<String,String> | 否       |
| NotifyHTTPParam.body   | String             | 否       |

```json
{
    "app": "est ad amet tempor proident",
    "name": "test1",
    "status": 1,
    "cron": "ex exercitation in sed",
    "notifyHTTPParam": "{\"method\": \"GET\", \"url\": \"https://golangcode.cn/\", \"header\": {\"Duis_c8\": \"culpa ex in\"}, \"body\": \"consectetur dolore consequat irure in\"}"
}
```

返回响应

```go
type CreateTimerReply struct {
	Code    int32                 `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Message string                `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	Data    *CreateTimerReplyData `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
}

// CreateTimerReplyData
type CreateTimerReplyData struct {
	TimerId int64 `protobuf:"varint,1,opt,name=timerId,proto3" json:"timerId,omitempty"`
}
```

### 2.激活 /xtimer/enableTimer GET

请求部分

| 字段 | 类型   | 是否必填 |
| ---- | ------ | -------- |
| id   | long   | 是       |
| app  | string | 是       |

返回响应

```go
type EnableTimerReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code    int32  `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Message string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
}
```

![image-20250308162737298](C:/Users/xz317/AppData/Roaming/Typora/typora-user-images/image-20250308162737298.png)

## 调度流程



## 迁移数据

