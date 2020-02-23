package main

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	wkhtmltopdf "github.com/SebastiaanKlippert/go-wkhtmltopdf"
	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/jinzhu/gorm"

	pdfcore "github.com/unidoc/unidoc/pdf/core"
	pdfmodel "github.com/unidoc/unidoc/pdf/model"
)

type PageDefinition struct {
	PageSize     string `json:"PageSize"`
	Orientation  string `json:"Orientation"`
	MarginBottom string `json:"MarginBottom"`
	MarginLeft   string `json:"MarginLeft"`
	MarginTop    string `json:"MarginTop"`
	MarginRight  string `json:"MarginRight"`
}

type SignaturePageDefinition struct {
	Page string `json:"Page"`
	X    string `json:"X"`
	Y    string `json:"Y"`
	W    string `json:"W"`
	H    string `json:"H"`
}

type SignaturePagesDefinition struct {
	Collection []SignaturePageDefinition
}

type PDFTemplate struct {
	gorm.Model
	Name                    string
	Description             string
	Tags                    string
	CompanyID               uint   `gorm:"index"`
	UUID                    string `gorm:"index"`
	SampleData              []byte
	FormDefinition          []byte
	PageDefinition          []byte
	Header                  []byte
	Content                 []byte
	Footer                  []byte
	AddSignatureWidget      bool
	SigAcroField            string
	SignaturePageDefinition []byte
}

func partPatchPDFTemplate(template *PDFTemplate, patch *PDFTemplate) {
	if patch.Name != "" {
		template.Name = patch.Name
	}
	if patch.Description != "" {
		template.Description = patch.Description
	}
	if patch.Tags != "" {
		template.Tags = patch.Tags
	}
	if len(patch.FormDefinition) > 0 {
		template.FormDefinition = patch.FormDefinition
	}
	if len(patch.PageDefinition) > 0 {
		template.PageDefinition = patch.PageDefinition
	}
	if len(patch.Header) > 0 {
		template.Header = patch.Header
	}
	if len(patch.Content) > 0 {
		template.Content = patch.Content
	}
	if len(patch.Footer) > 0 {
		template.Footer = patch.Footer
	}
	if len(patch.SampleData) > 0 {
		template.SampleData = patch.SampleData
	}
	if len(patch.SignaturePageDefinition) > 0 {
		template.SignaturePageDefinition = patch.SignaturePageDefinition
		template.AddSignatureWidget = true
		template.SigAcroField = patch.SigAcroField
	} else {
		template.AddSignatureWidget = false
	}

}

func getPDFTemplateHandler(c *gin.Context) {
	id := c.Param("id")
	var template PDFTemplate
	db.First(&template, id)
	if template.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": template})
}

func editPDFTemplateHandler(c *gin.Context) {
	var patch PDFTemplate
	var template PDFTemplate
	if err := c.ShouldBindJSON(&patch); err == nil {
		db.First(&template, patch.ID)
		if template.ID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Edit template with invalid ID"})
			return
		}
		partPatchPDFTemplate(&template, &patch)
		if err := db.Save(&template).Error; err != nil {
			printError(err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": err})
			return
		}
	} else {
		printError(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
	}
}

func createPDFTemplateHandler(c *gin.Context) {
	var patch PDFTemplate
	var template PDFTemplate
	if err := c.ShouldBindJSON(&patch); err == nil {
		claims := jwt.ExtractClaims(c)
		cid := uint(claims["companyID"].(float64))
		template.CompanyID = cid
		template.UUID = NewUUID()
		template.Name = patch.Name
		defaultPageDefinition := `{
			"PageSize": "A4",
			"Orientation": "portrait",			
			"MarginBottom": "0",
			"MarginLeft": "0",
			"MarginTop": "0",
			"MarginRight": "0"
			}`
		template.PageDefinition = []byte(defaultPageDefinition)
		template.SampleData = []byte("{}")
		template.AddSignatureWidget = false
		if err := db.Create(&template).Error; err != nil {
			printError(err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": err})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": template})
	} else {
		printError(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
	}
}

func getPDFTemplatesHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	cid := uint(claims["companyID"].(float64))
	var pdftemplates []PDFTemplate
	db.Where("company_id = ?", cid).Find(&pdftemplates)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": pdftemplates})
}

func deletePDFTemplateHandler(c *gin.Context) {
	id := c.Param("id")
	var template PDFTemplate
	db.First(&template, id)
	if template.ID == 0 {
		// printError(Errors.error("attempting to delete template that do not exist"))
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Not found"})
		return
	}
	if err := db.Delete(&template).Error; err != nil {
		printError(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}
}

type PDFQuery struct {
	Email  string                 `json:"email" binding:"required"`
	ApiKey string                 `json:"api_key" binding:"required"`
	Data   map[string]interface{} `json:"data" binding:"required"`
}

type QueryString map[string]interface{}

func getPDFTemplateByUUID(uuid string) (*PDFTemplate, error) {
	var template PDFTemplate
	db.Where("uuid = ?", uuid).First(&template)
	if template.ID > 0 {
		return &template, nil
	} else {
		return nil, errors.New("Template not found")
	}
}

func preflightCheck(c *gin.Context) (*PDFTemplate, PDFQuery, *User, error) {
	var query PDFQuery
	if err := c.MustBindWith(&query, binding.JSON); err != nil {
		return &PDFTemplate{}, PDFQuery{}, &User{}, err
	}
	user, err := checkUserApiKey(query.Email, query.ApiKey)
	if err != nil {
		return &PDFTemplate{}, PDFQuery{}, &User{}, err
	}
	if !strings.Contains(user.Roles, "ADMIN") {
		return &PDFTemplate{}, PDFQuery{}, &User{}, errors.New("role not provisionned")
	}
	template, err := getPDFTemplateByUUID(c.Param("id"))
	if err != nil {
		return &PDFTemplate{}, PDFQuery{}, &User{}, err
	}
	if template.CompanyID != user.CompanyID {
		return &PDFTemplate{}, PDFQuery{}, &User{}, errors.New("preflight check \"belong to\" failed")
	}
	return template, query, user, nil
}

func serveHeader(c *gin.Context) {
	// check serve only / 127.1 ?
	var pdfTemplate *PDFTemplate
	pdfTemplate, err := getPDFTemplateByUUID(c.Param("id"))
	if err != nil {
		printError(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}
	t, err := template.New(pdfTemplate.Name).Parse(string(pdfTemplate.Header))
	if err != nil {
		printError(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}
	q := c.Request.URL.Query()
	serveTemplate(c, t, q)
}

func serveFooter(c *gin.Context) {
	// check serve only / 127.1 ?
	var pdfTemplate *PDFTemplate
	pdfTemplate, err := getPDFTemplateByUUID(c.Param("id"))
	if err != nil {
		printError(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}
	q := c.Request.URL.Query()
	page := q["page"][0]

	log.Println("====== Bind By Query String ======")
	fmt.Printf("query string: %v %v\n", q, page)

	t, err := template.New(pdfTemplate.Name).Parse(string(pdfTemplate.Footer))
	if err != nil {
		printError(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}
	serveTemplate(c, t, c.Request.URL.Query())
}

func serveTemplate(c *gin.Context, t *template.Template, data map[string][]string) {
	buf := new(bytes.Buffer)
	fmt.Printf("params %v\n", data)
	err = t.Execute(buf, data)
	if err != nil {
		printError(err)
		return
	}
	c.Data(200, "text/html; charset=utf-8", buf.Bytes())
}

func buildPDF(pdfTemplate *PDFTemplate, data map[string]interface{}) ([]byte, error) {

	t, err := template.New(pdfTemplate.Name).Parse(string(pdfTemplate.Content))
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	fmt.Printf("params %v\n", data)
	err = t.Execute(buf, data)
	if err != nil {
		return nil, err
	}

	// Create new PDF generator
	pdfgen, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return nil, err
	}

	p := wkhtmltopdf.NewPageReader(buf)

	if len(pdfTemplate.Header) > 0 {
		p.PageOptions.HeaderHTML.Set("http://127.0.0.1:" + port + "/document/" + pdfTemplate.UUID + "/header")
	}
	if len(pdfTemplate.Footer) > 0 {
		p.PageOptions.FooterHTML.Set("http://127.0.0.1:" + port + "/document/" + pdfTemplate.UUID + "/footer")
	}

	var pageDef PageDefinition

	err = json.Unmarshal(pdfTemplate.PageDefinition, &pageDef)
	if err != nil {
		return nil, err
	}

	pdfgen.Orientation.Set(pageDef.Orientation)
	pdfgen.PageSize.Set(pageDef.PageSize)
	pdfgen.MarginTop.Set(strToInt(pageDef.MarginTop))
	pdfgen.MarginBottom.Set(strToInt(pageDef.MarginBottom))
	pdfgen.MarginLeft.Set(strToInt(pageDef.MarginLeft))
	pdfgen.MarginRight.Set(strToInt(pageDef.MarginRight))
	pdfgen.AddPage(p)

	err = pdfgen.Create()
	if err != nil {
		return nil, err
	}
	return pdfgen.Bytes(), nil
}

func serveDocument(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			printDebug("serveDocument panic", err)
			c.JSON(http.StatusBadRequest, gin.H{"message": err})
		}
	}()
	var query PDFQuery
	var pdfTemplate *PDFTemplate
	var user *User
	if pdfTemplate, query, user, err = preflightCheck(c); err != nil {
		printError(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	content, err := buildPDF(pdfTemplate, query.Data)
	if err != nil {
		println("problem building PDF template")
		printError(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	if pdfTemplate.AddSignatureWidget == true {
		//if false {
		var jsonSigPagesDef SignaturePagesDefinition
		err = json.Unmarshal(pdfTemplate.SignaturePageDefinition, &jsonSigPagesDef)

		sigPagesDef := make([]SignaturePageDefinition, 0)
		err = json.Unmarshal(pdfTemplate.SignaturePageDefinition, &sigPagesDef)
		if err != nil {
			println("probleme decoding Signature page definition")
			printError(err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		// create two tmp file
		infile, err := ioutil.TempFile(os.TempDir(), "prefix")
		defer os.Remove(infile.Name())
		if err != nil {
			printError(err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		outfile, err := ioutil.TempFile(os.TempDir(), "prefix")
		defer os.Remove(outfile.Name())
		if err != nil {
			printError(err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		// dump templated to tmp file
		err = ioutil.WriteFile(infile.Name(), []byte(content), 0644)

		// add pdfAnnotationWidget for Signature
		content, err = addAnnotationWidget(infile.Name(), outfile.Name(), sigPagesDef, pdfTemplate.SigAcroField)
		if err != nil {
			printError(err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
	}

	// serve file
	c.Data(200, "application/pdf; charset=utf-8", content)

	// make a log entry
	d, err := json.Marshal(query.Data)
	data := []byte(d)
	checksum := md5.Sum(data)
	e := LogEntry{
		CompanyID:    pdfTemplate.CompanyID,
		UserID:       user.ID,
		Action:       "PDF",
		DocumentUUID: pdfTemplate.UUID,
		Data:         d,
		DataMD5:      string(checksum[:])}
	// save a log entry
	if err = db.Save(&e).Error; err != nil {
		printError(err)
		//c.JSON(http.StatusInternalServerError, gin.H{"SENT but, unable to save  log entry": err})
		return
	}
}

func strToInt(s string) uint {
	i, err := strconv.Atoi(s)
	if err != nil {
		printError(err)
	}
	return uint(i)
}

// Annotate pdf file.
func addAnnotationWidget(inputPath string, outputPath string, sigPagesDef []SignaturePageDefinition, fieldName string) ([]byte, error) {

	pdfWriter := pdfmodel.NewPdfWriter()

	// Read the input pdf file.
	f, err := os.Open(inputPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	pdfReader, err := pdfmodel.NewPdfReader(f)
	if err != nil {
		return nil, err
	}

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return nil, err
	}

	var pagesToSign []uint
	if strToInt(sigPagesDef[0].Page) == 0 {
		// all
		pagesToSign = make([]uint, numPages)
		for i := 0; i < numPages; i++ {
			pagesToSign[i] = uint(i)
		}
	} else {
		// some
		for _, sigDef := range sigPagesDef {
			if strToInt(sigDef.Page) > 0 {
				pagesToSign = append(pagesToSign, strToInt(sigDef.Page)-1)
			}
		}
	}
	// save the items of slice in map
	m := make(map[uint]bool)
	for i := 0; i < len(pagesToSign); i++ {
		m[pagesToSign[i]] = true
	}
	// save the items of slice in map
	m2 := make(map[uint]int)
	for i := 0; i < len(pagesToSign); i++ {
		m2[pagesToSign[i]] = i
	}

	acroForm := pdfmodel.NewPdfAcroForm()
	sigFlag := pdfcore.MakeInteger(3)
	acroForm.SigFlags = sigFlag
	dr := pdfcore.MakeString("/Helv 0 Tf 0 g")
	acroForm.DA = dr

	var acroFields []*pdfmodel.PdfField

	acroField := pdfmodel.NewPdfField()

	acroField.T = pdfcore.MakeName(fieldName)
	acroField.FT = pdfcore.MakeName("Sig")
	acroField.DA = pdfcore.MakeString("/MyriadPro-Regular 0 Tf 0 Tz 0 g")

	s := append(acroFields, acroField)
	acroForm.Fields = &s

	err = pdfWriter.SetForms(acroForm)

	for i := 0; i < numPages; i++ {
		// Read the page.
		page, err := pdfReader.GetPage(i + 1)
		if err != nil {
			return nil, err
		}

		if m[uint(i)] {
			a := pdfmodel.NewPdfAnnotationWidget()
			a.F = pdfcore.MakeInteger(132)
			x, _ := strconv.ParseFloat(sigPagesDef[m2[uint(i)]].X, 64)
			y, _ := strconv.ParseFloat(sigPagesDef[m2[uint(i)]].Y, 64)
			h, _ := strconv.ParseFloat(sigPagesDef[m2[uint(i)]].H, 64)
			w, _ := strconv.ParseFloat(sigPagesDef[m2[uint(i)]].W, 64)
			a.Rect = pdfcore.MakeArrayFromFloats([]float64{x, y, h, w})
			a.P = page.ToPdfObject()
			//a.primitive = container
			page.Annotations = append(page.Annotations, a.PdfAnnotation)

			acroField.KidsA = append(acroField.KidsA, a.PdfAnnotation)
		}

		err = pdfWriter.AddPage(page)
		if err != nil {
			return nil, err
		}

	}

	fWrite, err := os.Create(outputPath)
	if err != nil {
		return nil, err
	}

	defer fWrite.Close()

	err = pdfWriter.Write(fWrite)
	if err != nil {
		return nil, err
	}

	// read tmp file
	content, err := ioutil.ReadFile(outputPath)

	return content, nil
}
