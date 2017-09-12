using System;
using System.Net;
using System.Net.Sockets;
using GTech.Log;
using GTech.IO;
using GTech.Net.Protocol;

namespace GTech.Net
{

    /**
     * 
     * @author wangyanqing
     *
     */
    public class GSocketClient
    {
        private Socket socket = null;
        private int readTimeout = 30000;
        private int connectTimeout = 30000;
        private string ip;
        private int port;
        InputStream inputStream = null;
        SocketOutputStream outputStream = null;
        LittleEndianDataOutputStream dos = null;
        LittleEndianDataInputStream dis = null;

        public string toString()
        {
            // TODO Auto-generated method stub
            return "ip=" + ip + " port=" + port;
        }

        public GSocketClient(string ip, int port)
        {
            this.ip = ip;
            this.port = port;

            //socket = new Socket();
        }

        public bool isConnected()
        {
            return socket != null && socket.Connected;
        }

        public void close()
        {
            if (dos != null)
            {
                dos.close();
                dos = null;
            }

            if (dis != null)
            {
                dis.close();
                dis = null;
            }

            if (null != inputStream)
            {
                inputStream.Close();
                inputStream = null;
            }

            if (null != outputStream)
            {
                outputStream.close();
                outputStream = null;
            }

            if (socket != null)
            {
                socket.Close();
                socket = null;
            }
        }

        public bool connect()
        {
            if (isConnected())
                return true;

            socket = new Socket(AddressFamily.InterNetwork, SocketType.Stream, ProtocolType.Tcp);
            socket.NoDelay = true;
            //socket.setTcpNoDelay(true);
            socket.ReceiveTimeout = readTimeout;
            socket.SendTimeout = readTimeout;
            //socket.setSoTimeout(VoiceMessageParam.SOCKET_READ_TIMEOUT);
            //SocketAddress address = new InetSocketAddress(VoiceMessageParam.voiceIp, VoiceMessageParam.voicePort);

            IPAddress ipadr = IPAddress.Parse(ip);
            IPEndPoint iep = new IPEndPoint(ipadr, port);
            GLog.d("---------------------开始连接上传服务器 ip:" + ip + " port:" + port + "---------------------");
            try
            {
                socket.Connect(iep);
            }
            catch (Exception e)
            {
                socket.Close();
                socket = null;
                GLog.d("---------------------上传服务器2连接失败---------------------");
                return false;
            }

            GLog.d("连接服务器成功");

            inputStream = new SocketInputStream(socket);
            outputStream = new SocketOutputStream(socket);
           // inputStream = socket.getInputStream();
            //outputStream = socket.getOutputStream();
            dos = new LittleEndianDataOutputStream(outputStream);
            dis = new LittleEndianDataInputStream(inputStream);
            return true;
        }

        public bool sendCmd(GServerCmd cmd)
        {
            return cmd.write(dos);
        }

        public bool receiveCmd(GServerCmd cmd)
        {
            return cmd.read(dis);
        }

        public int read(byte[] buffer)
        {
            return read(buffer, 0, buffer.Length);
        }

        public int read(byte[] buffer, int byteOffset, int byteCount)
        {
            return inputStream.Read(buffer);
        }

        public void write(byte[] buffer)
        {
            write(buffer, 0, buffer.Length);
        }

        public void write(byte[] buffer, int offset, int count)
        {
            outputStream.write(buffer);
            outputStream.flush();
        }

        public void setConnectTimeout(int timeoutMillis)
        {
            connectTimeout = timeoutMillis;
            //urlConnection.setConnectTimeout(timeoutMillis);
        }

        public void setReadTimeout(int timeoutMillis)
        {
            readTimeout = timeoutMillis;
            //urlConnection.setReadTimeout(timeoutMillis);
        }
    }
}