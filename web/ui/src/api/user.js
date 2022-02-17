import { WithoutAuth, Service } from '../utils/interceptor'

export const register = ({
    user_name,
    name,
    email,
    password
}) => {
    return WithoutAuth({
        url: '/users/register',
        method: 'post',
        data: {
            user_name: user_name,
            name: name,
            email: email,
            password: password
        }
    })
}

export const login = (email, password) => {
    return WithoutAuth({
        url: '/users/login',
        method: 'post',
        data: {
            email: email,
            password: password
        }
    })
}

export const logout = () => {
    return WithoutAuth({
        url: '/users/logout',
        method: 'post'
    })
}

export const getMe = () => {
    return WithoutAuth({
        url: '/users/me',
        method: 'get'
    })
}

export const getUserByID = (user_id) => {
    return WithoutAuth({
        url: '/users/' + user_id,
        method: 'get'
    })
}

export const getUserByName = (user_name) => {
    return WithoutAuth({
        url: '/users/username/' + user_name,
        method: 'get'
    })
}

export const getAllUser = (page) => {
    let pageQuery = ''
    if (page !== '') {
        pageQuery = 'page=' + page
    }
    return WithoutAuth({
        url: '/users?' + pageQuery,
        method: 'get'
    })
}

export const updateUser = (user_id, {
    user_name,
    name,
    gender,
    email,
    about,
    phone_number,
    country,
    birthday
}) => {
    return Service({
        url: '/users/' + user_id,
        method: 'patch',
        data: {
            user_name: user_name,
            name: name,
            gender: gender,
            email: email,
            about: about,
            phone_number: phone_number,
            country: country,
            birthday: birthday
        }
    })
}

export const updateUserRole = (user_id, role) => {
    return Service({
        url: '/users/' + user_id + '/role',
        method: 'patch',
        data: {
            role: role
        }
    })
}

export const deleteUser = (user_id) => {
    return Service({
        url: '/users/' + user_id,
        method: 'delete'
    })
}

export const follow = (follower_id, following_id) => {
    return Service({
        url: '/users/' + follower_id + '/following',
        method: 'post',
        data: {
            user_id: following_id
        }
    })
}

export const getFollowing = (user_id, page) => {
    if (page === '') {
        page = ''
    } else {
        page = '?page=' + page
    }
    return Service({
        url: '/users/' + user_id + '/following' + page + '&size=10',
        method: 'get'
    })
}

export const getFollowers = (user_id, page) => {
    if (page === '') {
        page = ''
    } else {
        page = '?page=' + page
    }
    return Service({
        url: '/users/' + user_id + '/followers' + page + '&size=10',
        method: 'get'
    })
}

export const deleteFollowing = (follower_id, following_id) => {
    return Service({
        url: '/users/' + follower_id + '/following/' + following_id,
        method: 'delete'
    })
}

export const getTweetsByUserID = (user_id, page) => {
    if (page === '') {
        page = ''
    } else {
        page = '?page=' + page
    }
    return WithoutAuth({
        url: '/users/' + user_id + '/tweets' + page + '&size=10',
        method: 'get'
    })
}

export const getLikedTweets = (user_id, page) => {
    if (page === '') {
        page = ''
    } else {
        page = '?page=' + page
    }
    return WithoutAuth({
        url: '/users/' + user_id + '/liked_tweets' + page + '&size=10',
        method: 'get'
    })
}

export const likeTweet = (user_id, tweet_id) => {
    return Service({
        url: '/users/' + user_id + '/liked',
        method: 'post',
        data: {
            tweet_id: tweet_id
        }
    })
}

export const deleteLike = (user_id, tweet_id) => {
    return Service({
        url: '/users/' + user_id + '/liked/' + tweet_id,
        method: 'delete'
    })
}

export const getCSRFToken = () => {
    return WithoutAuth({
        url: '/users/token',
        method: 'get'
    })
}