package main

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"net/mail"
	"net/smtp"
	"net/textproto"
	"strings"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/jinzhu/gorm"
	"github.com/jordan-wright/email"
)

type MailTemplate struct {
	gorm.Model
	//Company    Company
	CompanyID uint `gorm:"index"`
	//User       User
	UserID     uint   `gorm:"index"`
	UUID       string `gorm:"index"`
	Subject    string
	Content    []byte
	SampleData []byte
}

func getMailTemplateHandler(c *gin.Context) {
	id := c.Param("id")
	var mail MailTemplate
	db.First(&mail, id)
	if mail.ID == 0 {
		printError(errors.New("Mail template not found " + id))
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Mail template not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": mail})
}

func editMailTemplateHandler(c *gin.Context) {
	var patch MailTemplate
	var mail MailTemplate
	id := c.Param("id")
	if err := c.ShouldBindJSON(&patch); err == nil {
		db.First(&mail, id)
		if mail.ID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"unable to find mail template": patch.ID})
			return
		}
		if patch.Subject != "" {
			mail.Subject = patch.Subject
		}
		if len(patch.Content) > 0 {
			mail.Content = patch.Content
		}
		if len(patch.SampleData) > 0 {
			mail.SampleData = patch.SampleData
		}
		if err := db.Save(&mail).Error; err != nil {
			printError(err)
			c.JSON(http.StatusInternalServerError, gin.H{"unable to save data": err})
			return
		}
	} else {
		printError(err)
		c.JSON(http.StatusBadRequest, gin.H{"unable to bind json": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": mail})
}

func getMailTemplatesHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	cid := uint(claims["companyID"].(float64))
	var mails []MailTemplate
	db.Where("company_id = ?", cid).Find(&mails)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": mails})
}

func deleteMailTemplateHandler(c *gin.Context) {
	id := c.Param("id")
	var mail MailTemplate
	db.First(&mail, id)
	if mail.ID == 0 {
		printError(errors.New("attempting to delete Mail Template that do not exist"))
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Not found"})
		return
	}
	if err := db.Delete(&mail).Error; err != nil {
		printError(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}
}

func createMailTemplateHandler(c *gin.Context) {
	var patch MailTemplate
	var mail MailTemplate
	if err := c.ShouldBindJSON(&patch); err == nil {
		claims := jwt.ExtractClaims(c)
		cid := uint(claims["companyID"].(float64))
		uid := uint(claims["uid"].(float64))
		mail.CompanyID = cid
		mail.UserID = uid
		mail.UUID = NewUUID()
		mail.Subject = patch.Subject
		mail.SampleData = []byte("{}")
		if err := db.Create(&mail).Error; err != nil {
			printError(err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": err})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": mail})
	} else {
		printError(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	}
}

type MailQuery struct {
	Email  string                 `json:"email" binding:"required"`
	ApiKey string                 `json:"api_key" binding:"required"`
	From   string                 `json:"from" binding:"required"`
	To     string                 `json:"to" binding:"required"`
	Data   map[string]interface{} `json:"data" binding:"required"`
}

func sendMailPreflightCheck(c *gin.Context) (*MailTemplate, MailQuery, error) {
	var query MailQuery
	if err := c.MustBindWith(&query, binding.JSON); err != nil {
		return &MailTemplate{}, MailQuery{}, err
	}
	user, err := checkUserApiKey(query.Email, query.ApiKey)
	if err != nil {
		return &MailTemplate{}, MailQuery{}, err
	}
	if !strings.Contains(user.Roles, "ADMIN") {
		return &MailTemplate{}, MailQuery{}, errors.New("role not provisionned")
	}
	template, err := getMailTemplateByUUID(c.Param("id"))
	if err != nil {
		return &MailTemplate{}, MailQuery{}, err
	}
	if template.CompanyID != user.CompanyID {
		return &MailTemplate{}, MailQuery{}, errors.New("preflight check \"belong to\" failed")
	}
	return template, query, nil
}

func getMailTemplateByUUID(uuid string) (*MailTemplate, error) {
	var template MailTemplate
	db.Where("uuid = ?", uuid).First(&template)
	if template.ID > 0 {
		return &template, nil
	} else {
		return nil, errors.New("Mail Template not found")
	}
}

func checkCompanyMailType(company Company) error {
	switch company.MailType {
	case "mailchimp":
		if company.SmtpApiKey == "" {
			return errors.New("Please setup your company Mailchimp Api Key")
		}
	case "smtp":
		if company.SmtpHost == "" {
			return errors.New("Please setup your company Smtp Host")
		}
		if company.SmtpPort == "" {
			return errors.New("Please setup your company Smtp Port")
		}
		if company.SmtpUsername == "" {
			return errors.New("Please setup your company Smtp Username")
		}
		if company.SmtpPassword == "" {
			return errors.New("Please setup your company Smtp Password")
		}
	default:
		return errors.New("Please setup your company SMTP settings")
	}
	return nil
}

func sendMailHandler(c *gin.Context) {
	var mailTemplate *MailTemplate
	var query MailQuery
	if mailTemplate, query, err = sendMailPreflightCheck(c); err != nil {
		printError(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}
	query.To = "youcall.co@gmail.com"
	query.From = "youcall.co@gmail.com"
	var company Company
	db.First(&company, mailTemplate.CompanyID)
	if company.ID == 0 {
		printError(errors.New("Company not found "))
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Company not found"})
		return
	}
	var user User
	db.First(&user, mailTemplate.UserID)
	if user.ID == 0 {
		printError(errors.New("User not found "))
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "User not found"})
		return
	}

	if company.ID != mailTemplate.ID {
		printError(errors.New("Company check ????? "))
	}
	err := checkCompanyMailType(company)
	if err != nil {
		printError(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}
	_, err = mail.ParseAddressList(query.To)
	if err != nil {
		fmt.Printf("parsin address list %v \n", query.To)
		printError(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}
	subject, err := buildSubject(mailTemplate, query.Data)
	if err != nil {
		printError(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}
	htmlContent, err := buildHtml(mailTemplate, query.Data)
	if err != nil {
		printError(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}
	// func (e *Email) Attach(r io.Reader, filename string, c string) (a *Attachment, err error)
	e := &email.Email{
		To:      strings.Split(query.To, ","),
		From:    query.From,
		Subject: string(subject),
		//Text:    []byte("Text Body"),
		HTML:    htmlContent,
		Headers: textproto.MIMEHeader{},
	}

	switch company.MailType {
	case "mailchimp":
		// mime, err := e.Bytes()
		c.JSON(http.StatusNotFound, gin.H{"message": "mailchimp integration under construction"})
		return
	case "smtp":
		err = sendMail(company, e)
		if err != nil {
			printError(err)
			c.JSON(http.StatusBadRequest, gin.H{"message": err})
			return
		} else {
			d, err := json.Marshal(query.Data)
			data := []byte(d)
			checksum := md5.Sum(data)
			e := LogEntry{
				CompanyID:    company.ID,
				UserID:       user.ID,
				Action:       "MAIL",
				DocumentUUID: mailTemplate.UUID,
				Data:         d,
				DataMD5:      string(checksum[:])}
			if err = db.Save(&e).Error; err != nil {
				printError(err)
				c.JSON(http.StatusInternalServerError, gin.H{"SENT but, unable to save  log entry": err})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": "smtp mail sent"})
			return
		}
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"message": "smtp setting not set."})
		return
	}

}

func buildSubject(mailTemplate *MailTemplate, data map[string]interface{}) ([]byte, error) {
	t, err := template.New(mailTemplate.Subject).Parse(string(mailTemplate.Subject))
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	err = t.Execute(buf, data)
	return buf.Bytes(), nil
}

func buildHtml(mailTemplate *MailTemplate, data map[string]interface{}) ([]byte, error) {
	t, err := template.New(mailTemplate.Subject).Parse(string(mailTemplate.Content))
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	err = t.Execute(buf, data)
	return buf.Bytes(), nil
}

func sendMail(company Company, e *email.Email) error {
	decoded, err := Encryptable(company.SmtpPassword).Decrypt()
	if err != nil {
		return err
	}
	// Set up authentication information.
	// hostname is used by PlainAuth to validate the TLS certificate.
	auth := smtp.PlainAuth(
		"",                   // identity
		company.SmtpUsername, // username
		string(decoded),      // password
		company.SmtpHost,     // hostname
	)
	err = e.Send(company.SmtpHost+":"+company.SmtpPort, auth)
	return err
}
