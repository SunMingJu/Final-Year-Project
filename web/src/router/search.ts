export default [{
    path: "/search/:text",
    name: "Search",
    meta: {
        title: 'search',
        requireAuth: false,
        keepAlive: false
    },
    component: () => import('@/views/search/Layout.vue'),
    children: [
        {
            path: 'video',
            name: 'VideoSearch',
            meta: {
                title: 'Video search',
                requireAuth: false,
                keepAlive: false
            },
            component: () => import('@/views/search/video.vue')
        },
        {
            path: 'user',
            name: 'UserSearch',
            meta: {
                title: 'User search',
                requireAuth: false,
                keepAlive: false
            },
            component: () => import('@/views/search/user.vue')
        },
    ]
}];