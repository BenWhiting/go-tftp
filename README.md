In-memory TFTP Server
=====================

This is a simple in-memory TFTP server, implemented in Go.  It is
RFC1350-compliant, but doesn't implement the additions in later RFCs.  In
particular, options are not recognized.

Usage
-----
### Build
Build Command
``` 
go build -o go-tftp.exe ./cmd 
```

### Run 
Command-Line Flags
```
Usage of go-tftp:
  -f int
        TFTP server flush period (Seconds). (default 1)
  -h string
        TFTP server address (default "127.0.0.1:69")
  -l int
        TFTP server log flag. [1-7] (default 2).
  -p string
        TFTP server request log file path. (default "./go-tftp.log")
  -r int
        TFTP server transfer retry time (Seconds). (default 30)
  -t int
        TFTP server transfer timeout (Seconds). (default 30)
```

### Directory Structure
This repository is structured as such
```
    /cmd - Main tfpt server deployment

    /internal - Utilities for this repository
        /constants - None changing variables for this repository
        /logger - Logger "constructors"

    /pkg - Library code that could be used by external projects
        /communication - Transfers of files
        /wire - Packet parsing library, that was provided by igneous.io 
```

Testing
-------
Run all unit tests
```
go test ./...
```