import { EShellType } from '../lib/ts-js-utils/enum/Eshell-type';
import { EItemType } from '../enum/Eitem-type';

export interface IJsonItem {
    name?: string;
    item: string;
    description?: string;
    type?: EItemType;
    shell?: EShellType;
}