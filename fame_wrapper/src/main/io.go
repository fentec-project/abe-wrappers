/*
 * Copyright (c) 2019 ATOS
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
 
 package main

import (
	"bytes"
	"fmt"
	"encoding/gob"
	"log"
	"os"
)

func WriteGob_pn(filePath string, fileName string, object interface{}) error {
	os.MkdirAll(filePath, 0777)
	file, err := os.Create(filePath + "/" + fileName)
	if err == nil {
		encoder := gob.NewEncoder(file)
		err = encoder.Encode(object)
		if err != nil {
			fmt.Errorf("Error writting file: %v", err)
		}
	}
	file.Sync()
	file.Close()

	return err
}

func WriteGob_p(filePath string, object interface{}) error {
	file, err := os.Create(filePath)
	if err == nil {
		encoder := gob.NewEncoder(file)
		err = encoder.Encode(object)
                if err != nil {
                        fmt.Errorf("Error writting file: %v", err)
                }

	}
	file.Close()

	return err
}

func ReadGob_pn(filePath string, fileName string, object interface{}) error {
	file, err := os.Open(filePath + "/" + fileName)
	if err == nil {
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(object)
	}
	file.Close()

	return err
}

func ReadGob_p(filePath string, object interface{}) error {
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
