<template>
  <v-container fluid>
      <v-card style="margin-top: 6px;">    
        <v-card-title><h1>Signatures</h1></v-card-title>
        <v-card-text v-if="signatures.length > 0">
          <v-layout row wrap>
            <v-flex sm12>
                <v-list>
                  <template v-for="(item, index) in signatures">
                    <v-list-tile
                      @click.stop="open(item.ID)"
                      avatar
                      ripple
                      :key="item.ID">
                      <v-list-tile-content>
                          <h4>{{ item.Name }}</h4>                 
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
            <v-btn color="primary" raised @click="addSignature">
      add signature
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
  name: 'P12Signatures',
  data () {
    return {
      user: Persistent.getters.user,
      company: Persistent.getters.company,
      signatures: []
    }
  },
  created () {
    this.loadSignatures()
  },
  methods: {
    open: function (id) {
      this.$router.push({name: 'P12Signature', params: { id: id }})
    },
    addSignature: function () {
      let count = this.signatures.length + 1
      let url = Persistent.getters.baseUrl + API_VERSION + 'p12signature/'
      let contact = this.user.Name + '<' + this.user.Email + '>'
      let data = {
        name: this.company.Name + ' - ' + this.user.Name + ' ' + count,
        DefaultContact: contact
      }
      this.$http.put(url, data).then(function (res) {
        // @todo check this response
        // this.signatures = res.body.data
        this.loadSignatures()
      }, response => {
        console.log(response)
        this.notifyError(response.status, response.body.message)
      })
    },
    loadSignatures: function () {
      let url = Persistent.getters.baseUrl + API_VERSION + 'p12signatures/'
      this.$http.get(url).then(function (res) {
        // @todo check this response
        if (res.body.data !== null) {
          this.signatures = res.body.data
        }
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
