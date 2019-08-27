# 依赖安装
```
go get -u -v gopkg.in/alecthomas/kingpin.v2
go get -u -v golang.org/x/crypto/ssh/terminal
```
# 使用
1. 直接运行
```
> go run foreach.go -h
> echo -e 'a\nb\nc' | go run foreach.go "echo #1"
> ls |go run foreach.go "du -sh #1"
```
2. 编译后运行
```
> go build foreach.go
> ls |./foreach "du -sh #1"
```
