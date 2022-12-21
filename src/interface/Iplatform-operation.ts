import { ENotifyType } from '../../vendor/utils/typescript/src/enum/Enotify-type';
import { IItemInfo } from './Iitem-info';
import { IJsonItem } from './Ijson-item';

export interface IPlatformOperations {
    notify(message: string, type: ENotifyType): void;
    getInfo(jsonItem: IJsonItem): IItemInfo | null;
    run(file: IItemInfo): void;
}
