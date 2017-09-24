using GTech.Log;
using GTech.Net;
using System;
using System.Collections.Generic;
using System.Configuration;
using System.Data;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using System.Windows;
using System.Windows.Navigation;

namespace ChatClient_Win
{
    public class UserData
    {
        public ulong ID { get; set; }
        public string Nickname { get; set; }
        public string Desc { get; set; }
        public string HeadID { get; set; }
        public bool IsOnline { get; set; }
    }

    public class GroupData
    {
        public string Name { get; set; }
        public int OnlineCount { get; set; }
        public int AllCount { get; set; }

        public List<UserData> UserList { get; set; }
    }

    public class PlayerData
    {
        public ulong ID { get; set; }
        public string Nickname { get; set; }
        public string Desc { get; set; }
        public string HeadID { get; set; }

        public List<GroupData> GroupList { get; set; }
    }

    /// <summary>
    /// App.xaml 的交互逻辑
    /// </summary>
    public partial class App : Application
    {
        ChatClient chatClient;
        PlayerData playerData;
        Timer tickTimer;

        protected override void OnStartup(StartupEventArgs e)
        {
            GLog.d("OnStartup");
            chatClient = new ChatClient("127.0.0.1:9090");
            chatClient.LoginedHandler += OnLogined;
            chatClient.ErrorHandler += OnError;
            chatClient.CloseHandler += OnClose;
        }

        protected override void OnLoadCompleted(NavigationEventArgs e)
        {
            GLog.d("OnLoadCompleted");
        }

        protected override void OnExit(ExitEventArgs e)
        {
            GLog.d("OnExit");
        }

        public void Login(ulong uid, string password)
        {
            chatClient.Login(uid, password);
        }

        private void SendTick(Object data)
        {
            chatClient.SendTick();
        }

        private void OnLogined()
        {
            GLog.d("OnLogined");
            playerData = new PlayerData();

            tickTimer = new Timer(new TimerCallback(SendTick), null, 0, 30000);
        }

        private void OnError(int errorcode)
        {
            GLog.e("error, errorcode:" + errorcode);
        }

        private void OnClose()
        {
            GLog.e("OnClose");

            if(tickTimer != null)
            {
                tickTimer.Dispose();
                tickTimer = null;
            }
        }
    }
}
