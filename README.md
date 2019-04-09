# backup-and-sync

backup-and-sync is a simple wrapper application around restic and rclone. I wrote it to automate backups on my NAS.

At the moment the wrapper only includes the functionality needed for my use case.

## Installation

To install backup-and-sync from source run the following commands

``` bash
git clone https://github.com/th3noname/backup-and-sync.git

cd backup-and-sync

go run build.go
```

The application is automatically cross-compiled for windows and linux (386 and amd64). The binaries are stored in the bin directory.

If you're building outside of the GOPATH and have module support enabled the application should build without a problem. If you're building without module support you have to fetch the dependencies manually:

``` bash
# fetch all dependencies
go get -u ./src/...
```