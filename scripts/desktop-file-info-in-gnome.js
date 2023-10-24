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

function getTypeArg(args) {
    if (args && args.length > 1) {
        if (args[1] === 'name') return 'name';
        else if (args[1] === 'name-start') return 'name-start';
        else if (args[1] === 'system') return 'system';
        else if (args[1] === 'system-name-start') return 'system-name-start';
    }
    return '';
}

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

function getById(nameOrFile, isStart = false) {
    const applicationsAll = Gio.AppInfo.get_all();
    for (const key in applicationsAll) {
        const id = applicationsAll[key].get_id() || "";
        if (isStart) {
            if (id.startsWith(nameOrFile)) getApplicationInfo(applicationsAll[key]);
        } else {
            if (id == nameOrFile) getApplicationInfo(applicationsAll[key]);
        }
    }
}

function getByName(nameOrFile, isStart = false) {
    const applicationsAll = Gio.AppInfo.get_all();
    for (const key in applicationsAll) {
        const name = applicationsAll[key].get_name() || "";
        const displayName = applicationsAll[key].get_display_name() || "";
        if (isStart) {
            if (name.startsWith(nameOrFile) || displayName.startsWith(nameOrFile)) getApplicationInfo(applicationsAll[key]);
        } else {
            if (nameApp == applicationsAll[key].get_name() || nameApp == applicationsAll[key].get_display_name()) getApplicationInfo(applicationsAll[key]);
        }
    }
}

function main(args) {
    if (args && args.length > 0) {
        const nameOrFile = args[0] || "";
        const typeArg = getTypeArg(args);
        if (typeArg == "name" || typeArg == "name-start") {
            getByName(nameOrFile, typeArg == "name-start");
        } else if (typeArg == "system" || typeArg == "system-name-start") {
            getById(nameOrFile, typeArg == "system-name-start");
        }
    }
}
main(ARGV);

