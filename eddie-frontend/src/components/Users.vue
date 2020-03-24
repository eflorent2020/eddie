<template>
  <v-container fluid>
      <v-card>
        <v-card-title><h1>Users</h1></v-card-title>        
        <v-card-text>
          <v-layout row wrap>
            <v-flex sm12>
                <v-list>
                  <template v-for="(item, index) in users">
                    <v-list-tile
                      @click.stop="open(item.ID)"
                      avatar
                      ripple
                      :key="item.Email">
                      <v-list-tile-content>
                          <h4>{{index}} {{ item.Name }}</h4>                 
                      </v-list-tile-content>
                      <v-list-tile-action>
                        <v-list-tile-action-text>
                          <v-icon>details</v-icon>
                        </v-list-tile-action-text>
                      </v-list-tile-action>                
                    </v-list-tile>              
                    <v-divider :key="'divider-' + item.ID"></v-divider>    
                  </template>
                </v-list>
            </v-flex>      
          </v-layout>
      </v-card-text>
      <v-card-actions>
    <v-btn color="primary" raised
      disabled>
      invite user
    </v-btn>
      </v-card-actions>
      </v-card>

  </v-container>
</template>

<script>
import Persistent from '../store/persistent'
import { EventBus } from '../event-bus.js'
import Events from '../store/event-api.js'

const API_VERSION = '/api/v1/rest/'

export default {
  name: 'Users',
  data () {
    return {
      users: []
    }
  },
  created () {
    this.loadUsers()
  },
  methods: {
    open: function (id) {
      console.log('navigating to ', id)
      this.$router.push({name: 'User', params: { id: id }})
    },
    loadUsers: function () {
      let url = Persistent.getters.baseUrl + API_VERSION + 'users/'
      this.$http.get(url).then(function (res) {
        // @todo check this response
        this.users = res.body.data
      }, response => {
        console.log(response)
        this.notifyError(response.status, response.body.message)
      })
    },
    notifyError: function (code, message) {
      EventBus.$emit(Events.apiError, code, message)
      EventBus.$emit(Events.loadingEnd)
    }
  }
}
</script>
