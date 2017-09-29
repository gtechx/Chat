using GTech.Log;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Windows;
using System.Windows.Controls;
using System.Windows.Data;
using System.Windows.Documents;
using System.Windows.Input;
using System.Windows.Media;
using System.Windows.Media.Imaging;
using System.Windows.Navigation;
using System.Windows.Shapes;

namespace ChatClient_Win
{
    public class FData
    {
        public string Name { get; set; }
        public string Desc { get; set; }
    }

    public class Data
    {
        public string Name { get; set; }

        List<FData> _FList = new List<FData>();
        public List<FData> FList { get { return _FList; } }
    }

    /// <summary>
    /// Friends.xaml 的交互逻辑
    /// </summary>
    public partial class Friends : Page
    {
        public Friends()
        {
            InitializeComponent();
        }

        private void Page_Loaded(object sender, RoutedEventArgs e)
        {
            List<Data> datalist = new List<Data>();
            for (int i = 0; i < 5; i++)
            {
                Data data = new Data();
                data.Name = "aaa" + i;

                for (int j = 0; j < 3; j++)
                {
                    FData fdata = new FData();
                    fdata.Name = "bbb" + j;
                    fdata.Desc = "desc" + j;
                    data.FList.Add(fdata);
                }
                datalist.Add(data);
            }

            treeView.ItemsSource = datalist;
        }

        private void Button_Click(object sender, RoutedEventArgs e)
        {
            FriendAdd fawin = new FriendAdd();

            if (fawin.ShowDialog() == true)
            {
                GLog.d("friend sending...");
            }
        }
    }
}
