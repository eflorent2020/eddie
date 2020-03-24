<template>
  <v-app-bar       :clipped-left="$vuetify.breakpoint.lgAndUp"
      app dense dark color="indigo darken-3">
     <v-app-bar-nav-icon @click="toggleDrawer()"></v-app-bar-nav-icon>
    <v-toolbar-title class="white--text">{{companyname}}</v-toolbar-title>
    <v-spacer></v-spacer>
    <v-btn icon>
    </v-btn>
        <v-btn color="primary"  v-if="!isLoggedIn" slot="activator">Login</v-btn>
          <v-btn v-if="isLoggedIn" icon @click="logout">
         <v-icon>exit_to_app</v-icon>
        </v-btn>   
    <v-dialog v-model="openLogin" persistent max-width="500px">
        <v-card>
          <v-card-title>
            <span class="headline">Login</span>
          </v-card-title>
          <v-card-text>
            <v-container grid-list-md>
              <v-layout wrap>
                <v-flex xs12>
                  <v-text-field :rules="emailRules" label="Email" v-model="login" required></v-text-field>
                </v-flex>
                <v-flex xs12>
                  <v-text-field label="Password" v-model="password" :rules="ruleReq" type="password" required></v-text-field>
                </v-flex>
              </v-layout>
            </v-container>
            <small>*indicates required field</small>
          </v-card-text>
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn color="blue darken-1" flat @click.native="closeLogin">Cancel</v-btn>
            <v-btn color="blue darken-1" raised @click.native="proceedToLogin">Login</v-btn>
          </v-card-actions>
        </v-card>
  </v-dialog>
  </v-app-bar>
</template>

<script>

import { EventBus } from '../event-bus.js'
import Events from '../store/event-api.js'
import Persistent from '../store/persistent'
import JwtDecode from 'jwt-decode'
import Vue from 'vue'

const API_VERSION = '/api/v1/rest/'

export default {
  name: 'toolbar',
  data () {
    return {
      openLogin: false,
      login: null,
      password: null,
      ruleReq: [
        (v) => !!v || 'Field is required'
      ],
      emailRules: [
        (v) => !!v || 'E-mail is required',
        /* eslint-disable no-useless-escape */
        (v) => /^\w+([\.-]?\w+)*@\w+([\.-]?\w+)*(\.\w{2,3})+$/.test(v) || 'E-mail must be valid'
      ]
    }
  },
  created () {
    if (Persistent.getters.authToken != null) {
      try {
        let data = JwtDecode(Persistent.getters.authToken)
        if (data.exp > new Date().getTime() / 1000) {
          this.setUpHttpHeader(Persistent.getters.authToken)
        }
      } catch (e) {
        console.log(e)
      }
    }
    this.openLogin = !this.isLoggedIn
  },
  methods: {
    logout: function () {
      Persistent.commit('deleteAuthToken')
      // localStorage.clear()
      window.location = '/'
    },
    closeLogin: function () {
      this.openLogin = false
      EventBus.$emit(Events.loadingEnd)
    },
    proceedToLogin: function () {
      if (!this.login || !this.password) {
        console.log('empty req.canceled')
        return
      }
      EventBus.$emit(Events.loadingStart)
      let url = Persistent.getters.baseUrl + '/login'
      var data = {
        username: this.login,
        password: this.password
      }
      this.$http.post(url, data).then(function (res) {
        if (res.body.token) {
          Persistent.commit('addAuthToken', res.body.token)
          let data = JwtDecode(res.body.token)
          this.setUpHttpHeader(res.body.token)
          this.preloadCompany(data.companyID, data.uid)
          this.openLogin = false
        } else {
          this.notifyError(res.status, res.body)
        }
        EventBus.$emit(Events.toggleDrawer)
      }, response => {
        this.notifyError(response.status, response.body.message)
      })
    },
    setUpHttpHeader: function (token) {
      Vue.http.interceptors.push((request, next) => {
        request.headers.set('Authorization', 'Bearer ' + token)
        request.headers.set('Accept', 'application/json')
        next()
      })
    },
    preloadUser: function (id) {
      let url = Persistent.getters.baseUrl + API_VERSION + 'user/' + id
      this.$http.get(url).then(function (res) {
        Persistent.commit('setUser', res.body.data)
        EventBus.$emit(Events.loadingEnd)
      }, response => {
        this.notifyError(response.status, response.body.message)
      })
    },
    preloadCompany: function (cid, uid) {
      let url = Persistent.getters.baseUrl + API_VERSION + 'company/' + cid
      this.$http.get(url).then(function (res) {
        Persistent.commit('setUserCompany', res.body.data)
        this.preloadUser(uid)
      }, response => {
        this.notifyError(response.status, response.body.message)
      })
    },
    notifyError: function (code, message) {
      EventBus.$emit(Events.apiError, code, message)
      EventBus.$emit(Events.loadingEnd)
    },
    toggleDrawer: function () {
      EventBus.$emit(Events.toggleDrawer)
    }
  },
  computed: {
    companyname () {
      try {
        return Persistent.getters.company.Name
      } catch (e) {
        return ''
      }
    },
    isLoggedIn () {
      return Persistent.getters.authToken !== null
    }
  }
}
</script>
<style>

</style>
