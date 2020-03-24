<template>

<v-navigation-drawer
  dark
  temporary
  v-model="drawer"
  absolute
  >
           <v-list
          dense
          nav
          class="py-0"
        >

        <v-list-item @click.stop="toggle">
            <v-list-item-icon>
              <v-icon>chevron_left</v-icon>
            </v-list-item-icon>

            <v-list-item-content>
              <v-list-item-title></v-list-item-title>
            </v-list-item-content>
        </v-list-item>

          <v-list-item
            v-for="item in items"
            :key="item.title"
            @click="navigate(item.route)"
          >
            <v-list-item-icon>
              <v-icon>{{ item.icon }}</v-icon>
            </v-list-item-icon>

            <v-list-item-content>
              <v-list-item-title>{{ item.title }}</v-list-item-title>
            </v-list-item-content>
          </v-list-item>
        </v-list>

    </v-navigation-drawer>
</template>
<script>

import { EventBus } from '../event-bus.js'
import Events from '../store/event-api.js'
import router from '../router'
import Persistent from '../store/persistent'

export default {
  created: function () {
    EventBus.$on(Events.toggleDrawer, this.toggle)
  },
  methods: {
    toggle () {
      this.drawer = !this.drawer
    },
    navigate (location) {
      router.push(location)
    }
  },
  computed: {
    username () {
      try {
        return Persistent.getters.user.Name
      } catch (e) {
        return ''
      }
    },
    companyname () {
      try {
        return Persistent.getters.company.Name
      } catch (e) {
        return ''
      }
    }
  },
  data () {
    return {
      drawer: null,
      items: [
        { title: 'Users', icon: 'people', route: '/users' },
        { title: 'Company', icon: 'account_balance', route: '/company' },
        { title: 'PDF Templates', icon: 'content_paste', route: '/pdftemplates' },
        { title: 'P12 Signatures', icon: 'vpn_key', route: '/p12signatures' },
        // { title: 'Documents', icon: 'storage', route: '/documents' },
        // { title: 'PGP Signatures', icon: 'lock', route: '/pgpsignatures' },
        { title: 'Mail Templates', icon: 'mail_outline', route: '/mailtemplates' },
        // { title: 'Sent mails', icon: 'mail', route: '/mails' },
        { title: 'Log', icon: 'sort', route: '/logentries' },
        { title: 'Help', icon: 'help', route: '/help' },
        { title: 'About', icon: 'question_answer', route: '/about' }
      ],
      mini: false,
      right: null
    }
  }
}
</script>