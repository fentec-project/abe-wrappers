# fame_wrapper

This repository is a wrapper for the https://github.com/fentec-project/gofe/blob/master/abe/gpsw.go implementation.
It serialize and save cryptographic data structures into the file system.

In order to use this wrapper from applications developed in other programming languages different of GOLANG it can be compiled
into a shared object library:

````bash
$ go build -o gpsw_adaptor.so -buildmode=c-shared gpsw_adaptor.go io.go
````
