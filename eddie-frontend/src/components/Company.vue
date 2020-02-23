<template>
  <v-card>
    <v-card-title>
    <h1>Company {{company.Name}}</h1>
      </v-card-title>        
    <v-card-text>
   <v-form v-model="valid" ref="form" lazy-validation autocomplete="off">
    <v-text-field
      label="Name"
      v-model="company.Name"
      :rules="nameRules"
      required
    ></v-text-field>
  
      <v-btn-toggle v-model="company.MailType">
              <v-btn color="blue-grey" raised value="none">
                None
              </v-btn>
              <v-btn color="blue-grey" raised value="smtp">
                Smtp
              </v-btn>
              <v-btn color="blue-grey" raised value="mailchimp">
                MailChimp
              </v-btn>           
    </v-btn-toggle>
<br><br>
        <div v-show="company.MailType !== 'none'">
        <v-layout v-show="company.MailType === 'smtp'"  row wrap>
        <v-flex xs5 >
          <v-text-field pa-4
            label="SMTP Host"
            v-model="company.SmtpHost"
          ></v-text-field>
        </v-flex>
        <v-flex xs1>
        </v-flex>
        <v-flex xs2>
          <v-text-field
          numeric
            label="SMTP Port"
            v-model="company.SmtpPort"
          ></v-text-field>
        </v-flex>
        </v-layout>
        <v-layout v-show="company.MailType === 'smtp'"  row wrap>        
        <v-flex xs5>
          <v-text-field
            type="text"
            label="SMTP Username"
            v-model="company.SmtpUsername"
          ></v-text-field>
        </v-flex>   
        <v-flex xs1>
        </v-flex>       
        <v-flex xs5>
          <v-text-field
          autocomplete="new-password"
            type="password"
            :label="SMTPPasswordLabel"
            v-model="company.SmtpPassword"
          ></v-text-field>
        </v-flex>
        </v-layout>


        <v-layout v-show="company.MailType === 'mailchimp'" row wrap>
        <v-flex xs10>
          <v-text-field
            disabled
            :label="smtpApiKeyLabel"
            v-model="company.SmtpApiKey"
          ></v-text-field>
        </v-flex>
        </v-layout>
        </div>
  </v-form>
    </v-card-text>
    <v-card-actions>
    <v-btn color="primary"
      @click="submit"
      :disabled="!valid"
    >
      submit
    </v-btn>
    <v-btn  @click="close">Close</v-btn>
    </v-card-actions>
  </v-card>
</template>

<script>
import Persistent from '../store/persistent'
import { EventBus } from '../event-bus.js'
import Events from '../store/event-api.js'

const API_VERSION = '/api/v1/rest/'

export default {
  name: 'Company',
  data: () => ({
    smtpApiKeyLabel: 'Not set. Please register your Api Key.',
    SMTPPasswordLabel: 'Please set your SMTP password.',
    valid: true,
    nameRules: [
      (v) => !!v || 'Name is required',
      /* eslint-disable no-mixed-operators */
      (v) => v && v.length >= 3 || 'Name must be more than 3 characters'
    ],
    company: {}
  }),
  created () {
    this.loadCompany()
  },
  methods: {
    close () {
      this.$router.push({name: 'HelloWorld'})
    },
    submit () {
      if (this.$refs.form.validate()) {
        this.saveCompany()
      }
    },
    loadCompany: function () {
      let url = Persistent.getters.baseUrl + API_VERSION + 'company/' + Persistent.getters.company.ID
      this.$http.get(url).then(function (res) {
        // @todo check this response
        this.company = res.body.data
        if (this.company.SmtpApiKey) {
          this.smtpApiKeyLabel = 'Your Api Key is registerd'
        }
        this.company.SmtpApiKey = null
        if (this.company.SmtpPassword) {
          this.SMTPPasswordLabel = 'Your password is registerd'
        }
        this.company.SmtpPassword = null
      }, response => {
        console.log(response)
        this.notifyError(response.status, response.body.message)
      })
    },
    saveCompany: function () {
      let url = Persistent.getters.baseUrl + API_VERSION + 'company/' + Persistent.getters.company.ID
      this.$http.post(url, this.company).then(function (res) {
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
