using System;
using System.Collections.Generic;
using System.Linq;
using System.Net.Sockets;
using System.Text;
using System.Threading.Tasks;

namespace GTech.IO
{
    public class SocketOutputStream : OutputStream
    {
        Socket socket;
        public SocketOutputStream(Socket sock)
        {
            socket = sock;
        }

        public override void write(byte[] buffer)
        {
            socket.Send(buffer);
        }

        public override void write(byte[] buffer, int offset, int count)
        {
            byte[] sendbuffer = buffer;
            if(offset != 0 || count != buffer.Length)
            {
                int size = buffer.Length - offset > count ? count : buffer.Length - offset;
                sendbuffer = new byte[size];
                Array.Copy(buffer, offset, sendbuffer, 0, size);
            }
            write(sendbuffer);
        }

        public override void flush()
        {
            //
        }

        public override void close()
        {
            //
        }
    }
}
