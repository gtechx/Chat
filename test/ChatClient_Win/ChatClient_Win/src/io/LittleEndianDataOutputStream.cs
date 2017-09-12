using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace GTech.IO
{
    public class LittleEndianDataOutputStream
    {

        OutputStream os = null;

        public LittleEndianDataOutputStream(OutputStream os)
        {
            this.os = os;
        }

        public void close()
        {
            os = null;
        }

        public void flush()
        {
            os.flush();
        }

        public void writeDouble(double val)
        {
            byte[] bytes = System.BitConverter.GetBytes(val);
            os.write(bytes);
        }

        public void writeFloat(float val)
        {
            byte[] bytes = System.BitConverter.GetBytes(val);
            os.write(bytes);
        }

        public void writeInt(int data)
        {
            byte[] bytes = System.BitConverter.GetBytes(data);
            os.write(bytes);
        }

        public void writeInt(uint data)
        {
            byte[] bytes = System.BitConverter.GetBytes(data);
            os.write(bytes);
        }

        public void writeLong(long data)
        {
            byte[] bytes = System.BitConverter.GetBytes(data);
            os.write(bytes);
        }

        public void writeLong(ulong data)
        {
            byte[] bytes = System.BitConverter.GetBytes(data);
            os.write(bytes);
        }

        public void writeShort(short data)
        {
            byte[] bytes = System.BitConverter.GetBytes(data);
            os.write(bytes);
        }

        public void writeUShort(ushort data)
        {
            byte[] bytes = System.BitConverter.GetBytes(data);
            os.write(bytes);
        }

        public void writeChar(char data)
        {
            byte[] bytes = System.BitConverter.GetBytes(data);
            os.write(bytes);
        }

        public void writeChars(char[] data)
        {
            byte[] bytes = new byte[data.Length];
            Buffer.BlockCopy(data, 0, bytes, 0, data.Length);
            //Array.Copy(data, 0, bytes, 0, data.Length);
            //   for (int i = 0; i<data.Length; ++i){
            //	bytes[i] = (byte)data[i];
            //}

            os.write(bytes);
        }

        public void writeByte(byte data)
        {
            byte[] bytes = new byte[1];
            bytes[0] = data;

            os.write(bytes);
        }

        public void write(byte[] buffer)
        {
            if (buffer == null)
            {
                throw new Exception("buffer == null");
            }
            write(buffer, 0, buffer.Length);
        }

        public void write(byte[] buffer, int offset, int count)
        {
            if (buffer == null)
            {
                throw new Exception("buffer == null");
            }
            os.write(buffer, offset, count);
        }
    }
}
