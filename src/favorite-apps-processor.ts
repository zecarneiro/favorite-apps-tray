import { NativeImage, Notification, app, nativeImage } from 'electron';
import path from 'path';
import { FileUtils } from './lib/ts-js-utils/file-utils';
import { ConsoleUtils } from './lib/ts-js-utils/console-utils';
import { SystemUtils } from './lib/ts-js-utils/system-utils';
import { LoggerUtils } from './lib/ts-js-utils/logger-utils';
import { EShellType } from './lib/ts-js-utils/enum/Eshell-type';
import { IBoot } from './interface/Iboot';
import { ICommandInfo } from './lib/ts-js-utils/interface/Icommand-info';
import { IDirectories } from './interface/Idirectories';

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
    protected fileUtils: FileUtils;
    protected consoleUtils: ConsoleUtils;
    protected systemUtils: SystemUtils;
    private alreadySHowInvalidPlatform = false;
    constructor() {
        this.fileUtils = new FileUtils();
        this.consoleUtils = new ConsoleUtils();
        this.systemUtils = new SystemUtils();
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

    protected get directories(): IDirectories {
        const scriptsDir = path.resolve(this.resourceDir, 'scripts');
        return {
            bashScript: path.resolve(scriptsDir, 'bash-utils'),
            powershellScript: path.resolve(scriptsDir, 'powershell-utils'),
            script: scriptsDir
        };
    }


    protected addBootApp(info: IBoot): boolean {
        let cmdInfo: ICommandInfo | undefined;
        if (SystemUtils.isWindows && (info.command || info.isDelete)) {
            cmdInfo = { cmd: `. "${this.directories.powershellScript}\\MainUtils.ps1"`, shellType: EShellType.powershell };
            if (info.command) {
                cmdInfo = { ...cmdInfo,
                    args: [
                        `;AddBootApplication -name "${info.name}" -command "${info.command}" ${info.hidden ? ' -hidden' : ''}`,
                    ],
                };
            } else {
                cmdInfo = { ...cmdInfo,
                    args: [
                        `;DelBootApplication -name "${info.name}"`,
                    ],
                };
            }
        } else if (SystemUtils.isLinux && (info.command || info.isDelete)) {
            cmdInfo = { cmd: `. "${this.directories.bashScript}/main-utils.sh"`, shellType: EShellType.bash };
            if (info.command) {
                cmdInfo = { ...cmdInfo,
                    args: [
                        `;add_boot_application "${info.name}" "${info.command}" ${info.hidden ? '1' : ''}`,
                    ],
                };
            } else {
                cmdInfo = { ...cmdInfo,
                    args: [
                        `;del_boot_application "${info.name}"`,
                    ],
                };
            }
        }
        if (cmdInfo) {
            const response = this.consoleUtils.execSync({ ...cmdInfo, verbose: true, isThrow: false, verboseOnlyCommand: true });
            return response.hasError ? false : true;
        }
        return false;
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
}