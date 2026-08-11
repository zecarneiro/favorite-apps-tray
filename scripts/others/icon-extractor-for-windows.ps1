# Source code from https://www.powershellgallery.com/packages/IconExport/1.0.1/Content/IconExport.psm1
$file = "${env:FAT_ICON_EXTRACTOR_FILE}"
$dest = "${env:FAT_ICON_EXTRACTOR_DEST}"
$code = '
using System;
using System.Drawing;
using System.Runtime.InteropServices;
using System.IO;

namespace System {
    public class IconExtractor {
        public static Icon Extract(string file, int number, bool largeIcon) {
            IntPtr large;
            IntPtr small;
            ExtractIconEx(file, number, out large, out small, 1);
            try { return Icon.FromHandle(largeIcon ? large : small); }
            catch { return null; }
        }
        [DllImport("Shell32.dll", EntryPoint = "ExtractIconExW", CharSet = CharSet.Unicode, ExactSpelling = true, CallingConvention = CallingConvention.StdCall)]
        private static extern int ExtractIconEx(string sFile, int iIndex, out IntPtr piLargeVersion, out IntPtr piSmallVersion, int amountIcons);
    }
}

public class PngIconConverter
{
    public static bool Convert(System.Drawing.Bitmap input_bit, string output_icon, int size, bool keep_aspect_ratio = false)
    {
        System.IO.Stream output_stream = new System.IO.FileStream(output_icon, System.IO.FileMode.OpenOrCreate);
        if (input_bit != null)
        {
            int width, height;
            if (keep_aspect_ratio)
            {
                width = size;
                height = input_bit.Height / input_bit.Width * size;
            }
            else
            {
                width = height = size;
            }
            System.Drawing.Bitmap new_bit = new System.Drawing.Bitmap(input_bit, new System.Drawing.Size(width, height));
            if (new_bit != null)
            {
                System.IO.MemoryStream mem_data = new System.IO.MemoryStream();
                new_bit.Save(mem_data, System.Drawing.Imaging.ImageFormat.Png);

                System.IO.BinaryWriter icon_writer = new System.IO.BinaryWriter(output_stream);
                if (output_stream != null && icon_writer != null)
                {
                    icon_writer.Write((byte)0);
                    icon_writer.Write((byte)0);
                    icon_writer.Write((short)1);
                    icon_writer.Write((short)1);
                    icon_writer.Write((byte)width);
                    icon_writer.Write((byte)height);
                    icon_writer.Write((byte)0);
                    icon_writer.Write((byte)0);
                    icon_writer.Write((short)0);
                    icon_writer.Write((short)32);
                    icon_writer.Write((int)mem_data.Length);
                    icon_writer.Write((int)(6 + 16));
                    icon_writer.Write(mem_data.ToArray());
                    icon_writer.Flush();
                    return true;
                }
            }
            return false;
        }
        return false;
    }
}'
Add-Type -TypeDefinition $code -ReferencedAssemblies System.Drawing, System.IO -ErrorAction SilentlyContinue
$icon=[System.Drawing.Icon]::ExtractAssociatedIcon("${file}")
[PngIconConverter]::Convert($icon.ToBitmap(),"${dest}",32,$true) | Out-Null
$icon.Dispose()
