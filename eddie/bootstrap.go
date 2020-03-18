package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
)

var db *gorm.DB
var err error

var APP_NAME = "Codename Eddy"

// for db encryption replaced by docker secret if any
// must be 32 bytes long if any
// overwritted with env variable "DB_KEY" if any
var dbkey = []byte("8f9bfeff0cb4ffebe93bff9bfeff0cbc")

// this secret is regenerated at each boot
var JWT_SECRET = "mysupersecret"

func IsDebugging() bool {
	return true
}

func NewUUID() string {
	uuid := uuid.NewV4()
	return uuid.String()
}

func bootstrap() (*gorm.DB, error) {
	// overting db key from env at boot
	if os.Getenv("DB_KEY") != "" {
		dbkey = []byte(os.Getenv("DB_KEY"))
		if len(dbkey) != 32 {
			log.Fatal("encryption key must be 32 bytes long")
		}
	} else {
		printError(errors.New("************ security warning *************"))
		printError(errors.New("using \"in binary\" hardcoded encryption key"))
		printError(errors.New("please consider using env variable for this"))
	}

	JWT_SECRET = NewUUID()
	dbType := os.Getenv("DB_TYPE")
	if dbType == "" {
		dbType = "sqlite3"
	}
	// db, err := gorm.Open("postgres", "host=myhost user=gorm dbname=gorm sslmode=disable password=mypassword")

	dbConnection := os.Getenv("DB_CONNECTION")
	if dbConnection == "" {
		dbConnection = "./samples/gorm.db"
	}
	db, err = gorm.Open(dbType, dbConnection)
	if err != nil {
		fmt.Println(err)
		log.Fatal("cannot open db connection" + dbType + " " + dbConnection)
	}

	db.AutoMigrate(&User{})
	db.AutoMigrate(&Company{})
	db.AutoMigrate(&P12Signature{})
	db.AutoMigrate(&PDFTemplate{})
	db.AutoMigrate(&MailTemplate{})
	db.AutoMigrate(&LogEntry{})
	initACompanyAndAUser("My Tenant company", "Administrator", "admin@eddie.onthewifi.com", "aze123", "test")
	return db, err
}

func initACompanyAndAUser(companyName string, name string, email string, password string, smtpApiKey string) {
	var company Company
	db.Where("Name = ?", companyName).First(&company)
	if company.ID == 0 {
		apiKey, err := Encryptable(smtpApiKey).Encrypt()
		if err != nil {
			log.Fatal(err)
		}
		company = Company{Name: companyName,
			SmtpApiKey: string(apiKey)}
		db.NewRecord(company)
		db.Create(&company)
		printDebug("created company \"" + company.Name + "\"")
	} else {
		_, err := Encryptable(company.SmtpApiKey).Decrypt()
		if err != nil {
			log.Fatal(err)
			printError(errors.New("Unable to decode this database, did you provide the DB_KEY ?"))
		}

		printDebug("*********************************************")
		printDebug("*** db encryption test pass *****************")
		printDebug("*********************************************")
	}
	var user User
	db.Where("Email = ?", email).First(&user)
	if user.ID == 0 {
		pwd, _ := HashPassword(password)
		apiKey, err := Encryptable(NewUUID()).Encrypt()
		if err != nil {
			log.Fatal(err)
		}
		user = User{Name: name,
			Email:    email,
			Password: pwd,
			Company:  company,
			ApiKey:   string(apiKey),
			Roles:    strings.Join([]string{SUPERADMIN.String(), ADMIN.String()}, ", ")}
		db.NewRecord(user)
		db.Create(&user)
		fmt.Print(user)
		printDebug("created User \"" + user.Email + "\"")
	} else {
		printDebug("cool User \"" + user.Email + "\" exists")
	}
}
