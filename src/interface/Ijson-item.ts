export interface IJsonItem {
    name?: string;
    item: string;
    description?: string;
    type?: 'item' | 'custom-item' | 'command';
    shell?: 'powershell' | 'cmd' | 'bash' | 'terminal-osx';
}
