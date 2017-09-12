using System;
using System.Collections.Generic;
using System.Linq;
using System.Security.Cryptography;
using System.Text;
using System.Threading.Tasks;

namespace GTech.Utils
{
    public class MD5Utils
    {
        /***
	 * MD5加码 生成32位md5码
	 */
        public static string StringToMD5(string inStr)
        {
            char[] charArray = inStr.ToCharArray();
            byte[] byteArray = new byte[charArray.Length];

            for (int i = 0; i < charArray.Length; i++)
                byteArray[i] = (byte)charArray[i];

            //Buffer.BlockCopy(charArray, 0, byteArray, 0, charArray.Length);
            //Array.Copy(charArray, byteArray, charArray.Length);

            byte[] md5Bytes = BytesToMD5(byteArray);

            //string md5str = ByteUtils.ToString(md5Bytes);
            //string hexValue = "";
            //for (int i = 0; i < md5Bytes.Length; i++)
            //{
            //    int val = ((int)md5Bytes[i]) & 0xff;
            //    if (val < 16)
            //        hexValue += "0";
            //    hexValue += val.ToString("x");
            //    //hexValue.append(Integer.toHexString(val));
            //}

            return ByteUtils.ToString(md5Bytes).Replace("-", "").ToLower();
        }

        /***
         * MD5加码 生成32位md5码
         */
        public static byte[] BytesToMD5(byte[] bytes)
        {
            MD5 md5 = new MD5CryptoServiceProvider();
            byte[] output = md5.ComputeHash(bytes);

            return output;

            //MessageDigest md5 = null;
            //try
            //{
            //    md5 = MessageDigest.getInstance("MD5");
            //}
            //catch (Exception e)
            //{
            //    e.printStackTrace();
            //    return "MD5 Exception:" + e.toString();
            //}

            //byte[] md5Bytes = md5.digest(bytes);
            //StringBuffer hexValue = new StringBuffer();
            //for (int i = 0; i < md5Bytes.length; i++)
            //{
            //    int val = ((int)md5Bytes[i]) & 0xff;
            //    if (val < 16)
            //        hexValue.append("0");
            //    hexValue.append(Integer.toHexString(val));
            //}
            //return hexValue.toString();
        }
    }
}
