<template>
  <v-card>
    <v-card-title>
    <h1>User {{user.Name}}</h1>
      </v-card-title>        
    <v-card-text>
    <p><v-icon>mail</v-icon> {{user.Email}}</p>      
    <p><v-icon>link</v-icon><code>api_key</code>{{user.ApiKey}}</p>
   <v-form v-model="valid" ref="form" lazy-validation>
    <v-text-field
      label="Name"
      v-model="user.Name"
      :rules="nameRules"
      required
    ></v-text-field>
    <h2>Roles</h2>
      <v-layout row wrap>
        <v-flex xs6 md3>
            <v-checkbox
              label="Admin"
              v-model="roleAdmin"
              disabled
              ></v-checkbox>
        </v-flex>
        <v-flex xs6 md3>
            <v-checkbox
              label="Edit templates"
              v-model="roleEditTemplates"
              disabled
              ></v-checkbox>
        </v-flex>
        <v-flex xs6 md3>
            <v-checkbox
              label="Generate documents"
              v-model="roleGenerateDocuments"
              disabled
              ></v-checkbox>      
        </v-flex>
        <v-flex xs6 md3>
            <v-checkbox
              label="Sign documents"
              v-model="roleSignDocuments"
              disabled
              ></v-checkbox>      
        </v-flex>
        <v-flex xs6 md3>
            <v-checkbox
              label="Read documents"
              v-model="roleReadDocuments"
              disabled
              ></v-checkbox>
        </v-flex>
        <v-flex xs6 md3>
            <v-checkbox
              label="Mail documents"
              v-model="roleReadDocuments"
              disabled
              ></v-checkbox>
        </v-flex>        
        <v-flex xs6 md3>      
            <v-checkbox
              label="Manage Users"
              v-model="roleManageUsers"
              disabled
              ></v-checkbox>          
        </v-flex>
        <v-flex xs6 md3>
            <v-checkbox
              label="Manage Api Keys"
              v-model="roleManageApiKeys"
              disabled
              ></v-checkbox> 
    </v-flex></v-layout>
        <v-btn class="warning" @click.stop="dialog = true">reset api key</v-btn>
    <v-btn color="primary" raised
      @click="submit"
      :disabled="!valid">
      submit
    </v-btn>
    <v-btn  raised
      @click="close">
      close
    </v-btn>    
  </v-form>
  </v-card-text>
    <v-confirm v-on:cancel="cancel" v-on:ok="resetApiKey" 
      title="Delete api key and get a new one ?" :modal="dialog">
    </v-confirm>

  </v-card>
</template>
<script>
import Persistent from '../store/persistent'
import { EventBus } from '../event-bus.js'
import Events from '../store/event-api.js'
import OKCancelDialog from '@/components/OKCancelDialog'

const API_VERSION = '/api/v1/rest/'

export default {
  name: 'User',
  data: () => ({
    dialog: false,
    company: Persistent.getters.company,
    id: null,
    roleAdmin: true,
    roleEditTemplates: true,
    roleGenerateDocuments: true,
    roleSignDocuments: true,
    roleReadDocuments: true,
    roleManageUsers: true,
    roleManageApiKeys: true,
    valid: true,
    nameRules: [
      (v) => !!v || 'Name is required',
      /* eslint-disable no-mixed-operators */
      (v) => v && v.length >= 3 || 'Name must be more than 3 characters'
    ],
    user: {}
  }),
  components: {
    'v-confirm': OKCancelDialog
  },
  created () {
    this.id = this.$route.params.id
    this.loadUser(this.id)
    console.log('navigated to ', this.id)
  },
  methods: {
    cancel () {
      this.dialog = false
    },
    close () {
      this.$router.push({name: 'HelloWorld'})
    },
    resetApiKey () {
      EventBus.$emit(Events.loadingStart)
      let url = Persistent.getters.baseUrl + API_VERSION + 'user/' + this.id + '/resetapikey'
      this.$http.post(url).then(function (res) {
        this.loadUser(this.id)
        this.dialog = false
        this.notifyError(res.status, 'Api Key reset')
      }, response => {
        this.notifyError(response.status, response.body.message)
      })
    },
    submit () {
      if (this.$refs.form.validate()) {
        this.saveUser()
      }
    },
    loadUser: function () {
      let url = Persistent.getters.baseUrl + API_VERSION + 'user/' + this.id
      this.$http.get(url).then(function (res) {
        // @todo check this response
        this.user = res.body.data
      }, response => {
        console.log(response)
        this.notifyError(response.status, response.body.message)
      })
    },
    saveUser: function () {
      let url = Persistent.getters.baseUrl + API_VERSION + 'user/' + this.id
      this.$http.post(url, this.user).then(function (res) {
        console.log(res)
        // @todo check this response
        // this.company = res.body.data
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
