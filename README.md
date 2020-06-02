# gotop builder

This program creates the files necessary to compile gotop with selected
extensions. You need:

1. The tagged version of gotop to compile (e.g., "v3.5.1")
2. One (or more) extensions to enable (e.g. "github.com/xxxserxxx/gotop-nvidia")
3. Go. Since gotop requires Go >= 1.14, that's what you'll need.

Run with `-h` to get help text, including a mostly complete copy of this 
file.

## Example

```
$ go run ./build.go -r v4.0.0 github.com/xxxserxxx/gotop-nvidia
$ go build -o gotop ./gotop.go
```

## Binaries

In the releases section are all-inclusive binaries for Linux.  These binaries have all extensions compiled into them:

- [NVidia GPU support](https://github.com/xxxserxxx/gotop-nvidia)
- [Remote server monitoring](https://github.com/xxxserxxx/gotop-remote)

