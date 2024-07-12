# easy-go
go lang toolkit

## 依赖
- go get -u gopkg.in/natefinch/lumberjack.v2

## 运行
#### 作为socket原始服务端
```shell
 go run ./cmd/easy_entry.go --run srs --port 8282
```
#### 作为socket原始客户端
```shell
 go run ./cmd/easy_entry.go --run src --host 127.0.0.1 --port 8282 -nc 3
```

## 构建
#### PowerShell
```shell
 # 读取当前环境变量中的GOARCH与GOOS设置 [powershell]
 $env:GOARCH
 $env:GOOS 
```

```shell
 # 设置GOOS=windows [powershell]
 $env:GOOS = "windows"
```

```shell
 # 设置GOOS=linux [powershell]
 $env:GOOS = "linux"
```
#### 生成可执行文件

```shell
## 构建windows可执行文件:
 go build -o ./build/easy-go.exe ./cmd/easy_entry.go
```

```shell
## 构建linux可执行文件:
 go build -o ./build/easy-go ./cmd/easy_entry.go
```

#### 运行
```shell
 nohup ./excel-go > output.log 2>&1 &
```