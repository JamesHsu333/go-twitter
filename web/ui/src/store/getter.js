const getters = {
    user: state => state.user.user,
    session: state => state.user.session,
    csrf_token: state => state.user.csrf_token,
    token: state => state.user.token,
}

export default getters