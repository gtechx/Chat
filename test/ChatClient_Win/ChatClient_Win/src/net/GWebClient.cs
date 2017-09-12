using System;
using System.IO;
using System.Net;
using System.Net.Sockets;
using GTech.IO;

namespace GTech.Net
{
    public class GWebClient
    {
        string url;
        HttpStream inputStream = null;
        StreamReader streamReader;

        private int readTimeout = 30000;
        private int connectTimeout = 30000;

        public GWebClient(string addr)
        {
            url = addr;
        }

        public bool close()
        {
            if (streamReader != null)
            {
                streamReader.Close();
                streamReader = null;
            }

            if (inputStream != null)
            {
                inputStream.Close();
                inputStream = null;
            }

            return true;
        }

        public InputStream openStream()
        {
            HttpWebRequest req = WebRequest.Create(url) as HttpWebRequest;
            HttpWebResponse res = req.GetResponse() as HttpWebResponse;
            Stream st = res.GetResponseStream();
            inputStream = new HttpStream(st);
            streamReader = new StreamReader(st);

            return inputStream;
        }

        public int read(byte[] buffer)
        {
            return read(buffer, 0, buffer.Length);
        }

        public int read(byte[] buffer, int byteOffset, int byteCount)
        {
            return inputStream.Read(buffer, 0, buffer.Length);
        }

        public string readLine()
        {
            return streamReader.ReadLine();
        }

        //	public InputStream getInputStream(){
        //		try {
        //			return urlConnection.getInputStream();
        //		}catch(Exception e){
        //			e.printStackTrace();
        //			return null;
        //		}
        //    }
        //	
        //	public OutputStream getOutputStream(){
        //		try {
        //			return urlConnection.getOutputStream();
        //		}catch(Exception e){
        //			e.printStackTrace();
        //			return null;
        //		}
        //    }

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