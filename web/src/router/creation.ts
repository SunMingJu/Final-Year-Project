export default [

    {
        path: "/creation",
        name: "Creation",
        component: () => import('@/views/creation/Layout.vue'),
        children: [
            {
                path: "",
                name: "Contribute",
                component: () => import('@/views/creation/contribute/contribute.vue'),
                children: [
                ]
            },
            {
                path: 'videoManagement',
                name: 'VideoManagement',
                meta: {
                    title: 'Video manuscript',
                    requireAuth: false,
                    keepAlive: false
                },
                component: () => import('@/views/creation/manuscript/video.vue')
            },
            {
                path: 'articleManagement',
                name: 'ArticleManagement',
                meta: {
                    title: 'Column manuscript',
                    requireAuth: false,
                    keepAlive: false
                },
                component: () => import('@/views/creation/manuscript/article.vue')
            },
            {
                path: 'barrageDiscuss',
                name: 'BarrageDiscuss',
                meta: {
                    title: 'Barrage management',
                    requireAuth: false,
                    keepAlive: false
                },
                component: () => import('@/views/creation/discuss/barrage.vue')
            },
            {
                path: 'commentDiscuss',
                name: 'CommentDiscuss',
                meta: {
                    title: 'Comment management',
                    requireAuth: false,
                    keepAlive: false
                },
                component: () => import('@/views/creation/discuss/comment.vue')
            },
        ]
    }

]  