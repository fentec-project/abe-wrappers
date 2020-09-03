package main

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

import "C"

import (
	"encoding/hex"
	"fmt"
	"github.com/fentec-project/gofe/abe"
	//"github.com/serialization"
	"strings"
)

// This is an adaptor for the fame ABE scheme on github.com/fentec-project/gofe/abe
// in order to save and load keys from files

//export GenerateMasterKeys
func GenerateMasterKeys(path string, debug string) *C.char {
	// create a new FAME struct with the universe of attributes
	// denoted by integer
	printMsg("GenerateMasterKeys: path: "+path+" debug: "+debug, debug)
	printMsg("GenerateMasterKeys: init", debug)
	a := abe.NewFAME()
	err := WriteGob_pn(path, "fame.gob", a)
	if err != nil {
		printMsg("ha habido un error en Generate master keys",debug)
		fmt.Println(err)
		fmt.Errorf("Failed to write fame file: %v", err)
		output := fmt.Sprintf("%s: %s", "Error", err)
		return C.CString(output)
	}

	printMsg("GenerateMasterKeys: write fame ok", debug)

	pubKey, secKey, err := a.GenerateMasterKeys()
	if err != nil {
		fmt.Errorf("Failed to generate master key: %v", err)
		output := fmt.Sprintf("%s: %s", "Error", err)
		return C.CString(output)

	}
	err = WriteGob_pn(path, "publicKey.gob", pubKey)
if err != nil {
		fmt.Errorf("Failed to write PK file: %v", err)
		output := fmt.Sprintf("%s: %s", "Error", err)
		return C.CString(output)
	}				
	printMsg("GenerateMasterKeys: write pubkey ok", debug)

	WriteGob_pn(path, "secretKey.gob", secKey)
	if err != nil {
		fmt.Errorf("Failed to write MK file: %v", err)
		output := fmt.Sprintf("%s: %s", "Error", err)
		return C.CString(output)
	}			
	printMsg("GenerateMasterKeys: write secKey ok", debug)

	return C.CString("ok")
}

//export Encrypt
func Encrypt(path string, msg string, policy string, debug string) *C.char {

	printMsg("Encrypt: policy: "+policy, debug)
	msp, err := abe.BooleanToMSP(policy, false)
	if err != nil {
		fmt.Errorf("Failed to generate MSP: %v", err)
		output := fmt.Sprintf("%s: %s", "Error", err)
		return C.CString(output)
	}
	printMsg("Encrypt 1", debug)
	a := new(abe.FAME)
	ReadGob_pn(path, "fame.gob", a)
	pk := new(abe.FAMEPubKey)
	ReadGob_pn(path, "publicKey.gob", pk)

	printMsg("Encrypt 2", debug)

	cipher, err := a.Encrypt(msg, msp, pk)
	if err != nil {
		fmt.Errorf("Failed to encrypt: %v", err)
		output := fmt.Sprintf("%s: %s", "Error", err)
		return C.CString(output)
	}

	cypherS := WriteTString(cipher)
	s := WriteTCharA(cipher)
	printMsg("Encrypt result in golang: ", debug)
	printMsg(cypherS, debug)
	printMsg("Encrypt end in golang", debug)
	sEncoded := hex.EncodeToString(s)
	printMsg("Encrypted encoded: ", debug)
	printMsg(sEncoded, debug)
	printMsg("Encrypted end in golang", debug)

	return C.CString(sEncoded)

}

//export GenerateAttribKeys
func GenerateAttribKeys(path string, pathGen string, gamma string, debug string) *C.char {
	gammastr := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(gamma)), ","), "[]")

	printMsg("generateAttribKeys gamma: "+gammastr, debug)
	a := new(abe.FAME)
	ReadGob_pn(path, "fame.gob", a)
	printMsg("GenerateAttribKeys 1", debug)

	sk := new(abe.FAMESecKey)
	ReadGob_pn(path, "secretKey.gob", sk)
	printMsg("GenerateAttribKeys 2", debug)
	
	gammaSlice := strings.Split(gamma, ",")
	keys, err := a.GenerateAttribKeys(gammaSlice, sk)
	if err != nil {
		fmt.Errorf("Failed to generate keys: %v", err)
		output := fmt.Sprintf("%s: %s", "Error", err)
		return C.CString(output)
	}

	printMsg("GenerateAttribKeys 2.5", debug)
	WriteGob_pn(pathGen, "genKey.gob", keys)
	printMsg("GenerateAttribKeys 3", debug)

	return C.CString("ok")
}

// Decrypt takes as an input a cipher and an FAMEAttribKeys and tries to decrypt
// the cipher. This is possible only if the set of possessed attributes (and
// corresponding keys FAMEAttribKeys) suffices the encryption policy of the
// cipher. If this is not possible, an error is returned.

//export Decrypt
func Decrypt(path string, pathGen string, cypherS string, debug string) *C.char {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("into revover")
		}
        
    	}()
	printMsg("fame decrypt 0",debug)
	
	a := new(abe.FAME)
	err := ReadGob_pn(path, "fame.gob", a)
	
	
	if err != nil {
		printMsg("ha habido un error en retrieving FAME",debug)
		fmt.Println(err)
		fmt.Errorf("Failed to write fame file: %v", err)
		output := fmt.Sprintf("%s: %s", "Error", err)
		return C.CString(output)
	}
	printMsg("decrypt: cypherS:\n"+cypherS, debug)

	cipherDecoded, err := hex.DecodeString(cypherS)
	if err != nil {
		fmt.Errorf("Failed to decode b64: %v", err)
		output := fmt.Sprintf("%s: %s", "Error", err)
		return C.CString(output)
	}
	cipher := new(abe.FAMECipher)
	ReadFCharA(cipherDecoded, cipher)

	key := new(abe.FAMEAttribKeys)
	err = ReadGob_pn(pathGen, "genKey.gob", key)
	if err != nil {
		printMsg("ha habido un error en retrieving FAME",debug)
		fmt.Println(err)
		fmt.Errorf("Failed to write fame file: %v", err)
		output := fmt.Sprintf("%s: %s", "Error", err)
		return C.CString(output)
	}
	pk := new(abe.FAMEPubKey)
	err = ReadGob_pn(path, "publicKey.gob", pk)
	if err != nil {
		printMsg("ha habido un error en retrieving FAME",debug)
		fmt.Println(err)
		fmt.Errorf("Failed to write fame file: %v", err)
		output := fmt.Sprintf("%s: %s", "Error", err)
		return C.CString(output)
	}
	
	msgCheck, err := a.Decrypt(cipher, key, pk)
	if err != nil {
		fmt.Errorf("Failed to decrypt: %v", err)
		output := fmt.Sprintf("%s: %s", "Error", err)
		return C.CString(output)
	}

	return C.CString(msgCheck)
}

func printMsg(msg string, debug string) {
	if debug == "ok" {
		fmt.Println(msg + "\n")
	}
}

func main() {}

