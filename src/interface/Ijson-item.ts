import { EShellType } from '../../vendor/utils/typescript/shell-utils';
import { EItemType } from '../enum/Eitem-type';

export interface IJsonItem {
    name?: string;
    item: string;
    description?: string;
    type?: EItemType;
    shell?: EShellType;
}