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
	//	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/fentec-project/gofe/abe"
	"github.com/fentec-project/gofe/data"
	"strings"
)

// This is an adaptor for the KP ABE scheme (https://eprint.iacr.org/2006/309.pdf) on github.com/fentec-project/gofe/abe
// in order to save and load keys from files

//export GenerateMasterKeys
func GenerateMasterKeys(path string, ownerName string, name string, numAtt int, debug string) *C.char {
	b := string(numAtt)

	printMsg("GenerateMasterKeys with num attr: "+b, debug)
	a := abe.NewGPSW(numAtt)
	WriteGob(path+"gpsw"+ownerName+"_"+name+".gob", a)

	printMsg("GenerateMasterKeys: write gpsw ok", debug)

	pubKey, secKey, err := a.GenerateMasterKeys()
	if err != nil {
		fmt.Errorf("Failed to generate master key: %v", err)
		output := fmt.Sprintf("%s: %s", "Error", err)
		return C.CString(output)

	}

	WriteGob(path+"pubgpsw"+ownerName+"_"+name+".gob", pubKey)
	printMsg("GenerateMasterKeys: write pubkey ok", debug)

	WriteGob(path+"masterKeygpsw"+ownerName+"_"+name+".gob", secKey)
	printMsg("GenerateMasterKeys: write secKey ok", debug)

	return C.CString("ok")
}

//export Encrypt
func Encrypt(path string, ownerName string, name string, msg string, gamma []int, debug string) *C.char {
	gammastr := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(gamma)), ","), "[]")

	printMsg("Encrypt gamma: "+gammastr, debug)
	a := new(abe.GPSW)
	ReadGob(path+"gpsw"+ownerName+"_"+name+".gob", a)

	pubKey := new(abe.GPSWPubKey)
	ReadGob(path+"pubgpsw"+ownerName+"_"+name+".gob", pubKey)

	cipher, err := a.Encrypt(msg, gamma, pubKey)

	if err != nil {
		fmt.Errorf("Failed to encrypt: %v", err)
		output := fmt.Sprintf("%s: %s", "Error", err)
		return C.CString(output)
	}

	//	cypherS := WriteTString(cipher)

	s := WriteTCharA(cipher)

	mis := WriteTString(cipher)
	//	printMsg("Encrypt result in golang: ", debug)
	//	printMsg(cypherS, debug)
	//	printMsg("Encrypt end in golang", debug)
	//	sEncoded := base64.StdEncoding.Strict().EncodeToString(s)
	//	printMsg("Encrypted encoded: ", debug)
	//	printMsg(sEncoded, debug)
	//	printMsg("Encrypted end in golang", debug)

	sEncodedHex := hex.EncodeToString(s)

	printMsg("encript cipher as string without encode:\n"+mis, debug)
	printMsg("encript sEncodedHex: \n"+sEncodedHex, debug)

	return C.CString(sEncodedHex)

}

//export GeneratePolicyKeys
func GeneratePolicyKeys(path string, ownerName string, name string, clientName string, policy string, attrib []int, debug string) *C.char {

	a := new(abe.GPSW)
	ReadGob(path+"gpsw"+ownerName+"_"+name+".gob", a)
	printMsg("GeneratePolicyKeys", debug)
	//check existence of msp and private keys

	msProgram, err := abe.BooleanToMSP(policy, true)

	sk := new(data.Vector)
	ReadGob(path+"masterKeygpsw"+ownerName+"_"+name+".gob", sk)
	keyVector, err := a.GeneratePolicyKeys(msProgram, *sk)
	if err != nil {
		fmt.Errorf("Failed to create keyVector: %v", err)
		output := fmt.Sprintf("%s: %s", "Error", err)
		return C.CString(output)
	}

	abeKey := a.DelegateKeys(keyVector, msProgram, attrib)
	WriteGob(path+"abeKey"+ownerName+"_"+name+"_"+clientName+".gob", abeKey)

	return C.CString("ok")
}

//export Decrypt
func Decrypt(path string, ownerName string, name string, clientName string, cypherS string, debug string) *C.char {

	a := new(abe.GPSW)
	ReadGob(path+"gpsw"+ownerName+"_"+name+".gob", a)

	printMsg("decrypt: cypherS:\n"+cypherS, debug)

	decodedHex, err := hex.DecodeString(cypherS)
	if err != nil {
		fmt.Errorf("decript. decoding hex error: %v", err)
		output := fmt.Sprintf("%s: %s", "Error", err)
		return C.CString(output)

	}
	cipher := new(abe.GPSWCipher)

	ReadFCharA(decodedHex, cipher)
	mis := WriteTString(cipher)
	printMsg("decript cipher as string without encode:\n"+mis, debug)

	abekey := new(abe.GPSWKey)
	ReadGob(path+"abeKey"+ownerName+"_"+name+"_"+clientName+".gob", abekey)

	msgCheck, err := a.Decrypt(cipher, abekey)
	if err != nil {
		fmt.Errorf("Failed to decrypt: %v", err)
		output := fmt.Sprintf("%s: %s", "Error", err)
		return C.CString(output)
	}
	//printMsg("goland result final: "+msgCheck, debug)
	return C.CString(msgCheck)
}

func printMsg(msg string, debug string) {
	if debug == "ok" {
		fmt.Println(msg + "\n")
	}
}

func main() {}
