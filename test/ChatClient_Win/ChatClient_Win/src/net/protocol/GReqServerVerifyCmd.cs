using System;
using TestGCloud.io;
using TestGCloud.utils;

namespace TestGCloud.net.protocol
{
    public class GReqServerVerifyCmd : GServerCmd
    {
        private long userID;
        private int gameID;
        private int zoneID;
        private byte[] extensionName = new byte[16];
        private int fileSize;
        private byte[] token = new byte[33];

        public GReqServerVerifyCmd(long UID, int GID, int ZID, byte[] extName, int fsize, byte[] tok)
        {
            userID = UID;
            gameID = GID;
            zoneID = ZID;
            //extensionName = extName;
            Array.Clear(extensionName, 0, extensionName.Length);
            //Arrays.fill(extensionName, '\0'); 
            Array.Copy(extName, 0, extensionName, 0, extName.Length);
            fileSize = fsize;
            //token = tok;
            Array.Clear(token, 0, token.Length);
            //Arrays.fill(token, '\0');
            Array.Copy(tok, 0, token, 0, tok.Length);

            scmd = (byte)(eVMP_ReqServerVerifyMessage & 0xff);
        }

        public override bool write(LittleEndianDataOutputStream dos)
        {
            // TODO Auto-generated method stub
            base.write(dos);

            //LittleEndianDataOutputStream dos = new LittleEndianDataOutputStream(os);
            int bytecount = getLength() - base.getLength();
            byte[] buffer = new byte[bytecount];

            byte[] tmp = ByteUtils.GetBytes(userID);
            Array.Copy(tmp, 0, buffer, 0, tmp.Length);
            tmp = ByteUtils.GetBytes(gameID);
            Array.Copy(tmp, 0, buffer, 8, tmp.Length);
            tmp = ByteUtils.GetBytes(zoneID);
            Array.Copy(tmp, 0, buffer, 12, tmp.Length);
            //tmp = ByteUtils.GetBytes(extensionName);
            Array.Copy(extensionName, 0, buffer, 16, extensionName.Length);
            //Array.Copy(extensionName, 0, buffer, 16, tmp.Length);
            tmp = ByteUtils.GetBytes(fileSize);
            Array.Copy(tmp, 0, buffer, 32, tmp.Length);
            //tmp = GByteUtils.charsToBytes(token);
            Array.Copy(token, 0, buffer, 36, token.Length);
            //Array.Copy(token, 0, buffer, 36, tmp.Length);
            //GByteUtils.dumpMemory("buffer before crypt", buffer);
            CryptUtils.CryptBuff(buffer, buffer.Length);
            //GByteUtils.dumpMemory("buffer after crypt", buffer);
            dos.write(buffer);
            //		dos.writeLong(userID);
            //		dos.writeInt(gameID);
            //		dos.writeInt(zoneID);
            //		dos.writeChars(extensionName);
            //		dos.writeInt(fileSize);
            //		dos.writeChars(token);
            dos.flush();
            return true;
        }

        public override int getLength()
        {
            // TODO Auto-generated method stub
            return base.getLength() + 8 + 4 + 4 + 16 + 4 + 33;
        }

        public override string toString()
        {
            // TODO Auto-generated method stub
            return base.toString() + " userID=" + userID + " gameID=" + gameID + " zoneID=" + zoneID
                    + " extensionName=" + ByteUtils.ToString(extensionName) + " fileSize=" + fileSize + " token=" + ByteUtils.ToString(token);
        }
    }
}