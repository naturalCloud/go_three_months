## 目录
```text
-db
--student.sqllist 数据文件
-model
--student 模型文件
errorHandle_test.go 测试文件
 
```

`sql.ErrNoRows` 表示查不到数据,不应返回上层调用方,返回空切片