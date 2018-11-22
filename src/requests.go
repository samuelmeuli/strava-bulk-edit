package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"golang.org/x/net/html"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

const baseUrl = "https://www.strava.com"
const loginUrl = baseUrl + "/login"
const sessionUrl = baseUrl + "/session"
const activitiesUrl = baseUrl + "/athlete/training_activities"

var client *http.Client
var csrfToken string

func update(email string, password string, update Activity) {
	// Set up HTTP client which stores cookies and does not allow redirects
	cookieJar, _ := cookiejar.New(nil)
	client = &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Jar: cookieJar,
	}

	// Obtain CSRF token and Strava website cookie
	csrfToken = getCsrfToken()
	logIn(email, password)

	// Fetch list of activities
	activities := getActivities()

	// Update activities with the new value
	updateActivities(activities, update)
}

func getCsrfToken() string {
	// Fetch login page
	res, err := client.Get(loginUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	// Extract CSRF token from login page
	var csrfToken string
	htmlTokens := html.NewTokenizer(res.Body)
loop:
	for {
		tt := htmlTokens.Next()
		switch tt {
		// Exit on errors or end of file if CSRF token hasn't been found yet
		case html.ErrorToken:
			log.Fatal("CSRF token not found on login page")
			// Inspect meta tags for CSRF token
		case html.SelfClosingTagToken:
			t := htmlTokens.Token()
			if t.Data == "meta" {
				attr := t.Attr
				if attr[0].Val == "csrf-token" {
					csrfToken = attr[1].Val
					break loop
				}
			}
		}
	}
	return csrfToken
}

func logIn(email string, password string) {
	// Log in to obtain Strava client cookies (stored in http.Client's cookie jar)
	res, err := client.PostForm(
		sessionUrl,
		url.Values{"email": {email}, "password": {password}, "authenticity_token": {csrfToken}},
	)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusFound {
		log.Fatal("Login failed")
	}
}

func getActivities() []Activity {
	// TODO fetch more than one page
	// TODO only fetch specified date range

	// Set up request
	req, _ := http.NewRequest("GET", activitiesUrl, nil)
	req.Header.Add("X-CSRF-Token", csrfToken)
	req.Header.Add("X-Requested-With", "XMLHttpRequest")

	// Send request
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusMovedPermanently {
		log.Fatal("Incorrect Strava login credentials")
	}
	if res.StatusCode != http.StatusOK {
		log.Fatalf("Request to %q returned status code %d", activitiesUrl, res.StatusCode)
	}

	activityList := new(ActivityList)
	json.NewDecoder(res.Body).Decode(activityList)
	return activityList.Activities
}

func updateActivities(activities []Activity, update Activity) {
	// Loop over activities and update the specified attribute
	for _, activity := range activities {
		// Set up request
		updateUrl := fmt.Sprintf("%s/%d", activitiesUrl, activity.Id)
		body, _ := json.Marshal(update)
		req, _ := http.NewRequest("PUT", updateUrl, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("X-CSRF-Token", csrfToken)
		req.Header.Add("X-Requested-With", "XMLHttpRequest")

		// Send request
		res, err := client.Do(req)
		res.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
		if res.StatusCode != http.StatusOK {
			log.Fatalf("Request to %q returned status code %d", updateUrl, res.StatusCode)
		}
	}
}
