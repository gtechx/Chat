using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace GTech.Utils
{
    public class ByteUtils
    {
        public static long DoubleToInt64Bits(double value)
        {
            return BitConverter.DoubleToInt64Bits(value);
        }

        public static byte[] GetBytes(string value)
        {
            byte[] byteArray = new byte[value.Length];

            for (int i = 0; i < value.Length; i++)
                byteArray[i] = (byte)value[i];

            return byteArray;
        }

        public static byte[] GetBytes(char[] value)
        {
            byte[] byteArray = new byte[value.Length];

            for (int i = 0; i < value.Length; i++)
                byteArray[i] = (byte)value[i];

            return byteArray;
        }

        public static byte[] GetBytes(bool value)
        {
            return BitConverter.GetBytes(value);
        }

        public static byte[] GetBytes(char value)
        {
            return BitConverter.GetBytes(value);
        }

        public static byte[] GetBytes(short value)
        {
            return BitConverter.GetBytes(value);
        }

        public static byte[] GetBytes(int value)
        {
            return BitConverter.GetBytes(value);
        }

        public static byte[] GetBytes(long value)
        {
            return BitConverter.GetBytes(value);
        }

        public static byte[] GetBytes(ushort value)
        {
            return BitConverter.GetBytes(value);
        }

        public static byte[] GetBytes(uint value)
        {
            return BitConverter.GetBytes(value);
        }

        public static byte[] GetBytes(ulong value)
        {
            return BitConverter.GetBytes(value);
        }

        public static byte[] GetBytes(float value)
        {
            return BitConverter.GetBytes(value);
        }

        public static byte[] GetBytes(double value)
        {
            return BitConverter.GetBytes(value);
        }

        public static double Int64BitsToDouble(long value)
        {
            return BitConverter.Int64BitsToDouble(value);
        }

        public static bool ToBoolean(byte[] value, int startIndex)
        {
            return BitConverter.ToBoolean(value, startIndex);
        }

        public static char ToChar(byte[] value, int startIndex)
        {
            return BitConverter.ToChar(value, startIndex);
        }

        public static char[] ToChar(byte[] value)
        {
            char[] charArray = new char[value.Length];

            for (int i = 0; i < value.Length; i++)
                charArray[i] = (char)value[i];

            return charArray;
        }

        public static char[] GetAsciiChar(string value)
        {
            char[] charArray = new char[value.Length];

            for (int i = 0; i < value.Length; i++)
                charArray[i] = (char)value[i];

            return charArray;
        }

        public static double ToDouble(byte[] value, int startIndex)
        {
            return BitConverter.ToDouble(value, startIndex);
        }

        public static short ToInt16(byte[] value, int startIndex)
        {
            return BitConverter.ToInt16(value, startIndex);
        }

        public static int ToInt32(byte[] value, int startIndex)
        {
            return BitConverter.ToInt32(value, startIndex);
        }

        public static long ToInt64(byte[] value, int startIndex)
        {
            return BitConverter.ToInt64(value, startIndex);
        }

        public static float ToSingle(byte[] value, int startIndex)
        {
            return BitConverter.ToSingle(value, startIndex);
        }

        public static string ToString(byte[] value)
        {
            return BitConverter.ToString(value);
        }

        public static string ToString(byte[] value, int startIndex)
        {
            return BitConverter.ToString(value, startIndex);
        }

        public static string ToString(byte[] value, int startIndex, int length)
        {
            return BitConverter.ToString(value, startIndex);
        }

        public static ushort ToUInt16(byte[] value, int startIndex)
        {
            return BitConverter.ToUInt16(value, startIndex);
        }

        public static uint ToUInt32(byte[] value, int startIndex)
        {
            return BitConverter.ToUInt32(value, startIndex);
        }

        public static ulong ToUInt64(byte[] value, int startIndex)
        {
            return BitConverter.ToUInt64(value, startIndex);
        }

        public static void dumpMemory(string item, byte[] buf)
        {
            string byteStr = "\r\n";
            int i = 0;
            foreach (var bytes in buf)
            {
                byteStr += string.Format("0x{0:X2} ", bytes);
                i++;

                if (i % 8 == 0)
                    byteStr += "\r\n";
            }
        }
    }
}
