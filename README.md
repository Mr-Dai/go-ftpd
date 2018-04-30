# go-ftpd

Simple FTP server based on [github.com/goftp/server](https://github.com/goftp/server).

## Usage

```
$ go build
$ ./go-ftpd -h
Usage of ./go-ftpd:
  -a, --authdb string     LevelDB file for auth info of the FTP server (default "/tmp/go-ftpd/auth.db")
  -d, --datapath string   data directory for the FTP server (default "/tmp/go-ftpd/data")
  -n, --name string       name of the FTP server (default "my-ftpd")
  -p, --port int          port to listen to (default 21)
```

## Docker


