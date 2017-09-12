using System;
using TestGCloud.io;
using TestGCloud.utils;

namespace TestGCloud.net.protocol
{
    public class GRtnServerUploadCmd : GRtnServerCmd
    {
        private byte byRet = 1;
        private short size;
        private byte[] url;

        public override bool read(LittleEndianDataInputStream dis)
        {
            // TODO Auto-generated method stub
            base.read(dis);

            //LittleEndianDataInputStream dis = new LittleEndianDataInputStream(is);

            byRet = dis.readByte();
            size = dis.readShort();

            byte[] databytes = new byte[size];
            dis.read(databytes);
            CryptUtils.CryptBuff(databytes, databytes.Length);
            //url = GByteUtils.bytesToChars(databytes);
            url = databytes;// new byte[databytes.Length];
            //Buffer.BlockCopy(databytes, 0, url, 0, databytes.Length);
            //Array.Copy(databytes, url, databytes.Length);
            //GCloudLog.d(toString());
            return true;
        }

        public override int getLength()
        {
            // TODO Auto-generated method stub
            if (url == null)
                return base.getLength() + 1 + 2;
            else
                return base.getLength() + 1 + 2 + url.Length;
        }

        public override bool isSuccess()
        {
            // TODO Auto-generated method stub
            return byRet == 0;
        }

        public override string toString()
        {
            // TODO Auto-generated method stub
            if (url == null)
                return base.toString() + " byRet=" + byRet + " size=" + size + " url=null";
            else
                return base.toString() + " byRet=" + byRet + " size=" + size + " url=" + System.Text.Encoding.Default.GetString(url);
        }

        public string getUrl()
        {
            if (url == null)
                return "url is null";
            else
                return System.Text.Encoding.Default.GetString(url);
        }
    }
}