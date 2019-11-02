# abe-Wrappers

This repository contains wrappers for the chryptographic schemes implementations developed in FENTEC procject which can be found at: https://github.com/fentec-project/gofe/blob/master/abe/
They serialize and save cryptographic data structures into the file system.

In order to use these wrappers from applications developed in other programming languages different of GOLANG it can be compiled
into a shared object library:

*GPSW scheme (KP-ABE)
````bash
$ go build -o gpsw_adaptor.so -buildmode=c-shared gpsw_adaptor.go io.go
````
*FAME shcme (CP-ABE)
````bash
$ go build -o fame_adaptor.so -buildmode=c-shared fame_adaptor.go io.go
````


