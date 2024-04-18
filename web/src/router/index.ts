import { useGlobalStore } from "@/store/main";
import { createRouter, createWebHistory } from "vue-router";

import creation from "@/router/creation";
import liveRouter from "@/router/live";
import personalRouter from "@/router/personal";
import articleShow from "@/router/show/article";
import videoShow from "@/router/show/video";
import search from "@/router/search"
import space from "@/router/space";
const routes = [
    {
        path: '',
        name: 'Home',
        meta: {
            title: 'front pageont page',
            requireAuth: false,
            keepAlive: false
        },
        component: () => import('@/views/home/home.vue')
    },
    {
        path: '/column',
        name: 'Column',
        meta: {
            title: 'Column',
            requireAuth: false,
            keepAlive: false
        },
        component: () => import('@/views/home/column.vue')
    },
    {
        path: '/live',
        name: 'Live',
        meta: {
            title: 'Column',
            requireAuth: false,
            keepAlive: false
        },
        component: () => import('@/views/home/live.vue')
    },
    {
        path: "/",
        name: "Index",
        component: () => import('@/views/Layout.vue'),
        children: [
            ...personalRouter,
            ...liveRouter,
            ...videoShow,
            ...space,
        ]
    },
    //login
    {
        path: "/login",
        name: "Login",
        meta: {
            title: 'login',
            requireAuth: false,
            keepAlive: false
        },
        component: () => import('@/views/login/login.vue'),
    },
    ...search,
    //Creation Center
    ...creation,
    //Column display
    ...articleShow
    //Route not matched 404
    ,
    {
        path: '/:pathMatch(.*)*',
        name: 'NotFound',
        component: () => import('@/views/404.vue'),

    }

];

const router = createRouter({
    history: createWebHistory(),
    routes,
});


router.beforeEach((to, from, next) => {
    const globalStore = useGlobalStore()
    globalStore.globalData.router.currentRouter = to.path

    next()
})

export default router;
