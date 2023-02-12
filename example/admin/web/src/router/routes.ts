import {RouteRecordRaw} from 'vue-router';

const routes: RouteRecordRaw[] = [
    {
        path: '/',
        name: 'Layouts',
        component: () => import('layouts/main.vue'),
        children: [
            {
                path: '',
                name: 'Home',
                component: () => import('pages/index.vue'),
                meta: {requireAuth: true, keepAlive: true}
            }, {
                path: '/console',
                name: 'Console',
                component: () => import('pages/console/index.vue'),
                meta: {requireAuth: true, keepAlive: true}
            }
        ]
    },
    // Always leave this as last one,
    // but you can also remove it
    {
        path: '/login',
        name: 'Login',
        component: () => import('pages/login.vue')
    },
    {
        name: 'Replace',
        path: '/replace',
        component: () => import('pages/replace.vue')
    },
    {
        path: '/:catchAll(.*)*',
        component: () => import('pages/404.vue')
    }
];


export default routes;
