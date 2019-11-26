# 確認方法
- dd if=/dev/zero of=giga1 bs=1M count=1000
- dd if=/dev/zero of=giga2 bs=1M count=2000
でファイルを作っておいて、
```
go run main.go -s giga1
go run main.go -s giga2
```
で出力を見る

