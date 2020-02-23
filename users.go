package main

import (
	"errors"
	"net/http"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type Role int

const (
	SUPERADMIN Role = 1 + iota
	ADMIN
	EDIT_TEMPLATES
	GENERATE_DOCUMENTS
	SIGN_DOCUMENT
	READ_DOCUMENTS
	MANAGE_USERS
	MANAGE_API_KEYS
)

var roles = [...]string{
	"SUPERADMIN",
	"ADMIN",
	"EDIT_TEMPLATES",
	"GENERATE_DOCUMENTS",
	"SIGN_DOCUMENT",
	"READ_DOCUMENTS",
	"MANAGE_USERS",
	"MANAGE_API_KEYS"}

func (r Role) String() string { return roles[r-1] }

type SmtpType int

const (
	STANDARD SmtpType = 1 + iota
	MAILCHIMP
)

var SmtpTypes = [...]string{
	"with stmp",
	"with mailchimp"}

type Company struct {
	gorm.Model
	Name         string
	MailType     string
	SmtpHost     string
	SmtpPort     string
	SmtpUsername string
	SmtpPassword string
	SmtpApiKey   string
}

type TransientCompany struct {
	gorm.Model
	ID           uint
	Name         string
	MailType     string
	SmtpHost     string
	SmtpUsername string
	SmtpPort     string
	SmtpApiKey   bool
	SmtpPassword bool
}

func (c Company) Transient() TransientCompany {
	t := TransientCompany{
		ID:           c.ID,
		Name:         c.Name,
		MailType:     c.MailType,
		SmtpHost:     c.SmtpHost,
		SmtpUsername: c.SmtpUsername,
		SmtpPort:     c.SmtpPort,
		SmtpApiKey:   c.SmtpApiKey != "",
		SmtpPassword: c.SmtpPassword != ""}
	return t
}

type User struct {
	gorm.Model
	Name      string
	Email     string `gorm:"unique_index"`
	Password  string
	Company   Company
	CompanyID uint `gorm:"index"`
	ApiKey    string
	Roles     string
}

type TransientUser struct {
	ID      uint
	Name    string
	Email   string
	Company Company
	ApiKey  string
	Roles   string
}

func (u User) Transient() TransientUser {
	apiKey, err := Encryptable(u.ApiKey).Decrypt()
	if err != nil {
		printError(err)
	}
	t := TransientUser{
		ID:      u.ID,
		Name:    u.Name,
		Email:   u.Email,
		Company: u.Company,
		ApiKey:  string(apiKey),
		Roles:   u.Roles}
	return t
}

func getUsersHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	cid := uint(claims["companyID"].(float64))
	var users []User
	var _users []TransientUser
	db.Where("company_id = ?", cid).Find(&users)
	for _, user := range users {
		_users = append(_users, user.Transient())
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": _users})
}

func getUserHandler(c *gin.Context) {
	id := c.Param("id")
	var user User
	db.First(&user, id)
	if user.ID == 0 {
		printError(errors.New("getting a user unknow " + id))
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "User found"})
		return
	}
	_user := user.Transient()
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": _user})
}

func editUserHandler(c *gin.Context) {
	var patch User
	var user User
	if err := c.ShouldBindJSON(&patch); err == nil {
		db.First(&user, patch.ID)
		if user.ID == 0 {
			printError(errors.New("User id not found at patching "))
			c.JSON(http.StatusBadRequest, gin.H{"error": patch.ID})
			return
		}
		if patch.Name != "" {
			user.Name = patch.Name
		}
		if patch.Roles != "" {
			user.Roles = patch.Roles
		}
		if err := db.Save(&user).Error; err != nil {
			printError(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
		}
	} else {
		printError(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func getCompanyHandler(c *gin.Context) {
	id := c.Param("id")
	var company Company
	db.First(&company, id)
	if company.ID == 0 {
		printError(errors.New("Company not found " + id))
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Company not found"})
		return
	}
	_company := company.Transient()
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": _company})
}

func editCompanyHandler(c *gin.Context) {
	var patch Company
	var company Company
	if err := c.ShouldBindJSON(&patch); err == nil {
		db.First(&company, patch.ID)
		if company.ID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"unable to find company": patch.ID})
			return
		}
		if patch.Name != "" {
			company.Name = patch.Name
		}
		company.MailType = patch.MailType
		switch patch.MailType {
		case "mailchimp":
			if patch.SmtpApiKey != "" {
				apiKey, err := Encryptable(patch.SmtpApiKey).Encrypt()
				if err != nil {
					printError(err)
				}
				company.SmtpApiKey = string(apiKey)
			}
		case "smtp":
			if patch.SmtpHost != "" {
				company.SmtpHost = patch.SmtpHost
				company.SmtpPort = patch.SmtpPort
				company.SmtpUsername = patch.SmtpUsername
				if patch.SmtpPassword != "" {
					password, err := Encryptable(patch.SmtpPassword).Encrypt()
					if err != nil {
						printError(err)
						c.JSON(http.StatusInternalServerError, gin.H{"unable to save data": err})
						return
					}
					company.SmtpPassword = string(password)
				}
			}
		}
		if err := db.Save(&company).Error; err != nil {
			printError(err)
			c.JSON(http.StatusInternalServerError, gin.H{"unable to save data": err})
		}
	} else {
		printError(err)
		c.JSON(http.StatusBadRequest, gin.H{"unable to bind json": err.Error()})
	}
}

func checkUserApiKey(username string, apiKey string) (*User, error) {
	var user User
	db.Where("email = ?", username).First(&user)

	if user.ID == 0 {
		return nil, errors.New("User not found")
	}
	decoded, err := Encryptable(user.ApiKey).Decrypt()
	if err != nil {
		return nil, err
	}
	if string(decoded) != apiKey {
		return nil, errors.New("Api key doen't match")
	}
	return &user, nil
}

func resetApiKeyUserHandler(c *gin.Context) {
	id := c.Param("id")
	var user User
	db.First(&user, id)
	if user.ID == 0 {
		printError(errors.New("User not found " + id))
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "User not found"})
		return
	}
	apiKey, err := Encryptable(NewUUID()).Encrypt()
	if err != nil {
		printError(err)
	}
	user.ApiKey = apiKey

	if err := db.Save(&user).Error; err != nil {
		printError(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": user})
}
