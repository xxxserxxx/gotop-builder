# gotop builder

This program creates the files necessary to compile gotop with selected
extensions. You need:

1. The tagged version of gotop to compile (e.g., "v3.5.1")
2. One (or more) extensions to enable (e.g. "github.com/xxxserxxx/gotop-nvidia")
3. Go. Since gotop requires Go >= 1.14, that's what you'll need.

Example:

```
$ go run ./build.go -r v3.5.1 github.com/xxxserxxx/gotop-nvidia
$ go build -o gotop ./gotop.go
$ sudo cp gotop /usr/local/bin
```

