import router from './'
import store from '../store'
import {
    ElMessage
} from 'element-plus'
import NProgress from 'nprogress'
import 'nprogress/nprogress.css'

NProgress.configure({
    showSpinner: false
})

const whitelist = ['/login', '/register']

router.beforeEach(async(to, from, next) => {
    // Start Progress Bar
    NProgress.start()

    const session = localStorage.getItem('session.token')

    if (session) {
        if (to.path === '/login') {
            // If is logged in, redirect to the home page
            next({
                path: '/'
            })
            NProgress.done()
        } else {
            if (store.getters.user.user_id) {
                if (to.path === '/' + store.getters.user.user_name) {
                    next('/profile')
                }
                next()
            } else {
                try {
                    await store.dispatch('getMe')
                    if (to.path === '/' + store.getters.user.user_name) {
                        next('/profile')
                    }
                    next()
                } catch (err) {
                    await store.dispatch('removeUserInfo')
                    next('/login')
                }
            }
            NProgress.done()
        }
    } else {
        // Has no token
        if (whitelist.indexOf(to.path) !== -1) {
            next()
            NProgress.done()
        } else {
            // other pages that do not have permission to access are redirected to login page
            next('/login')
            ElMessage.error('Please Log In')
            NProgress.done()
        }
    }
    window.document.title = to.meta.title + ' / Twitter';
})

router.afterEach(() => {
    NProgress.done()
})