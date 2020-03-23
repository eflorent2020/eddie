package main

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var port string

func main() {
	port = os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	// TODO check if port is busy
	env := os.Getenv("ENV")
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	if env == "development" {
		printDebug("Disabled CORS check")
		config := cors.Config{
			AllowMethods:     []string{"GET", "POST", "PUT", "HEAD", "OPTION", "DELETE"},
			AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "authorization"},
			AllowCredentials: false,
			AllowAllOrigins:  true,
			MaxAge:           12 * time.Hour}
		r.Use(cors.New(config))
	}

	db, err = bootstrap()
	defer db.Close()

	// the jwt middleware
	authMiddleware := setupAuth()

	//router.Static("/assets", "./eddie-frontend/dist/static")
	r.StaticFS("/static", http.Dir("./eddie-frontend/dist/static"))
	r.StaticFile("/", "./eddie-frontend/dist/index.html")

	r.POST("/login", authMiddleware.LoginHandler)
	r.OPTIONS("/login", optionsHandler)
	r.GET("/refresh_token", authMiddleware.RefreshHandler)

	r.POST("/document/:id", serveDocument)
	r.GET("/document/:id/header", serveHeader)
	r.GET("/document/:id/footer", serveFooter)

	r.POST("/mail/:id", sendMailHandler)

	api := r.Group("/api/v1/rest")
	api.Use(authMiddleware.MiddlewareFunc())
	{
		api.GET("/users/", getUsersHandler)
		api.GET("/user/:id", getUserHandler)
		api.POST("/user/:id/resetapikey", resetApiKeyUserHandler)
		api.POST("/user/:id", editUserHandler)

		api.OPTIONS("/company/:id", optionsHandler)
		api.GET("/company/:id", getCompanyHandler)
		api.POST("/company/:id", editCompanyHandler)

		api.GET("/p12signatures/", getP12SignaturesHandler)
		api.OPTIONS("/p12signature/", optionsHandler)
		api.PUT("/p12signature/", createP12SignatureHandler)
		api.OPTIONS("/p12signature/:id", optionsHandler)
		api.DELETE("/p12signature/:id", deleteP12SignatureHandler)
		api.GET("/p12signature/:id", getP12SignatureHandler)
		api.POST("/p12signature/:id", editP12SignatureHandler)

		api.GET("/pdftemplates/", getPDFTemplatesHandler)
		api.OPTIONS("/pdftemplate/", optionsHandler)
		api.PUT("/pdftemplate/", createPDFTemplateHandler)
		api.OPTIONS("/pdftemplate/:id", optionsHandler)
		api.DELETE("/pdftemplate/:id", deletePDFTemplateHandler)
		api.GET("/pdftemplate/:id", getPDFTemplateHandler)
		api.POST("/pdftemplate/:id", editPDFTemplateHandler)

		api.GET("/mailtemplates/", getMailTemplatesHandler)
		api.OPTIONS("/mailtemplate/", optionsHandler)
		api.PUT("/mailtemplate/", createMailTemplateHandler)
		api.OPTIONS("/mailtemplate/:id", optionsHandler)
		api.DELETE("/mailtemplate/:id", deleteMailTemplateHandler)
		api.GET("/mailtemplate/:id", getMailTemplateHandler)
		api.POST("/mailtemplate/:id", editMailTemplateHandler)

		api.GET("/logentries/", getLogEntriesHandler)
	}

	http.ListenAndServe(":"+port, r)
}
