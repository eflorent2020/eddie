import Vue from 'vue'
import Vuex from 'vuex'
import createPersistedState from 'vuex-persistedstate'

Vue.use(Vuex)

export default new Vuex.Store({
  plugins: [createPersistedState()],
  state: {
    company: null,
    companyID: null,
    user: null,
    userID: null,
    authToken: null,
    baseUrl: null
  },
  mutations: {
    addAuthToken: (state, authToken) => {
      state.authToken = authToken
    },
    deleteAuthToken: (state) => {
      state.authToken = null
    },
    setUser: (state, user) => {
      state.user = user
    },
    setUserCompany: (state, company) => {
      state.company = company
    },
    setBaseUrl: (state, url) => {
      state.baseUrl = url
    }
  },
  getters: {
    authToken: state => state.authToken,
    user: state => state.user,
    baseUrl: state => state.baseUrl,
    company: state => state.company,
    userEmail: state => state.userEmail,
    isLoggedIn: state => {
      return state.authToken !== null
    }
  }
})
