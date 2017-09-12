using System;
using GTech.IO;
using GTech.Utils;

namespace GTech.Net.Protocol
{
    public class GServerCmd : IServerCmd
    {

        //大命令号
        protected static int VOICE_MESSAGE_USERCMD = 252;

        //小命令号
        protected static int eVMP_ReqServerTimeMessage = 1; //获取服务器时间
        protected static int eVMP_RtnServerTimeMessage = 2;
        protected static int eVMP_ReqServerVerifyMessage = 3; //验证
        protected static int eVMP_RtnServerVerifyMessage = 4;
        protected static int eVMP_ReqServerUploadMessage = 5; //上传
        protected static int eVMP_RtnServerUploadMessage = 6;

        protected int header;
        protected byte lcmd = (byte)VOICE_MESSAGE_USERCMD;//large cmd
        protected byte scmd;//small cmd
        protected int timeStamp;

        public GServerCmd()
        {
            //header = getLength() - 4;
        }

        public virtual bool read(LittleEndianDataInputStream dis)
        {
            header = dis.readInt();
            lcmd = dis.readByte();
            scmd = dis.readByte();
            timeStamp = dis.readInt();

            return true;
        }

        public virtual bool write(LittleEndianDataOutputStream dos)
        {
            header = getLength() - 4;
            timeStamp = (int)(TimeUtils.GetTimeStamp() / 1000); //转化为unix时间戳
                                                                //GCloudLog.d(toString());
            dos.writeInt(header);
            dos.writeByte(lcmd);
            dos.writeByte(scmd);
            dos.writeInt(timeStamp);

            return true;
        }

        public virtual int getLength()
        {
            // TODO Auto-generated method stub
            return 4 + 1 + 1 + 4;
        }

        public virtual byte[] toBytes()
        {
            throw new NotImplementedException();
        }

        public virtual string toString()
        {
            // TODO Auto-generated method stub
            return " header=" + header + " scmd=" + scmd;
        }

    }
}
