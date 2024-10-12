#!/usr/bin/gjs

imports.gi.versions.Gtk = "3.0";
const { Gio, Gtk } = imports.gi;

var dataArr = [];

function fileExist(file) {
  const gfile = Gio.File.new_for_path(file);
  return gfile.query_exists(null);
}

function loadIcon(iconStr) {
  if (iconStr !== "" && !fileExist(iconStr)) {
    const icon = Gio.Icon.new_for_string(iconStr);
    const iconTheme = new Gtk.IconTheme();
    const info = iconTheme.lookup_by_gicon(
      icon,
      32,
      Gtk.IconLookupFlags.FORCE_SIZE,
    );
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
      command = command.replace(/%\w/g, "");
      canStop = false;
    }
    if (command.match(/%./g)) {
      command = command.replace(/%./g, "");
      canStop = false;
    }
    if (command.match(/@.{2}/g)) {
      command = command.replace(/@.{2}/g, "");
      canStop = false;
    }
    if (command.match(/@@/g)) {
      command = command.replace(/@@/g, "");
      canStop = false;
    }
    if (command.match(/^"/g)) {
      command = command.replace(/^"/g, "");
      canStop = false;
    }
    if (command.match(/" \*\$/g)) {
      command = command.replace(/" \*\$/g, "");
      canStop = false;
    }
    if (command.match(/"\$/g)) {
      command = command.replace(/"\$/g, "");
      canStop = false;
    }
    if (command.match(/"/g)) {
      command = command.replace(/"/g, "");
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
  print(JSON.stringify(dataArr, null, 2));
}

function getApplicationInfo(application) {
  if (application) {
    const command = application.get_id() ? `gtk-launch ${application.get_id().substring(0, application.get_id().lastIndexOf(".")) || application.get_id()}`: null;
    let data = {
      displayName: application.get_display_name() ? application.get_display_name() : application.get_name(),
      shortcut: application.get_id(),
      icon: loadIcon(application.get_icon() ? application.get_icon().to_string() : ""),
      command: command ? command : application.get_commandline() ? trimComamnd(application.get_commandline()) : "",
    };
    dataArr.push(data);
  }
}

function main() {
  const applicationsAll = Gio.AppInfo.get_all();
  for (let key in applicationsAll) {
    getApplicationInfo(applicationsAll[key]);
  }
  printData();
}
main();
