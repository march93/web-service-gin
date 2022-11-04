package tests

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"web-service-gin/controllers"
	"web-service-gin/database"
	"web-service-gin/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type ReturnType struct {
	oid  string
	size int
}

func Test_UploadDocument_BadRequest(t *testing.T) {
	a := assert.New(t)
	database.InitDB()
	db := database.GetDB()

	// Create new document for the request body
	document := models.Document{
		Content: "",
	}
	reqBody, err := json.Marshal(document)
	if err != nil {
		a.Error(err)
	}

	req, w, err := setPutDocumentRouter(db, "/data/repo", bytes.NewBuffer(reqBody))
	if err != nil {
		a.Error(err)
	}
	defer closeDB(db)

	a.Equal(http.MethodPut, req.Method, "HTTP request method error")
	a.Equal(http.StatusBadRequest, w.Code, "HTTP request status code error")
	database.ClearTable()
}

func Test_UploadDocument_Success(t *testing.T) {
	a := assert.New(t)
	database.InitDB()
	db := database.GetDB()

	// Create new document for the request body
	document := models.Document{
		Content: "123",
	}
	reqBody, err := json.Marshal(document)
	if err != nil {
		a.Error(err)
	}

	req, w, err := setPutDocumentRouter(db, "/data/repo", bytes.NewBuffer(reqBody))
	if err != nil {
		a.Error(err)
	}
	defer closeDB(db)

	a.Equal(http.MethodPut, req.Method, "HTTP request method error")
	a.Equal(http.StatusCreated, w.Code, "HTTP request status code error")

	body, err := ioutil.ReadAll(w.Body)
	if err != nil {
		a.Error(err)
	}

	// Decode the real value returned
	var actual ReturnType
	if err := json.Unmarshal(body, &actual); err != nil {
		a.Error(err)
	}

	// Fetch the newly created document
	err = db.Where("oid = ?", actual.oid).First(&document).Error
	if err != nil {
		a.Error(err)
	}
	a.NotNil(document)

	expected := ReturnType{oid: actual.oid, size: 3}
	actual = ReturnType{oid: actual.oid, size: len(document.Content)}
	a.Equal(expected, actual)
	database.ClearTable()
}

func setPutDocumentRouter(db *gorm.DB, url string, body *bytes.Buffer) (*http.Request, *httptest.ResponseRecorder, error) {
	r := gin.New()
	controller := &controllers.DocumentController{DB: db}
	r.PUT("/data/:repository", controller.UploadDocument)

	req, err := http.NewRequest(http.MethodPut, url, body)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return req, w, nil
}
