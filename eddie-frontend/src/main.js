// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import App from './App'
import router from './router'
import Vuetify from 'vuetify'
import VueResource from 'vue-resource'
import '../node_modules/vuetify/dist/vuetify.min.css'
import persistent from './store/persistent'

Vue.config.productionTip = false
Vue.use(VueResource)

Vue.use(Vuetify)
/* eslint-disable no-new */
new Vue({
  el: '#app',
  router,
  persistent,
  template: '<App/>',
  components: { App }
})
