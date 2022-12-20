#!/usr/bin/env node
import SysTray from 'systray2';
import { Logger } from '../vendor/utils/typescript/src/lib/logger';
import { ILogger } from '../vendor/utils/typescript/src/interface/Ilogger';
import { Console } from '../vendor/utils/typescript/src/lib/console/console';
import { FileSystem } from '../vendor/utils/typescript/src/lib/file-system';
import { Ui } from '../vendor/utils/typescript/src/lib/ui';
import { Windows } from './lib/windows';
import { MenuItemClickable } from './interface/Imenu-item-clickable';
import { IJsonItem } from './interface/Ijson-item';
import { Functions } from '../vendor/utils/typescript/src/lib/functions';
import { IJson } from './interface/Ijson';
import { ENotifyType } from '../vendor/utils/typescript/src/enum/Enotify-type';

export class FavoriteAppsTray {
    public static readonly APP_NAME: string = 'Favorite Applications Indicator';
    public static readonly DIALOG_TIMEOUT: number = 15;
    private imagesDir: string;
    private icon: string;
    private jsonConfigFile: string;
    private systray: SysTray | undefined;
    private fileSystem: FileSystem;
    private console: Console;
    private ui: Ui;
    private logger: ILogger;
    private notifyUpdated: boolean;
    private windows: Windows;

    constructor() {
        this.logger = new Logger();
        this.fileSystem = new FileSystem(this.logger);
        this.console = new Console(this.fileSystem, this.logger);
        this.ui = new Ui(this.logger, this.fileSystem, this.console);
        this.imagesDir = this.fileSystem.resolvePath(`${__dirname}/../images`);
        this.jsonConfigFile = this.fileSystem.resolvePath(`${this.fileSystem.systemInfo.homeDir}/favorite-apps-tray.json`);
        this.icon = this.fileSystem.resolvePath(`${this.imagesDir}/logo/app.ico`);
        this.windows = new Windows(this.fileSystem, this.console);
        this.notifyUpdated = false;
    }

    private getBaseMenu(items: MenuItemClickable[] = []): MenuItemClickable[] {
        if (items && items.length > 0) {
            items.push(SysTray.separator);
        }
        items.push({
            title: 'Update',
            tooltip: 'Update all items in the application',
            checked: false,
            enabled: true,
            click: () => {
                if (this.systray && this.systray.process && this.systray.process.pid) {
                    process.kill(this.systray.process.pid);
                }
                this.process();
                this.notifyUpdated = true;
            },
        });
        items.push({
            title: 'Exit',
            tooltip: 'Exit of the application',
            checked: false,
            enabled: true,
            click: () => {
                this.systray?.kill(true);
            },
        });
        return items;
    }

    private mapJsonItemToItem(data: IJsonItem): MenuItemClickable {
        let item: MenuItemClickable = {} as MenuItemClickable;
        const info = this.windows.getInfo(data);
        if (info) {
            item = {
                title: info.name,
                tooltip: data.description,
                checked: false,
                enabled: true,
                icon: info.icon,
                click: () => {
                    this.windows.run(info);
                },
            } as MenuItemClickable;
        }
        return item;
    }

    private buildSubmenu(data: IJsonItem[], sort?: boolean): MenuItemClickable[] {
        const items: MenuItemClickable[] = data.map((value) => this.mapJsonItemToItem(value)).filter((value) => value !== null && Object.keys(value).length > 0);
        if (sort && items && items.length > 0) {
            return Functions.sort(items, 'title', 'string');
        }
        return items;
    }

    private buildMainMenu(): MenuItemClickable[] {
        if (!this.fileSystem.isWindows) {
            this.logger.error(new Error('Invalid platform'));
            return this.getBaseMenu([]);
        }
        const jsonRes = this.fileSystem.readJsonFile<IJson>(this.jsonConfigFile);
        if (jsonRes.hasError) {
            this.logger.error(jsonRes.error);
            return this.getBaseMenu([]);
        }
        const items: MenuItemClickable[] = [];
        const jsonData: IJson = jsonRes.data;
        if (jsonData.OTHERS) {
            for (const key in jsonData.OTHERS) {
                if (Object.prototype.hasOwnProperty.call(jsonData.OTHERS, key)) {
                    const submenu = this.buildSubmenu(jsonData.OTHERS[key], jsonData.SORT);
                    if (submenu.length > 0) {
                        items.push({
                            title: key,
                            tooltip: key,
                            enabled: true,
                            checked: false,
                            hidden: false,
                            items: submenu,
                        });
                    }
                }
            }
        }
        if (jsonData.NO_MENU && jsonData.NO_MENU.length > 0) {
            if (items.length > 0) {
                items.push(SysTray.separator);
            }
            // eslint-disable-next-line prefer-spread
            items.push.apply(items, this.buildSubmenu(jsonData.NO_MENU, jsonData.SORT));
        }
        return this.getBaseMenu(items);
    }

    public process() {
        const items: MenuItemClickable[] = this.buildMainMenu();
        this.systray = new SysTray({
            menu: {
                icon: this.icon,
                title: '',
                tooltip: FavoriteAppsTray.APP_NAME,
                isTemplateIcon: true,
                items: items,
            },
            debug: false,
            copyDir: false, // copy go tray binary to outside directory, useful for packing tool like pkg.
        });
        this.systray.onClick((action) => {
            const item: MenuItemClickable = Functions.convert<MenuItemClickable>(action.item);
            if (item.click != null) {
                item.click();
            }
        });

        // Systray.ready is a promise which resolves when the tray is ready.
        this.systray.ready().then(() => {
            let message = 'systray started!';
            if (this.notifyUpdated) {
                this.notifyUpdated = false;
                message = 'Update, done.';
                this.ui.notify(FavoriteAppsTray.APP_NAME, message, ENotifyType.info, FavoriteAppsTray.DIALOG_TIMEOUT);
            }
            this.logger.info(message);
        }).catch((err) => {
            this.ui.notify(FavoriteAppsTray.APP_NAME, 'systray failed to start: ' + err.message, ENotifyType.error, FavoriteAppsTray.DIALOG_TIMEOUT);
        });
    }

    public static start() {
        const favoriteAppsTray = new FavoriteAppsTray();
        favoriteAppsTray.process();
    }
}
FavoriteAppsTray.start();
