# goforeach
```
usage: foreach [<flags>] <execute>

run the command which you want by goroutine.

Flags:
  -h, --help       Show context-sensitive help (also try --help-long and --help-man).
  -c, --circle=1   The circle times of command run
  -n, --cpu_num=2  Sets the maximum number of CPUs that can be executing simultaneously
  -f, --fork=10    Specify the number of concurrent goroutine
  -s, --show       Show infos of localhost
      --version    Show application version.

Args:
  <execute>  The command to be executed
```
# usage
```
❯ seq 8 |goforeach 'echo #1'
3
8
4
2
6
1
5
7
```
```
❯ goforeach -c 3 'tcping 8.8.8.8 53'
8.8.8.8 port 53 open.
8.8.8.8 port 53 open.
8.8.8.8 port 53 open.
```
```
❯ echo '8.8.8.8 53'|goforeach 'tcping #1 #2'
8.8.8.8 port 53 open.
```
