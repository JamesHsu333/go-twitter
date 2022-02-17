import { createRouter, createWebHistory } from "vue-router";
import Layout from '../layout/index.vue';


const routes = [{
        path: '/login',
        name: 'login',
        hidden: true,
        component: () =>
            import ('../views/login/index.vue'),
        meta: {
            title: 'Login'
        }
    },
    {
        path: '/register',
        name: 'register',
        hidden: true,
        component: () =>
            import ('../views/register/index.vue'),
        meta: {
            title: 'Register'
        }
    },
    {
        path: '/403',
        component: () =>
            import ('../views/error/forbidden.vue'),
        hidden: true,
        meta: {
            title: '403'
        }
    },
    {
        path: '/404',
        component: () =>
            import ('../views/error/notfound.vue'),
        hidden: true,
        meta: {
            title: '404'
        }
    },
    {
        path: '/',
        component: Layout,
        redirect: '/home',
        children: [{
            path: 'home',
            name: 'home',
            component: () =>
                import ('../views/home/index.vue'),
            meta: {
                title: 'Home',
                icon: 'house'
            }
        }]
    },
    {
        path: '/:user',
        component: Layout,
        hidden: true,
        children: [{
            path: '',
            name: 'user',
            component: () =>
                import ('../views/user/index.vue'),
            meta: {
                title: 'User',
                icon: 'user'
            }
        }]
    },
    {
        path: '/:user/followers',
        component: Layout,
        hidden: true,
        children: [{
            path: '',
            name: 'follow',
            component: () =>
                import ('../views/follow/index.vue'),
            meta: {
                title: 'Follow'
            },
            alias: ['/:user/following']
        }]
    },
    {
        path: '/:user/status/:tweet',
        component: Layout,
        hidden: true,
        children: [{
            path: '',
            name: 'tweet',
            component: () =>
                import ('../views/tweet/index.vue'),
            meta: {
                title: 'Tweet'
            }
        }]
    },
    {
        path: '/profile',
        component: Layout,
        children: [{
            path: '',
            name: 'profile',
            component: () =>
                import ('../views/profile/index.vue'),
            meta: {
                title: 'Profile',
                icon: 'user'
            }
        }]
    },
    {
        path: '/configuration',
        component: Layout,
        redirect: '/configuration/profile',
        children: [{
                path: 'profile',
                component: () =>
                    import ('../views/configuration/index.vue'),
                meta: {
                    title: 'Configuration',
                    icon: 'setting'
                },
                children: [{
                    path: '',
                    name: 'Profile',
                    component: () =>
                        import ('../views/configuration/profile/index.vue'),
                    meta: {
                        title: 'Profile',
                        icon: 'document',
                    },
                }]
            },
            {
                path: 'users',
                component: () =>
                    import ('../views/configuration/index.vue'),
                meta: {
                    title: 'Configuration',
                    icon: 'user'
                },
                children: [{
                    path: '',
                    name: 'Users',
                    component: () =>
                        import ('../views/configuration/users/index.vue'),
                    meta: {
                        title: 'Users',
                        icon: 'user',
                    },
                }]
            }
        ]
    },
    { path: '/:pathMatch(.*)*', redirect: '/404', hidden: true }
];

export default createRouter({
    history: createWebHistory(),
    routes
})