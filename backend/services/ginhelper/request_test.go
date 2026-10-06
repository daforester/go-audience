package ginhelper

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http/httptest"
	"strconv"
	"testing"
)

var jsonStr = "{\"string\": \"string\",\"int\": 5,\"bool\": true}"
var badJsonStr = "{string\": \"string\",\"int\": 5,\"bool\": true}"

func TestGetRequestBodyWithJSON(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/", nil)

	bodyReader := ioutil.NopCloser(bytes.NewReader([]byte(jsonStr)))
	r.Body = bodyReader

	gin.SetMode(gin.ReleaseMode)
	c, _ := gin.CreateTestContext(w)
	c.Request = r

	m, err := GetRequestBody(c)

	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	if m == nil {
		t.Fatalf("Expected map of values but got nil")
	}

	if m["bool"] == nil {
		t.Fatalf("Expected value for bool")
	}

	if m["int"] == nil {
		t.Fatalf("Expected value for int")
	}

	if m["string"] == nil {
		t.Fatalf("Expected value for string")
	}

	if m["bool"].(bool) != true {
		t.Fatalf("Expected map key bool to be \"true\" but got : %b", m["bool"])
	}

	if m["int"].(float64) != 5 {
		t.Fatalf("Expected map key int to be \"5\" but got : %d", m["float64"])
	}

	if m["string"].(string) != "string" {
		t.Fatalf("Expected map key string to be \"string\" but got : %s", m["string"])
	}

	bodyReader = ioutil.NopCloser(bytes.NewReader([]byte(badJsonStr)))
	r.Body = bodyReader

	_, err = GetRequestBody(c)

	if err == nil {
		t.Fatalf("Expected error for malformed JSON")
	}
}

func TestGetRequestBodyWithForm(t *testing.T) {
	boundary := "GoLangTestFormBoundary"
	disposition := "Content-Disposition: form-data;"
	input := "--" + boundary + "\n" + disposition + " name=\"string\"\n\nstring\n" +
		"--" + boundary + "\n" + disposition + " name=\"int\"\n\n5\n" +
		"--" + boundary + "\n" + disposition + " name=\"bool\"\n\ntrue\n" +
		"--" + boundary + "\n" + disposition + " name=\"onlyarray[]\"\n\noa1\n" +
		"--" + boundary + "\n" + disposition + " name=\"array[]\"\n\na1\n" +
		"--" + boundary + "\n" + disposition + " name=\"array[]\"\n\na2\n" +
		"--" + boundary + "\n" + disposition + " name=\"arrays\"\n\nas1\n" +
		"--" + boundary + "\n" + disposition + " name=\"arrays\"\n\nas2\n" +
		"--" + boundary + "\n" + disposition + " name=\"onlyfile\"; filename=\"testonlyfile\"\nContent-Type: application/octet-stream\n\nb2N0ZXQ=\n" +
		"--" + boundary + "\n" + disposition + " name=\"onlyfilearray[]\"; filename=\"testonlyfilearray\"\nContent-Type: application/octet-stream\n\nb2N0ZXQ=\n" +
		"--" + boundary + "\n" + disposition + " name=\"file[]\"; filename=\"testfile1\"\nContent-Type: application/octet-stream\n\nb2N0ZXQ=\n" +
		"--" + boundary + "\n" + disposition + " name=\"file[]\"; filename=\"testfile2\"\nContent-Type: application/octet-stream\n\nb2N0ZXQ=\n" +
		"--" + boundary + "\n" + disposition + " name=\"files\"; filename=\"testfiles1\"\nContent-Type: application/octet-stream\n\nb2N0ZXQ=\n" +
		"--" + boundary + "\n" + disposition + " name=\"files\"; filename=\"testfiles2\"\nContent-Type: application/octet-stream\n\nb2N0ZXQ=\n" +
		"--" + boundary + "--"

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/", nil)

	bodyReader := ioutil.NopCloser(bytes.NewReader([]byte(input)))
	r.Body = bodyReader

	gin.SetMode(gin.ReleaseMode)
	c, _ := gin.CreateTestContext(w)
	c.Request = r
	c.Request.Header.Add("Content-Length", strconv.Itoa(len(input)))
	c.Request.Header.Add("Content-Type","multipart/form-data; boundary=" + boundary)

	m, err := GetRequestBody(c)

	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	if m == nil {
		t.Fatalf("Expected map of values but got nil")
	}

	if m["bool"] == nil {
		t.Fatalf("Expected value for bool")
	}

	if m["int"] == nil {
		t.Fatalf("Expected value for int")
	}

	if m["string"] == nil {
		t.Fatalf("Expected value for string")
	}

	if *m["bool"].(*string) != "true" {
		t.Fatalf("Expected map key bool to be \"true\" but got : %s", m["bool"])
	}

	if *m["int"].(*string) != "5" {
		t.Fatalf("Expected map key int to be \"5\" but got : %s", m["float64"])
	}

	if *m["string"].(*string) != "string" {
		t.Fatalf("Expected map key string to be \"string\" but got : %s", m["string"])
	}
}

func TestGetRequestBodyWithJSONForm(t *testing.T) {
	boundary := "GoLangTestFormBoundary"
	disposition := "Content-Disposition: form-data;"
	input := "--" + boundary + "\n" + disposition + " name=\"_json\"\n\n" + jsonStr + "\n" +
		"--" + boundary + "--"

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/", nil)

	bodyReader := ioutil.NopCloser(bytes.NewReader([]byte(input)))
	r.Body = bodyReader

	gin.SetMode(gin.ReleaseMode)
	c, _ := gin.CreateTestContext(w)
	c.Request = r
	c.Request.Header.Add("Content-Length", strconv.Itoa(len(input)))
	c.Request.Header.Add("Content-Type","multipart/form-data; boundary=" + boundary)

	m, err := GetRequestBody(c)

	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	if m == nil {
		t.Fatalf("Expected map of values but got nil")
	}

	if m["bool"] == nil {
		t.Fatalf("Expected value for bool")
	}

	if m["int"] == nil {
		t.Fatalf("Expected value for int")
	}

	if m["string"] == nil {
		t.Fatalf("Expected value for string")
	}

	if m["bool"].(bool) != true {
		t.Fatalf("Expected map key bool to be \"true\" but got : %b", m["bool"])
	}

	if m["int"].(float64) != 5 {
		t.Fatalf("Expected map key int to be \"5\" but got : %d", m["float64"])
	}

	if m["string"].(string) != "string" {
		t.Fatalf("Expected map key string to be \"string\" but got : %s", m["string"])
	}
}
