using GTech.IO;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace GTech.Net.Protocol
{
    public class MsgTick : GServerCmd
    {
        public MsgTick()
        {
            msgId = (ushort)MsgId.MsgId_Tick;
        }
    }

    public class MsgEcho : GServerCmd
    {
        public byte[] Data;

        public MsgEcho()
        {
            msgId = (ushort)MsgId.MsgId_Echo;
        }

        public MsgEcho(byte[] buff) : base(buff)
        {
            int datalength = buff.Length - 2;
            this.Data = new byte[datalength];
            Array.Copy(buff, 2, this.Data, 0, datalength);
        }

        public override bool read(LittleEndianDataInputStream dis)
        {
            base.read(dis);

            return true;
        }

        public override bool write(LittleEndianDataOutputStream dos)
        {
            base.write(dos);
            dos.write(Data);

            return true;
        }

        public override int getLength()
        {
            // TODO Auto-generated method stub
            return base.getLength() + Data.Length;
        }

        public override byte[] toBytes()
        {
            throw new NotImplementedException();
        }

        public override string toString()
        {
            // TODO Auto-generated method stub
            return base.toString();
        }
    }

    public class MsgReqLogin : GServerCmd
    {
        public UInt64 Uid;
        public byte[] Password;

        public MsgReqLogin()
        {
            msgId = (ushort)MsgId.MsgId_ReqLogin;
        }

        public MsgReqLogin(byte[] buff) : base(buff)
        {
            this.Uid = System.BitConverter.ToUInt64(buff, 0);
            int datalength = buff.Length - 10;
            this.Password = new byte[datalength];
            Array.Copy(buff, 2 + 8, this.Password, 0, datalength);
        }

        public override bool read(LittleEndianDataInputStream dis)
        {
            base.read(dis);

            return true;
        }

        public override bool write(LittleEndianDataOutputStream dos)
        {
            base.write(dos);
            dos.writeULong(Uid);
            dos.write(Password);

            return true;
        }

        public override int getLength()
        {
            // TODO Auto-generated method stub
            return base.getLength() + 8 + Password.Length;
        }

        public override byte[] toBytes()
        {
            ushort length = (ushort)getLength();
            byte[] data = new byte[length + 2];
            System.Array.Copy(System.BitConverter.GetBytes(length), 0, data, 0, 2);
            System.Array.Copy(System.BitConverter.GetBytes(msgId), 0, data, 2, 2);
            System.Array.Copy(System.BitConverter.GetBytes(Uid), 0, data, 4, 8);
            System.Array.Copy(Password, 0, data, 12, Password.Length);

            return data;
        }

        public override string toString()
        {
            // TODO Auto-generated method stub
            return base.toString();
        }
    }
}
