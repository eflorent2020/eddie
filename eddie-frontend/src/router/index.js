import Vue from 'vue'
import Router from 'vue-router'
import HelloWorld from '@/components/HelloWorld'
import Company from '@/components/Company'
import Users from '@/components/Users'
import User from '@/components/User'
import P12Signatures from '@/components/P12Signatures'
import P12Signature from '@/components/P12Signature'
import PDFTemplates from '@/components/PDFTemplates'
import PDFTemplate from '@/components/PDFTemplate'
import MailTemplates from '@/components/MailTemplates'
import MailTemplate from '@/components/MailTemplate'
import Help from '@/components/Help'
import LogEntries from '@/components/LogEntries'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      name: 'HelloWorld',
      component: HelloWorld
    },
    {
      path: '/company',
      name: 'Company',
      component: Company
    },
    {
      path: '/users',
      name: 'Users',
      component: Users
    },
    {
      path: '/user/:id',
      name: 'User',
      component: User
    },
    {
      path: '/p12signatures',
      name: 'P12Signatures',
      component: P12Signatures
    },
    {
      path: '/p12signature/:id',
      name: 'P12Signature',
      component: P12Signature
    },
    {
      path: '/pdftemplates',
      name: 'PDFTemplates',
      component: PDFTemplates
    },
    {
      path: '/pdftemplate/:id',
      name: 'PDFTemplate',
      component: PDFTemplate
    },
    {
      path: '/mailtemplates',
      name: 'MailTemplates',
      component: MailTemplates
    },
    {
      path: '/mailtemplate/:id',
      name: 'MailTemplate',
      component: MailTemplate
    },
    {
      path: '/help',
      name: 'Help',
      component: Help
    },
    {
      path: '/logentries',
      name: 'LogEntries',
      component: LogEntries
    }
  ]
})
