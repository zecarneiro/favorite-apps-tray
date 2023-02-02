import { MenuItemConstructorOptions } from 'electron';
import { IJsonItem } from './Ijson-item';

export interface IPlatform {
    getItem(jsonItem: IJsonItem): MenuItemConstructorOptions | undefined;
}