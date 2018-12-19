[toc]

## 最最最基础功能

> 这是一个可以将 Telegram 表情转发到微信的小机器人

- tg 表情转到出到 wx
- 可以设置表情包为组

## 基础功能

- tg 绑定微信帐号 （user 表 ，tg 帐号和 wx 帐号绑定）
- 图片处理为 gif 图片
- 小段 mp4 抽取成 gif
- 可保存别人的图片到本人表情包
- 设置私有表情包
- 可分享表情
- 可在微信小程序上传 gif 或者小段 mp4
- 随机浏览表情包
- 最热表情包

## 扩展功能

- 制作文字表情动图 前台跑（模板+gif.js）

## tg_bot 命令参考

```
/start_send - 开始发送，非表情包组
/stop_send - 结束发送，开始后台处理
/start_group - 开始发送表情包组
/bind_wx - 绑定微信号
/un_bind_wx - 解绑微信号
```

## 数据库结构

- 文件数据库 ID 主键
  - FileID ----> 文件 fileID
  - GroupID--->Group.ID
  - UserID--->User.UserID

| ID  | FileID           | GroupID | UserID |
| --- | ---------------- | ------- | ------ |
| 12  | d1jhkhd1U_D21SDH | 123     | 23123  |

- wxID 微信用户(wx 用户每次打开小程序的时候会查询用户名是否改变，如果改变则修改数据库？？？？)

| ID   | openID                       | NickName |
| ---- | ---------------------------- | -------- |
| 1223 | o2cb64in4yFf9YS-KMQcEYkiPtJ0 | lozzow   |

- tgID Telegram 用户
  |ID|
  |-----|
  |123|

- User 用户表,可扩展其他程序

| UserID | wxID | tgID |
| ------ | ---- | ---- |
| 12     | 33   | 44   |

- Group 表情包组表

| ID  | Name |
| --- | ---- |
| 123 | ha   |

## TgBot 部分挂在国外，数据库放在国内，数据上传到 oss

- 国外访问国内稳定性大于国内访问国外
- 国内只放置前端 API 服务，对带宽要求不高，随便找个羊毛机都可以跑
