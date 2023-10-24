import { IJsonItem } from '../../interface/Ijson-item';
import { IPlatform } from '../../interface/Iplatform';
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
import { IFileInfo } from '../ts-js-utils/interface/Ifile-info';
import { FunctionUtils } from '../ts-js-utils/function-utils';

interface IAppInfo {
    name: string;
    command: string;
    exec: string;
}

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

    private extractAppInfo(appName: string, isStart?: boolean): IAppInfo | undefined {
        let appInfo: IAppInfo | undefined;
        let cmdInfo: ICommandInfo | undefined;
        if (SystemUtils.isWindows) {
            const script = path.resolve(this.directories.script, 'app-info-in-windows.ps1');
            cmdInfo = {
                cmd: path.normalize(script),
                args: ['-name', `"${appName}"`, isStart ? '-isStart' : ''],
                shellType: EShellType.powershell,
            };
        }
        if (cmdInfo) {
            const resp = this.consoleUtils.execSync({ ...cmdInfo, isThrow: false, verbose: false });
            if (!resp.hasError) {
                appInfo = FunctionUtils.stringToObject<IAppInfo>(resp.data);
            }
        }
        return appInfo && Object.keys(appInfo).length > 0 ? appInfo : undefined;
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

    private getIcon(fileInfo: IFileInfo): string | undefined {
        const icon = path.resolve(`${this.iconsDirectory}/${fileInfo.basenameWithoutExtension}.ico`);
        if (!FileUtils.fileExist(icon)) {
            this.extractIconLnkExe(fileInfo.original, icon);
        }
        return icon;
    }

    private getCommandData(jsonItem: IJsonItem): MenuItemConstructorOptions | undefined {
        const fileInfo = FileUtils.getFileInfo(jsonItem.nameOrFile);
        const displayName = jsonItem.displayName ? jsonItem.displayName : fileInfo?.basenameWithoutExtension || jsonItem.nameOrFile;
        return {
            label: displayName,
            icon: undefined,
            click: () => {
                this.consoleUtils.exec({ cmd: jsonItem.nameOrFile, shellType: jsonItem.shell ? EShellType[jsonItem.shell] : undefined, verbose: false, isThrow: false });
            }
        };
    }

    private getSystemOrOtherData(jsonItem: IJsonItem): MenuItemConstructorOptions | undefined {
        const fileInfo = FileUtils.getFileInfo(jsonItem.nameOrFile);
        const icon = this.getIcon(fileInfo);
        const displayName = jsonItem.displayName ? jsonItem.displayName : fileInfo?.basenameWithoutExtension || jsonItem.nameOrFile;
        const item: MenuItemConstructorOptions = {
            label: displayName,
            icon: FileUtils.fileExist(icon) ? icon : undefined,
            click: () => {
                this.consoleUtils.exec({ cmd: 'Start-Process', args: ['-FilePath', `'${jsonItem.nameOrFile}'`], verbose: false, isThrow: false, shellType: EShellType.powershell });
            }
        };
        if (item.icon) {
            item.icon = nativeImage.createFromPath(icon);
            item.icon.setTemplateImage(true);
            item.icon = item.icon.resize(this.iconsSize);
        }
        return item;
    }

    private getSystemData(jsonItem: IJsonItem): MenuItemConstructorOptions | undefined {
        let item: MenuItemConstructorOptions | undefined;
        for (const infoDir of this.startMenuDirectoryInfo) {
            const index = infoDir.files.findIndex((file) => {
                const fileData = FileUtils.getFileInfo(file);
                if (jsonItem.type == 'system') {
                    return fileData.basename === jsonItem.nameOrFile;
                } else if (jsonItem.type == 'system-name-start') {
                    return fileData.basenameWithoutExtension.startsWith(jsonItem.nameOrFile);
                }
            });
            if (index >= 0) {
                jsonItem.nameOrFile = infoDir.files[index];
                item = this.getSystemOrOtherData(jsonItem);
            }
        }
        return item;
    }

    private getNameData(jsonItem: IJsonItem): MenuItemConstructorOptions | undefined {
        let item: MenuItemConstructorOptions | undefined;
        const appInfo = this.extractAppInfo(jsonItem.nameOrFile, jsonItem.type == 'name-start');
        if (appInfo) {
            const displayName = jsonItem.displayName ? jsonItem.displayName : appInfo.name || jsonItem.nameOrFile;
            item = {
                label: displayName,
                icon: undefined,
                click: () => {
                    this.consoleUtils.exec({ cmd: appInfo.command, verbose: true, isThrow: false, shellType: EShellType.powershell });
                }
            }
        }
        return item;
    }

    getItem(jsonItem: IJsonItem): MenuItemConstructorOptions | undefined {
        if (jsonItem.type == 'command') {
            return this.getCommandData(jsonItem);
        } else if (jsonItem.type == 'other' && FileUtils.fileExist(jsonItem.nameOrFile)) {
            return this.getSystemOrOtherData(jsonItem);
        } else if (jsonItem.type == 'system' || jsonItem.type == 'system-name-start') {
            return this.getSystemData(jsonItem);
        } else if (jsonItem.type == 'name' || jsonItem.type == 'name-start') {
            return this.getNameData(jsonItem);
        }
        return undefined;
    }
}