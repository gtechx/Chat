using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using GTech.IO;

namespace GTech.Utils
{
    public class FileUtils
    {
        static string FILE_EXTENSION_SEPARATOR = ".";
        static string FILE_SEPARATOR = "/";

        public static string GetFileNameWithoutExtension(string filePath)
        {
            return Path.GetFileNameWithoutExtension(filePath);
        }

        public static string GetFilePath(string filePath)
        {
            string filename = GetFileName(filePath);
            return filePath.Substring(0, filePath.Length - filename.Length);
        }

        public static string GetFileName(string filePath)
        {
            return Path.GetFileName(filePath);
        }

        public static string GetExtension(string filePath)
        {
            return Path.GetExtension(filePath);
        }

        public static bool MakeDir(string filePath)
        {
            try {
                Directory.CreateDirectory(filePath);
                return true;
            }catch(Exception)
            {
                return false;
            }
        }

        public static bool DeleteDir(string filePath)
        {
            DirectoryInfo dir = new DirectoryInfo(filePath);
            if (dir.Exists)
            {
                DirectoryInfo[] childs = dir.GetDirectories();
                foreach (DirectoryInfo child in childs)
                {
                    child.Delete(true);
                }
                dir.Delete(true);
            }

            return true;
        }

        public static bool IsFileExist(string filePath)
        {
            return File.Exists(filePath);
        }

        public static bool IsDirExist(string Path)
        {
            return Directory.Exists(Path);
        }

        public static bool DeleteFile(string path)
        {
            try
            {
                File.Delete(path);
                return true;
            }
            catch (Exception)
            {
                return false;
            }
        }

        public static bool CopyFile(string oldPath, string newPath)
        {
            try
            {
                File.Copy(oldPath, newPath);
                return true;
            }
            catch (Exception)
            {
                return false;
            }
        }

        public static bool RenameFile(string srcfile, string desfile)
        {
            try
            {
                File.Move(srcfile, desfile);
                return true;
            }
            catch (Exception)
            {
                return false;
            }
        }

        public static long GetFileSize(string file)
        {
            FileInfo fileInfo = new FileInfo(file);
            return fileInfo.Length;
        }

        public static string ReadFile(InputStream stream)
        {
            byte[] buffer = new byte[1024];
            int len = -1;
            string result = "";

            while ((len = stream.Read(buffer)) > 0)
            {
                result += System.Text.Encoding.Default.GetString(buffer, 0, len);
            }

            return result;
        }

        public static byte[] ReadFile(string file)
        {
            byte[] buffer;
            long filesize = GetFileSize(file);
            buffer = new byte[filesize];

            FileStream fs = new FileStream(file, FileMode.Open, FileAccess.Read);

            if(fs != null)
            {
                fs.Read(buffer, 0, (int)filesize);
                fs.Close();
            }

            return buffer;
        }

        public static string ReadFileAsText(string file)
        {
            byte[] buffer;
            long filesize = GetFileSize(file);
            buffer = new byte[filesize];

            FileStream fs = new FileStream(file, FileMode.Open, FileAccess.Read);

            if (fs != null)
            {
                fs.Read(buffer, 0, (int)filesize);
                fs.Close();
            }

            return System.Text.Encoding.Default.GetString(buffer, 0, (int)filesize);
        }

        public static bool WriteFile(string file, byte[] buffer)
        {
            FileStream fs = new FileStream(file, FileMode.OpenOrCreate, FileAccess.Write);
            if (fs != null)
            {
                fs.Write(buffer, 0, buffer.Length);
                fs.Close();
            }

            return true;
        }

        public static bool WriteFile(string file, InputStream stream)
        {
            FileStream fs = new FileStream(file, FileMode.OpenOrCreate, FileAccess.Write);
            if (fs != null)
            {
                byte[] buffer = new byte[1024];
                int len = -1;

                while((len = stream.Read(buffer)) > 0)
                {
                    fs.Write(buffer, 0, len);
                }
                fs.Flush();
                fs.Close();
            }
            return true;
        }

        public static bool WriteFile(string file, string content, bool bAppend = false)
        {
            FileMode type = bAppend ? FileMode.OpenOrCreate | FileMode.Append : FileMode.OpenOrCreate;

            if (IsFileExist(file))
            {
                type = bAppend ? FileMode.Append : FileMode.Truncate;
            } 
            else
            {
                type = FileMode.Create;
            }

            FileStream fs = new FileStream(file, type, FileAccess.Write);

            if (fs != null)
            {
                StreamWriter sw = new StreamWriter(fs);
                sw.Write(content);
                sw.Flush();
                sw.Close();
                fs.Close();
            }
            return true;
        }
    }
}
