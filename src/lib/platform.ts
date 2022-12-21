import { ENotifyType } from '../../vendor/utils/typescript/src/enum/Enotify-type';
import { ESystem } from '../../vendor/utils/typescript/src/enum/Esystem';
import { Console } from '../../vendor/utils/typescript/src/lib/console/console';
import { FileSystem } from '../../vendor/utils/typescript/src/lib/file-system';
import { Ui } from '../../vendor/utils/typescript/src/lib/ui';
import { IItemInfo } from '../interface/Iitem-info';
import { IJsonItem } from '../interface/Ijson-item';
import { FavoriteAppsTray } from '../main';
import { IPlatformOperations } from '../interface/Iplatform-operation';

export abstract class Platform implements IPlatformOperations {
    protected ui: Ui;

    constructor(protected fileSystem: FileSystem, protected console: Console) {
        this.ui = new Ui(this.fileSystem, this.console);
    }

    public notify(message: string, type: ENotifyType = ENotifyType.none) {
        this.ui.notify(FavoriteAppsTray.APP_NAME, message, ENotifyType.info, FavoriteAppsTray.DIALOG_TIMEOUT, ESystem.gnome);
    }

    public abstract getInfo(jsonItem: IJsonItem): IItemInfo | null
    public abstract run(info: IItemInfo): void
}
