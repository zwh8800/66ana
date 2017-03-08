# API 定义

## 通用结构

### 返回值

```json
{
  "code": "OK",   // 英文纯大写，下划线分割
  "message": "操作成功",
  "data": null
}
```


## /working-room/count 

### 描述

获取当前抓取中的房间数量

### url

`https://BASE/working-room/count`

### method 

GET

### request

无

### response

```json
{
  "data": {
    "count": 537
  }
}
```

## /working-room/list

### 描述

获取当前正在抓取的房间详细信息

### url

`https://BASE/working-room/list`

### method 

GET

### request

- keyword: string 搜索关键字【暂时不做】
- offset: integer
- limit: integer

### response

```json
{
  "data": {
    "total": 32768,
    "workingRoomList": [
      {
        "id": 28, 
        "createdAt": "2017-02-06T22:51:18.130572+08:00", 
        "updatedAt": "2017-03-08T22:21:27.743961+08:00", 
        "deletedAt": null, 
        "rid": 522423, 
        "cateId": 2, 
        "name": "2017LCK 春季赛 KDM VS ROX", 
        "status": 1, 
        "thumb": "https://rpic.douyucdn.cn/a1703/08/21/522423_170308214144.jpg", 
        "avatar": "https://apic.douyucdn.cn/upload/avanew/face/201612/22/11/d78b969ee15585976403feee9f246d51_big.jpg?rltime", 
        "fansCount": 791307, 
        "onlineCount": 2479914, 
        "ownerName": "Riot 丶 LCK", 
        "weight": 4770000, 
        "lastLiveTime": "2017-02-22T03:10:00+08:00"
      }
    ]
  }
}
```
