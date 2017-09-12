using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace GTech.Threading
{
    public abstract class Runnable
    {
        public static void Execute(object obj)
        {
            Runnable runnable = obj as Runnable;
            runnable.run();
        }

        abstract public void run();
    }
}
