import { Console } from '../../vendor/utils/typescript/src/lib/console/console';
import { FileSystem } from '../../vendor/utils/typescript/src/lib/file-system';
import { IItemInfo } from '../interface/Iitem-info';
import { IJsonItem } from '../interface/Ijson-item';
import { Platform } from './platform';

export class NoPlatform extends Platform {
    constructor(protected fileSystem: FileSystem, protected console: Console) {
        super(fileSystem, console);
    }

    public getInfo(jsonItem: IJsonItem): IItemInfo | null {
        return null;
    }

    public run(file: IItemInfo) {
      // TODO document why this method 'run' is empty
     }
}
