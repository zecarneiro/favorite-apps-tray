import { FileSystem } from '../../vendor/utils/typescript/src/lib/file-system';
import { IDirectoryInfo } from '../../vendor/utils/typescript/src/interface/Idirectory-info';
import * as os from 'os';
import { Console } from '../../vendor/utils/typescript/src/lib/console/console';
import { IItemInfo } from '../interface/Iitem-info';
import { IJsonItem } from '../interface/Ijson-item';
import { ICommandInfo } from '../../vendor/utils/typescript/src/interface/Icomand-info';
import { EShellType } from '../../vendor/utils/typescript/src/enum/Eshell-type';
import { EJsonType } from '../enum/Ejson-type';
import { ESystem } from '../../vendor/utils/typescript/src/enum/Esystem';
import { FavoriteAppsTray } from '../main';
import { ENotifyType } from '../../vendor/utils/typescript/src/enum/Enotify-type';
import { Ui } from '../../vendor/utils/typescript/src/lib/ui';
import { Platform } from './platform';

export class Windows extends Platform {
    private startMenuDirectory: string[];
    private _startMenuDirectoryInfo: IDirectoryInfo[] | undefined;
    constructor(protected fileSystem: FileSystem, protected console: Console) {
        super(fileSystem, console);
        this.startMenuDirectory = [
            fileSystem.resolvePath(`${os.homedir}\\AppData\\Roaming\\Microsoft\\Windows\\Start Menu`),
            fileSystem.resolvePath(`C:\\ProgramData\\Microsoft\\Windows\\Start Menu\\Programs`),
        ];
        this.ui = new Ui(this.fileSystem, this.console);
    }

    private get startMenuDirectoryInfo(): IDirectoryInfo[] {
        if (!this._startMenuDirectoryInfo) {
            this._startMenuDirectoryInfo = [];
            this.startMenuDirectory.forEach((directory) => {
                const info = this.fileSystem.readDirRecursive(directory);
                if (info) {
                    this._startMenuDirectoryInfo?.push(info);
                }
            });
        }
        return this._startMenuDirectoryInfo;
    }

    private getItemInfo(jsonItem: IJsonItem): IItemInfo {
        const fileInfo = this.fileSystem.getFileInfo(jsonItem.item);
        const info: IItemInfo = { ...jsonItem, name: jsonItem.name ? jsonItem.name : fileInfo.data.basenameWithoutExtension };
        if (jsonItem.type === EJsonType.customItem || jsonItem.type === EJsonType.item) {
            const shurtcutsInfo = this.fileSystem.getShurtcutInfo(fileInfo.data.filename as string, ESystem.windows);
            info.item = fileInfo.data.filename as string;
            info.icon = shurtcutsInfo.icon;
        }
        return info;
    }

    public getInfo(jsonItem: IJsonItem): IItemInfo | null {
        let info: IItemInfo | null = null;
        if ((jsonItem.type === 'custom-item' && this.fileSystem.fileExist(jsonItem.item)) || jsonItem.type === 'command') {
            info = this.getItemInfo(jsonItem);
        } else {
            for (const infoDir of this.startMenuDirectoryInfo) {
                const index = infoDir.files.findIndex((file) => {
                    const fileData = this.fileSystem.getFileInfo(file);
                    return !fileData.hasError && fileData.data.basename === jsonItem.item;
                });
                if (index >= 0) {
                    jsonItem.item = infoDir.files[index];
                    info = this.getItemInfo(jsonItem);
                }
            }
        }
        return info;
    }

    public async run(info: IItemInfo) {
        let cmd: ICommandInfo = info.type === 'command' ? { cmd: info.item } : { cmd: 'Start-Process', args: ['-FilePath', `'${info.item}'`]};
        cmd = { ...cmd, isThrow: false, shellType: info.shell === 'cmd' ? EShellType.cmd : EShellType.system, cwd: this.fileSystem.systemInfo.homeDir };
        this.console.exec(cmd);
    }

    public notify(message: string, type: ENotifyType = ENotifyType.none) {
        this.ui.notify(FavoriteAppsTray.APP_NAME, message, ENotifyType.info, FavoriteAppsTray.DIALOG_TIMEOUT, ESystem.windows);
    }
}
