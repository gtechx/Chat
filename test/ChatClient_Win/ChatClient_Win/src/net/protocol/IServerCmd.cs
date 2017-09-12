using System;
using GTech.IO;

namespace GTech.Net.Protocol
{
    public interface IServerCmd
    {
        bool read(LittleEndianDataInputStream dis);
        bool write(LittleEndianDataOutputStream dos);
        int getLength();
        byte[] toBytes();
        string toString();
    }
}