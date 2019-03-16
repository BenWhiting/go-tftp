In-memory TFTP Server
=====================

This is a simple in-memory TFTP server, implemented in Go.  It is
RFC1350-compliant, but doesn't implement the additions in later RFCs.  In
particular, options are not recognized.

Usage
-----
TODO

### Directory structure
```
    /cmd - Main tfpt server deployment
    
    /internal - 

    /pkg - Library code that could be used by external projects
        /wire - Packet parsing library, that was provided by igneous.io 
```

Testing
-------
TODO

TODO: other relevant documentation
