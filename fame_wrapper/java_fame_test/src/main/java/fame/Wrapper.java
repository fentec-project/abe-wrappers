

/**
*
*FENTEC Functional Encryption Technologies
*Privacy-preserving and Auditable Digital Currency Use Case
*Copyright Â© 2019 Atos Spain SA
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
package main.java.fame;

import java.util.Arrays;
import java.util.List;

import com.sun.jna.Library;
import com.sun.jna.Memory;
import com.sun.jna.Native;
import com.sun.jna.Pointer;
import com.sun.jna.Structure;

/*
 * Wrapper to shared "so" library of FAME in Golang
 */
public class Wrapper {

	/*
	 * JNA adaptor interface
	 */
	public interface FAMEWrapper extends Library {
		// GoSlice class maps to:
		// C type struct { void *data; GoInt len; GoInt cap; }
		public class GoSlice extends Structure {
			public static class ByValue extends GoSlice implements Structure.ByValue {}
			public Pointer data;
			public long len;
			public long cap;
			protected List getFieldOrder(){
				return Arrays.asList(new String[]{"data","len","cap"});
			}
		}

		// GoString class maps to:
		// C type struct { const char *p; GoInt n; }
		public class GoString extends Structure {
			public static class ByValue extends GoString implements Structure.ByValue {}
			public String p;
			public long n;
			protected List getFieldOrder(){
				return Arrays.asList(new String[]{"p","n"});
			}

		}

		// Foreign functions
		public String GenerateMasterKeys (GoString.ByValue pathStr, GoString.ByValue debugStr);
		public String Encrypt(GoString.ByValue path, GoString.ByValue msg, GoString.ByValue policy, GoString.ByValue debugStr);
		//public String GenerateAttribKeys(GoString.ByValue path, GoString.ByValue pathGen, GoSlice.ByValue gamma, GoString.ByValue debugStr);
		public String GenerateAttribKeys(GoString.ByValue path, GoString.ByValue pathGen, GoString.ByValue gamma, GoString.ByValue debugStr);
		public String Decrypt(GoString.ByValue path, GoString.ByValue pathGen,  GoString.ByValue cypherS, GoString.ByValue debugStr);
		
		
	}

	/*
	 * Wrapper object
	 */
	FAMEWrapper fameWrapper = null;
	String debug;

	/*
	 * Constructor
	 */
	public Wrapper (String libPath, String debug) {
		fameWrapper = (FAMEWrapper) Native.loadLibrary(libPath, FAMEWrapper.class);
		this.debug = debug;
	}

	/*
	 * Generate Master Key
	 */
	public String genMaster (String path) {
		FAMEWrapper.GoString.ByValue pathStr = getGoString(path);
		FAMEWrapper.GoString.ByValue debugStr = getGoString(debug);

		String p = fameWrapper.GenerateMasterKeys(pathStr, debugStr);

		return p;
	}


	/*
	 * Encrypt
	 */
	public String encrp (String path, String name, String msg, String policy) {
		
		FAMEWrapper.GoString.ByValue pathStr = getGoString (path);
		FAMEWrapper.GoString.ByValue msgStr = getGoString (msg);
		FAMEWrapper.GoString.ByValue policyStr = getGoString (policy);
		FAMEWrapper.GoString.ByValue debugStr = getGoString(debug);
		
		//String result = fameWrapper.Encrypt(pathStr, ownerNameStr, nameStr, msgStr, policyStr, debugStr);
		String r = fameWrapper.Encrypt(pathStr, msgStr, policyStr, debugStr);

		
		//byte[] resultAsUTF8 = result.getBytes(StandardCharsets.UTF_8);
		return r;
	}
	
	/*
	 * Generate Keys
	 */
	public String GenAttrKeys(String path, String pathGen, String[] gamma){
		
		StringBuffer sb = new StringBuffer();
		int iMax = gamma.length -1;
		for (int i=0 ; i<gamma.length ; i++) {
			sb.append(gamma[i]);
			if ( i == iMax)
				break;
			sb.append(",");
		}
		
		FAMEWrapper.GoString.ByValue pathStr = getGoString (path);
		FAMEWrapper.GoString.ByValue pathGenStr = getGoString (pathGen);
		FAMEWrapper.GoString.ByValue gammaValStr = getGoString(sb.toString());
		FAMEWrapper.GoString.ByValue debugStr = getGoString(debug);
		
		String result= fameWrapper.GenerateAttribKeys(pathStr, pathGenStr, gammaValStr, debugStr);
		return result;
	}
	
	
	/*
	 * Decrypt	
	 */
	public String Decrypt(String path, String pathGen, String cypherS) {
		FAMEWrapper.GoString.ByValue pathStr = getGoString (path);
		FAMEWrapper.GoString.ByValue pathGenStr = getGoString (pathGen);
		FAMEWrapper.GoString.ByValue cypherSStr = getGoString (cypherS);
		FAMEWrapper.GoString.ByValue debugStr = getGoString(debug);
		
		String result = fameWrapper.Decrypt(pathStr, pathGenStr, cypherSStr, debugStr);
		return result;
	}
	
	/*
	 * JNA data struts creators
	 */

	private FAMEWrapper.GoString.ByValue getGoString (String s){
		FAMEWrapper.GoString.ByValue gs = new FAMEWrapper.GoString.ByValue();
		gs.p = s;
		gs.n = gs.p.length();
		return gs;
	}
	
	private FAMEWrapper.GoSlice.ByValue getGoByValue (long[] vals){
		Memory arr = new Memory(vals.length * Native.getNativeSize(Long.TYPE));
		arr.write(0,  vals, 0, vals.length);
		
		FAMEWrapper.GoSlice.ByValue slice = new FAMEWrapper.GoSlice.ByValue();
		slice.data = arr;
		slice.len = vals.length;
		slice.cap = vals.length;
		return slice;
	}
}

