import axios from 'axios'
import store from '../store'
import {
    tip,
    toLogin,
    to403Page
} from './error'

/**
 * Handling Request Fail 
 * @param {Number} status 
 */

const errorHandle = (status, msg) => {
    switch (status) {
        case 400:
            tip(msg)
            break
        case 401:
            tip('Log in has expired, please login again')
            setTimeout(() => {
                toLogin()
            }, 1000)
            break
        case 403:
            to403Page()
            break
        case 404:
            tip(msg)
            break
        default:
            tip('Unknown Error: ' + msg)
    }
}

//Create API axios instance

const WithoutAuth = axios.create({
    baseURL: '/api/v1',
    timeout: 7000,
    headers: {
        'Content-Type': 'application/json'
    }
})

WithoutAuth.interceptors.request.use(async(config) => {
    return config
}, (error) => {
    return Promise.reject(error)
})

WithoutAuth.interceptors.response.use((response) => {
    return response
}, (error) => {
    if (error) {
        // Get error response
        console.log(error.response)
        errorHandle(error.response.status, error.response.data.error)
        return Promise.reject(error.response.data)
    } else {
        // Error but not get response
        if (!window.navigator.onLine) {
            tip('Internet has been offline, please check')
        } else {
            return Promise.reject(error.response.data)
        }
    }
})

//Create API axios instance

const Service = axios.create({
    baseURL: '/api/v1',
    timeout: 7000,
    headers: {
        'Content-Type': 'application/json'
    }
})

//Service Request interceptor

Service.interceptors.request.use(async(config) => {
    if (!store.getters.csrf_token) {
        await store.dispatch('getCSRFToken')
    }
    config.headers['X-CSRF-Token'] = store.getters.csrf_token
    return config
}, (error) => {
    return Promise.reject(error)
})

//Service Response interceptor

Service.interceptors.response.use((response) => {
    return response
}, (error) => {
    if (error) {
        // Get error response
        console.log(error.response)
        errorHandle(error.response.status, error.response.data.error)
        return Promise.reject(error.response.data)
    } else {
        // Error but not get response
        if (!window.navigator.onLine) {
            tip('Internet has been offline, please check')
        } else {
            return Promise.reject(error.response.data)
        }
    }
})


export { WithoutAuth, Service }