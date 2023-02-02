import { IJsonItem } from './Ijson-item';

export interface IJson {
    noMenu: IJsonItem[],
    others: { [key: string]: IJsonItem[]; },
}