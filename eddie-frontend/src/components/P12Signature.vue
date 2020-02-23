<template>
  <v-layout row wrap >
    <v-flex xs12> 
      <v-form v-model="valid" ref="form" lazy-validation>

  <v-card style="margin-top: 6px;">    
    <v-card-title primary-title>
      <div class="headline">P12 Signature - {{signature.Name}}</div>      
    </v-card-title>
    <v-card-text>      
        <v-icon>link</v-icon>
        <code>sig_id</code>
        {{signature.UUID}}
      <v-layout fluid>
        <v-flex xs11>
          <v-upload v-model="p12file" id="P12File"
          v-on:formData="uploadP12Ev"></v-upload>
        </v-flex>
        <v-flex xs1 class="text-xs-center">
          <p></p>
          <v-icon v-if="signature.P12File">details</v-icon>
          <v-icon v-if="!signature.P12File" color="red">error</v-icon>
        </v-flex>
      </v-layout>
    
      <v-layout fluid>
        <v-flex xs11>
          <v-upload v-model="imageFile" id="SignatureImage" 
          v-on:formData="uploadImageEv"></v-upload>
        </v-flex>
        <v-flex xs1 class="text-xs-center">
          <p></p>
          <v-icon v-if="signature.SignatureImage">details</v-icon>
          <v-icon v-if="!signature.SignatureImage" color="red">error</v-icon>
        </v-flex>
      </v-layout>
    </v-card-text>
    <v-card-actions>
      <v-btn color="warning" @click.stop="dialog = true">Delete</v-btn> 
      <v-btn color="primary" raised
        @click="submit"
        :disabled="!valid">
        submit
      </v-btn>
      <v-btn  @click="close">Close</v-btn>
      <v-btn  @click="test" disabled>Test</v-btn>      
    </v-card-actions>     
  </v-card>

  <v-card style="margin-top: 24px;">    
    <v-card-text>
      <v-text-field
        label="Reference"
        v-model="signature.Name"
        :rules="nameRules"
        required
      ></v-text-field>
      <v-text-field
        label="Description"
        v-model="signature.Description"
      ></v-text-field>
    </v-card-text>
  </v-card>

  <v-card style="margin-top: 24px;">    
    <v-card-text>
    <v-text-field type="password"
      label="P12 Pin protection"
      v-model="signature.PinProtection"
    ></v-text-field>
    <v-text-field
      label="Default sign location"
      v-model="signature.DefaultLocation"
    ></v-text-field>
    <v-text-field
      label="Default signing contact"
      v-model="signature.DefaultContact"
    ></v-text-field>
    <v-text-field
      label="Default signing reason"
      v-model="signature.DefaultReason"
    ></v-text-field>
    </v-card-text>  
  </v-card>

        </v-form>
      </v-flex>
    <v-confirm v-on:cancel="cancel" v-on:ok="deleteSignature" 
      title="Do you confirm ?" :modal="dialog">
    </v-confirm>
  </v-layout>
</template>
<script>
import Persistent from '../store/persistent'
import Upload from '@/components/Upload'
import { EventBus } from '../event-bus.js'
import Events from '../store/event-api.js'
import OKCancelDialog from '@/components/OKCancelDialog'

const API_VERSION = '/api/v1/rest/'

export default {
  name: 'P12Signature',
  data: () => ({
    company: Persistent.getters.company,
    id: null,
    p12file: 'click to attach a p12 file',
    imageFile: 'click to attach an image signature',
    formData: null,
    signature: {},
    valid: true,
    dialog: false,
    nameRules: [
      (v) => !!v || 'Name is required',
      /* eslint-disable no-mixed-operators */
      (v) => v && v.length >= 3 || 'Name must be more than 3 characters'
    ]
  }),
  components: {
    'v-confirm': OKCancelDialog,
    'v-upload': Upload
  },
  created () {
    this.id = this.$route.params.id
    this.loadSignature(this.id)
  },
  methods: {
    cancel () {
      this.dialog = false
    },
    test () {
    },
    close: function () {
      this.$router.push({name: 'P12Signatures'})
    },
    deleteSignature () {
      EventBus.$emit(Events.loadingStart)
      let url = Persistent.getters.baseUrl + API_VERSION + 'p12signature/' + this.id
      this.$http.delete(url).then(function (res) {
        this.$router.push({name: 'P12Signatures'})
        EventBus.$emit(Events.loadingEnd)
      }, response => {
        this.notifyError(response.status, response.body.message)
      })
    },
    uploadP12Ev (data) {
      let url = Persistent.getters.baseUrl + API_VERSION + 'p12signature/' + this.id
      data.set('ID', this.id)
      this.$http.post(url, data).then(function (res) {
        console.log(res)
      }, response => {
        console.log(response)
        this.notifyError(response.status, response.body.message)
      })
    },
    uploadImageEv (data) {
      let url = Persistent.getters.baseUrl + API_VERSION + 'p12signature/' + this.id
      data.set('ID', this.id)
      this.$http.post(url, data).then(function (res) {
        console.log(res)
      }, response => {
        console.log(response)
        this.notifyError(response.status, response.body.message)
      })
    },
    submit () {
      if (this.$refs.form.validate()) {
        this.saveSignature()
      }
    },
    loadSignature: function () {
      let url = Persistent.getters.baseUrl + API_VERSION + 'p12signature/' + this.id
      this.$http.get(url).then(function (res) {
        // @todo check this response
        this.signature = res.body.data
        if (!this.signature.P12File) {
          this.p12file = 'p12 file is absent click to attach'
        } else {
          this.p12file = 'click to set a new p12 file'
        }
        if (!this.signature.SignatureImage) {
          this.imageFile = 'signature image file is absent click to attach'
        } else {
          this.imageFile = 'click to set a new signature image file'
        }
        console.log('response', res.body.data.PinProtection)
      }, response => {
        console.log(response)
        this.notifyError(response.status, response.body.message)
      })
    },
    saveSignature: function () {
      let url = Persistent.getters.baseUrl + API_VERSION + 'p12signature/' + this.id
      // let postData = JSON.parse(JSON.stringify(this.signature))
      delete this.signature.P12File
      delete this.signature.SignatureImage
      this.$http.post(url, this.signature).then(function (res) {
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
