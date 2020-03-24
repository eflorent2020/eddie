<template>
  <v-card>
    <v-card-title>
    <h1>Log entries</h1>
      </v-card-title>        
    <v-card-text>
        <table id="datatable"></table>
    </v-card-text>
  </v-card>
</template>

<script>
import { EventBus } from '../event-bus.js'
import Events from '../store/event-api.js'
import Persistent from '../store/persistent'

const API_VERSION = '/api/v1/rest/'

export default {
  name: 'LogEntries',
  data () {
    return {
      entries: [],
      headers: [
        {title: 'Created At', data: 'CreatedAt'},
        {title: 'Action', data: 'Action'},
        {title: 'Document', data: 'DocumentUUID'}
      ],
      dtHandle: null
    }
  },
  created () {
    this.loadEntries()
  },
  methods: {
    close: function () {
      this.$router.push({name: 'HelloWord'})
    },
    prepareTable: function () {
      /* eslint-disable no-undef */
      this.dtHandle = $('#datatable').DataTable({
        columns: this.headers,
        data: this.entries,
        searching: false,
        paging: true,
        info: false
      })
    },
    loadEntries: function () {
      let url = Persistent.getters.baseUrl + API_VERSION + 'logentries/'
      this.$http.get(url).then(function (res) {
        // @todo check this response
        if (res.body.data !== null) {
          this.entries = res.body.data
          console.log(this.entries)
          this.prepareTable()
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
<style>
#datatable {
  width: 100%;
}
</style>
