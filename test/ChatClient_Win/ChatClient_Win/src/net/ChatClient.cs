using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace GTech.Net
{
    class ChatClient : IMsgParse, IConnListener
    {
        TcpClient tcpClient;

        public ChatClient(string addr)
        {
            tcpClient = new TcpClient(addr);
        }

        public void Login(string account, string password)
        {

        }

        public void LogOut()
        {

        }

        public void OnClose()
        {
            throw new NotImplementedException();
        }

        public void OnError(int code, string desc)
        {
            throw new NotImplementedException();
        }

        public void OnPostSend(byte[] buff, int num)
        {
            throw new NotImplementedException();
        }

        public void OnPreSend(byte[] buff)
        {
            throw new NotImplementedException();
        }

        public void OnRecvBusy(byte[] buff)
        {
            throw new NotImplementedException();
        }

        public void OnSendBusy(byte[] buff)
        {
            throw new NotImplementedException();
        }

        public int ParseHeader(byte[] buff)
        {
            throw new NotImplementedException();
        }

        public void ParseMsg(byte[] buff)
        {
            throw new NotImplementedException();
        }
    }
}
