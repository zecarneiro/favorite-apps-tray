import { MenuItemConstructorOptions, nativeImage } from 'electron';
import { IJsonItem } from '../../interface/Ijson-item';
import { IPlatform } from '../../interface/Iplatform';
import { EItemType } from '../../enum/Eitem-type';
import { FileUtils } from '../../../vendor/utils/typescript/file-utils';
import { ConsoleUtils } from '../../../vendor/utils/typescript/console-utils';
import { EShellType } from '../../../vendor/utils/typescript/shell-utils';
import { IIconsSize } from './favorite-apps-tray';
import { LoggerUtils } from '../../../vendor/utils/typescript/logger-utils';
import path from 'path';

export class Gnome implements IPlatform {
    constructor(private fileUtils: FileUtils, private consoleUtils: ConsoleUtils, private iconsDirectory: string, private iconsSize: IIconsSize) {
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
            const info = this.fileUtils.extractDesktopFileInfo(jsonItem.item);
            if (info) {
                let icon: string | undefined = info.icon;
                if (info.icon) {
                    if (FileUtils.validateFileExtension(icon, ['.svg'])) {
                        const iconInfo = FileUtils.getFileInfo(icon);
                        icon = path.resolve(this.iconsDirectory, `${iconInfo.basenameWithoutExtension}.png`);
                        if (!FileUtils.fileExist(icon) && !this.fileUtils.svgToPng(info.icon, icon)) {
                            icon = undefined;
                        }
                    }
                }
                item = {
                    label: info.displayName || jsonItem.name,
                    click: () => {
                        this.consoleUtils.exec({ cmd: info.command, verbose: false, isThrow: false, shellType: EShellType.bash });
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

    getItem(jsonItem: IJsonItem): MenuItemConstructorOptions | undefined {
        if (jsonItem.type === EItemType.command) {
            return this.getItemCommandData(jsonItem);
        } else if (jsonItem.type === EItemType.item) {
            return this.getItemData(jsonItem);
        }
        return undefined;
    }
}