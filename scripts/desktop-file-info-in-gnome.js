#!/usr/bin/gjs

imports.gi.versions.Gtk = '3.0';
const { Gio, Gtk } = imports.gi;

const data = {
    name: '',
    icon: '',
    command: '',
    commandByGtk: '',
    displayName: '',
    desktopFile: '',
};

function fileExist(file) {
    const gfile = Gio.File.new_for_path(file);
    return gfile.query_exists(null);
}

function loadIcon(iconStr) {
    if (iconStr !== '' && !fileExist(iconStr)) {
        const icon = Gio.Icon.new_for_string(iconStr);
        const iconTheme = new Gtk.IconTheme();
        const info = iconTheme.lookup_by_gicon(icon, 32, Gtk.IconLookupFlags.FORCE_SIZE);
        if (info) {
            return info.get_filename();
        }
    }
    return iconStr;
}

function trimComamnd(command) {
    let canStop = true;
    // removes any arguments - %u, %f, @U, @@u, @@ etc
    // removes " around command (if present)
    while (true) {
        if (command.match(/%\w/g)) {
            command = command.replace(/%\w/g, '');
            canStop = false;
        }
        if (command.match(/%./g)) {
            command = command.replace(/%./g, '');
            canStop = false;
        }
        if (command.match(/@.{2}/g)) {
            command = command.replace(/@.{2}/g, '');
            canStop = false;
        }
        if (command.match(/@@/g)) {
            command = command.replace(/@@/g, '');
            canStop = false;
        }
        if (command.match(/^"/g)) {
            command = command.replace(/^"/g, '');
            canStop = false;
        }
        if (command.match(/" \*\$/g)) {
            command = command.replace(/" \*\$/g, '');
            canStop = false;
        }
        if (command.match(/"\$/g)) {
            command = command.replace(/"\$/g, '');
            canStop = false;
        }
        if (command.match(/"/g)) {
            command = command.replace(/"/g, '');
            canStop = false;
        }
        if (canStop) {
            break;
        }
        canStop = true;
    }
    return command.trim();
}

function isDesktopFile(desktopFile) {
    if (desktopFile.lastIndexOf('.') < 0) {
        return false;
    }
    const extension = desktopFile.substring(desktopFile.lastIndexOf('.') + 1, desktopFile.length) || '';
    return extension === 'desktop';
}

function printData() {
    print(JSON.stringify(data));
}

function getApplicationInfo(application) {
    if (application) {
        const command = application.get_commandline();
        data.name = application.get_name();
        data.icon = loadIcon(application.get_icon().to_string());
        data.command = command ? trimComamnd(command) : '';
        data.displayName = application.get_display_name();
        data.desktopFile = application.get_id();
        data.commandByGtk = `gtk-launch ${data.desktopFile.substring(0, data.desktopFile.lastIndexOf('.')) || data.desktopFile}`;
        printData();
    }
}

function getAppByDesktopExtension(desktopFile) {
    const applicationsAll = Gio.AppInfo.get_all();
    for (const key in applicationsAll) {
        if (applicationsAll[key].get_id() == desktopFile) {
            getApplicationInfo(applicationsAll[key]);
        }
    }
}

function getAppByName(nameApp) {
    const applicationsAll = Gio.AppInfo.get_all();
    for (const key in applicationsAll) {
        if (applicationsAll[key].get_name() == nameApp || applicationsAll[key].get_display_name() == nameApp) {
            getApplicationInfo(applicationsAll[key]);
        }
    }
}

function getAppByDesktopExtension(desktopFile) {
    if (desktopFile) {
        const applicationsAll = Gio.AppInfo.get_all();
        for (const key in applicationsAll) {
            if (applicationsAll[key].get_id() == desktopFile) {
                getApplicationInfo(applicationsAll[key]);
            }
        }
    }
}

function main(args) {
    if (args && args.length == 1) {
        if (!isDesktopFile(args[0])) {
            getAppByName(args[0]);
        } else {
            getAppByDesktopExtension(args[0]);
        }
    }
}
main(ARGV);

