import { IJsonItem } from './Ijson-item';

export interface IJson {
    SORT: boolean,
    NO_MENU: IJsonItem[],
    OTHERS: { [key: string]: IJsonItem[]; },
}
