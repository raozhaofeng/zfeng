import { MutationTree } from 'vuex';
import { UserStateInterface } from './state';
import { UserToken, UserTokenKey, UserInfo, UserMenu, UserRouteList, MenuListInterface } from './state';

const mutation: MutationTree<UserStateInterface> = {
    // 更新用户Token
    updateToken(state: UserStateInterface, token: string) {
        state.token = token;
        localStorage.setItem(UserToken, token);
    },
    // 更新用户TokenKey
    updateTokenKey(state: UserStateInterface, tokenKey: string){
        state.tokenKey = tokenKey
        localStorage.setItem(UserTokenKey, tokenKey)
    },
    // 更新用户信息
    updateInfo(state: UserStateInterface, info: object) {
        state.info = info;
        localStorage.setItem(UserInfo, JSON.stringify(info));
    },
    // 更新用户路由
    updateRouteList(state: UserStateInterface, routeList: []) {
        state.routeList = routeList
        localStorage.setItem(UserRouteList, JSON.stringify(routeList))
    },
    // 更新用户菜单
    updateMenu(state: UserStateInterface, menuList: MenuListInterface[]) {
        state.menuList = menuList;
        localStorage.setItem(UserMenu, JSON.stringify(menuList));
    }
};

export default mutation;
