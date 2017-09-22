using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace GTech.Log
{
    public class GLog
    {
        static string sTag = "GLog:";
        static bool sIsDebug = true;

        public static bool isDebug()
        {
            return sIsDebug;
        }

        public static void setDebug(bool isDebug)
        {
            sIsDebug = isDebug;
        }

        public static void setLogTag(string tag)
        {
            sTag = tag;
        }

        public static void d(object obj)
        {
            string msg = null;
            if (obj == null)
            {
                msg = "Obj is null.";
            }
            else {
                msg = obj.ToString();
            }
            d(msg);
        }

        public static void d(string msg)
        {
            //return;
            if (sIsDebug)
            {
#if DEBUG
                System.Diagnostics.Debug.WriteLine(sTag + msg);
#endif
                Console.WriteLine(sTag + msg);
            }
        }

        public static void i(string msg)
        {
            if (sIsDebug)
            {
#if DEBUG
                System.Diagnostics.Debug.WriteLine(sTag + msg);
#endif
                Console.WriteLine(sTag + msg);
            }
        }

        public static void i(object obj)
        {
            string msg = null;
            if (obj == null)
            {
                msg = "Obj is null.";
            }
            else {
                msg = obj.ToString();
            }
            i(msg);
        }

        public static void w(string msg)
        {
            if (sIsDebug)
            {
#if DEBUG
                System.Diagnostics.Debug.WriteLine(sTag + msg);
#endif
                Console.WriteLine(sTag + msg);
            }
        }

        public static void w(object obj)
        {
            string msg = null;
            if (obj == null)
            {
                msg = "Obj is null.";
            }
            else {
                msg = obj.ToString();
            }
            w(msg);
        }

        // 	public static void i(string format, Object... args) {
        // 		StringBuffer sb = new StringBuffer();
        // 		Formatter f = new Formatter(sb, Locale.getDefault());
        // 		f.format(format, args);
        // 		f.close();
        // 		i(sb);
        // 	}

        // 	public static void m(Class<?> clazz, string method) {
        // 		m(clazz.getName(), method);
        // 	}
        // 
        // 	public static void m(Object obj, string method) {
        // 		if (sIsDebug) {
        // 			int len = obj.toString().length();
        // 			StringBuffer sb = new StringBuffer();
        // 			if (len < LOG_WIDTH_LENGTH) {
        // 				for (int i = 0; i < LOG_WIDTH_LENGTH - len; i++) {
        // 					sb.append(".");
        // 				}
        // 			}
        // 			i(obj + sb.toString() + method);
        // 		}
        // 	}
        // 
        // 	public static void md(Class<?> clazz, string method) {
        // 		md(clazz.getName(), method);
        // 	}
        // 
        // 	public static void md(Object obj, string method) {
        // 		if (sIsDebug) {
        // 			int len = obj.toString().length();
        // 			StringBuffer sb = new StringBuffer();
        // 			if (len < LOG_WIDTH_LENGTH) {
        // 				for (int i = 0; i < LOG_WIDTH_LENGTH - len; i++) {
        // 					sb.append(".");
        // 				}
        // 			}
        // 			d(obj + sb.toString() + method);
        // 		}
        // 	}
        // 
        // 	public static void logInObj(Object obj, string msg) {
        // 		Class<?> cls = obj.getClass();
        // 		string name = cls.getSimpleName();
        // 		i(name + "--->%s", msg);
        // 	}
        // 
        // 	public static void logInObj(Object obj, string format, Object... args) {
        // 		Class<?> cls = obj.getClass();
        // 		string name = cls.getSimpleName();
        // 		format = name + "----->" + format;
        // 		i(format, args);
        // 	}

        public static void e(string msg)
        {
#if DEBUG
            System.Diagnostics.Debug.WriteLine(sTag + msg);
#endif
            Console.WriteLine(sTag + msg);
        }

        public static void log(string msg)
        {
#if DEBUG
            System.Diagnostics.Debug.WriteLine(msg);
#endif
            Console.WriteLine(msg);
        }
    }
}
