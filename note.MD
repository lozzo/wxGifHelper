### 数据库错误
- 执行插入时报：`Incorrect string value: '\xF0\x9F\x8C\xB4' for column 'nickName' at row 1` 是因为`按照四个字节一个单位进行编码的，而我们通常使用的utf-8编码在mysql数据库中默认是按照3个字节一个单位进行编码的`
- 解决办法:`1c7cedb014`更改，数据库执行`ALTER TABLE wxUser CONVERT TO CHARACTER SET utf8mb4;`

### golang使用通道控制并发数
- 结束最好采用`defer` 来释放通道，见`97f82166c5`