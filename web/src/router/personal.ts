export default [{
  path: 'personal',
  name: 'Personal',
  meta: {
    title: 'User Info',
    requireAuth: false,
    keepAlive: false
  },
  component: () => import('@/views/personal/Layout.vue'),
  children: [
    {
      path: '',
      name: 'UserShow',
      meta: {
        title: 'front page',
        requireAuth: false,
        keepAlive: false
      },
      component: () => import('@/views/personal/userInfo/userShow.vue')
    },
    {
      path: 'userinfo',
      name: 'UserInfo',
      meta: {
        title: 'User Info',
        requireAuth: false,
        keepAlive: false
      },
      component: () => import('@/views/personal/userInfo/userInfo.vue')
    },{
      path: 'picturesetting',
      name: 'PictureSetting',
      meta: {
        title: 'User Info',
        requireAuth: false,
        keepAlive: false
      },
      component: () => import('@/views/personal/userInfo/pictureSetting.vue')
    } ,{
      path: 'safety',
      name: 'Safety',
      meta: {
        title: 'Safety',
        requireAuth: false,
        keepAlive: false
      },
      component: () => import('@/views/personal/userInfo/safety.vue')
    },{
      path: 'livesetup',
      name: 'LiveSetUp',
      meta: {
        title: 'Live broadcast settings',
        requireAuth: false,
        keepAlive: false
      },
      component: () => import('@/views/personal/live/setUp.vue')
    },
    {
      path: 'favorites',
      name: 'Favorites',
      meta: {
        title: 'my collection',
        requireAuth: false,
        keepAlive: false
      },
      component: () => import('@/views/personal/collect/favorites.vue')
    },
    {
      path: 'collectList/:id',
      name: 'CollectList',
      meta: {
        title: 'Favorites',
        requireAuth: false,
        keepAlive: false
      },
      component: () => import('@/views/personal/collect/collectList.vue')
    },{
      path: 'record',
      name: 'Record',
      meta: {
        title: 'history record',
        requireAuth: false,
        keepAlive: false
      },
      component: () => import('@/views/personal/record/record.vue')
    },
  ]
}] 