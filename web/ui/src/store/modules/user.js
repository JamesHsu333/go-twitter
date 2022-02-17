import { register, login, getCSRFToken, getMe } from '../../api/user'

const state = {
    user: {
        user_id: '',
        user_name: '',
        name: '',
        email: '',
        role: '',
        about: '',
        avatar: '',
        header: '',
        phone_number: '',
        country: '',
        gender: '',
        birthday: '',
        create_at: '',
        updated_at: '',
        login_date: '',
        following: 0,
        followers: 0,
    },
    token: '',
    csrf_token: '',
    session: ''
}

const mutations = {
    SET_USER: (state, user) => {
        state.user = user
    },
    SET_USER_FOLLOWING: (state, increment) => {
        state.user.following += increment
    },
    SET_TOKEN: (state, token) => {
        state.token = token
    },
    SET_CSRF_TOKEN: (state, csrf_token) => {
        state.csrf_token = csrf_token
    },
    SET_SESSION: (state, session) => {
        state.session = session
    }
}

const actions = {
    register({ commit }, {
        user_name,
        name,
        email,
        password
    }) {
        return new Promise((resolve, reject) => {
            register({
                user_name,
                name,
                email,
                password
            }).then((res) => {
                localStorage.setItem('session.token', res.data.token);
                commit('SET_TOKEN', res.data.token)
                commit('SET_USER', res.data.user)
                resolve()
            }).catch((err) => {
                reject(err)
            })
        })
    },
    login({ commit }, { email, password }) {
        return new Promise((resolve, reject) => {
            login(email, password).then((res) => {
                localStorage.setItem('session.token', res.data.token);
                commit('SET_TOKEN', res.data.token)
                commit('SET_USER', res.data.user)
                resolve()
            }).catch((err) => {
                reject(err)
            })
        })
    },
    getCSRFToken({ commit }) {
        return new Promise((resolve, reject) => {
            getCSRFToken().then((res) => {
                commit('SET_CSRF_TOKEN', res.headers['x-csrf-token'])
                resolve()
            }).catch((err) => {
                reject(err)
            })
        })
    },
    getMe({ commit }) {
        return new Promise((resolve, reject) => {
            getMe().then((res) => {
                commit('SET_USER', res.data)
                resolve()
            }).catch((err) => {
                reject(err)
            })
        })
    },
    removeUserInfo({
        commit
    }) {
        return new Promise(
            resolve => {
                localStorage.clear()
                commit('SET_TOKEN', '')
                commit('SET_USER', '')
                commit('SET_CSRF_TOKEN', '')
                resolve()
            }
        )
    },
    updateUserInfo({
        commit
    }, user) {
        return new Promise(
            resolve => {
                commit('SET_USER', user)
                resolve()
            }
        )
    },
    updateUserFollowing({ commit }, increment) {
        return new Promise(
            resolve => {
                commit('SET_USER_FOLLOWING', increment)
                resolve
            }
        )
    }
}

export default {
    state,
    mutations,
    actions
}