using System;
using TestGCloud.io;
using TestGCloud.utils;

namespace TestGCloud.net.protocol
{
    public class GReqServerUploadCmd : GServerCmd
    {
        private short size;
        private byte byFlag;
        private byte[] data;

        public GReqServerUploadCmd(short siz, byte flag, byte[] dat)
        {
            size = siz;
            byFlag = flag;
            data = dat;
            scmd = (byte)(eVMP_ReqServerUploadMessage & 0xff);
        }

        public override bool write(LittleEndianDataOutputStream dos)
        {
            // TODO Auto-generated method stub
            base.write(dos);

            //LittleEndianDataOutputStream dos = new LittleEndianDataOutputStream(os);

            dos.writeShort(size);
            dos.writeByte(byFlag);
            dos.write(data);
            dos.flush();
            return true;
        }

        public override int getLength()
        {
            // TODO Auto-generated method stub
            if (data == null)
                return base.getLength() + 2 + 1;
            else
                return base.getLength() + 2 + 1 + data.Length;
        }

        public override string toString()
        {
            // TODO Auto-generated method stub
            if (data == null)
                return base.toString() + " size=" + size + " byFlag=" + byFlag + " data=null";
            else
                return base.toString() + " size=" + size + " byFlag=" + byFlag + " data=" + ByteUtils.ToString(data);
        }
    }
}