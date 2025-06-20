package main

import (
	"crypto/tls"
	"fmt"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/tebeka/selenium"
)

const (
	seleniumURL = "http://localhost:4444/wd/hub"
	testURL     = "https://127.0.0.1/"
	apiURL      = "https://127.0.0.1/api/devices"
	username    = "admin"
	password    = "moxa"
	token       = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImFkbWluIiwiaWF0IjoxNzQ5MDE2Mjg3LCJqdGkiOiIxYTk3Zjk1MmU4ODg2ODU1ZWU2ZDQ5ZWRlMmQ2OTJjZTI5NWM1YWQ0In0.nblWedEDfQYZdKJTmKIoBBu8dlymawIPB_c2hp8piFQ"
)

func TestLoginAndCheckAPI(t *testing.T) {
	driver := createWebDriver(t)
	defer driver.Quit()

	performLogin(t, driver)

	verifyAPIResponse(t)
}

func createWebDriver(t *testing.T) selenium.WebDriver {
	// 加入 ChromeOptions
	chromeCaps := map[string]interface{}{
		"args": []string{
			"--ignore-certificate-errors", // 這是關鍵
			"--disable-web-security",      // 可選：禁用一些瀏覽器安全限制
			"--allow-insecure-localhost",  // 可選：允許 localhost 的自簽憑證
			"--headless",                  // 無頭模式（可選）
			"--no-sandbox",                // 避免沙箱錯誤（有時需用）
		},
	}

	// 設定 Chrome capabilities
	caps := selenium.Capabilities{
		"browserName":        "chrome",
		"goog:chromeOptions": chromeCaps,
	}

	driver, err := selenium.NewRemote(caps, seleniumURL)
	if err != nil {
		t.Fatalf("failed to start selenium session: %v", err)
	}
	return driver
}

func performLogin(t *testing.T, driver selenium.WebDriver) {
	if err := driver.Get(testURL); err != nil {
		t.Fatalf("failed to open login page: %v", err)
	}

	time.Sleep(2 * time.Second)

	usernameField, err := driver.FindElement(selenium.ByID, "input-userName")
	if err != nil {
		t.Fatalf("failed to find username field: %v", err)
	}
	passwordField, err := driver.FindElement(selenium.ByID, "input-password")
	if err != nil {
		t.Fatalf("failed to find password field: %v", err)
	}
	loginButton, err := driver.FindElement(selenium.ByID, "button-login")
	if err != nil {
		t.Fatalf("failed to find login button: %v", err)
	}

	usernameField.SendKeys(username)
	passwordField.SendKeys(password)
	loginButton.Click()

	time.Sleep(3 * time.Second)
}

func verifyAPIResponse(t *testing.T) {
	client := resty.New().
		SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+token).
		Get(apiURL)

	if err != nil {
		t.Fatalf("API call failed: %v", err)
	}
	if resp.StatusCode() != 200 {
		t.Fatalf("unexpected API status code: %d", resp.StatusCode())
	}

	fmt.Println(string(resp.Body()))
}
