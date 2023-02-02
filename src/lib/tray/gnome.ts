import { IJsonItem } from '../../interface/Ijson-item';
import { IPlatform } from '../../interface/Iplatform';

export class Gnome implements IPlatform {
    getItem(jsonItem: IJsonItem): Electron.MenuItemConstructorOptions {
        throw new Error('Method not implemented.');
    }
}