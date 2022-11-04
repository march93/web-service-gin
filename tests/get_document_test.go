package tests

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"web-service-gin/controllers"
	"web-service-gin/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (suite *TestSuiteEnv) Test_GetDocument_NotFound(t *testing.T) {
	req, w := setGetDocumentRouter(suite.db, "/data/repo/123")

	a := suite.Assert()
	a.Equal(http.MethodGet, req.Method, "HTTP request method error")
	a.Equal(http.StatusNotFound, w.Code, "HTTP request status code error")

	body, err := ioutil.ReadAll(w.Body)
	if err != nil {
		a.Error(err)
	}

	var actual models.Document
	if err := json.Unmarshal(body, &actual); err != nil {
		a.Error(err)
	}

	var expected models.Document
	a.Equal(expected, actual)
}

func (suite *TestSuiteEnv) Test_GetDocument_Success(t *testing.T) {
	a := suite.Assert()

	document, err := insertTestDocument(suite.db)
	if err != nil {
		a.Error(err)
	}

	req, w := setGetDocumentRouter(suite.db, "/data/TestRepo/"+document.Oid)

	a.Equal(http.MethodGet, req.Method, "HTTP request method error")
	a.Equal(http.StatusOK, w.Code, "HTTP request status code error")

	body, err := ioutil.ReadAll(w.Body)
	if err != nil {
		a.Error(err)
	}

	var actual models.Document
	if err := json.Unmarshal(body, &actual); err != nil {
		a.Error(err)
	}

	expected := document
	a.Equal(expected, actual)
}

func setGetDocumentRouter(db *gorm.DB, url string) (*http.Request, *httptest.ResponseRecorder) {
	r := gin.New()
	controller := &controllers.DocumentController{DB: db}
	r.GET("/data/:repository/:oid", controller.GetDocument)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return req, w
}

func insertTestDocument(db *gorm.DB) (models.Document, error) {
	repository := models.Repository{
		Name: "TestRepo",
	}
	repository, err := models.CreateRepository(db, &repository)
	if err != nil {
		return models.Document{}, err
	}

	document := models.Document{
		Content:        "123456",
		RepositoryName: repository.Name,
	}
	document, err = models.CreateDocument(db, &document)
	if err != nil {
		return models.Document{}, err
	}

	return document, nil
}
