package main

import (
	"bou.ke/monkey"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestPingRoute(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	resBody := gin.H{}
	err := json.Unmarshal(w.Body.Bytes(), &resBody)
	assert.Equal(t, err, nil)
	assert.Equal(t, gin.H{"message": "pong"}, resBody)
}

func TestGetPath(t *testing.T) {
	router := setupRouter()

	reqParam := "test"

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/path/"+reqParam, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	resBody := w.Body.String()
	assert.Equal(t, fmt.Sprintf("param: %s", reqParam), resBody)
}

func TestGetURI(t *testing.T) {
	router := setupRouter()

	reqParentId, reqChildId := "25", "3"
	reqPath := fmt.Sprintf("/parent/%s/children/%s", reqParentId, reqChildId)
	expBody := gin.H{
		"ParentID": reqParentId,
		"ChildID":  reqChildId,
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", reqPath, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	resBody := gin.H{}
	err := json.Unmarshal(w.Body.Bytes(), &resBody)
	assert.Equal(t, err, nil)
	assert.Equal(t, expBody, resBody)
}

func TestPostJson(t *testing.T) {
	router := setupRouter()

	reqBody := gin.H{
		"name": "test.name",
		"desc": "test.desc",
	}
	reqBytes, err := json.Marshal(&reqBody)
	if err != nil {
		assert.Error(t, err)
	}
	reqBuf := bytes.NewBuffer(reqBytes)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/body", reqBuf)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	resBody := gin.H{}
	err = json.Unmarshal(w.Body.Bytes(), &resBody)
	assert.Equal(t, err, nil)
	assert.Equal(t, reqBody, resBody)
}

func TestPostJsonMissingName(t *testing.T) {
	router := setupRouter()

	reqBody := gin.H{
		"desc": "test.desc",
	}
	reqBytes, err := json.Marshal(&reqBody)
	if err != nil {
		assert.Error(t, err)
	}
	reqBuf := bytes.NewBuffer(reqBytes)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/body", reqBuf)
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
}

func TestGetQueryTest(t *testing.T) {
	router := setupRouter()

	reqTestQ := "valid"

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/query?test="+reqTestQ, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	resBody := gin.H{}
	err := json.Unmarshal(w.Body.Bytes(), &resBody)
	assert.Equal(t, err, nil)
	assert.Equal(t, reqTestQ, resBody["Test"])
}

func TestGetQueryDefaultTimes(t *testing.T) {
	now := time.Now()
	patch := monkey.Patch(time.Now, func() time.Time { return now })
	defer patch.Unpatch()

	router := setupRouter()

	reqTestQ := "valid"
	reqQuery := fmt.Sprintf("?test=%s", reqTestQ)
	expBody := gin.H{
		"Test":     "valid",
		"Earliest": time.Time{}.UTC().Format(time.RFC3339),
		"Latest":   time.Now().UTC().Format(time.RFC3339),
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/query"+reqQuery, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	resBody := gin.H{}
	err := json.Unmarshal(w.Body.Bytes(), &resBody)
	assert.Equal(t, err, nil)
	assert.Equal(t, expBody, resBody)
}

func TestGetQueryLatestGTEarliest(t *testing.T) {
	router := setupRouter()

	now := time.Now().UTC()
	test, earliest, latest := "valid", now.Format(time.RFC3339), now.Add(-time.Hour).Format(time.RFC3339)
	reqQuery := fmt.Sprintf("?test=%s&earliest=%s&latest=%s", test, earliest, latest)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/query"+reqQuery, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
}

func TestGetQueryLatestGTNow(t *testing.T) {
	now := time.Now().UTC()
	patch := monkey.Patch(time.Now, func() time.Time { return now })
	defer patch.Unpatch()

	router := setupRouter()

	test, latest := "valid", now.Add(time.Hour).Format(time.RFC3339)
	reqQuery := fmt.Sprintf("?test=%s&latest=%s", test, latest)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/query"+reqQuery, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
}

func TestGetQueryLatestEqNow(t *testing.T) {
	now := time.Now().UTC()
	patch := monkey.Patch(time.Now, func() time.Time { return now })
	defer patch.Unpatch()

	router := setupRouter()

	test, latest := "valid", now.Format(time.RFC3339)
	reqQuery := fmt.Sprintf("?test=%s&latest=%s", test, latest)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/query"+reqQuery, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestPutAll(t *testing.T) {
	router := setupRouter()
	path, query, body := "path", "query", "body"
	reqUri := fmt.Sprintf("/all/%s?q=%s", path, query)
	reqBody := gin.H{"data": body}
	reqBytes, err := json.Marshal(&reqBody)
	if err != nil {
		assert.Error(t, err)
	}
	reqBuf := bytes.NewBuffer(reqBytes)
	expBody := reqBody
	expBody["Path"] = path
	expBody["Query"] = query

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", reqUri, reqBuf)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	resBody := gin.H{}
	err = json.Unmarshal(w.Body.Bytes(), &resBody)
	assert.Equal(t, err, nil)
	assert.Equal(t, expBody, resBody)
}
