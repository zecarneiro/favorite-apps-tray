import { NativeImage, Notification, app, nativeImage } from 'electron';
import path from 'path';
import { FileUtils } from '../vendor/utils/typescript/file-utils';
import { ConsoleUtils } from '../vendor/utils/typescript/console-utils';
import { SystemUtils } from '../vendor/utils/typescript/system-utils';
import { FunctionUtils } from '../vendor/utils/typescript/function-utils';
import { LoggerUtils } from '../vendor/utils/typescript/logger-utils';

export interface IConfigurations {
    isStartup: boolean;
    isSortMenuEntries: boolean;
}

export enum ETheme {
    dark = 'dark',
    light = 'light',
}

export enum EEnvironment {
    development = 'development',
    production = 'production'
}

export abstract class FavoriteAppsProcessor {
    public readonly APP_NAME: string = 'Favorite Apps Tray';
    protected fileUtils: FileUtils ;
    protected consoleUtils: ConsoleUtils;
    protected systemUtils: SystemUtils;
    private alreadySHowInvalidPlatform = false;
    constructor() {
        this.fileUtils = new FileUtils();
        this.consoleUtils = new ConsoleUtils();
        this.systemUtils = new SystemUtils();
        this.setConfigurations();
    }

    protected _theme: ETheme = ETheme.dark;
    protected get theme(): ETheme {
        return this._theme;
    }

    private get homeDir(): string {
        return SystemUtils.systemInfo.homeDir;
    }

    protected get nodeEnvironment(): EEnvironment {
        return SystemUtils.getEnv('NODE_ENV') === 'development' ? EEnvironment.development : EEnvironment.production;
    }

    protected get rootDir(): string {
        return app.getAppPath();
    }

    protected get resourceDir(): string {
        return this.nodeEnvironment === EEnvironment.production ? path.resolve(app.getAppPath(), '..') : app.getAppPath();
    }

    protected get assetsDir(): string {
        return path.resolve(this.resourceDir, 'assets');
    }

    private _logo: NativeImage;
    protected get logo(): NativeImage {
        if (!this._logo) {
            let iconPath = path.resolve(this.assetsDir, 'image/logo/icon');
            if (SystemUtils.isWindows) {
                iconPath += '.ico';
            } else {
                iconPath += '.png';
            }
            this._logo = nativeImage.createFromPath(iconPath);
            this._logo.setTemplateImage(true);
        }
        return this._logo;
    }

    protected get configDir(): string {
        const dir = path.resolve(this.homeDir, '.config', app.getName());
        if (!FileUtils.fileExist(dir)) {
            FileUtils.createDir(dir);
        }
        return dir;
    }

    private get configurationFile(): string {
        return path.resolve(this.configDir, 'configurations.json');
    }

    private _configurations: IConfigurations;
    protected get configurations(): IConfigurations {
        if (!this._configurations) {
            this._configurations = FileUtils.readJsonFile<IConfigurations>(this.configurationFile);
            if (!this._configurations) {
                this._configurations = { isSortMenuEntries: false, isStartup: false };
                FileUtils.writeJsonFile(this.configurationFile, this._configurations);
            }
        }
        return this._configurations;
    }

    protected updateConfigurations(options: IConfigurations) {
        if (options) {
            FileUtils.writeJsonFile(this.configurationFile, options);
            this._configurations = options;
        }
    }

    protected isNodeEnvironment(env: EEnvironment): boolean {
        return this.nodeEnvironment === env;
    }

    protected isValidPlatform(forceNotify?: boolean): boolean {
        const valid = SystemUtils.isWindows || SystemUtils.isGnome;
        if ((!valid && forceNotify) || (!this.alreadySHowInvalidPlatform && !valid)) {
            this.showNotification('Invalid platform');
            this.alreadySHowInvalidPlatform = true;
        }
        return valid;
    }

    protected showProcessing(isDone?: boolean) {
        if (isDone) {
            this.showNotification(`${this.APP_NAME} is ready.`);
        } else {
            this.showNotification('Processing...');
        }
    }

    protected showNotification(message: string) {
        new Notification({ title: this.APP_NAME, body: message, icon: this.logo }).show();
    }

    protected abstract process(showNotify: boolean, forceReadJsonMenuConfig?: boolean): void;
    public abstract distroy(): void;
    public start() {
        try {
            this.process(true);
        } catch (error) {
            LoggerUtils.error(error);
            this.showNotification(error.message);
        }
    }

    private setConfigurations() {
        global.bashUtilsDir = this.isNodeEnvironment(EEnvironment.development) ? path.resolve(this.resourceDir, 'vendor/utils/bash') : path.resolve(this.resourceDir, 'bash');
        global.powershellUtilsDir = this.isNodeEnvironment(EEnvironment.development) ? path.resolve(this.resourceDir, 'vendor/utils/powershell') : path.resolve(this.resourceDir, 'powershell');
        global.appsUtilsDir = this.isNodeEnvironment(EEnvironment.development) ? path.resolve(this.resourceDir, 'vendor/utils/apps') : path.resolve(this.resourceDir, 'apps');
    }
}