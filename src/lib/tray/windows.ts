import { IJsonItem } from '../../interface/Ijson-item';
import { IPlatform } from '../../interface/Iplatform';
import { EItemType } from '../../enum/Eitem-type';
import { MenuItemConstructorOptions, nativeImage } from 'electron';
import path from 'path';
import { IIconsSize } from './favorite-apps-tray';
import { FileUtils } from '../ts-js-utils/file-utils';
import { ConsoleUtils } from '../ts-js-utils/console-utils';
import { SystemUtils } from '../ts-js-utils/system-utils';
import { IDirectoryInfo } from '../ts-js-utils/interface/Idirectory-info';
import { EShellType } from '../ts-js-utils/enum/Eshell-type';
import { ICommandInfo } from '../ts-js-utils/interface/Icommand-info';
import { IDirectories } from '../../interface/Idirectories';

export class Windows implements IPlatform {
    private startMenuDirectory: string[];
    private _startMenuDirectoryInfo: IDirectoryInfo[] | undefined;
    constructor(private fileUtils: FileUtils, private consoleUtils: ConsoleUtils, private iconsDirectory: string, private iconsSize: IIconsSize, private directories: IDirectories) {
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

    private getItemCommandData(jsonItem: IJsonItem): MenuItemConstructorOptions | undefined {
        let item: MenuItemConstructorOptions | undefined;
        if (jsonItem.type && jsonItem.type === EItemType.command) {
            const fileInfo = FileUtils.getFileInfo(jsonItem.item);
            item = {
                label: jsonItem.name || fileInfo.basenameWithoutExtension,
                icon: undefined,
                click: () => {
                    this.consoleUtils.exec({ cmd: jsonItem.item, shellType: jsonItem.shell ? EShellType[jsonItem.shell] : undefined, verbose: false, isThrow: false });
                }
            };
        }
        return item;
    }

    private getItemData(jsonItem: IJsonItem): MenuItemConstructorOptions | undefined {
        let item: MenuItemConstructorOptions | undefined;
        if (jsonItem.type) {
            const fileInfo = FileUtils.getFileInfo(jsonItem.item);
            const icon = path.resolve(`${this.iconsDirectory}/${fileInfo.basenameWithoutExtension}.ico`);
            if (!FileUtils.fileExist(icon)) {
                this.extractIconLnkExe(jsonItem.item, icon);
            }
            item = {
                label: jsonItem.name || fileInfo.basenameWithoutExtension,
                icon: FileUtils.fileExist(icon) ? icon : undefined,
                click: () => {
                    this.consoleUtils.exec({ cmd: 'Start-Process', args: ['-FilePath', `'${jsonItem.item}'`], verbose: false, isThrow: false, shellType: EShellType.powershell });
                }
            };
            if (item.icon) {
                item.icon = nativeImage.createFromPath(icon);
                item.icon.setTemplateImage(true);
                item.icon = item.icon.resize(this.iconsSize);
            }
        }
        return item;
    }

    private extractIconLnkExe(file: string, dest: string) {
        if (SystemUtils.isWindows && FileUtils.fileExist(file) && FileUtils.validateFileExtension(file, ['.lnk', '.exe'])) {
            const cmdInfo: ICommandInfo = {
                cmd: `. "${path.normalize(this.directories.powershellScript)}\\MainUtils.ps1"`,
                args: [
                    `;IconExtractor -file "${file}" -dest "${dest}" -display`,
                ],
                isThrow: false,
                verbose: false,
                shellType: EShellType.powershell,
            };
            this.consoleUtils.execSync(cmdInfo);
        }
    }


    getItem(jsonItem: IJsonItem): MenuItemConstructorOptions | undefined {
        if (jsonItem.type === EItemType.command) {
            return this.getItemCommandData(jsonItem);
        } else if ((jsonItem.type === EItemType.customItem && FileUtils.fileExist(jsonItem.item))) {
            return this.getItemData(jsonItem);
        } else {
            for (const infoDir of this.startMenuDirectoryInfo) {
                const index = infoDir.files.findIndex((file) => {
                    const fileData = FileUtils.getFileInfo(file);
                    if (fileData.basename === jsonItem.item) {
                        return true;
                    } else {
                        const jsonItemInfo = FileUtils.getFileInfo(jsonItem.item);
                        if (jsonItemInfo.extension !== '.lnk' && jsonItemInfo.extension !== 'lnk') {
                            return fileData.basename.startsWith(jsonItem.item);
                        }
                    }
                    return false;
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