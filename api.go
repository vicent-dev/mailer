package main

import (
	"encoding/json"
	"net/url"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func startApi() {
	r := gin.New()

	r.Use(jsonMiddleware)
	r.POST("/", sendEmailHandler)
	r.GET("/", func(c *gin.Context) {
		json.NewEncoder(c.Writer).Encode("pong")
	})


	if os.Getenv("ENV") == "prod" {
		r.RunTLS(":" + os.Getenv("API_PORT"), os.Getenv("CERT_PATH"),os.Getenv("KEY_PATH"))
	} else {
		r.Run(":" + os.Getenv("API_PORT"))
	}
}

func writErrorResponse(w http.ResponseWriter, err error, code int) {
	if code == 0 {
		code = http.StatusInternalServerError
	}
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}

func jsonMiddleware(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Next()
}

func recaptcha(token string) (bool, interface{}) {
	if os.Getenv("ENV") != "prod" && os.Getenv("RECAPTCHA_ON") == "true"{
		return true, nil
	}

	secret := os.Getenv("GOOGLE_RECAPTCHA_SECRET")

	apiURL := "https://www.google.com"
	resource := "/recaptcha/api/siteverify"
	data := url.Values{}
	data.Set("secret", secret)
	data.Set("response", token)

	u, _ := url.ParseRequestURI(apiURL)
	u.Path = resource
	urlStr := u.String()

	client := &http.Client{}
	r, _ := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode())) // URL-encoded payload
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, _ := client.Do(r)
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	return result["success"].(bool), result["error-codes"]
}
