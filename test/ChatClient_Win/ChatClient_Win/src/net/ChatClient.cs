using GTech.Log;
using GTech.Net.Protocol;
using GTech.Utils;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace GTech.Net
{
    public enum GTChatType
    {
        ChatType_Chat,
        ChatType_GroupChat
    }

    public delegate void GTLoginedHandler();
    public delegate void GTCloseHandler();
    public delegate void GTErrorHandler(int errorcode);
    //public delegate void GTMessageHandler(string name, GTChatType chattype, IMessage message);
    //public delegate void GTRoomJoinHandler(string name, GTErrorCode code);
    //public delegate void GTMessageSendHandler(IMessage message, GTErrorCode code);

    public partial class ChatClient : IMsgParse, IConnListener
    {
        public GTLoginedHandler LoginedHandler;
        public GTCloseHandler CloseHandler;
        public GTErrorHandler ErrorHandler;
        //public GTMessageHandler OnMessageHandler;

        TcpClient tcpClient;
        string addr;

        public ChatClient(string addr)
        {
            this.addr = addr;
        }

        public void OnClose()
        {
            GLog.d("socket closed");
        }

        public void OnError(int code, string desc)
        {
            throw new NotImplementedException();
        }

        public void OnPostSend(byte[] buff, int num)
        {
            //throw new NotImplementedException();
        }

        public void OnPreSend(byte[] buff)
        {
            //throw new NotImplementedException();
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
            GLog.d("ParseHeader " + System.BitConverter.ToUInt16(buff, 0));
            return (int)System.BitConverter.ToUInt16(buff, 0);
        }

        public void ParseMsg(byte[] buff)
        {
            ushort msgid = System.BitConverter.ToUInt16(buff, 0);

            switch (msgid)
            {
                case (ushort)MsgId.MsgId_Tick:
                    GLog.d("recv tick from server");
                    break;
                case (ushort)MsgId.MsgId_Echo:
                    GLog.d("recv echo from server:" + Encoding.UTF8.GetString(buff, 2, buff.Length - 2));
                    break;
                case (ushort)MsgId.MsgId_RetLogin:
                    ushort result = System.BitConverter.ToUInt16(buff, 2);
                    if(result == 0)
                    {
                        GLog.d("login success");

                        if (LoginedHandler != null)
                        {
                            LoginedHandler.Invoke();
                        }
                    }
                    else
                    {
                        GLog.d("login failed:" + result);

                        if (ErrorHandler != null)
                        {
                            ErrorHandler.Invoke(result);
                        }
                    }
                    break;
                default:
                    break;
            }
        }
    }
}
