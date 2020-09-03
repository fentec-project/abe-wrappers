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
package main.java.gpsw;

import java.util.Arrays;

public class GpswTest1 {
                         
	static String libPath = "../src/main/gpsw_adaptor.so";
    
	
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
    
        System.out.println("txt sin cifrar: " + msg);
    
        String result = wrapper.genMaster(path,31);
        System.out.println ("\n#################\n\n");
		System.out.println("Create Master Key result: " + result);
    
        long[] gamma = {1,2,3,4,5,6,7,8,9,10,11,12,13,14,29,30};
        String result3 = wrapper.encrp (path, msg, gamma);
        System.out.println ("\n#################\n\n");
        System.out.println("encrypt result in java: \n"+ result3);
        System.out.println("encrypt result in java end.");
        System.out.println("result length = " + result3.length());
        System.out.println("Encrypted with attributes: " + Arrays.toString(gamma));
        
 
        String policy = "((15 AND 13) OR 28)";
        
        System.out.println ("\n#################\n\n");
		System.out.println("Generate decryption key with policy: " + policy);
		result = wrapper.GeneratePolicyK(path, path, policy);
        System.out.println("GeneratePolicyK result: "+ result);
 
    	String result2 = wrapper.Decrypt(path, path, result3);
    	System.out.println ("\n#################\n\n");
    	System.out.println("decrypt result: "+ result2);

        
        policy = "((((10 OR 11 OR 12) AND 13) OR (1 AND 14)) AND 29)";
    
        System.out.println ("\n#################\n\n");
		System.out.println("Generate decryption key with policy: " + policy);
		result = wrapper.GeneratePolicyK(path, path, policy);
        System.out.println("GeneratePolicyK result: "+ result);
 
    	result2 = wrapper.Decrypt(path, path, result3);
    	System.out.println ("\n#################\n\n");
    	System.out.println("decrypt result: "+ result2);
	}
}
