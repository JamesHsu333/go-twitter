import store from '../store'

export const getPermission = () => {
    let role = store.getters.user.role
    if (role === 'admin') {
        return true
    }
    return false
}