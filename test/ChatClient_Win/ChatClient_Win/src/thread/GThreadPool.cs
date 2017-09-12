using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading;
using System.Threading.Tasks;

namespace GTech.Threading
{
    public class GThreadPool
    {
        int threadCount;
        static bool isInit = false;

        public GThreadPool(int threadnum)
        {
            threadCount = threadnum;

            if (!isInit)
            {
                int workerThreadnum;
                int portThreadnum;
                ThreadPool.GetMaxThreads(out workerThreadnum, out portThreadnum);
                ThreadPool.SetMaxThreads(workerThreadnum, portThreadnum);
                isInit = true;
            }
        }

        public void Execute(Runnable runnable)
        {
            //ThreadPool.QueueUserWorkItem(new WaitCallback(Runnable.Execute), runnable);
            Thread thread = new Thread(new ParameterizedThreadStart(Runnable.Execute));
            thread.IsBackground = true;
            thread.Start(runnable);  //启动异步线程 
        }

        public void Destroy()
        {
            //TODO
        }
    }
}
