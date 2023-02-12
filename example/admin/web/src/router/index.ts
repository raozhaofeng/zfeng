// import { route } from 'quasar/wrappers';
import {
    createMemoryHistory,
    createRouter,
    createWebHashHistory,
    createWebHistory
} from 'vue-router';
import store from 'src/store';
// import { StateInterface } from 'src/store';

import routes from 'src/router/routes';
import {MenuListInterface} from 'src/store/user/state';

const routePageList = import.meta.glob('../pages/**/*.vue')

/*
 * If not building with SSR mode, you can
 * directly export the Router instantiation;
 *
 * The function below can be async too; either use
 * async/await or return a Promise which resolves
 * with the Router instance.
 */

// export default route<StateInterface>(function({ /* store, ssrContext */ }) {
//     const createHistory = process.env.SERVER
//         ? createMemoryHistory
//         : (process.env.VUE_ROUTER_MODE === 'history' ? createWebHistory : createWebHashHistory);
//
//     const Router = createRouter({
//         scrollBehavior: () => ({ left: 0, top: 0 }),
//         routes,
//
//         // Leave this as is and make changes in quasar.conf.js instead!
//         // quasar.conf.js -> build -> vueRouterMode
//         // quasar.conf.js -> build -> publicPath
//         history: createHistory(process.env.VUE_ROUTER_BASE)
//     });
//
//     return Router;
// });


const createHistory = process.env.SERVER
    ? createMemoryHistory
    : (process.env.VUE_ROUTER_MODE === 'history' ? createWebHistory : createWebHashHistory);

const Router = createRouter({
    scrollBehavior: () => ({left: 0, top: 0}),
    routes,

    // Leave this as is and make changes in quasar.conf.js instead!
    // quasar.conf.js -> build -> vueRouterMode
    // quasar.conf.js -> build -> publicPath
    history: createHistory(process.env.VUE_ROUTER_BASE)
});

// 路由前置守卫
Router.beforeEach((to, form, next) => {
    //  如果是进入登录页面并且
    if ((to.name === 'Login' || to.name === 'Register') && store.getters['user/isLogin']) {
        next({name: 'Home'});
    } else {
        // 验证是否跳转到登录页面
        if (to.matched.some(record => record.meta.requireAuth) && !store.getters['user/isLogin']) {
            next({name: 'Login', query: {next: to.fullPath}});
        } else {
            next();
        }
    }
});

// 加载动态路由
export const dynamicRouter = (routerList: MenuListInterface[]) => {
    routerList.forEach((item) => {
        if (item.route !== '' && !Router.hasRoute(item.route)) {
            Router.addRoute('Layouts', {
                path: item.route,
                component: routePageList['../pages' + item.data.temp + '.vue'],
                meta: {requireAuth: true, keepAlive: true}
            });
        }
        if (item.hasOwnProperty('children') && item.children !== null && item.children.length > 0) {
            dynamicRouter(item.children);
        }
    });
};

// 加载动态路由
dynamicRouter(store.state.user.menuList);

export default Router;
