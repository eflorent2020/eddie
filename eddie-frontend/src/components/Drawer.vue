<template>

<v-navigation-drawer
dark
      temporary
      v-model="drawer"
      absolute
    >
      <v-list class="pa-1">
        <v-list-tile avatar tag="div">
          <v-list-tile-avatar>
            <v-icon>perm_identity</v-icon>
            <!--
            <img src="https://randomuser.me/api/portraits/men/85.jpg" />
          -->
          </v-list-tile-avatar>
          <v-list-tile-content>
            <v-list-tile-title>{{ username }}</v-list-tile-title>
            <v-list-tile-title class="caption">{{ companyname }}</v-list-tile-title>
          </v-list-tile-content>
          <v-list-tile-action>
            <v-btn icon @click.stop="toggle">
              <v-icon>chevron_left</v-icon>
            </v-btn>
          </v-list-tile-action>
        </v-list-tile>
      </v-list>
      <v-list class="pt-0" dense>
        <v-divider light></v-divider>
        <v-list-tile v-for="item in items" :key="item.title" @click="navigate(item.route)">
          <v-list-tile-action>
            <v-icon>{{ item.icon }}</v-icon>
          </v-list-tile-action>
          <v-list-tile-content>
            <v-list-tile-title>{{ item.title }}</v-list-tile-title>
          </v-list-tile-content>
        </v-list-tile>
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