<template>
  <v-layout row wrap >
    <v-flex xs12> 
      <v-form v-model="valid" ref="form" lazy-validation>
        
        <v-card>    
<v-card-text>
          <v-text-field
            type="text"
            class="headline"
            label="PDF Template"
            v-model="template.Name"
          ></v-text-field>
<v-icon>link</v-icon><code>document_id</code>{{template.UUID}}
<div id=ici></div><br>
<v-btn-toggle v-model="card">
              <v-btn color="blue-grey" raised value="showPageDefinition">
                Page Definition
              </v-btn>
              <v-btn color="blue-grey" raised value="showHeader">
                Header
              </v-btn>
              <v-btn color="blue-grey" raised value="showContent">
                Content
              </v-btn>
              <v-btn color="blue-grey" raised value="showFooter">
                Footer
              </v-btn>
              <v-btn color="blue-grey" raised value="showSigWidget">
                Sig. Widget
              </v-btn>              
              <v-btn color="blue-grey" raised value="showSampleData">
                Sample data
              </v-btn>              
</v-btn-toggle>


</v-card-text>


          <v-card-actions> 
            <v-btn color="warning" @click.stop="dialog = true">Delete</v-btn> 
            <v-btn color="primary" @click="submit" :disabled="!valid">Save</v-btn>             
<!--            <v-btn color="secondary">Log file</v-btn>   -->
            <v-btn color="green"  @click="test">Test</v-btn>
             <v-btn  @click="close">Close</v-btn>
          </v-card-actions>
        </v-card>

        <v-card style="margin-top: 12px;"  v-show="card === 'showContent'">    
          <v-card-title primary-title>
            <div class="headline">Content</div>            
          </v-card-title>
            <v-card-text>
              <v-text-field
                label="Content HTML"
                multi-line
                v-model="template.Content"
              ></v-text-field> 
            </v-card-text>
        </v-card>

        <v-card style="margin-top: 12px;"  v-show="card === 'showHeader'">    
          <v-card-title primary-title>
            <div class="headline">Header</div>            
          </v-card-title>
          <v-slide-y-transition>
            <v-card-text>
              <v-text-field
                label="Header HTML"
                multi-line
                v-model="template.Header"
              ></v-text-field> 
            </v-card-text>
          </v-slide-y-transition>
        </v-card>

        <v-card style="margin-top: 12px;"  v-show="card === 'showFooter'">    
          <v-card-title primary-title>
            <div class="headline">Footer</div>            
          </v-card-title>
          <v-slide-y-transition>
            <v-card-text>
              <v-text-field
                label="Footer HTML"
                multi-line
                v-model="template.Footer"
              ></v-text-field> 
            </v-card-text>
          </v-slide-y-transition>
        </v-card>

        <v-card style="margin-top: 12px;"  v-show="card === 'showSigWidget'">    
          <v-card-title primary-title>
            <div class="headline">Signature Widget</div>            
          </v-card-title>
            <v-card-text>
              <v-layout>
            <v-flex xs3 class="pa-3">
              <v-checkbox @click.stop="checkAddPageSig" 
                v-model="template.AddSignatureWidget"
                label="Add signature widget"
                type="checkbox"></v-checkbox>
            </v-flex>
            <v-flex xs3 >
                <v-text-field 
                v-show="template.AddSignatureWidget"
                  label="Acrofield"
                  type="text"
                  v-model="template.SigAcroField"
                ></v-text-field> 
            </v-flex>
            <v-flex xs-3 v-show="template.AddSignatureWidget" class="text-xs-right">

            </v-flex>
              <v-flex xs3 v-show="template.AddSignatureWidget"  class="text-xs-right pa-3">
                <v-btn large raised  icon  color="primary"  @click.stop="addPageSig" >
                   <v-icon>add</v-icon>
                  </v-btn> 
              </v-flex> 

              </v-layout>
            <v-layout :key="item.Page" row wrap v-for="(item,index) in template.SignaturePageDefinition"  
            v-show="template.AddSignatureWidget">
              <v-flex xs2 class="pa-3">
              <v-text-field
                label="Page"
                type="number"

                v-model="template.SignaturePageDefinition[index].Page"
              ></v-text-field> 
              </v-flex>
              <v-flex xs2  class="pa-3">
              <v-text-field
                label="X"
                type="number"
                v-model="template.SignaturePageDefinition[index].X"
              ></v-text-field> 

              </v-flex>              
              <v-flex xs2  class="pa-3">
              <v-text-field
                label="Y"
                type="number"
                v-model="template.SignaturePageDefinition[index].Y"
              ></v-text-field>                 
              </v-flex>
              <v-flex xs2  class="pa-3">
              <v-text-field
                label="W"
                type="number"
                v-model="template.SignaturePageDefinition[index].W"
              ></v-text-field>                 
              </v-flex>
              <v-flex xs2  class="pa-3">
              <v-text-field
                label="H"
                type="number"
                v-model="template.SignaturePageDefinition[index].H"
              ></v-text-field>                 
              </v-flex>   
              <v-flex xs2  class="pa-3 text-xs-right">
                <v-btn large raised  icon  color="warning"  @click.stop="deletePageSig(index)">
                   <v-icon>delete</v-icon>
                  </v-btn> 
              </v-flex>                
            </v-layout>
            <p v-show="template.AddSignatureWidget">Use Page 0 (zero) for repeat all pages or add pages using +</p>

            </v-card-text>
        </v-card>

        <v-card style="margin-top: 12px;" v-show="card === 'showPageDefinition'">    
          <v-card-title primary-title>
            <div class="text-xs-right">
            </div>
            <div class="headline">Page Definition</div>
          </v-card-title>
        
            <v-card-text>


  <v-container fill-height fluid>
    <v-layout fill-height>
      <v-flex xs5 align-end flexbox>
        <h3>Margins</h3>
        <v-container fluid>
          <v-layout >
          <v-flex xs4 align-end flexbox>
            
          </v-flex> 
          <v-flex xs4 align-end flexbox>
              <v-text-field
                label="top"
                type="number"
                v-model="template.PageDefinition.MarginTop"
              ></v-text-field>          
          </v-flex> 
          <v-flex xs4 align-end flexbox>
            
          </v-flex>        
          </v-layout>
          <v-layout >
          <v-flex xs4 align-end flexbox>
              <v-text-field
                label="left"
                type="number"
                v-model="template.PageDefinition.MarginLeft"
              ></v-text-field> 
          </v-flex> 
          <v-flex xs4 text-xs-center flexbox>
            <p></p>
            <v-icon>crop</v-icon>
          </v-flex> 
          <v-flex xs4 align-end flexbox>
            <v-text-field
                label="right"
                type="number"
                v-model="template.PageDefinition.MarginRight"
              ></v-text-field> 
          </v-flex>        
          </v-layout>
          <v-layout >
          <v-flex xs4 align-end flexbox>
            
          </v-flex> 
          <v-flex xs4 align-end flexbox>
            <v-text-field
                label="bottom"
                type="number"
                v-model="template.PageDefinition.MarginBottom"
              ></v-text-field> 
          </v-flex> 
          <v-flex xs4 align-end flexbox>
          
          </v-flex>        
          </v-layout>          
        </v-container>

      </v-flex>
            <v-flex xs1 align-end flexbox>
            </v-flex>
      <v-flex xs5 align-end flexbox>
        <h3>Format</h3>        
      <v-select
              prepend-icon="photo_size_select_small"
              v-bind:items="pageSizes"
              v-model="template.PageDefinition.PageSize"
              label="Page size"
              single-line
              bottom
      ></v-select>
      <v-select
              prepend-icon="crop_rotate"
              v-bind:items="orientations"
              v-model="template.PageDefinition.Orientation"
              label="Orientation"
              single-line
              bottom
      ></v-select>
      </v-flex>
    </v-layout>
</v-container>




            </v-card-text>        
        </v-card>

        <v-card style="margin-top: 12px;"  v-show="card === 'showSampleData'">    
          <v-card-title primary-title>
            <div class="headline">Sample data</div>            
          </v-card-title>
            <v-card-text style="min-height: 480px;">
              <code>curl -X POST 'http://{{host}}/document/{{template.UUID}}' \
                -d '{"api_key":"{{user.ApiKey}}", "email":"{{user.Email}}", "data": {{ sampleData }} }' -o output.pdf</code>

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
const DEFAULT_PAGE_DEFINITION = JSON.parse('{"PageSize": "A4","MarginBottom": 0,"MarginLeft": 0,"MarginTop": 0,"MarginRight": 0}')

export default {
  name: 'PDFTemplate',
  data: () => ({
    user: Persistent.getters.user,
    orientations: ['portrait', 'landscape'],
    pageSizes: ['A4', 'Letter'],
    dialog: false,
    card: 'showContent',
    id: null,
    company: Persistent.getters.company,
    formData: null,
    template: { PageDefinition: DEFAULT_PAGE_DEFINITION, SampleData: {}, SignaturePageDefinition: [] },
    valid: true,
    nameRules: [
      (v) => !!v || 'Name is required',
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
    deletePageSig (index) {
      this.template.SignaturePageDefinition.splice(index, 1)
      if (this.template.SignaturePageDefinition.length === 0) {
        this.template.AddSignatureWidget = 0
      }
    },
    checkAddPageSig () {
      this.template.AddSignatureWidget = !this.template.AddSignatureWidget
      if (this.template.AddSignatureWidget === true) {
        if (this.template.SignaturePageDefinition.length === 0) {
          this.addPageSig()
        }
        if (typeof this.template.SigAcroField === 'undefined') {
          this.template.SigAcroField = 'Signature'
        }
      }
    },
    addPageSig () {
      this.template.SignaturePageDefinition.push(this.nextSigPageDef())
    },
    nextSigPageDef () {
      if (this.template.SignaturePageDefinition.length === 0) {
        return { Page: '1', X: '150', Y: '50', W: '200', H: '50' }
      } else {
        let obj = JSON.parse(JSON.stringify(this.template.SignaturePageDefinition[this.template.SignaturePageDefinition.length - 1]))
        obj.Page = obj.Page + 1
        return obj
      }
    },
    cancel () {
      this.dialog = false
    },
    deleteTemplate () {
      EventBus.$emit(Events.loadingStart)
      let url = Persistent.getters.baseUrl + API_VERSION + 'pdftemplate/' + this.id
      this.$http.delete(url).then(function () {
        this.$router.push({name: 'PDFTemplates'})
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
      let url = Persistent.getters.baseUrl + API_VERSION + 'pdftemplate/' + this.id
      this.$http.get(url).then(function (res) {
        this.template = res.body.data
        this.decodeBase64()
        try {
          this.template.PageDefinition = JSON.parse(this.template.PageDefinition)
        } catch (e) {
          console.log('error getting page definition', this.template.PageDefinition, e)
          this.template.PageDefinition = DEFAULT_PAGE_DEFINITION
          this.notifyError('500', 'error getting page definition, It \'ve been reset.')
        }
        try {
          this.template.SampleData = JSON.parse(this.template.SampleData)
        } catch (e) {
          console.log('error getting page definition', this.template.SampleData, e)
          this.template.SampleData = {}
          this.notifyError('500', 'error getting sample data, It \'ve been reset.')
        }
        try {
          this.template.SignaturePageDefinition = JSON.parse(this.template.SignaturePageDefinition)
        } catch (e) {
          console.log('error getting signature page definition', this.template.SignaturePageDefinition, e)
          this.template.SignaturePageDefinition = []
          this.notifyError('500', 'error getting signature page definition data, It \'ve been reset.')
        }
        if (this.template.AddSignatureWidget === true && !this.template.SigAcroField) {
          this.template.SigAcroField = 'Signature'
        }
        if (this.template.SignaturePageDefinition === null) {
          this.template.SignaturePageDefinition = []
        }
        console.log('sig def', this.template.SampleData)
      }, response => {
        console.log(response)
        this.notifyError(response.status, response.body.message)
      })
    },
    saveTemplate: function () {
      // shallow copy
      let copyData = JSON.parse(JSON.stringify(this.template))
      copyData.PageDefinition = JSON.stringify(this.template.PageDefinition)
      const editor = this.$refs.editor.editor
      this.template.SampleData = editor.get()
      copyData.SampleData = JSON.stringify(this.template.SampleData)
      copyData.SignaturePageDefinition = JSON.stringify(this.template.SignaturePageDefinition)
      copyData = this.encodeBase64(copyData)
      let url = Persistent.getters.baseUrl + API_VERSION + 'pdftemplate/' + this.id
      this.$http.post(url, copyData).then(function () {
        EventBus.$emit(Events.loadingEnd)
        this.loadTemplate()
        // @todo check this response
        // this.company = res.body.data
      }, response => {
        console.log(response)
        this.notifyError(response.status, response.body.message)
      })
    },
    onChange (data) {
      this.template.SampleData = data
    },
    decodeBase64 () {
      let props = ['FormDefinition', 'PageDefinition', 'Header', 'Content', 'Footer', 'SampleData', 'SignaturePageDefinition']
      for (let i = 0; i < props.length; i++) {
        if (this.template[props[i]] !== null) {
          this.template[props[i]] = this.Base64Decode(this.template[props[i]])
        }
      }
    },
    encodeBase64 (copyData) {
      let props = ['FormDefinition', 'PageDefinition', 'Header', 'Content', 'Footer', 'SampleData', 'SignaturePageDefinition']
      for (let i = 0; i < props.length; i++) {
        if (copyData[props[i]] !== null) {
          copyData[props[i]] = this.Base64Encode(copyData[props[i]])
        }
      }
      return copyData
    },
    Base64Encode (str, encoding = 'utf-8') {
      var bytes = new (TextEncoder)(encoding).encode(str)
      return base64js.fromByteArray(bytes)
    },
    Base64Decode (str, encoding = 'utf-8') {
      var bytes = base64js.toByteArray(str)
      return new (TextDecoder)(encoding).decode(bytes)
    },
    notifyError: function (code, message) {
      EventBus.$emit(Events.apiError, code, message)
      EventBus.$emit(Events.loadingEnd)
    },
    close: function () {
      this.$router.push({name: 'PDFTemplates'})
    },
    test: function () {
      EventBus.$emit(Events.loadingStart)
      let data = {
        email: this.user.Email,
        api_key: this.user.ApiKey,
        data: this.template.SampleData
      }
      let url = Persistent.getters.baseUrl + '/document/' + this.template.UUID
      var me = this
      var xhr = new XMLHttpRequest()
      xhr.onreadystatechange = function () {
        if (this.readyState === 4) {
          if (this.status === 200) {
            EventBus.$emit(Events.loadingEnd)
            var blob = this.response
            var link = document.createElement('a')
            link.href = window.URL.createObjectURL(blob)
            document.getElementById('ici').appendChild(link)
            // link.innerHTML = 'OK'
            // link.target = me.template.UUID
            link.click()
            // window.URL.revokeObjectURL(link.href)
          } else {
            // resend the data to get the error (not a blob)
            me.$http.post(url, data).then(function () {
              EventBus.$emit(Events.loadingEnd)
              // window.URL.revokeObjectURL(url);
            }, response => {
              me.notifyError(response.status, response.bodyText + ' ' + data)
            })
          }
        }
      }
      xhr.open('POST', url)
      xhr.setRequestHeader('Content-Type', 'application/json; charset=utf-8')
      xhr.responseType = 'blob' // the response will be a blob and not text
      xhr.send(JSON.stringify(data))
      /*
      this.$http.post(url, data).then(function (res) {
        EventBus.$emit(Events.loadingEnd)
        var blob = new Blob([res.body], {type: 'application/pdf'})
        var link = document.createElement('a')
        link.href = window.URL.createObjectURL(blob)
        // link.download = 'myFileName.pdf'
        link.innerHTML = 'OK'
        document.getElementById('ici').appendChild(link)
        link.click()
        // window.URL.revokeObjectURL(url);
      }, response => {
        console.log(response)
        this.notifyError(response.status, response.body.message)
      })
      */
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