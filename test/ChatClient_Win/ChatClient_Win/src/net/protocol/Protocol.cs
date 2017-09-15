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
            msgId = 0;
        }

        public override bool read(LittleEndianDataInputStream dis)
        {
            base.read(dis);

            return true;
        }

        public override bool write(LittleEndianDataOutputStream dos)
        {
            base.write(dos);

            return true;
        }

        public override int getLength()
        {
            // TODO Auto-generated method stub
            return base.getLength();
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

    public class MsgEcho : GServerCmd
    {
        public MsgEcho()
        {
            msgId = 0;
        }

        public override bool read(LittleEndianDataInputStream dis)
        {
            base.read(dis);

            return true;
        }

        public override bool write(LittleEndianDataOutputStream dos)
        {
            base.write(dos);

            return true;
        }

        public override int getLength()
        {
            // TODO Auto-generated method stub
            return base.getLength();
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
}
