using System;
using System.Collections.Generic;
using System.Linq;
using System.Net.Sockets;
using System.Text;
using System.Threading.Tasks;

namespace GTech.IO
{
    public class SocketInputStream : InputStream
    {
        Socket socket;
        public SocketInputStream(Socket sock)
        {
            socket = sock;
        }

        public override int Read()
        {
            return -1;
        }

        //public override void write(byte[] buffer)
        //{
        //    socket.Send(buffer);
        //}

        public override int Read(byte[] buffer)
        {
            return socket.Receive(buffer);
        }

        public override int Read(byte[] buffer, int byteOffset, int byteCount)
        {
            return socket.Receive(buffer, byteOffset, byteCount, SocketFlags.None);
        }

        public override void Flush()
        {
            //
        }

        public override void Close()
        {
            //
        }
    }
}
