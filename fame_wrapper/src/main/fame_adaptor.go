package main

/**
*
*FENTEC Functional Encryption Technologies
*Privacy-preserving and Auditable Digital Currency Use Case
*Copyright © 2019 Atos Spain SA
*
*This program is free software: you can redistribute it and/or modify
*it under the terms of the GNU General Public License as published by
*the Free Software Foundation, either version 3 of the License, or
*(at your option) any later version.
*
*This program is distributed in the hope that it will be useful,
*but WITHOUT ANY WARRANTY; without even the implied warranty of
*MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
*GNU General Public License for more details.
*
*You should have received a copy of the GNU General Public License
*along with this program.  If not, see <http://www.gnu.org/licenses/>.
**/

import "C"

import (
	"encoding/base64"
	"fmt"
	"github.com/fentec-project/gofe/abe"
	//"github.com/serialization"
	"strings"
)

// This is an adaptor for the fame ABE scheme on github.com/fentec-project/gofe/abe
// in order to save and load keys from files

//export GenerateMasterKeys
func GenerateMasterKeys(path string, ownerName string, name string, debug string) *C.char {
	// create a new FAME struct with the universe of attributes
	// denoted by integer
	printMsg("GenerateMasterKeys: path: "+path+" ownerName: "+ownerName+" name: "+name+" debug: "+debug, debug)
	printMsg("GenerateMasterKeys: init", debug)
	a := abe.NewFAME()
	WriteGob(path+"fame"+ownerName+"_"+name+".gob", a)

	printMsg("GenerateMasterKeys: write fame ok", debug)

	pubKey, secKey, err := a.GenerateMasterKeys()
	if err != nil {
		fmt.Errorf("Failed to generate master key: %v", err)
		output := fmt.Sprintf("%s: %s", "Error", err)
		return C.CString(output)

	}
	WriteGob(path+"pubParams"+ownerName+"_"+name+".gob", pubKey)
	printMsg("GenerateMasterKeys: write pubkey ok", debug)

	WriteGob(path+"masterKey"+ownerName+"_"+name+".gob", secKey)
	printMsg("GenerateMasterKeys: write secKey ok", debug)

	return C.CString("ok")
}

//export Encrypt
func Encrypt(path string, ownerName string, name string, msg string, policy string, debug string) *C.char {

	printMsg("Encrypt: policy: "+policy, debug)
	msp, err := abe.BooleanToMSP(policy, false)
	if err != nil {
		fmt.Errorf("Failed to generate MSP: %v", err)
		output := fmt.Sprintf("%s: %s", "Error", err)
		return C.CString(output)
	}
	printMsg("Encrypt 1", debug)
	a := new(abe.FAME)
	ReadGob(path+"fame"+ownerName+"_"+name+".gob", a)
	pk := new(abe.FAMEPubKey)
	ReadGob(path+"pubParams"+ownerName+"_"+name+".gob", pk)

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
	sEncoded := base64.StdEncoding.EncodeToString(s)
	printMsg("Encrypted encoded: ", debug)
	printMsg(sEncoded, debug)
	printMsg("Encrypted end in golang", debug)

	return C.CString(sEncoded)

}

//export GenerateAttribKeys
func GenerateAttribKeys(path string, ownerName string, name string, gamma []int, debug string) *C.char {
	gammastr := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(gamma)), ","), "[]")

	printMsg("generateAttribKeys gamma: "+gammastr, debug)
	a := new(abe.FAME)
	ReadGob(path+"fame"+ownerName+"_"+name+".gob", a)
	printMsg("GenerateAttribKeys 1", debug)

	sk := new(abe.FAMESecKey)
	ReadGob(path+"masterKey"+ownerName+"_"+name+".gob", sk)
	printMsg("GenerateAttribKeys 2", debug)

	keys, err := a.GenerateAttribKeys(gamma, sk)
	if err != nil {
		fmt.Errorf("Failed to generate keys: %v", err)
		output := fmt.Sprintf("%s: %s", "Error", err)
		return C.CString(output)
	}

	WriteGob(path+ownerName+name+".gob", keys)
	printMsg("GenerateAttribKeys 3", debug)

	return C.CString("ok")
}

// Decrypt takes as an input a cipher and an FAMEAttribKeys and tries to decrypt
// the cipher. This is possible only if the set of possessed attributes (and
// corresponding keys FAMEAttribKeys) suffices the encryption policy of the
// cipher. If this is not possible, an error is returned.

//export Decrypt
func Decrypt(path string, ownerName string, name string, cypherS string, debug string) *C.char {
	a := new(abe.FAME)
	ReadGob(path+"fame"+ownerName+"_"+name+".gob", a)
	printMsg("Decrypt 1 \n", debug)

	cipher := new(abe.FAMECipher)
	cipherDecoded, err := base64.StdEncoding.DecodeString(cypherS)

	if err != nil {
		fmt.Errorf("Failed to decode b64: %v", err)
		output := fmt.Sprintf("%s: %s", "Error", err)
		return C.CString(output)
	}
	ReadFCharA(cipherDecoded, cipher)
	printMsg("Decrypt 2 \n", debug)

	key := new(abe.FAMEAttribKeys)
	ReadGob(path+ownerName+name+".gob", key)
	printMsg("Decrypt 3 \n", debug)

	pk := new(abe.FAMEPubKey)
	ReadGob(path+"pubParams"+ownerName+"_"+name+".gob", pk)
	printMsg("Decrypt 4 \n", debug)

	msgCheck, err := a.Decrypt(cipher, key, pk)
	printMsg("Decrypt 5 \n", debug)

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
