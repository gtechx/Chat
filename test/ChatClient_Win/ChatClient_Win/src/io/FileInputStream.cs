using Microsoft.Win32.SafeHandles;
using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Security.AccessControl;
using System.Text;
using System.Threading.Tasks;

namespace GTech.IO
{
    public class FileInputStream : FileStream
    {
        public FileInputStream(string file) : base(file, FileMode.Open, FileAccess.Read) { }

        public FileInputStream(string path, FileMode mode): base(path, mode){}

        public FileInputStream(SafeFileHandle handle, FileAccess access) : base(handle, access) { }

        public FileInputStream(string path, FileMode mode, FileAccess access) : base(path, mode, access) { }

        public FileInputStream(SafeFileHandle handle, FileAccess access, int bufferSize) : base(handle, access, bufferSize) { }

        public FileInputStream(string path, FileMode mode, FileAccess access, FileShare share) : base(path, mode, access, share) { }

        public FileInputStream(SafeFileHandle handle, FileAccess access, int bufferSize, bool isAsync) : base(handle, access, bufferSize, isAsync) { }

        public FileInputStream(string path, FileMode mode, FileAccess access, FileShare share, int bufferSize) : base(path, mode, access, share, bufferSize) { }

        public FileInputStream(string path, FileMode mode, FileAccess access, FileShare share, int bufferSize, bool useAsync) : base(path, mode, access, share, bufferSize, useAsync) { }

        public FileInputStream(string path, FileMode mode, FileAccess access, FileShare share, int bufferSize, FileOptions options) : base(path, mode, access, share, bufferSize, options) { }

        public FileInputStream(string path, FileMode mode, FileSystemRights rights, FileShare share, int bufferSize, FileOptions options) : base(path, mode, rights, share, bufferSize, options) { }

        public FileInputStream(string path, FileMode mode, FileSystemRights rights, FileShare share, int bufferSize, FileOptions options, FileSecurity fileSecurity) : base(path, mode, rights, share, bufferSize, options, fileSecurity) { }
    }
}
