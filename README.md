# myapp_server
服务于自己的测试app的后台

sh build.sh 编译bin文件
sh service_control.sh start 运行服务

## 开发说明

### API 文档生成

1. 安装 swaggo：
```bash
go install github.com/swaggo/swag/cmd/swag@latest
```
2. 生成文档：
```bash
swag init -g ./cmd/main.go
```  
3. 添加新API请遵循swagger流程
