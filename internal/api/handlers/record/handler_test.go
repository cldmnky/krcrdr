// BEGIN: yz9d8f4g5h6j
package record

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/require"
)

func TestMount(t *testing.T) {
	// Create a new gin engine
	r := gin.New()

	// setup the validator
	fa, err := NewFakeAuthenticator()
	require.NoError(t, err)

	// Mount the API on the gin engine
	if err := Mount(r, fa); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Create a JWT token with the fake authenticator that has the write and read premissions
	wJWT, err := fa.CreateJWSWithClaims([]string{"records:w", "records:r"})
	require.NoError(t, err)

	/*
		// Create a JWT token with the fake authenticator that no permissions
		rJWT, err := fa.CreateJWSWithClaims([]string{})
		require.NoError(t, err)

		// create a new HTTP request to the /records endpoint, add bearer authenticaion
		req, err := http.NewRequest("GET", "/record", nil)
		require.NoError(t, err)
		bearer := "Bearer " + string(rJWT)
		req.Header.Add("Authorization", bearer)
		// create a new HTTP response recorder
		w := httptest.NewRecorder()
		// Dispatch the HTTP request
		r.ServeHTTP(w, req)
		// Check the HTTP response status code
		if w.Code != http.StatusUnauthorized {
			t.Errorf("unexpected status code: %v", w.Code)
		}
	*/

	// Create a new HTTP request to the /records endpoint, add bearer authenticaion
	req, err := http.NewRequest("GET", "/record", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	bearer := "Bearer " + string(wJWT)
	req.Header.Add("Authorization", bearer)

	// Create a new HTTP response recorder
	w := httptest.NewRecorder()

	// Dispatch the HTTP request
	r.ServeHTTP(w, req)

	// Check the HTTP response status code
	if w.Code != http.StatusOK {
		t.Errorf("unexpected status code: %v", w.Code)
	}
}

func TestRecordImpl_AddRecord(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	// Create a new gin context and RecordImpl instance
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	recordApi := RecordImpl{}

	// Call the AddRecord function
	recordApi.AddRecord(c)

	// Check that the response status code is 200 and the response body is "AddRecord"
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d but got %d", http.StatusOK, w.Code)
	}
	gomega.Expect(w.Body.String()).To(gomega.ContainSubstring("AddRecord"))

}

func TestRecordImpl_ListRecords(t *testing.T) {
	// Create a new gin context and RecordImpl instance
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	recordApi := RecordImpl{}

	// Call the ListRecords function
	recordApi.ListRecords(c)

	// Check that the response status code is 200 and the response body is "ListRecords"
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d but got %d", http.StatusOK, w.Code)
	}
	gomega.Expect(w.Body.String()).To(gomega.ContainSubstring("ListRecords"))
}
