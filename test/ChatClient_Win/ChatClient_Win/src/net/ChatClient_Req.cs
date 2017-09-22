using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using GTech.Log;
using GTech.Net.Protocol;
using GTech.Utils;

namespace GTech.Net
{
    public partial class ChatClient : IMsgParse, IConnListener
    {
        public void Login(ulong account, string password)
        {
            tcpClient = new TcpClient(addr);
            tcpClient.Parser = this;
            tcpClient.Listener = this;

            try
            {
                tcpClient.Connect();
            }
            catch (Exception e)
            {
                GLog.d(e.ToString());
                return;
            }

            GLog.d("sending login req to server...");
            MsgReqLogin reqlogin = new MsgReqLogin();
            reqlogin.Uid = account;
            GLog.d("md5 password is " + MD5Utils.StringToMD5(password));
            reqlogin.Password = Encoding.UTF8.GetBytes(MD5Utils.StringToMD5(password));

            //reqlogin.write(tcpClient.SendStream);
            tcpClient.Send(reqlogin.toBytes());
        }

        public void LogOut()
        {

        }

        public void SendTick()
        {
            MsgTick tick = new MsgTick();
            tcpClient.Send(tick.toBytes());
        }
    }
}
