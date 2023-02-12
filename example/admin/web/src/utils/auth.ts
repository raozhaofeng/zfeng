import store from 'src/store'

// 判断是否有权限
export const isAuth = (url: string): boolean => {
    const indexStr = '/admin' + url
    return store.state.user.routeList.indexOf(<never>'*') > -1 || store.state.user.routeList.indexOf(<never>indexStr) > -1
}