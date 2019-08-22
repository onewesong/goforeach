# 依赖安装
```
go get -u -v github.com/akamensky/argparse
```
# 使用
1. 直接运行
```
> go run foreach.go -h
> echo -e 'a\nb\nc' | go run foreach.go -e "echo #1"
> ls |go run foreach.go -e "du -sh #1"
```
2. 编译后运行
```
> go build foreach.go
> ls |./foreach -e "du -sh #1"
```
