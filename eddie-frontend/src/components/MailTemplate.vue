<template>
  <v-layout row wrap >
    <v-flex xs12> 
      <v-form v-model="valid" ref="form" lazy-validation>
        
        <v-card>    
          <v-card-title primary-title>
              <div class="headline">Mail Template - {{template.Subject}}</div>      
          </v-card-title>
<v-card-text>
<v-icon>link</v-icon><code>mail_id</code>{{template.UUID}}
</v-card-text>


          <v-card-actions> 
            <v-btn color="warning" @click.stop="dialog = true">Delete</v-btn> 
            <v-btn color="primary" @click="submit" :disabled="!valid">Save</v-btn>             
<!--            <v-btn color="secondary">Log file</v-btn>   -->

             <v-btn  @click="close">Close</v-btn>

<v-btn-toggle v-model="card">
                <v-btn raised value="showContent">
                Content
              </v-btn>
              <v-btn color raised value="showSampleData">
                Sample data
              </v-btn>              
</v-btn-toggle>
            <v-btn color="green"  @click="test">Test</v-btn>
          </v-card-actions>
        </v-card>

        <v-card style="margin-top: 12px;"  v-show="card === 'showContent'">    
          <v-card-title primary-title>
            <div class="headline">Content</div>            
          </v-card-title>
            <v-card-text>
              <v-text-field
                label="Subject"
                v-model="template.Subject"
              ></v-text-field>               
              <v-text-field
                label="Content HTML"
                multi-line
                v-model="template.Content"
              ></v-text-field> 
            </v-card-text>
        </v-card>

        <v-card style="margin-top: 12px;"  v-show="card === 'showSampleData'">    
          <v-card-title primary-title>
            <div class="headline">Sample data</div>            
          </v-card-title>
            <v-card-text style="min-height: 480px;">
              <code class="hljs-meta">
curl -X POST https://{{host}}/mail/{{template.UUID}} \
     -d '{ "api_key":"{{user.ApiKey}}", 
           "email":"{{user.Email}}", \ 
           "from" : "{{user.Name}} <{{user.Email}}>", \           
           "to" : "{{user.Name}} <{{user.Email}}>", \
           "data": {{ sampleData }}', \
           "attachements": null
           </code>
                <p></p>
              <json-editor ref="editor" :json="template.SampleData" />
            </v-card-text>
        </v-card>

      </v-form>
    </v-flex>
    <v-confirm v-on:cancel="cancel" v-on:ok="deleteTemplate" 
        title="Do you confirm ?" :modal="dialog">
    </v-confirm>
  </v-layout>
</template>
<script>
import Persistent from '../store/persistent'
import { EventBus } from '../event-bus.js'
import Events from '../store/event-api.js'
import OKCancelDialog from '@/components/OKCancelDialog'
import JSONEditor from 'vue2-jsoneditor'
import base64js from '../../node_modules/base64-js/index.js'

const API_VERSION = '/api/v1/rest/'

export default {
  name: 'MailTemplate',
  data: () => ({
    template: {},
    user: Persistent.getters.user,
    dialog: false,
    card: 'showContent',
    id: null,
    company: Persistent.getters.company,
    formData: null,
    valid: true,
    nameRules: [
      (v) => !!v || 'Subject is required',
      /* eslint-disable no-mixed-operators */
      (v) => v && v.length >= 3 || 'Name must be more than 3 characters'
    ]
  }),
  components: {
    'v-confirm': OKCancelDialog,
    'json-editor': JSONEditor
  },
  created () {
    this.id = this.$route.params.id
    this.loadTemplate(this.id)
  },
  computed: {
    sampleData () {
      return JSON.stringify(this.template.SampleData)
    },
    host () {
      return window.document.location.host
    }
  },
  methods: {
    cancel () {
      this.dialog = false
    },
    deleteTemplate () {
      EventBus.$emit(Events.loadingStart)
      let url = Persistent.getters.baseUrl + API_VERSION + 'mailtemplate/' + this.id
      this.$http.delete(url).then(function (res) {
        this.$router.push({name: 'MailTemplates'})
        EventBus.$emit(Events.loadingEnd)
      }, response => {
        this.notifyError(response.status, response.body.message)
      })
    },
    submit () {
      if (this.$refs.form.validate()) {
        EventBus.$emit(Events.loadingStart)
        this.card = 'none'
        this.saveTemplate()
        this.card = 'showContent'
      }
    },
    loadTemplate: function () {
      let url = Persistent.getters.baseUrl + API_VERSION + 'mailtemplate/' + this.id
      this.$http.get(url).then(function (res) {
        // @todo check this response
        this.template = res.body.data
        this.decodeBase64()
        try {
          this.template.SampleData = JSON.parse(this.template.SampleData)
        } catch (e) {
          console.log('error getting sample data definition', this.template.SampleData, e)
          this.template.SampleData = {}
          this.notifyError('500', 'error getting sample data, It \'ve been reset.')
        }
      }, response => {
        console.log(response)
        this.notifyError(response.status, response.body.message)
      })
    },
    saveTemplate: function () {
      let copyData = JSON.parse(JSON.stringify(this.template))
      const editor = this.$refs.editor.editor
      this.template.SampleData = editor.get()
      copyData.SampleData = JSON.stringify(this.template.SampleData)
      copyData = this.encodeBase64(copyData)
      let url = Persistent.getters.baseUrl + API_VERSION + 'mailtemplate/' + this.id
      console.log(copyData)
      this.$http.post(url, copyData).then(function (res) {
        EventBus.$emit(Events.loadingEnd)
        this.loadTemplate()
      }, response => {
        console.log(response)
        this.notifyError(response.status, response.bodyText)
      })
    },
    Base64Encode (str, encoding = 'utf-8') {
      var bytes = new (TextEncoder)(encoding).encode(str)
      return base64js.fromByteArray(bytes)
    },
    Base64Decode (str, encoding = 'utf-8') {
      var bytes = base64js.toByteArray(str)
      return new (TextDecoder)(encoding).decode(bytes)
    },
    decodeBase64 () {
      let props = ['Content', 'SampleData']
      for (let i = 0; i < props.length; i++) {
        if (this.template[props[i]] !== null) {
          this.template[props[i]] = this.Base64Decode(this.template[props[i]])
        }
      }
    },
    encodeBase64 (copyData) {
      let props = ['Content', 'SampleData']
      for (let i = 0; i < props.length; i++) {
        if (copyData[props[i]] !== null) {
          copyData[props[i]] = this.Base64Encode(copyData[props[i]])
        }
      }
      return copyData
    },
    onChange (data) {
      this.template.SampleData = data
    },
    notifyError: function (code, message) {
      EventBus.$emit(Events.apiError, code, message)
      EventBus.$emit(Events.loadingEnd)
    },
    close: function () {
      this.$router.push({name: 'MailTemplates'})
    },
    test: function () {
      this.card = 'showSampleData'
      EventBus.$emit(Events.loadingStart)
      let data = {
        email: this.user.Email,
        api_key: this.user.ApiKey,
        from: this.user.Name + ' ' + '<' + this.user.Email + '>',
        to: this.user.Name + ' ' + '<' + this.user.Email + '>',
        data: this.template.SampleData
      }
      let url = Persistent.getters.baseUrl + '/mail/' + this.template.UUID
      this.$http.post(url, data).then(function (res) {
        EventBus.$emit(Events.loadingEnd)
        console.log(res)
      }, response => {
        console.log(response)
        try {
          let message = JSON.stringify(JSON.parse(response.bodyText).message)
          this.notifyError(response.status, message)
        } catch (e) {
          this.notifyError('', response.bodyText)
        }
      })
    }
  }
}
</script>
<style>
.jsoneditor {
  min-height: 240px;
}
.jsoneditor-tree {
  min-height: 200px;
}
</style>