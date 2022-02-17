import {
    ElMessage
} from 'element-plus'
import router from '../routers'
import store from '../store'
import { logout } from '../api/user'

/**
 * Remind User with Error message
 * 
 * @function tip
 * 
 * @param {string} msg - Error message
 */
export const tip = msg => {
    ElMessage.error(msg)
}

/**
 * Log out user and Back to Login Page
 * 
 * @function toLogin
 */
export const toLogin = async() => {
    await logout()
    await store.dispatch('removeUserInfo')
    router.go('/login')
}

/**
 * To 403 Page
 * 
 * @function to403Page
 */

export const to403Page = () => {
    router.go('/403')
}