using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace GTech.IO
{
    public class HttpStream : InputStream
    {
        Stream stream = null;

        public HttpStream(Stream st)
        {
            stream = st;
        }

        public override void Close()
        {
            stream = null;
        }

        public override int Read()
        {
            return -1;
        }

        public override int Read(byte[] buffer)
        {
            return stream.Read(buffer, 0, buffer.Length);
        }

        public override int Read(byte[] buffer, int byteOffset, int byteCount)
        {
            return stream.Read(buffer, byteOffset, byteCount);
        }
    }
}
