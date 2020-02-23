<template>
  <v-app id="app">
    <v-drawer></v-drawer>
    <v-toolbar></v-toolbar>
    <v-progress-linear v-if="loading" v-bind:indeterminate="true"></v-progress-linear>
      <v-alert dismissible v-model="alert" transition="scale-transition">
      {{alertMessage}}
    </v-alert>
      <v-content >
    <v-container light-blue accent-4 fluid style="height: 100%;">
    <router-view/>
    </v-container>
      </v-content>
  </v-app>
</template>

<script>
import Toolbar from '@/components/Toolbar'
import Drawer from '@/components/Drawer'
import {EventBus} from './event-bus.js'
import Persistent from './store/persistent'
import Events from './store/event-api.js'

const DEV_HOST = 'http://127.0.0.1:8000'

export default {
  name: 'app',
  data () {
    return {
      loading: false,
      alert: false,
      alertMessage: '',
      user: Persistent.getters.user,
      company: Persistent.getters.company
    }
  },
  components: {
    'v-toolbar': Toolbar,
    'v-drawer': Drawer
  },
  created: function () {
    var me = this
    EventBus.$on(Events.apiError, function (status, message) {
      if (status === 0) {
        status = ''
        message = 'Server is unreachable, is your network OK ?'
      }
      if (status === 401) {
        message = 'Authentification error'
        setTimeout(function () {
          Persistent.commit('deleteAuthToken')
          window.location = '/'
        }, 5500)
      }
      me.alertMessage = status + ' ' + message
      me.alert = true
      setTimeout(function () {
        me.alert = false
      }, 5000)
    })
    EventBus.$on(Events.loadingStart, function () {
      me.loading = true
    })
    EventBus.$on(Events.loadingEnd, function () {
      me.loading = false
    })
    if (process.env.NODE_ENV === 'development') {
      Persistent.commit('setBaseUrl', DEV_HOST)
    } else {
      Persistent.commit('setBaseUrl', '')
    }
  }
}
</script>

<style>
body {
  background-color:blue;
}
#app {

}
.alert.alert {
  width: 100%;
  margin: 0;
  position: absolute;
  top: 60px;
  z-index: 10000;
}
.progress-linear {
  margin-top: 0px;
  margin-bottom: 0px;
}
</style>
