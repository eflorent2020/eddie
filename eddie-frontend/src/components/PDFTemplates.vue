<template>
  <v-container fluid>
      <v-card style="margin-top: 6px;">    
        <v-card-title><h1>PDF Templates</h1></v-card-title>
        <v-card-text v-if="templates.length > 0">
          <v-layout row wrap>
            <v-flex sm12>
                <v-list>
                  <template v-for="(item, index) in templates">
                    <v-list-tile
                      @click.stop="open(item.ID)"
                      avatar
                      ripple
                      :key="item.ID">
                      <v-list-tile-content>
                          <h4>{{ index }}{{ item.Name }}</h4>                 
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
        <v-btn color="primary" raised @click="addTemplate">
      add template
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
  name: 'PDFTemplates',
  data () {
    return {
      user: Persistent.getters.user,
      company: Persistent.getters.company,
      templates: []
    }
  },
  created () {
    this.loadTemplates()
  },
  methods: {
    open: function (id) {
      this.$router.push({name: 'PDFTemplate', params: { id: id }})
    },
    addTemplate: function () {
      let count = this.templates.length + 1
      let url = Persistent.getters.baseUrl + API_VERSION + 'pdftemplate/'
      let data = {
        name: 'template ' + count
      }
      this.$http.put(url, data).then(function (res) {
        // @todo check this response
        console.log(res)
        this.loadTemplates()
      }, response => {
        console.log(response)
        this.notifyError(response.status, response.body.message)
      })
    },
    loadTemplates: function () {
      let url = Persistent.getters.baseUrl + API_VERSION + 'pdftemplates/'
      this.$http.get(url).then(function (res) {
        // @todo check this response
        if (res.body.data !== null) {
          this.templates = res.body.data
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
