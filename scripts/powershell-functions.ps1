# Author: José M. C. Noronha
# ---------------------------------------------------------------------------- #
#                                 APP INFO AREA                                #
# ---------------------------------------------------------------------------- #
$appName = "${env:FAT_INFO_APP_NAME}"
$typeInfo = "${env:FAT_INFO_TYPE}"
$TRAY_APP_CONFIG_DIR = "${home}\.config\${appName}"
$JSON_FILE = "${TRAY_APP_CONFIG_DIR}\apps-info-${typeInfo}.json"
$APPS_ARR = @()

function saveToFile {
    if (!([string]::IsNullOrEmpty("$JSON_FILE")) -and (Test-Path -Path "$JSON_FILE" -PathType Leaf)) {
        Remove-Item "$JSON_FILE" -Recurse -Force
    }
    Write-Host "INFO: ${appName} - Generating: $JSON_FILE"
    $content = (ConvertTo-Json $APPS_ARR -Depth 2 | Out-String)
    New-Item -Force "$JSON_FILE" -Value ($content | Out-String) | Out-Null
}

function appinfo {
    New-Item -Path "$TRAY_APP_CONFIG_DIR" -ItemType Directory -Force | Out-Null
    if ("windows-apps" -eq $typeInfo) {
        (New-Object -ComObject Shell.Application).NameSpace('shell:AppsFolder').Items() | ForEach-Object {
            $nameApp = $_.Name
            $pathApp = $_.Path
            $path = "shell:AppsFolder\'$pathApp'"
            $appInfo = @{
                displayName="$nameApp";
                command="Start-Process -FilePath $path";
                shortcut="$nameApp"
            }
            $APPS_ARR += $appInfo
        }
        saveToFile
    } elseif ("shortcuts" -eq $typeInfo) {
        $menuDirs = "C:\ProgramData\Microsoft\Windows\Start Menu\Programs", "$home\AppData\Roaming\Microsoft\Windows\Start Menu"
        $menuDirs | ForEach-Object {
            $menuDir = $_
            Get-ChildItem -Path "$menuDir" -Include "*.lnk" -File -Recurse | ForEach-Object {
                $shortcut = $_
                $basename = ([System.IO.Path]::GetFileName("$shortcut"))
                $filename = ([System.IO.Path]::GetFileNameWithoutExtension("$basename"))
                $appInfo = [PSCustomObject]@{
                    displayName="$filename";
                    command="Start-Process -FilePath '$shortcut'";
                    shortcut=$basename
                    icon="$shortcut"
                }
                $APPS_ARR += $appInfo
            }
        }
        saveToFile
    }
}

# ---------------------------------------------------------------------------- #
#                                    UI AREA                                   #
# ---------------------------------------------------------------------------- #
function notify {
    $appId = "${env:FAT_NOTIFY_APP_ID}"
    $title = "${env:FAT_NOTIFY_TITLE}"
    $message = "${env:FAT_NOTIFY_MESSAGE}"
    $icon = "${env:FAT_NOTIFY_ICON}"

    [Windows.UI.Notifications.ToastNotificationManager, Windows.UI.Notifications, ContentType = WindowsRuntime] > $null
    $Template = [Windows.UI.Notifications.ToastNotificationManager]::GetTemplateContent([Windows.UI.Notifications.ToastTemplateType]::ToastImageAndText02)

    $RawXml = [xml] $Template.GetXml()
    ($RawXml.toast.visual.binding.text|Where-Object {$_.id -eq "1"}).AppendChild($RawXml.CreateTextNode($title)) > $null
    ($RawXml.toast.visual.binding.text|Where-Object {$_.id -eq "2"}).AppendChild($RawXml.CreateTextNode($message)) > $null
    if (![string]::IsNullOrEmpty($icon)) {
        ($RawXml.toast.visual.binding.image|Where-Object {$_.id -eq "1"}).SetAttribute('src', $icon) > $null
    }

    $SerializedXml = New-Object Windows.Data.Xml.Dom.XmlDocument
    $SerializedXml.LoadXml($RawXml.OuterXml)
    Write-Host $SerializedXml.ToString()

    $Toast = [Windows.UI.Notifications.ToastNotification]::new($SerializedXml)
    $Toast.Tag = "$appId"
    $Toast.Group = "$appId"
    $Toast.ExpirationTime = [DateTimeOffset]::Now.AddMinutes(1)

    $Notifier = [Windows.UI.Notifications.ToastNotificationManager]::CreateToastNotifier("$appId")
    $Notifier.Show($Toast);
}

function selectfiledialog {
    # Load Assembly
    Add-Type -AssemblyName System.Windows.Forms

    $OpenFileDialog = New-Object System.Windows.Forms.OpenFileDialog
    $OpenFileDialog.ShowDialog() | Out-Null
    $filename = $OpenFileDialog.FileName
    $data = @{ selected="$filename"; }
    ConvertTo-Json $data -Depth 1 | Out-String
}

function messagedialog {
    Add-Type -AssemblyName System.Windows.Forms
    $title = "${env:FAT_DIALOG_TITLE}"
    $message = "${env:FAT_DIALOG_MESSAGE}"
    $icon = "${env:FAT_DIALOG_ICON}"
    [System.Windows.Forms.MessageBox]::Show($message, $title, [System.Windows.Forms.MessageBoxButtons]::OK, "Information")
}

# ---------------------------------------------------------------------------- #
#                                  OTHERS AREA                                 #
# ---------------------------------------------------------------------------- #
function iconextractor {
    $file = "${env:FAT_ICON_EXTRACTOR_FILE}"
    $dest = "${env:FAT_ICON_EXTRACTOR_DEST}"
    # Source code from https://www.powershellgallery.com/packages/IconExport/1.0.1/Content/IconExport.psm1
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
}
