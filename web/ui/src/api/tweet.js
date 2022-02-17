import { WithoutAuth, Service } from '../utils/interceptor'

export const createTweet = (
    form
) => {
    return Service({
        url: '/tweets',
        method: 'post',
        data: form
    })
}

export const createReplyTweet = (
    id,
    form
) => {
    return Service({
        url: '/tweets/' + id + '/reply',
        method: 'post',
        data: form
    })
}

export const getTweets = (page) => {
    if (page === '') {
        page = ''
    } else {
        page = '?page=' + page
    }
    return WithoutAuth({
        url: '/tweets' + page + '&size=10',
        method: 'get'
    })
}

export const getTweetByID = (id) => {
    return WithoutAuth({
        url: '/tweets/' + id,
        method: 'get'
    })
}

export const getReplyTweets = (id) => {
    return WithoutAuth({
        url: '/tweets/' + id + '/replys',
        method: 'get'
    })
}

export const deleteTweet = (id) => {
    return Service({
        url: '/tweets/' + id,
        method: 'delete'
    })
}


export const getLikedUser = (id, page) => {
    if (page === '') {
        page = ''
    } else {
        page = '?page=' + page
    }
    return WithoutAuth({
        url: '/tweets/' + id + '/liking_users' + page + '&size=10',
        method: 'get'
    })
}