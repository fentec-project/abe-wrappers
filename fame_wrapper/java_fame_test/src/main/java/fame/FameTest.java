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

package main.java.fame;

import java.util.Arrays;

public class FameTest {
	
	static String libPath = "../src/main/fame_adaptor.so";


	static public void main(String args[]) {
		
		StringBuilder sb = new StringBuilder();
		sb.append("msg: ");
		for (int i = 0 ; i < args.length ; i++)
	       		sb.append(" " + args[i]);
		
		System.out.println ("\n#################\n\n");
		System.out.println ("message to cypher: " + sb.toString());
		
		//Wrapper wrapper = new Wrapper (libPath, "ok");
		Wrapper wrapper = new Wrapper (libPath, "ko");

		String path = "./keys";

		String msg = sb.toString();

		String result = wrapper.genMaster(path);
		System.out.println ("\n#################\n\n");
		System.out.println("Create Master Key result: " + result);

		String policy = "( director AND financial )";
		String encoded = wrapper.encrp(path,"pepe",msg,policy);
		System.out.println ("\n#################\n\n");
		System.out.println("encrypt result in java: \n"+ encoded);
        System.out.println("encrypt result in java end.");
        System.out.println("result length = " + encoded.length());
        System.out.println("Encrypted with policiy: " + policy);

        String[] gamma={"analyst","industrial", "engines"};
        System.out.println ("\n#################\n\n");
		System.out.println("Generate decryption key with atributes: " + Arrays.toString(gamma));
		result = wrapper.GenAttrKeys(path, path,gamma);
		System.out.println ("genAttibKeys result: " + result);
		
		result = wrapper.Decrypt(path,path,encoded);
		System.out.println ("\n#################\n\n");
		System.out.println("Decrypt result : " + result + "\nEND\n");
		
		String[] gamma2={"director","industrial", "financial"};
		System.out.println ("\n#################\n\n");
		System.out.println("Generate decryption key with atributes: " + Arrays.toString(gamma2));
		result = wrapper.GenAttrKeys(path, path,gamma2);
		System.out.println ("genAttibKeys result: " + result);
		

		result = wrapper.Decrypt(path,path,encoded);
		System.out.println ("Decrypt result : " + result + "\nEND\n");
		
	}
}


