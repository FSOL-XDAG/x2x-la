# x2x LogAnalyzer (aka x2x-la)
x2x LogAnalyzer is designed to work on the log files produced by xmrig2xdag (x2x) and extract relevant information to improve proxy monitoring.

## How to setup XMRig2XDAG ?

- Download 'XMRig2XDAG'(x2x) : https://github.com/XDagger/xmrig2xdag/releases
- Setup config.json to generate log (set proxy.log is a good choice) 

![image](https://user-images.githubusercontent.com/128682335/229854580-cd949f06-876d-4404-99df-6b531612c53a.png)
- run 'x2x'

## How to use x2x-la ?

- Copy 'x2x-la' executable inside x2x folder
- run it !

## How to get help ? 

```
x2x-la --help
```

## How to compile ?

```
go mod init x2x-la
go get github.com/fatih/color
go mod tidy
go build x2x-la.go
```
