import { IJsonItem } from '../../interface/Ijson-item';
import { IPlatform } from '../../interface/Iplatform';
import { EItemType } from '../../enum/Eitem-type';
import { MenuItemConstructorOptions, nativeImage } from 'electron';
import path from 'path';
import { IIconsSize } from './favorite-apps-tray';
import { FileUtils, IDirectoryInfo } from '../../../vendor/utils/typescript/file-utils';
import { SystemUtils } from '../../../vendor/utils/typescript/system-utils';
import { ConsoleUtils } from '../../../vendor/utils/typescript/console-utils';
import { EShellType } from '../../../vendor/utils/typescript/shell-utils';

export class Windows implements IPlatform {
    private startMenuDirectory: string[];
    private _startMenuDirectoryInfo: IDirectoryInfo[] | undefined;
    constructor(private fileUtils: FileUtils, private consoleUtils: ConsoleUtils, private iconsDirectory: string, private iconsSize: IIconsSize) {
        this.startMenuDirectory = [
            path.resolve(`${SystemUtils.systemInfo.homeDir}\\AppData\\Roaming\\Microsoft\\Windows\\Start Menu`),
            path.resolve('C:\\ProgramData\\Microsoft\\Windows\\Start Menu\\Programs'),
        ];
    }

    private get startMenuDirectoryInfo(): IDirectoryInfo[] {
        if (!this._startMenuDirectoryInfo) {
            this._startMenuDirectoryInfo = [];
            this.startMenuDirectory.forEach((directory) => {
                const info = FileUtils.readDirRecursive(directory);
                if (info) {
                    this._startMenuDirectoryInfo?.push(info);
                }
            });
        }
        return this._startMenuDirectoryInfo;
    }

    private getItemData(jsonItem: IJsonItem): MenuItemConstructorOptions | undefined {
        let item: MenuItemConstructorOptions | undefined;
        if (jsonItem.type) {
            const fileInfo = FileUtils.getFileInfo(jsonItem.item);
            item = {
                label: jsonItem.name || fileInfo.basenameWithoutExtension,
                icon: undefined,
                click: () => {
                    if (jsonItem.type === EItemType.customItem || jsonItem.type === EItemType.item) {
                        this.consoleUtils.exec({ cmd: 'Start-Process', args: ['-FilePath', `'${jsonItem.item}'`], verbose: true, isThrow: false, shellType: EShellType.powershell });
                    } else if (jsonItem.type === EItemType.command) {
                        this.consoleUtils.exec({ cmd: jsonItem.item, shellType: jsonItem.shell ? EShellType[jsonItem.shell] : undefined, verbose: true, isThrow: false });
                    }
                }
            };
            if (jsonItem.type !== EItemType.command) {
                item.icon = path.resolve(`${this.iconsDirectory}/${fileInfo.basenameWithoutExtension}.ico`);
                if (!FileUtils.fileExist(item.icon)) {
                    this.fileUtils.extractIconLnkExe(jsonItem.item, item.icon);
                    if (!FileUtils.fileExist(item.icon)) {
                        item.icon = undefined;
                    }
                }
                if (item.icon) {
                    item.icon = nativeImage.createFromPath(item.icon as string);
                    item.icon.setTemplateImage(true);
                    item.icon = item.icon.resize(this.iconsSize);
                }
            }
        }
        return item;
    }


    getItem(jsonItem: IJsonItem): MenuItemConstructorOptions | undefined {
        if ((jsonItem.type === EItemType.customItem && FileUtils.fileExist(jsonItem.item)) || jsonItem.type === EItemType.command) {
            return this.getItemData(jsonItem);
        } else {
            for (const infoDir of this.startMenuDirectoryInfo) {
                const index = infoDir.files.findIndex((file) => {
                    const fileData = FileUtils.getFileInfo(file);
                    return fileData.basename === jsonItem.item;
                });
                if (index >= 0) {
                    jsonItem.item = infoDir.files[index];
                    return this.getItemData(jsonItem);
                }
            }
        }
        return undefined;
    }
}