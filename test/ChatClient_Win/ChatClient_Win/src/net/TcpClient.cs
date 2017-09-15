using System;
using System.Collections.Generic;
using System.Linq;
using System.Net;
using System.Net.Sockets;
using System.Runtime.Remoting.Messaging;
using System.Text;
using System.Threading;
using System.Threading.Tasks;

namespace GTech.Net
{
    interface IMsgParse
    {
        int ParseHeader(byte[] buff);
        void ParseMsg(byte[] buff);
    }

    interface IConnListener
    {
        void OnClose();
        void OnError(int code, string desc);
        void OnPreSend(byte[] buff);
        void OnPostSend(byte[] buff, int num);
        void OnRecvBusy(byte[] buff);
        void OnSendBusy(byte[] buff);
    }

    class TcpClient
    {
        public delegate void AsyncSendMethod();
        public delegate void AsyncReceiveMethod();

        public IMsgParse Parser { get; set; }
        public IConnListener Listener { get; set; }
        public string Addr { get; set; }

        private string ip = "";
        private int port = -1;
        private Socket socket = null;

        int _wpos = 0;              // 写入的数据位置
        int _spos = 0;              // 发送完毕的数据位置
        int _sending = 0;

        List<byte[]> buffList;

        AsyncCallback _asyncCallback = null;
        AsyncSendMethod _asyncSendMethod;
        private object lockobj = new object();

        class SyncObject
        {
            public bool flag = false;
        }
        SyncObject syncobj;

        public TcpClient(string addr)
        {
            this.Addr = addr;

            var arr = addr.Split(':');
            if(arr.Length == 2)
            {
                ip = arr[0];
                port = int.Parse(arr[1]);
            }

            Reset();
        }

        public void Reset()
        {
            _wpos = 0;
            _spos = 0;
            _sending = 0;
            buffList = new List<byte[]>();
            buffList.Capacity = 1024;
            _asyncSendMethod = new AsyncSendMethod(this._asyncSend);
            _asyncCallback = new AsyncCallback(_onSent);
            syncobj = new Net.TcpClient.SyncObject();
        }

        public bool Connect()
        {
            if(ip == "" || port == -1)
            {
                throw new Exception("ip:" + ip + " or port:" + port + " is invalid.");
            }

            if (socket != null)
                Close();

            socket = new Socket(AddressFamily.InterNetwork, SocketType.Stream, ProtocolType.Tcp);
            socket.NoDelay = true;
            //socket.ReceiveTimeout = readTimeout;
            //socket.SendTimeout = readTimeout;

            IPAddress ipadr = IPAddress.Parse(ip);
            IPEndPoint iep = new IPEndPoint(ipadr, port);
            try
            {
                socket.Connect(iep);
            }
            catch (Exception e)
            {
                socket.Close();
                socket = null;
                throw e;
            }

            startRecv();
            return true;
        }

        public void Close()
        {
            if (Interlocked.CompareExchange(ref _sending, 1, 0) == 1)
            {
                syncobj.flag = true;
                socket = null;
            }
            else if (socket != null)
            {
                socket.Close();
                socket = null;
            }
        }

        public void SetReadTimeout(int timeoutMillis)
        {
            socket.SendTimeout = timeoutMillis;
        }

        public void SetWriteTimeout(int timeoutMillis)
        {
            socket.ReceiveTimeout = timeoutMillis;
        }

        public void Send(byte[] buff)
        {
            if (socket != null)
                return;

            if (0 == Interlocked.Add(ref _sending, 0))
            {
                if (_wpos == _spos)
                {
                    _wpos = 0;
                    _spos = 0;
                }
            }

            int t_spos = Interlocked.Add(ref _spos, 0);
            int space = 0;
            int tt_wpos = _wpos % buffList.Count;
            int tt_spos = t_spos % buffList.Count;

            if (tt_wpos >= tt_spos)
                space = buffList.Count - tt_wpos + tt_spos - 1;
            else
                space = tt_spos - tt_wpos - 1;

            if (1 > space)
            {
                return;
            }

            buffList[tt_wpos] = buff;

            Interlocked.Add(ref _wpos, 1);
            lock (lockobj)
            {
                if (Interlocked.CompareExchange(ref _sending, 1, 0) == 0)
                {
                    _startSend();
                }
            }

            return;
        }

        void _startSend()
        {
            // 由于socket用的是非阻塞式，因此在这里不能直接使用socket.send()方法
            // 必须放到另一个线程中去做
            _asyncSendMethod.BeginInvoke(_asyncCallback, null);
        }

        void _asyncSend()
        {
            Socket tmpsocket = this.socket;
            IConnListener tmplistener = this.Listener;
            List<byte[]> tmplist = this.buffList;
            SyncObject tmpsyncobj = this.syncobj;

            if (tmpsocket == null)
            {
                return;
            }

            while (true)
            {
                //int sendSize = Interlocked.Add(ref _wpos, 0) - _spos;
                //int t_spos = _spos % buffList.Count;
                //if (t_spos == 0)
                //    t_spos = sendSize;

                //if (sendSize > buffList.Count - t_spos)
                //    sendSize = buffList.Count - t_spos;

                int bytesSent = 0;
                var buff = tmplist[_spos];

                if(buff != null)
                {
                    if (tmplistener != null)
                    {
                        tmplistener.OnPreSend(buff);
                    }

                    try
                    {
                        bytesSent = tmpsocket.Send(buff, 0, buff.Length, 0);
                    }
                    catch (SocketException se)
                    {
                        return;
                    }

                    if (tmplistener != null)
                    {
                        tmplistener.OnPostSend(buff, bytesSent);
                    }

                    tmplist[_spos] = null;
                }

                int spos = Interlocked.Add(ref _spos, 1);
                
                lock (lockobj)
                {
                    // 所有数据发送完毕了
                    if (spos == Interlocked.Add(ref _wpos, 0))
                    {
                        Interlocked.Exchange(ref _sending, 0);
                        if(tmpsyncobj.flag)
                        {
                            if(tmpsocket != null)
                            {
                                tmpsocket.Close();
                            }
                        }
                        return;
                    }
                }
            }
        }

        private static void _onSent(IAsyncResult ar)
        {
            AsyncResult result = (AsyncResult)ar;
            AsyncSendMethod caller = (AsyncSendMethod)result.AsyncDelegate;
            caller.EndInvoke(ar);
        }

        public void startRecv()
        {
            var v = new AsyncReceiveMethod(this._asyncReceive);
            v.BeginInvoke(new AsyncCallback(_onRecv), null);
        }

        private void _asyncReceive()
        {
            if (socket == null)
            {
                return;
            }
            Socket tmpsocket = this.socket;
            IConnListener tmplistener = this.Listener;
            IMsgParse tmpparse = this.Parser;
            int MsgHeaderSize = 2;
            byte[] headerbuf = new byte[MsgHeaderSize];

            while (true)
            {
                int num = 0;

                try
                {
                    num = tmpsocket.Receive(headerbuf, 0, 2, 0);
                }
                catch (SocketException se)
                {
                    break;
                }

                if(num >= MsgHeaderSize)
                {
                    int datasize = 0;
                    if(tmpparse != null)
                    {
                        datasize = tmpparse.ParseHeader(headerbuf);
                    }

                    if(datasize > 0)
                    {
                        byte[] databuff = new byte[datasize];
                        try
                        {
                            num = tmpsocket.Receive(databuff, 0, datasize, 0);
                        }
                        catch (SocketException se)
                        {
                            break;
                        }

                        if(tmpparse != null && num == datasize)
                        {
                            tmpparse.ParseMsg(databuff);
                        }
                    }
                }
            }

            if (tmplistener != null)
            {
                tmplistener.OnClose();
            }
        }

        private void _onRecv(IAsyncResult ar)
        {
            AsyncResult result = (AsyncResult)ar;
            AsyncReceiveMethod caller = (AsyncReceiveMethod)result.AsyncDelegate;
            caller.EndInvoke(ar);
        }
    }
}
