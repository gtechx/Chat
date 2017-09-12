using System;
using TestGCloud.io;

namespace TestGCloud.net.protocol
{
    public class GReqServerTimeCmd : GServerCmd
    {
        private int clientVer = 0;

        public GReqServerTimeCmd(int cv)
        {
            clientVer = cv;
            scmd = (byte)(eVMP_ReqServerTimeMessage & 0xff);
        }

        public override bool write(LittleEndianDataOutputStream dos)
        {
            // TODO Auto-generated method stub
            base.write(dos);
            //LittleEndianDataOutputStream dos = new LittleEndianDataOutputStream(os);

            dos.writeInt(clientVer);
            dos.flush();
            return true;
        }

        public override int getLength()
        {
            // TODO Auto-generated method stub
            return base.getLength() + 4;
        }

        public override string toString()
        {
            // TODO Auto-generated method stub
            return base.toString() + " clientVer=" + clientVer;
        }
    }
}