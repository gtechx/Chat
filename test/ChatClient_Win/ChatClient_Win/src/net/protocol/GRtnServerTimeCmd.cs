using System;
using TestGCloud.io;
using TestGCloud.utils;

namespace TestGCloud.net.protocol
{
    public class GRtnServerTimeCmd : GRtnServerCmd
    {
        private byte byRet = 1;
        private long curTime;

        public override bool read(LittleEndianDataInputStream dis)
        {
            // TODO Auto-generated method stub
            base.read(dis);

            //LittleEndianDataInputStream dis = new LittleEndianDataInputStream(is);

            byRet = dis.readByte();
            curTime = dis.readLong();
            //GCloudLog.d(toString());
            return true;
        }

        public override int getLength()
        {
            // TODO Auto-generated method stub
            return base.getLength() + 1 + 8;
        }

        public override bool isSuccess()
        {
            // TODO Auto-generated method stub
            return byRet == 0;
        }

        public long getServerTime()
        {
            return curTime;
        }

        public override string toString()
        {
            // TODO Auto-generated method stub
            return base.toString() + " byRet=" + byRet + " curTime=" + curTime;
        }
    }
}