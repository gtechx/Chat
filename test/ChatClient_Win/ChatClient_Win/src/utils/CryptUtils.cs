using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace GTech.Utils
{
    public class CryptUtils
    {
        static string key = "JKL:DUIEN(*%&^%&$%#KDFDSGDSDGCVD&*(EJKFDS&*S";
        static int PACKET_MASK = 65535;
        static char[] keyVec;
        static bool isInit = false;

        public static void Init()
        {
            if (!isInit)
            {
                SetKey(key);
                isInit = true;
            }
        }

        private static void SetKey(string key)
        {
            keyVec = new char[PACKET_MASK];
            char[] keybytes = ByteUtils.GetAsciiChar(key);

            for (int i = 0; i < PACKET_MASK; ++i)
            {
                keyVec[i] = keybytes[i % key.Length];
            }
        }

        public static bool CryptBuff(byte[] buffer, int length)
        {
            if (buffer == null || length > PACKET_MASK || keyVec == null)
            {
                return false;
            }

            for (int i = 0; i < length; ++i)
            {
                buffer[i] ^= (byte)keyVec[i];
            }

            return true;
        }
    }
}
