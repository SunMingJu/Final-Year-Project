export default [{
    path: 'space/:id',
    name: 'Space',
    meta: {
        title: 'User Info',
        requireAuth: false,
        keepAlive: false
    },
    component: () => import('@/views/space/Layout.vue'),
    children: [
        {
            path: 'individual',
            name: 'SpaceIndividual',
            meta: {
                title: 'personal space',
                requireAuth: false,
                keepAlive: false
            },
            component: () => import('@/views/space/space.vue')
        },
        {
            path: 'myAttention',
            name: 'MyAttention',
            meta: {
                title: 'my focus',
                requireAuth: false,
                keepAlive: false
            },
            component: () => import('@/views/space/myAttention.vue')
        },
        {
            path: 'myVermicelli',
            name: 'MyVermicelli',
            meta: {
                title: 'my fans',
                requireAuth: false,
                keepAlive: false
            },
            component: () => import('@/views/space/myVermicelli.vue')
        }
    ]
}]