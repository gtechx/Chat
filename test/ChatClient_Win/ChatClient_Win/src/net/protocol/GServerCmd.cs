using System;
using GTech.IO;
using GTech.Utils;

namespace GTech.Net.Protocol
{
    public class GServerCmd : IServerCmd
    {
        protected UInt16 size;
        protected UInt16 msgId;

        public GServerCmd()
        {
            //header = getLength() - 4;
        }

        public virtual bool read(LittleEndianDataInputStream dis)
        {
            size = dis.readUnsignedShort();
            msgId = dis.readUnsignedShort();

            return true;
        }

        public virtual bool write(LittleEndianDataOutputStream dos)
        {
            size = (ushort)(getLength() - 4);

            dos.writeUShort(size);
            dos.writeUShort(msgId);

            return true;
        }

        public virtual int getLength()
        {
            // TODO Auto-generated method stub
            return 2;
        }

        public virtual byte[] toBytes()
        {
            throw new NotImplementedException();
        }

        public virtual string toString()
        {
            // TODO Auto-generated method stub
            return " size=" + size + " msgId=" + msgId;
        }

    }
}
