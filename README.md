
A SaaS that takes a template, and some JSON data and produce a PDF,
Supercharger by GoLang, VueJS interface included

## Features

* Super-fast Golang SaaS
* User and ACL management
* Complete documents templating features
* PDF Template engine with header, footer, page numbering and format management
* P12 Files signature management
* Documents signatures
* RSA Database encryption
* End to end data encryption (SSL)
* ISO 27001 datacenter

## Usage:

`
go build && ENV=development ./eddie
`


## Todo

### small improvments
- api key change reset cookie store
- better relogin , jwt exp.

### roadmap
- OAuth2
- log entries
- qrcode
- add custom page size
- bind header/footer template data (current page of ...)
- mailchimp integration
- PDF P12 signing or workflow
- P12 cert. expiration
- sample docs
- DASHBOARD
- Integrated cache for fast document delivery
- edit roles
- invite user / create account
- document / mail / versioning


