export interface UserStateInterface {
    token: string;
    tokenKey: string;
    info: object;
    menuList: MenuListInterface[],
    routeList: []
}

export interface MenuListDataInterface {
    icon: string;
    temp: string;
}

export interface MenuListInterface {
    name: string;
    route: string;
    children: MenuListInterface[];
    data: MenuListDataInterface;
}

export const UserToken = '_token';
export const UserTokenKey = '_tokenKey';
export const UserInfo = '_userInfo';
export const UserMenu = '_menu';
export const UserRouteList = '_routeList'

function state(): UserStateInterface {
    return {
        token: localStorage.getItem(UserToken) ?? '',
        tokenKey: localStorage.getItem(UserTokenKey) ?? '',
        info: localStorage.getItem(UserInfo) ? JSON.parse(<string>localStorage.getItem(UserInfo)) : {},
        menuList: localStorage.getItem(UserMenu) ? JSON.parse(<string>localStorage.getItem(UserMenu)) : [],
        routeList: localStorage.getItem(UserRouteList) ? JSON.parse(<string>localStorage.getItem(UserRouteList)) : []
    };
}

export default state;
