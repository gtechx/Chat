using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace GTech.IO
{
    public class InputStream
    {
        //public virtual void Write(byte[] buffer)
        //{
        //}

        public virtual int Read()
        {
            return -1;
        }

        public virtual int Read(byte[] buffer)
        {
            return 0;
        }

        public virtual int Read(byte[] buffer, int byteOffset, int byteCount)
        {
            return 0;
        }

        public virtual void Flush()
        {
            //
        }

        public virtual void Close()
        {
            //
        }
    }
}
