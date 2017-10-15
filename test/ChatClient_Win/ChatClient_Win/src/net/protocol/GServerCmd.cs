using System;
using GTech.IO;
using GTech.Utils;

namespace GTech.Net.Protocol
{
    public enum MsgId : ushort
    {
        MsgId_Tick = 1000,
        MsgId_Error,

        MsgId_Echo,

        MsgId_ReqLogin,
        MsgId_RetLogin,

        MsgId_ReqAppLogin,
        MsgId_RetAppLogin,

        MsgId_ReqTokenLogin,
        MsgId_RetTokenLogin,

        MsgId_ReqToken,
        MsgId_RetToken,

        MsgId_ReqLoginOut,
        MsgId_ReqRetLoginOut,

        MsgId_ReqFriendList,
        MsgId_RetFriendList,

        MsgId_ReqFriendAdd,
        MsgId_RetFriendAdd,

        MsgId_FriendReqAgree,
        MsgId_FriendReq,

        MsgId_FriendReqResult,

        MsgId_ReqFriendDel,
        MsgId_RetFriendDel,

        MsgId_ReqUserToBlack,
        MsgId_RetUserToBlack,

        MsgId_ReqRemoveUserInBlack,
        MsgId_RetRemoveUserInBlack,

        MsgId_ReqMoveFriendToGroup,
        MsgId_RetMoveFriendToGroup,

        MsgId_ReqSetFriendVerifyType,
        MsgId_RetSetFriendVerifyType,

        MsgId_Message,
        MsgId_RetMessage,

        MsgId_ReqUserInfo,
        MsgId_RetUserInfo,

        MsgId_ReqUserSearch,
        MsgId_RetUserSearch,

        MsgId_End,
    };

    public class GServerCmd : IServerCmd
    {
        //protected UInt16 size;
        protected ushort msgId;

        public GServerCmd()
        {
            //header = getLength() - 4;
        }

        public GServerCmd(byte[] buff)
        {
            //header = getLength() - 4;
            msgId = System.BitConverter.ToUInt16(buff, 0);
        }

        public virtual bool read(LittleEndianDataInputStream dis)
        {
            //size = dis.readUShort();
            msgId = dis.readUShort();

            return true;
        }

        public virtual bool write(LittleEndianDataOutputStream dos)
        {
            ushort size = (ushort)(getLength());

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
            ushort length = (ushort)getLength();
            byte[] data = new byte[length + 2];
            System.Array.Copy(System.BitConverter.GetBytes(length), 0, data, 0, 2);
            System.Array.Copy(System.BitConverter.GetBytes(msgId), 0, data, 2, 2);

            return data;
        }

        public virtual string toString()
        {
            // TODO Auto-generated method stub
            return "msgId=" + msgId;
        }

    }
}
