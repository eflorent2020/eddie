package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/pkcs12"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type P12Signature struct {
	gorm.Model
	Name            string
	Description     string
	Company         Company
	CompanyID       uint   `gorm:"index"`
	UUID            string `gorm:"index"`
	User            User
	UserID          uint `gorm:"index"`
	P12File         []byte
	Expire          time.Time `orm:"auto_now_add;type(datetime)"`
	PinProtection   string
	SignatureImage  []byte
	DefaultLocation string
	DefaultReason   string
	DefaultContact  string
}

type TransientP12Signature struct {
	ID              uint
	Name            string
	UUID            string
	Description     string
	DefaultLocation string
	DefaultReason   string
	DefaultContact  string
	SignatureImage  []byte
	P12File         bool
}

func (s P12Signature) Transient() TransientP12Signature {
	t := TransientP12Signature{
		ID:              s.ID,
		Name:            s.Name,
		UUID:            s.UUID,
		Description:     s.Description,
		DefaultLocation: s.DefaultLocation,
		DefaultReason:   s.DefaultReason,
		SignatureImage:  s.SignatureImage,
		P12File:         len(s.P12File) > 0,
		DefaultContact:  s.DefaultContact}

	return t
}

func getP12SignatureHandler(c *gin.Context) {
	id := c.Param("id")
	var signature P12Signature
	db.First(&signature, id)
	if signature.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Not found"})
		return
	}
	_signature := signature.Transient()
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": _signature})
}

func partPatchP12Signature(signature *P12Signature, patch *P12Signature) {
	if patch.Name != "" {
		signature.Name = patch.Name
	}
	if patch.Description != "" {
		signature.Description = patch.Description
	}
	if patch.DefaultLocation != "" {
		signature.DefaultLocation = patch.DefaultLocation
	}
	if patch.DefaultReason != "" {
		signature.DefaultReason = patch.DefaultReason
	}
	if patch.DefaultContact != "" {
		signature.DefaultContact = patch.DefaultContact
	}
}

func checkP12(pfxData []byte) error {
	_, _, err := pkcs12.Decode(pfxData, "")
	if err != nil {
		if strings.HasPrefix((err.Error()), "pkcs12:") {
			return nil
		}
	}
	return err
	// switch t := err.(type) {
	// case *pkcs12.ErrIncorrectPassword:
	//	fmt.Println("ModelMissingError", t)
	//default:
	//	return err
	//}
}

func updateP12(signature P12Signature, c *gin.Context) {
	// TODO factorize form read ...
	file, _, err := c.Request.FormFile("P12File")
	defer file.Close()
	if err != nil {
		printError(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": err})
		return
	}
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		printError(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": err.Error()})
		return
	}
	err = checkP12(buf.Bytes())
	if err != nil {
		fmt.Printf("%T", err)
		printError(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": err.Error()})
		return
	}
	p12, err := Encryptable(buf.Bytes()).Encrypt()
	if err != nil {
		printError(err)
		return
	}

	signature.P12File = []byte(p12)
	if err := db.Save(&signature).Error; err != nil {
		printError(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
	}
}

func updateSignatureImage(signature P12Signature, c *gin.Context) {
	file, header, err := c.Request.FormFile("SignatureImage")
	defer file.Close()
	if err != nil {
		printError(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": err})
		return
	}
	if !strings.HasPrefix(header.Header["Content-Type"][0], "image") {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest,
			"message": "Please provide an image file not a " + header.Header["Content-Type"][0]})
		return
	}

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		printError(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": err})
		return
	}
	signature.SignatureImage = buf.Bytes()
	if err := db.Save(&signature).Error; err != nil {
		printError(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
	}
}

func handlePostMultipart(c *gin.Context) {
	id := c.Param("id")
	var signature P12Signature
	db.First(&signature, id)
	if signature.ID == 0 {
		printDebug("Handling multipart form signature, no entity found")
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Not found"})
		return
	}
	something := false
	_, err := c.FormFile("P12File") // try P12File
	if err == nil {
		something = true
		updateP12(signature, c)
	}
	_, err = c.FormFile("SignatureImage") // try Signature Image
	if err == nil {
		something = true
		updateSignatureImage(signature, c)
	}
	if something == false {
		printError(errors.New("P12Signature.handlePostMultipart: no form-data received"))
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "no form-data received"})
		return
	}
}

func editP12SignatureHandler(c *gin.Context) {
	if c.ContentType() == "application/json" {
		handlePostJSON(c)
	} else {
		// certainly multipart-form-data
		handlePostMultipart(c)
	}
}

func handlePostJSON(c *gin.Context) {
	var patch P12Signature
	var signature P12Signature
	if err := c.ShouldBindJSON(&patch); err == nil {
		db.First(&signature, patch.ID)
		if signature.ID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Edit signature with Invalid CID"})
			return
		}
		partPatchP12Signature(&signature, &patch)
		if err := db.Save(&signature).Error; err != nil {
			printError(err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		}
	} else {
		printError(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	}
}

func createP12SignatureHandler(c *gin.Context) {
	var patch P12Signature
	var signature P12Signature
	if err := c.ShouldBindJSON(&patch); err == nil {
		claims := jwt.ExtractClaims(c)
		cid := uint(claims["companyID"].(float64))
		uid := uint(claims["uid"].(float64))
		signature.CompanyID = cid
		signature.UserID = uid
		signature.UUID = NewUUID()
		signature.Name = patch.Name
		signature.DefaultContact = patch.DefaultContact
		if err := db.Create(&signature).Error; err != nil {
			printError(err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": err})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": signature})
	} else {
		printError(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	}
}

func getP12SignaturesHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	cid := uint(claims["companyID"].(float64)) // company
	var p12signatures []P12Signature
	var _p12signatures []TransientP12Signature
	db.Where("company_id = ?", cid).Find(&p12signatures)
	for _, p12signature := range p12signatures {
		_p12signatures = append(_p12signatures, p12signature.Transient())
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": _p12signatures})
}

func deleteP12SignatureHandler(c *gin.Context) {
	id := c.Param("id")
	var signature P12Signature
	db.First(&signature, id)
	if signature.ID == 0 {
		printError(errors.New("attempting to delete signature that do not exist"))
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Not found"})
		return
	}
	if err := db.Delete(&signature).Error; err != nil {
		printError(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}
}
