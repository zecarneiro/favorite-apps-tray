import { MenuItem } from 'systray2';

export interface MenuItemClickable extends MenuItem {
    click?: () => void;
    items?: MenuItemClickable[];
}
