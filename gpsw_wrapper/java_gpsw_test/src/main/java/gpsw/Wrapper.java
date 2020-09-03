

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
package main.java.gpsw;


import java.io.File;
import java.util.Arrays;
import java.util.List;


import com.sun.jna.Library;
import com.sun.jna.Memory;
import com.sun.jna.Native;
import com.sun.jna.Pointer;
import com.sun.jna.Structure;

/*
 * Wrapper to shared "so" library of GPSW in Golang
 */
public class Wrapper{


	/*
	 * JNA adaptor interface
	 */
	public interface GPSWWrapper extends Library {
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
		public String GenerateMasterKeys(GoString.ByValue pathStr, long numAtt, GoString.ByValue debugStr);
		//func GenerateMasterKeys(path string, numAtt int, debug string) *C.char {
		
		public String Encrypt(GoString.ByValue pathStr, GoString.ByValue msg, GoSlice.ByValue gamma, GoString.ByValue debugStr);
		//func Encrypt(path string, msg string, gamma []int, debug string) *C.char {
		
		              
		public String GeneratePolicyK(GoString.ByValue pathStr, GoString.ByValue userPathStr, GoString.ByValue policy, GoString.ByValue debugStr);
		//func GeneratePolicyK(path string, policy string, debug string) *C.char {
		
		public String Decrypt(GoString.ByValue pathStr, GoString.ByValue pathUserStr, GoString.ByValue  cypherS , GoString.ByValue debugStr);
		//func Decrypt(path string, path_user string, cypherS string, debug string) *C.char {
		
		
	}	
	
	/*****************************************************
	 * 
	 *****************************************************/
	GPSWWrapper gpswWrapper = null;
	String debug;

	/*
	 * Constructor
	 */
	public Wrapper (String libPath, String debug) {
		gpswWrapper = (GPSWWrapper) Native.loadLibrary(libPath, GPSWWrapper.class);
		this.debug = debug;
	}

	private boolean createDir(String path) {
		File file = new File(path);
		try {
			if (!file.exists()) {
				if (file.mkdirs()) {
					file.setWritable(true, false);
					file.setReadable(true, false);
				}else
					return false;
            }
			return true;
		}catch (Exception e) {
			e.printStackTrace();
			return false;
		}
	}

	public String genMaster (String path, long numAtt) {
		this.createDir(path);

		GPSWWrapper.GoString.ByValue pathStr = getGoString(path);
		GPSWWrapper.GoString.ByValue debugStr = getGoString(debug);



		String p = gpswWrapper.GenerateMasterKeys(pathStr, numAtt, debugStr);

		return p;
	}

	/*
	 * Encrypt
	 */
	public String encrp (String path, String msg, long[] gamma) {

		GPSWWrapper.GoString.ByValue pathStr = getGoString (path);
		GPSWWrapper.GoString.ByValue msgStr = getGoString (msg);
		GPSWWrapper.GoSlice.ByValue gammaValStr = getGoByValue(gamma);
		GPSWWrapper.GoString.ByValue debugStr = getGoString(debug);

		//String result = GPSWWrapper.Encrypt(pathStr, ownerNameStr, nameStr, msgStr, policyStr, debugStr);
		String r = gpswWrapper.Encrypt(pathStr, msgStr, gammaValStr, debugStr);


		//byte[] resultAsUTF8 = result.getBytes(StandardCharsets.UTF_8);
		return r;
	}

	/*
	 * Generate Keys
	 */


	public String GeneratePolicyK(String path, String pathUser, String policy){
		GPSWWrapper.GoString.ByValue pathStr = getGoString (path);
		GPSWWrapper.GoString.ByValue pathUserStr = getGoString (pathUser);

		GPSWWrapper.GoString.ByValue policyStr = getGoString (policy);
		GPSWWrapper.GoString.ByValue debugStr = getGoString(debug);

		this.createDir(path);
		System.out.println("GeneratePolicyK, path: " +  path);

		String result= gpswWrapper.GeneratePolicyK(pathStr, pathUserStr, policyStr, debugStr);

		return result;
	}






	/*
	 * Decrypt
	 */
	public String Decrypt(String path, String pathUser, String cypherS) {
		GPSWWrapper.GoString.ByValue pathStr = getGoString (path);
		GPSWWrapper.GoString.ByValue cypherSStr = getGoString (cypherS);
		GPSWWrapper.GoString.ByValue debugStr = getGoString(debug);

		GPSWWrapper.GoString.ByValue pathUserStr = getGoString (pathUser);

		String result = gpswWrapper.Decrypt(pathStr, pathUserStr, cypherSStr, debugStr);

		return result;
	}

	/*
	 * JNA data struts creators
	 */

	private GPSWWrapper.GoString.ByValue getGoString (String s){
		GPSWWrapper.GoString.ByValue gs = new GPSWWrapper.GoString.ByValue();
		gs.p = s;
		gs.n = gs.p.length();
		return gs;
	}

	private GPSWWrapper.GoSlice.ByValue getGoByValue (long[] vals){
		Memory arr = new Memory(vals.length * Native.getNativeSize(Long.TYPE));
		arr.write(0,  vals, 0, vals.length);

		GPSWWrapper.GoSlice.ByValue slice = new GPSWWrapper.GoSlice.ByValue();
		slice.data = arr;
		slice.len = vals.length;
		slice.cap = vals.length;
		return slice;
	}
}



