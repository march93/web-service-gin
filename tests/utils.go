package tests

import (
	"testing"
	"web-service-gin/database"

	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type TestSuiteEnv struct {
	suite.Suite
	db *gorm.DB
}

// Tests are run before they start
func (suite *TestSuiteEnv) SetupSuite() {
	database.InitDB()
	suite.db = database.GetDB()
}

// Running after each test
func (suite *TestSuiteEnv) TearDownTest() {
	database.ClearTable()
}

// Running after all tests are completed
func (suite *TestSuiteEnv) TearDownSuite() {
	dbInstance, _ := suite.db.DB()
	_ = dbInstance.Close()
}

// This gets run automatically by `go test` so we call `suite.Run` inside it
func TestSuite(t *testing.T) {
	// This is what actually runs our suite
	suite.Run(t, new(TestSuiteEnv))
}
