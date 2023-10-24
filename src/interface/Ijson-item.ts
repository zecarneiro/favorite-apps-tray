import { EShellType } from '../lib/ts-js-utils/enum/Eshell-type';

export interface IJsonItem {
    nameOrFile: string;
    type: 'name' | 'name-start' | 'system' | 'system-name-start' | 'other' | 'command';
    displayName?: string;
    shell?: EShellType;
}