import { MenuItemConstructorOptions, nativeImage } from 'electron';
import { IJsonItem } from '../../interface/Ijson-item';
import { IPlatform } from '../../interface/Iplatform';
import { EItemType } from '../../enum/Eitem-type';
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
    private getItemCommandData(jsonItem: IJsonItem): MenuItemConstructorOptions | undefined {
        let item: MenuItemConstructorOptions | undefined;
        if (jsonItem.type && jsonItem.type === EItemType.command) {
            item = {
                label: jsonItem.name,
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
            const info = this.extractDesktopFileInfo(jsonItem.item);
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
                    label: info.displayName || jsonItem.name,
                    click: () => {
                        this.consoleUtils.exec({ cmd: info.commandByGtk ? info.commandByGtk : info.command, verbose: true, isThrow: false, shellType: EShellType.bash });
                    }
                };
                if (icon) {
                    item.icon = nativeImage.createFromPath(icon);
                    item.icon.setTemplateImage(true);
                    item.icon = item.icon.resize(this.iconsSize);
                }
            }
        }
        return item;
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

    private extractDesktopFileInfo(desktopFile: string): IDesktopFileInfo | undefined {
        let desktopFileInfo: IDesktopFileInfo | undefined;
        let cmdInfo: ICommandInfo | undefined;
        if (SystemUtils.isGnome) {
            const script = path.resolve(this.directories.script, 'desktop-file-info-in-gnome.js');
            cmdInfo = {
                cmd: 'gjs',
                args: [
                    `"${path.normalize(script)}"`,
                    `"${desktopFile}"`,
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

    getItem(jsonItem: IJsonItem): MenuItemConstructorOptions | undefined {
        if (jsonItem.type === EItemType.command) {
            return this.getItemCommandData(jsonItem);
        } else if (jsonItem.type === EItemType.item) {
            return this.getItemData(jsonItem);
        }
        return undefined;
    }
}