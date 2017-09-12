using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace GTech.IO
{
    public class LittleEndianDataInputStream
    {

        InputStream inputStream = null;
	
	    public LittleEndianDataInputStream(InputStream ins)
        {
            inputStream = ins;
        }

        public void close()
        {

            inputStream = null;
        }

        public float readFloat()
        {
            byte[] bytes = new byte[4];
            return System.BitConverter.ToSingle(bytes, 0);
        }

        public double readDouble()
        {
            byte[] bytes = new byte[4];
            return System.BitConverter.ToDouble(bytes, 0);
        }

        public long readLong()
        {

        byte[] bytes = new byte[8];
        inputStream.Read(bytes);
        return (long)System.BitConverter.ToUInt64(bytes, 0);
    //        return (0xffL & (long) bytes[0]) | (0xff00L & ((long) bytes[1] << 8)) | (0xff0000L & ((long) bytes[2] << 16)) | (0xff000000L & ((long) bytes[3] << 24))
				//| (0xff00000000L & ((long) bytes[4] << 32)) | (0xff0000000000L & ((long) bytes[5] << 40)) | (0xff000000000000L & ((long) bytes[6] << 48))
				//| (long)(0xff00000000000000L & ((ulong) bytes[7] << 56));
	    }

        public int readInt()
        { 
            byte[] bytes = new byte[4];
            inputStream.Read(bytes);
            return System.BitConverter.ToInt32(bytes, 0);
        }

        public byte readByte()
        {
            byte[] bytes = new byte[1];

            inputStream.Read(bytes);
            return bytes[0];
        }

        public short readShort()
        {
            byte[] bytes = new byte[2];

            inputStream.Read(bytes);
            return (short)System.BitConverter.ToInt16(bytes, 0);
        }

        public char readChar()
        {
            byte[] bytes = new byte[1];

            inputStream.Read(bytes);
            return (char)bytes[0];
        }

        public int read(byte[] bytes)
        {
            return inputStream.Read(bytes);
        }

        public int read(byte[] bytes, int byteOffset, int byteCount)
        {
            return inputStream.Read(bytes, byteOffset, byteCount);
        }

        public ushort readUnsignedShort()
        {
            byte[] bytes = new byte[2];

            inputStream.Read(bytes);
            return System.BitConverter.ToUInt16(bytes, 0);
        }

        /**
         * Reads a single byte from this stream and returns it as an integer in the
         * range from 0 to 255. Returns -1 if the end of the stream has been
         * reached. Blocks until one byte has been read, the end of the source
         * stream is detected or an exception is thrown.
         *
         * @throws IOException
         *             if the stream is closed or another IOException occurs.
         */
        public sbyte readSByte()
        {
            byte[] bytes = new byte[1];

            inputStream.Read(bytes);
            return (sbyte)bytes[0];
        }
	
    }

}
