import { FileSystem } from '../../vendor/utils/typescript/src/lib/file-system';
import { IDirectoryInfo } from '../../vendor/utils/typescript/src/interface/Idirectory-info';
import { Console } from '../../vendor/utils/typescript/src/lib/console/console';
import { IItemInfo } from '../interface/Iitem-info';
import { IJsonItem } from '../interface/Ijson-item';
import { ESystem } from '../../vendor/utils/typescript/src/enum/Esystem';
import { Platform } from './platform';

export class Gnome extends Platform {
    private _directoryInfo: IDirectoryInfo[] | undefined;
    constructor(protected fileSystem: FileSystem, protected console: Console) {
        super(fileSystem, console);
    }

    private get desktopDirectoriesInfo(): IDirectoryInfo[] {
        if (!this._directoryInfo) {
            const desktopDirectories = [
                `${this.fileSystem.systemInfo.homeDir}/.local/share/applications`,
                '/usr/share/applications',
                '/var/lib/snapd/desktop/applications',
                '/var/lib/flatpak/exports/share/applications',
                `${this.fileSystem.systemInfo.homeDir}/.local/share/flatpak/exports/share/applications`,
            ];
            this._directoryInfo = [];
            desktopDirectories.forEach((directory) => {
                const info = this.fileSystem.readDirRecursive(directory);
                if (info) {
                    this._directoryInfo?.push(info);
                }
            });
        }
        return this._directoryInfo;
    }

    private getDesktopInfo(fileName: string): IItemInfo | null {
        let infoData: IItemInfo | null = null;
        const shurtcutsInfo = this.fileSystem.getShurtcutInfo(fileName, ESystem.gnome);
        if (shurtcutsInfo) {
            infoData = {
                name: shurtcutsInfo.name,
                item: shurtcutsInfo.command?.trim() as string,
                icon: shurtcutsInfo.icon,
            };
            if (infoData.icon) {
                if (!this.fileSystem.fileExist(infoData.icon)) {
                    const iconFilesPath = ['/usr/share/icons/32x32', '/usr/share/pixmaps'];
                    for (const path of iconFilesPath) {
                        const files = this.fileSystem.readDirRecursive(path).files;
                        const index = files.findIndex((iconFile) => iconFile.includes(`${infoData?.icon}.png`));
                        if (index >= 0) {
                            infoData.icon = this.fileSystem.getBase64File(files[index]).base;
                            break;
                        }
                    }
                }
            }
        }
        return infoData;
    }

    public getInfo(jsonItem: IJsonItem): IItemInfo | null {
        let info: IItemInfo | null = null;
        if ((jsonItem.type === 'custom-item' && this.fileSystem.fileExist(jsonItem.item)) || jsonItem.type === 'item') {
            info = this.getDesktopInfo(jsonItem.item);
        } else if (jsonItem.type === 'command') {
            return jsonItem as IItemInfo;
        }
        return info;
    }

    public run(info: IItemInfo) {
        this.console.exec({ cmd: info.item, isThrow: false, verbose: false });
    }
}
