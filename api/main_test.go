package api

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
