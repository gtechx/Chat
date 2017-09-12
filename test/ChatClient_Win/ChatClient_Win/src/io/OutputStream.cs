using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace GTech.IO
{
    public class OutputStream
    {
        public virtual void write(byte[] buffer)
        {
        }

        public virtual void write(byte[] buffer, int offset, int count)
        {
        }

        public virtual void flush()
        {
            //
        }

        public virtual void close()
        {
            //
        }
    }
}
