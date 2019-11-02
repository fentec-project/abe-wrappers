package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"os"
)

func WriteGob(filePath string, object interface{}) error {
	file, err := os.Create(filePath)
	if err == nil {
		encoder := gob.NewEncoder(file)
		encoder.Encode(object)
	}
	file.Close()

	return err
}

func ReadGob(filePath string, object interface{}) error {
	file, err := os.Open(filePath)
	if err == nil {
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(object)
	}
	file.Close()

	return err
}

func WriteTString(object interface{}) string {
	var myB bytes.Buffer
	enc := gob.NewEncoder(&myB)
	err := enc.Encode(object)
	if err != nil {
		log.Fatal("encode error: ", err)
	}

	return myB.String()
}

func ReadFString(s string, object interface{}) {
	b := []byte(s)
	buf := bytes.NewBuffer(b)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(object)
	if err != nil {
		log.Fatal("decode error: ", err)
	}
}

func WriteTCharA(object interface{}) []byte {
	var myB bytes.Buffer
	enc := gob.NewEncoder(&myB)
	err := enc.Encode(object)
	if err != nil {
		log.Fatal("encode error: ", err)
	}

	return myB.Bytes()
}

func ReadFCharA(s []byte, object interface{}) {
	b := []byte(s)
	buf := bytes.NewBuffer(b)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(object)
	if err != nil {
		log.Fatal("decode error: ", err)
	}
}
