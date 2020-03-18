package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	b64 "encoding/base64"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

func setupAuth() *jwt.GinJWTMiddleware {
	return &jwt.GinJWTMiddleware{
		Realm:         APP_NAME,
		Key:           []byte(JWT_SECRET),
		Timeout:       time.Hour,
		MaxRefresh:    time.Hour,
		Authenticator: authenticator,
		Authorizator:  authorizator,
		Unauthorized:  unauthorized,
		PayloadFunc:   payload,
		TokenLookup:   "header:Authorization",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	}
}

func authorizator(data interface{}, c *gin.Context) bool {
	// v, ok := data.(*User);
	// uid := v.ID
	claims := jwt.ExtractClaims(c)
	uid := int(claims["uid"].(float64))       // user
	cid := int(claims["companyID"].(float64)) // company
	roles := claims["roles"]
	resource, _ := strconv.Atoi(c.Param("id"))
	switch c.HandlerName() {
	case "main.getCompanyHandler":
		return resource == cid
	case "main.getUsersHandler":
		return true // ok anybody can see who is enrolled
	case "main.getUserHandler":
		return resource == uid
	case "main.editUserHandler":
		return resource == uid || canEditUser(uid, resource)
	case "main.resetApiKeyUserHandler":
		return resource == uid || canEditUser(uid, resource)
	case "main.editCompanyHandler":
		return resource == cid && hasRole(roles, ADMIN, SUPERADMIN)

	case "main.getPDFTemplatesHandler":
		return hasRole(roles, ADMIN, SUPERADMIN, EDIT_TEMPLATES) &&
			PDFTemplateInCompany(resource, cid)
	case "main.createPDFTemplateHandler":
		return hasRole(roles, ADMIN, SUPERADMIN, EDIT_TEMPLATES) &&
			PDFTemplateInCompany(resource, cid)
	case "main.deletePDFTemplateHandler":
		return hasRole(roles, ADMIN, SUPERADMIN, EDIT_TEMPLATES) &&
			PDFTemplateInCompany(resource, cid)
	case "main.editPDFTemplateHandler":
		return hasRole(roles, ADMIN, SUPERADMIN, EDIT_TEMPLATES) &&
			PDFTemplateInCompany(resource, cid)
	case "main.getPDFTemplateHandler":
		return hasRole(roles, ADMIN, SUPERADMIN, EDIT_TEMPLATES) &&
			PDFTemplateInCompany(resource, cid)

	case "main.getP12SignaturesHandler":
		return true // ok anybody can see the list
	case "main.createP12SignatureHandler":
		return hasRole(roles, ADMIN, SUPERADMIN, SIGN_DOCUMENT)
	case "main.getP12SignatureHandler":
		return hasRole(roles, ADMIN, SUPERADMIN, SIGN_DOCUMENT) &&
			OwnSignature(resource, uid)
	case "main.editP12SignatureHandler":
		return hasRole(roles, ADMIN, SUPERADMIN, SIGN_DOCUMENT) &&
			OwnSignature(resource, uid)
	case "main.deleteP12SignatureHandler":
		return hasRole(roles, ADMIN, SUPERADMIN, SIGN_DOCUMENT) &&
			OwnSignature(resource, uid)

	case "main.getMailTemplatesHandler":
		return hasRole(roles, ADMIN, SUPERADMIN, EDIT_TEMPLATES) &&
			MailTemplateInCompany(resource, cid)
	case "main.createMailTemplateHandler":
		return hasRole(roles, ADMIN, SUPERADMIN, EDIT_TEMPLATES) &&
			MailTemplateInCompany(resource, cid)
	case "main.deleteMailTemplateHandler":
		return hasRole(roles, ADMIN, SUPERADMIN, EDIT_TEMPLATES) &&
			MailTemplateInCompany(resource, cid)
	case "main.editMailTemplateHandler":
		return hasRole(roles, ADMIN, SUPERADMIN, EDIT_TEMPLATES) &&
			MailTemplateInCompany(resource, cid)
	case "main.getMailTemplateHandler":
		return hasRole(roles, ADMIN, SUPERADMIN, EDIT_TEMPLATES) &&
			MailTemplateInCompany(resource, cid)
	case "main.getLogEntriesHandler":
		return true // res read only filtered by cid or uid
	}
	printError(errors.New("ACL not granted " + c.HandlerName() + " " + strconv.Itoa(resource) + " " + strconv.Itoa(uid)))
	return false
}

func OwnSignature(resource int, uid int) bool {
	var signature P12Signature
	db.First(&signature, resource)
	if signature.ID == 0 {
		return true // new
	}
	if signature.UserID == uint(uid) {
		return true
	}
	printError(errors.New("Signature not owned " + strconv.Itoa(resource) + " " + strconv.Itoa(uid)))
	return false
}

func MailTemplateInCompany(resource int, cid int) bool {
	var template MailTemplate
	db.First(&template, resource)
	if template.ID == 0 {
		return true // new
	}
	// TODO check me
	return true
	if template.CompanyID == uint(cid) {
		return true
	}
	printError(errors.New("Mail Template not owned " + strconv.Itoa(resource) + " " + strconv.Itoa(cid)))
	return false
}

func PDFTemplateInCompany(resource int, cid int) bool {
	var template PDFTemplate
	db.First(&template, resource)
	if template.ID == 0 {
		return true // new
	}
	if template.CompanyID == uint(cid) {
		return true
	}
	printError(errors.New("Template not owned " + strconv.Itoa(resource) + " " + strconv.Itoa(cid)))
	return false
}

func hasRole(roles interface{}, needed ...Role) bool {
	for i := 0; i < len(needed); i++ {
		return strings.Contains(roles.(string), needed[i].String())
	}
	printError(errors.New("User does have one of needed role with " + roles.(string)))
	return false
}

func canEditUser(uid int, resource int) bool {
	// @TODO
	// users are in same company
	// uid has roles
	return false
}

func payload(data interface{}) jwt.MapClaims {
	var user User
	if v, ok := data.(*User); ok {
		return jwt.MapClaims{
			"id": v.ID,
			"uid": v.ID,
			"companyID": v.CompanyID,
			"roles": user.Roles,
		}
	} else {
		return jwt.MapClaims{
			"id": 0,
			"uid": 0,
			"companyID": 0,
			"roles": 0,
		}
	}
}

func authenticator(c *gin.Context) (interface{}, error) {
	// email string, password string,
	var loginVals login
	if err := c.ShouldBind(&loginVals); err != nil {
		return "", jwt.ErrMissingLoginValues
	}
	email := loginVals.Username
	password := loginVals.Password

	var user User
	db.Where("Email = ?", email).First(&user)
	if user.ID == 0 {
		return nil, jwt.ErrFailedAuthentication
	} else {
		if CheckPasswordHash(password, user.Password) {
			return user, nil
		}
		return nil, jwt.ErrFailedAuthentication
	}
}

func unauthorized(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"code":    code,
		"message": message,
	})
}

func printError(err error) {
	if err != nil {
		log.Printf("[EDDIE-ERROR] %v\n", err)
	}
}

func printErrors(errs []error) {
	for _, err := range errs {
		if err != nil {
			log.Printf("[EDDIE-ERROR] %v\n", err)
		}
	}
}

func printDebug(format string, values ...interface{}) {
	if IsDebugging() {
		log.Printf("[EDDIE-debug] "+format, values...)
	}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func optionsHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}

// define an alias type user for encrypting senssitive information
// stored in database
type Encryptable []byte

// encrypt a string example:
// text := Encryptable("My name is Joe")
// ciphertext, err := text.Encrypt()
func (plaintext Encryptable) Encrypt() (string, error) {
	c, err := aes.NewCipher(dbkey)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	//b64.URLEncoding.EncodeToString(
	return b64.StdEncoding.EncodeToString(gcm.Seal(nonce, nonce, plaintext, nil)), nil
}

// decrypt a string  example :
// 	ciphertext := Encryptable(ciphertext)
//	plaintext, err := ciphertext.Decrypt()
func (sEnc Encryptable) Decrypt() ([]byte, error) {
	ciphertext, err := b64.StdEncoding.DecodeString(string(sEnc))
	if err != nil {
		return nil, err
	}
	c, err := aes.NewCipher(dbkey)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

type LogEntry struct {
	gorm.Model
	CompanyID    uint
	UserID       uint
	Action       string
	DocumentUUID string
	Data         []byte
	DataMD5      string
}

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func getLogEntriesHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	cid := uint(claims["companyID"].(float64))
	var entries []LogEntry
	db.Where("company_id = ?", cid).Find(&entries)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": entries})
}
