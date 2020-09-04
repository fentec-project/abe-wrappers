# abe-wrappers 

This repository contains an example about how to use FENTEC ABE libraries implemented in GOLANG from applications developed in other programming languages, by building such libraries into Shared Object (SO) Libraries.
 


### Before using the library
Please note that the FENTEC libraries are a work in progress and have not yet
reached a stable release. Code organization and APIs are **not stable**.
You can expect them to change at any point.
The current implementation of the wrapper has been developed for https://github.com/fentec-project/gofe/commit/2ce1030d0723d8fd1abe9c86b817a034039a088f

### Project structure:

The project contains the implementation of wrappers for the two ABE implemented schemes: **GPSW** and **FAME**. Each wrapper contains two different part implemented in GOLANG and Java.

* The GOLANG implementations:
    * files:
        ````bash
	    $PROJECT_HOME/gpsw_wrapper/src/main/gpsw_adaptor.go
	    $PROJECT_HOME/gpsw_wrapper/src/main/io.go
        $PROJECT_HOME/fame_wrapper/src/main/fame_adaptor.go
	    $PROJECT_HOME/fame_wrapper/src/main/io.go
        ````
    * extends the functionality of the ABE libraries, providing saving to disk functionality of cryptographic keys created.
	* provide an API for the ABE libraries.
	* Methods of this API are exported to build the SO libraries. In order to mark which methods must be exported, labels are added just before their declaration. Labels has the format of comments with the pattern "//export " plus the name of the method to export.

        ````go
        //export Encrypt
        func Encrypt(path string, msg string, policy string, debug string) *C.char {
        ````

	* Exported functions are:

		* GPSW:
			````go
			func GenerateMasterKeys(path string, numAtt int, debug string) *C.char
			func Encrypt(path string, msg string, gamma []int, debug string) *C.char
			func GeneratePolicyK(path string, path_user string, policy string, debug string) *C.char
			func Decrypt(path string, path_user string, cypherS string, debug string) *C.char
			````
		
		* FAME:
		    ````go
			func GenerateMasterKeys(path string, debug string) *C.char
			func Encrypt(path string, msg string, policy string, debug string) *C.char
			func GenerateAttribKeys(path string, pathGen string, gamma string, debug string) *C.char
			func Decrypt(path string, pathGen string, cypherS string, debug string) *C.char
		    ````
		
            Propety **path** is used to indicate the location to save Master and Private keys, while **path_user** and **pathGen** are used to indicate the location to save Generated keys (or decryption keys).
            The action of saving keys to files is performed at io.go. Methods implemented create the required directories if they do not exist, although this behaviour may change depending on the OOSS.


* The Java implementation use JNA (https://github.com/java-native-access/jna) to access the SO library.
    * files:
        ````bash
        $PROJECT_HOME/gpsw_wrapper/gpsw_wrapper/java_gpsw_test/src/main/java/gpsw/GpswTest1.java
        $PROJECT_HOME/gpsw_wrapper/gpsw_wrapper/java_gpsw_test/src/main/java/gpsw/Wrapper.java
        $PROJECT_HOME/fame_wrapper/java_fame_test/src/main/java/main/java/fame/FameTest.java
        $PROJECT_HOME/fame_wrapper/java_fame_test/src/main/java/main/java/fame/Wrapper.java
        ````
	    
    * Wrapper.java: implements the jna binding with the SO Library. This is:
        * an interface to declare data structures and methods provided by the SO library.
	    * API to invoke the SO library.
	* ___Test.java: are example of use of the wrapper.
	 

### Pre-requirements: 

* Donload and compile **gofe** project (https://github.com/fentec-project/gofe)


### How to run examples:

* Clone the project:  
	````bash
	cd $GOPATH/src/github.com
	git clone https://github.com/fentec-project/abe-wrappers.git
	````
		
* Compile libraries
	````bash
	cd $GOPATH/src/github.com/abe-wrappers/gpsw_wrapper/src/main
	go install .
	go build -o gpsw_adaptor.so -buildmode=c-shared gpsw_adaptor.go io.go
		
	cd $GOPATH/src/github.com/abe-wrappers/fame_wrapper/src/main
	go install .
	go build -o fame_adaptor.so -buildmode=c-shared fame_adaptor.go io.go
	````
	This will create SO libraries of both wrappers plus their header files.
		
	* Compile java wrappers and test classes
		````bash
		cd $GOPATH/src/github.com/abe-wrappers/gpsw_wrapper/java_gpsw_test
		mvn clean compile assembly:single
		
		cd $GOPATH/src/github.com/abe-wrappers/fame_wrapper/java_fame_test
		mvn clean compile assembly:single
		````
		
	* Run examples
		````bash
		cd $GOPATH/src/github.com/abe-wrappers/gpsw_wrapper/java_gpsw_test
		java -jar ./target/java_gpsw_test-0.0.1-SNAPSHOT-jar-with-dependencies.jar [message to encrypt]
		
		cd $GOPATH/src/github.com/abe-wrappers/fame_wrapper/java_fame_test_test
		java -jar ./target/java_fame_test-0.0.1-SNAPSHOT-jar-with-dependencies.jar [message to encrypt]

