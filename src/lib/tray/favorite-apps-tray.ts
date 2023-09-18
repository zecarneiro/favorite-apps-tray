import { Menu, MenuItemConstructorOptions, Tray, app, dialog } from 'electron';
import { IJson } from '../../interface/Ijson';
import { IJsonItem } from '../../interface/Ijson-item';
import { IPlatform } from '../../interface/Iplatform';
import { Windows } from './windows';
import * as path from 'path';
import { FavoriteAppsProcessor } from '../../favorite-apps-processor';
import { Gnome } from './gnome';
import { SystemUtils } from '../ts-js-utils/system-utils';
import { FileUtils } from '../ts-js-utils/file-utils';
import { FunctionUtils } from '../ts-js-utils/function-utils';

export interface IIconsSize {
  height: number, width: number
}

export class FavoriteAppsTray extends FavoriteAppsProcessor {
  private readonly MENU_ICON_SIZE: IIconsSize = {
    height: 16, width: 16
  };
  private trayMenu: MenuItemConstructorOptions[] = [];
  
  private _tray: Tray | undefined;
  private get tray(): Tray {
    if (!this._tray || this._tray.isDestroyed()) {
      this._tray = new Tray(this.logo);
      this.tray.setIgnoreDoubleClickEvents(true);
      this.tray.on('click', () => {
        this.tray?.popUpContextMenu();
      });
    }
    return this._tray;
  }

  private get jsonConfigMenuEntries(): string {
    return path.resolve(this.configDir, 'tray_menu_entries.json');
  }

  private _platformProcess: IPlatform | undefined;
  private get platformProcess(): IPlatform | undefined {
    if (!this._platformProcess) {
      if (SystemUtils.isWindows) {
        this._platformProcess = new Windows(this.fileUtils, this.consoleUtils, this.iconsDir, this.MENU_ICON_SIZE, this.directories);
      } else if (SystemUtils.isGnome) {
        this._platformProcess = new Gnome(this.fileUtils, this.consoleUtils, this.iconsDir, this.MENU_ICON_SIZE, this.directories);
      }
    }
    return this._platformProcess;
  }

  private get iconsDir(): string {
    const dir = path.resolve(this.configDir, 'icons');
    if (!FileUtils.fileExist(dir)) {
      FileUtils.createDir(dir);
    }
    return dir;
  }

  private _menuConfigData: IJson | undefined;
  private getMenuConfigData(force?: boolean): IJson | undefined {
    if (this.isValidPlatform() && (!this._menuConfigData || force)) {
      const dataMenu = FileUtils.readJsonFile<IJson>(this.jsonConfigMenuEntries);
      if (dataMenu) {
        this._menuConfigData = dataMenu;
      }
    }
    return this._menuConfigData;
  }

  private clear() {
    this.tray.setContextMenu(Menu.buildFromTemplate(this.buildBaseMenu()));
    this.trayMenu = [];
  }

  private updateMenu() {
    this.clear();
    this.process(false, true);
    this.showProcessing(true);
  }

  private insetSeparator() {
    if (this.trayMenu.length > 0) {
      this.trayMenu.push(
        {
          type: 'separator'
        });
    }
  }

  private buildBaseMenu(): MenuItemConstructorOptions[] {
    this.insetSeparator();
    return [
      {
        id: '1',
        label: 'Settings',
        type: 'submenu',
        submenu: [
          {
            label: 'Select/Change JSON file',
            enabled: this.isValidPlatform(),
            click: () => {
              const file = dialog.showOpenDialogSync({
                properties: ['openFile'],
                title: this.APP_NAME,
                message: 'Select a JSON configuration file'
              });
              if (file && file.length > 0) {
                if (!FunctionUtils.isJsonParsable(FileUtils.readTextFile(file[0]))) {
                  this.showNotification('Invalid JSON file');
                } else {
                  FileUtils.copyFile(file[0], this.jsonConfigMenuEntries, true);
                  this.updateMenu();
                }
              }
            }
          },
          {
            label: 'Sort Menu Entries',
            toolTip: 'Sort Menu Entries',
            type: 'checkbox',
            enabled: this.isValidPlatform(),
            checked: this.configurations.isSortMenuEntries,
            click: () => {
              this.configurations.isSortMenuEntries = this.configurations.isSortMenuEntries ? false : true;
              this.updateConfigurations(this.configurations);
              this.updateMenu();
            }
          },
          {
            label: 'Add Startup',
            toolTip: 'Add to boot',
            type: 'checkbox',
            enabled: this.isValidPlatform(),
            checked: this.configurations.isStartup,
            click: () => {
              if (this.configurations.isStartup) {
                this.addBootApp({ name: app.getName(),  isDelete: true });
                this.configurations.isStartup = false;
              } else {
                this.addBootApp({ name: app.getName(), command: app.getPath('exe'), hidden: true });
                this.configurations.isStartup = true;
              }
              this.updateConfigurations(this.configurations);
              if (SystemUtils.isGnome) {
                this.updateMenu();
              } else {
                this.showProcessing(true);
              }
            }
          },
        ]
      },
      {
        id: '2',
        label: 'Update',
        toolTip: 'Update Entries',
        click: () => {
          this.updateMenu();
        }
      },
      {
        id: '3',
        label: 'Exit',
        toolTip: 'Exit of the application',
        click: () => {
          app.exit();
        }
      }
    ];
  }

  private buildSubmenu(data: IJsonItem[]): MenuItemConstructorOptions[] {
    const items: MenuItemConstructorOptions[] = [];
    if (data && data.length > 0) {
      for (const item of data) {
        const data = this.platformProcess?.getItem(item);
        if (data) {
          items.push(data);
        }
      }
    }
    if (this.configurations.isSortMenuEntries && items && items.length > 0) {
      return FunctionUtils.sort(items, 'label', 'string');
    }
    return items;
  }

  private buildMainMenu(jsonData: IJson) {
    if (jsonData.others) {
      for (const key in jsonData.others) {
        if (Object.prototype.hasOwnProperty.call(jsonData.others, key)) {
          const submenu = this.buildSubmenu(jsonData.others[key]);
          if (submenu.length > 0) {
            this.trayMenu.push({
              label: key,
              submenu: submenu
            });
          }
        }
      }
    }
    if (jsonData.noMenu) {
      this.insetSeparator();
      this.trayMenu = this.trayMenu.concat(this.buildSubmenu(jsonData.noMenu));
    }
  }

  protected process(showNotify: boolean, forceReadJsonMenuConfig?: boolean): void {
    const config = this.getMenuConfigData(forceReadJsonMenuConfig);
    if (config) {
      this.buildMainMenu(config);
      if (showNotify) {
        this.showProcessing(true);
      }
    }
    this.tray.setContextMenu(Menu.buildFromTemplate(this.trayMenu.concat(this.buildBaseMenu())));
  }

  public distroy() {
    if (this.tray && !this.tray.isDestroyed()) {
      this.tray.destroy();
      this._tray = undefined;
    }
  }
}