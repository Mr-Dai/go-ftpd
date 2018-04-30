# go-ftpd

Simple FTP server based on [github.com/goftp/server](https://github.com/goftp/server).

## Install

```
$ go install github.com/Mr-Dai/go-ftpd
```

## Basic Usage

Run `go-ftpd run` to start the FTP server on port 21 (requires root permission on Linux).

```
$ go-ftpd help
NAME:
   go-ftpd - A Simple FTP Server

USAGE:
   go-ftpd [global options] command [command options] [arguments...]

VERSION:
   0.1.0

AUTHOR:
   Robert Peng <robert.peng@foxmail.com>

COMMANDS:
     run, s   Runs FTP server
     help, h  Shows a list of commands or help for one command
   Permission Management:
     chown  Changes file owner and group
     chmod  Changes file mode bits
   User/Group Management:
     user   Manages users of the FTP server
     group  Manages user groups of the FTP server

GLOBAL OPTIONS:
   --authdb value, -a value  path for auth DB of the FTP server (default: "/tmp/go-ftpd/auth.db")
   --help, -h                show help
   --version, -v             print the version
```

You can manage user/group and file permissions through other subcommands:

```
$ go-ftpd user help
NAME:
   go-ftpd user - Manages users of the FTP server

USAGE:
   go-ftpd user [global options] command [command options] [arguments...]

COMMANDS:
     list    Lists users of the FTP server
     add     Adds user to the FTP server
     del     Deletes user from the FTP server
     passwd  Updates user's password
```

```
$ go-ftpd group help
NAME:
   go-ftpd group - Manages user groups of the FTP server

USAGE:
   go-ftpd group [global options] command [command options] [arguments...]

COMMANDS:
     list     Lists user groups of the FTP server
     add      Adds user group to the FTP server
     del      Deletes user group from the FTP server
     adduser  Adds user to user group
     deluser  Deletes user from user group
```

```
$ go-ftpd chmod --help
NAME:
   go-ftpd chmod - Changes file mod bits

USAGE:
   go-ftpd chmod mode file
```

```
$ ./go-ftpd chown --help
NAME:
   go-ftpd chown - Changes file owner and group

USAGE:
   go-ftpd chown owner:[group] file
```

## Docker

## TODO

- [ ] Support for configuration file.
- [ ] Support for environment variables.


