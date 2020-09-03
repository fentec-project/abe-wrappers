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
	"encoding/hex"
	"fmt"
	"github.com/fentec-project/gofe/abe"
	"github.com/fentec-project/gofe/data"
	"strings"
)

// This is an adaptor for the KP ABE scheme (https://eprint.iacr.org/2006/309.pdf) on github.com/fentec-project/gofe/abe
// in order to save and load keys from files

//export GenerateMasterKeys
func GenerateMasterKeys(path string, numAtt int, debug string) *C.char {
	b := string(numAtt)

	printMsg("GenerateMasterKeys with num attr: "+b, debug)
	a := abe.NewGPSW(numAtt)
	err := WriteGob_pn(path, "abe.gob", a)
	if err != nil {
		printMsg("There has been an error while creating GPSW",debug)
		fmt.Println(err)
		fmt.Errorf("Failed to write gpsw file: %v", err)
		output := fmt.Sprintf("%s: %s", "Error", err)
		return C.CString(output)
	}

	printMsg("GenerateMasterKeys: write gpsw ok", debug)

	pubKey, secKey, err := a.GenerateMasterKeys()
	if err != nil {
		fmt.Errorf("Failed to generate Master and Secret keys: %v", err)
		output := fmt.Sprintf("%s: %s", "Error", err)
		return C.CString(output)

	}

	err = WriteGob_pn(path, "PK.gob", pubKey)
	if err != nil {
		fmt.Errorf("Failed to write PK file: %v", err)
		output := fmt.Sprintf("%s: %s", "Error", err)
		return C.CString(output)
	}
	printMsg("GenerateMasterKeys: write PK ok", debug)

	err = WriteGob_pn(path, "MK.gob", secKey)
	if err != nil {
		fmt.Errorf("Failed to write MK file: %v", err)
		output := fmt.Sprintf("%s: %s", "Error", err)
		return C.CString(output)
	}
	printMsg("GenerateMasterKeys: write MK ok", debug)

	return C.CString("ok")
}

//export Encrypt
func Encrypt(path string, msg string, gamma []int, debug string) *C.char {
	//format gamma
	gammastr := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(gamma)), ","), "[]")

	printMsg("Encrypt: gamma="+gammastr, debug)
	a := new(abe.GPSW)
	ReadGob_pn(path, "abe.gob", a)

	pubKey := new(abe.GPSWPubKey)
	ReadGob_pn(path, "PK.gob", pubKey)

	cipher, err := a.Encrypt(msg, gamma, pubKey)

	if err != nil {
		fmt.Errorf("Failed to encrypt: %v", err)
		output := fmt.Sprintf("%s: %s", "Error", err)
		return C.CString(output)
	}


	s := WriteTCharA(cipher)

	mis := WriteTString(cipher)

	sEncodedHex := hex.EncodeToString(s)

	printMsg("encript cipher as string without encode:\n"+mis, debug)
	printMsg("encript sEncodedHex: \n"+sEncodedHex, debug)

	return C.CString(sEncodedHex)

}


//export GeneratePolicyK
func GeneratePolicyK(path string, path_user string, policy string, debug string) *C.char {
	a := new(abe.GPSW)
	ReadGob_pn(path, "abe.gob", a)
	printMsg("GeneratePolicyK", debug)

	msProgram, err := abe.BooleanToMSP(policy, true)

	sk := new(data.Vector)
	ReadGob_pn(path, "MK.gob", sk)
	abeKey, err := a.GeneratePolicyKey(msProgram, *sk)
	if err != nil {
		fmt.Errorf("Failed to generate Policy Key: %v", err)
		output := fmt.Sprintf("%s: %s", "Error", err)
		return C.CString(output)
	}
	err = WriteGob_pn(path_user, "delegatedKey.gob", abeKey)
	if err != nil {
		fmt.Errorf("Failed to write policy key file: %v", err)
		output := fmt.Sprintf("%s: %s", "Error", err)
		return C.CString(output)
	}
	return C.CString("ok")
}


//export Decrypt
func Decrypt(path string, path_user string, cypherS string, debug string) *C.char {
	defer func() {
                if err := recover(); err != nil {
			err = fmt.Errorf("Error: could not decrypt, panic in abe.decrypt")
                }

        }()
	a := new(abe.GPSW)
	ReadGob_pn(path, "abe.gob", a)

	printMsg("Decrypt: cypherS:\n"+cypherS, debug)

	decodedHex, err := hex.DecodeString(cypherS)
	if err != nil {
		fmt.Errorf("Decrypt: decoding hex error: %v", err)
		output := fmt.Sprintf("%s: %s", "Error", err)
		return C.CString(output)

	}
	cipher := new(abe.GPSWCipher)

	ReadFCharA(decodedHex, cipher)
	mis := WriteTString(cipher)
	printMsg("Decrypt: cipher as string without encode:\n"+mis, debug)

	abekey := new(abe.GPSWKey)
	ReadGob_pn(path_user, "delegatedKey.gob", abekey)

	msgCheck, err := a.Decrypt(cipher, abekey)
	if err != nil {
		fmt.Errorf("Decrypt: Failed to decrypt: %v", err)
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
