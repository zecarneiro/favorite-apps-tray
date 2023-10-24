import { MenuItemConstructorOptions, nativeImage } from 'electron';
import { IJsonItem } from '../../interface/Ijson-item';
import { IPlatform } from '../../interface/Iplatform';
import { IIconsSize } from './favorite-apps-tray';
import path from 'path';
import { FileUtils } from '../ts-js-utils/file-utils';
import { ConsoleUtils } from '../ts-js-utils/console-utils';
import { EShellType } from '../ts-js-utils/enum/Eshell-type';
import { ICommandInfo } from '../ts-js-utils/interface/Icommand-info';
import { SystemUtils } from '../ts-js-utils/system-utils';
import { FunctionUtils } from '../ts-js-utils/function-utils';
import { IDirectories } from '../../interface/Idirectories';

interface IDesktopFileInfo {
    name: string;
    icon?: string;
    command: string;
    displayName: string;
    commandByGtk?: string;
}

export class Gnome implements IPlatform {
    constructor(private fileUtils: FileUtils, private consoleUtils: ConsoleUtils, private iconsDirectory: string, private iconsSize: IIconsSize, private directories: IDirectories) {
    }

    private svgToPng(svg: string, outFile: string): boolean {
        let cmdInfo: ICommandInfo | undefined;
        if (SystemUtils.isLinux && FileUtils.fileExist(svg) && FileUtils.validateFileExtension(svg, ['.svg'])) {
            cmdInfo = {
                cmd: `source "${path.normalize(this.directories.bashScript)}/main-utils.sh"`,
                args: [
                    `;svg_to_png "${svg}" "${outFile}"`,
                ],
                shellType: EShellType.bash,
            };
        }
        if (cmdInfo) {
            this.consoleUtils.execSync({ ...cmdInfo, isThrow: false, verbose: false });
            return FileUtils.fileExist(outFile);
        }
        return false;
    }

    private extractDesktopFileInfo(desktopFile: string, getType: string): IDesktopFileInfo | undefined {
        let desktopFileInfo: IDesktopFileInfo | undefined;
        let cmdInfo: ICommandInfo | undefined;
        if (SystemUtils.isGnome) {
            const script = path.resolve(this.directories.script, 'desktop-file-info-in-gnome.js');
            cmdInfo = {
                cmd: 'gjs',
                args: [
                    `"${path.normalize(script)}"`,
                    `"${desktopFile}"`,
                    `"${getType}"`
                ],
                shellType: EShellType.bash,
            };
        }
        if (cmdInfo) {
            const resp = this.consoleUtils.execSync({ ...cmdInfo, isThrow: false, verbose: false });
            if (!resp.hasError) {
                desktopFileInfo = FunctionUtils.stringToObject<IDesktopFileInfo>(resp.data);
            }
        }
        return desktopFileInfo && Object.keys(desktopFileInfo).length > 0 ? desktopFileInfo : undefined;
    }

    private getCommandData(jsonItem: IJsonItem): MenuItemConstructorOptions | undefined {
        return {
            label: jsonItem.displayName || jsonItem.nameOrFile,
            icon: undefined,
            click: () => {
                this.consoleUtils.exec({ cmd: jsonItem.nameOrFile, shellType: jsonItem.shell ? EShellType[jsonItem.shell] : undefined, verbose: false, isThrow: false });
            }
        };
    }

    private getNameOrSystemData(jsonItem: IJsonItem): MenuItemConstructorOptions | undefined {
        let item: MenuItemConstructorOptions | undefined;
        const info = this.extractDesktopFileInfo(jsonItem.nameOrFile, jsonItem.type);
        if (info) {
            let icon: string | undefined = info.icon;
            if (info.icon) {
                if (FileUtils.validateFileExtension(icon, ['.svg'])) {
                    const iconInfo = FileUtils.getFileInfo(icon);
                    icon = path.resolve(this.iconsDirectory, `${iconInfo.basenameWithoutExtension}.png`);
                    if (!FileUtils.fileExist(icon) && !this.svgToPng(info.icon, icon)) {
                        icon = undefined;
                    }
                }
            }
            item = {
                label: jsonItem.displayName ? jsonItem.displayName : info.displayName || jsonItem.nameOrFile,
                click: () => {
                    this.consoleUtils.exec({ cmd: info.commandByGtk ? info.commandByGtk : info.command, verbose: false, isThrow: false, shellType: EShellType.bash });
                }
            };
            if (icon) {
                item.icon = nativeImage.createFromPath(icon);
                item.icon.setTemplateImage(true);
                item.icon = item.icon.resize(this.iconsSize);
            }
        }
        return item;
    }

    getItem(jsonItem: IJsonItem): MenuItemConstructorOptions | undefined {
        if (jsonItem.type == 'command') {
            return this.getCommandData(jsonItem);
        } else if (jsonItem.type == 'name' || jsonItem.type == 'name-start' || jsonItem.type == 'system' || jsonItem.type == 'system-name-start') {
            return this.getNameOrSystemData(jsonItem);
        }
        return undefined;
    }
}