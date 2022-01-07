package connect

import (
	"bytes"
	"crypto/tls"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/iwvelando/tesla-energy-stats-collector/config"
	"github.com/iwvelando/tesla-energy-stats-collector/model"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

const expectedHttpStatus = 200

func Auth(config *config.Configuration) (*http.Client, time.Time, error) {

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: config.TeslaGateway.SkipVerifySsl}
	client := &http.Client{}

	data := &model.AuthPayload{
		Username:     "customer",
		Password:     config.TeslaGateway.Password,
		Email:        config.TeslaGateway.Email,
		Force_Sm_Off: false,
	}
	dataJson, err := json.Marshal(data)
	if err != nil {
		return client, time.Now(), err
	}

	req, err := http.NewRequest("POST", config.TeslaGateway.Address+"/api/login/Basic", bytes.NewBuffer(dataJson))
	if err != nil {
		return client, time.Now(), err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return client, time.Now(), err
	}
	defer resp.Body.Close()

	bodyJson := &model.AuthResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return client, time.Now(), err
	}

	status := resp.StatusCode
	if status != expectedHttpStatus {
		err = fmt.Errorf("expected %s HTTP status code but got %s; raw body %s", expectedHttpStatus, resp.StatusCode, body)
		return client, time.Now(), err
	}

	err = json.Unmarshal(body, bodyJson)
	if err != nil {
		return client, time.Now(), err
	}
	err = bodyJson.ParseTime()
	if err != nil {
		return client, time.Now(), fmt.Errorf("error when parsing authentication time, %s", err)
	}

	jar, _ := cookiejar.New(nil)
	var cookies []*http.Cookie
	cookie := &http.Cookie{
		Name:    "AuthCookie",
		Value:   bodyJson.Token,
		Path:    "/",
		Domain:  "",
		Expires: bodyJson.LoginTime.Add(24 * time.Hour),
	}
	cookies = append(cookies, cookie)
	cookie = &http.Cookie{
		Name:    "UserRecord",
		Value:   b64.StdEncoding.EncodeToString(body),
		Path:    "/",
		Domain:  "",
		Expires: bodyJson.LoginTime.Add(24 * time.Hour),
	}
	cookies = append(cookies, cookie)
	u, _ := url.Parse(config.TeslaGateway.Address)
	jar.SetCookies(u, cookies)
	client.Jar = jar

	return client, bodyJson.LoginTime.Add(23*time.Hour + 55*time.Minute), nil
}

func GetEndpoint(config *config.Configuration, client *http.Client, endpoint string, data interface{}) error {
	req, err := http.NewRequest("GET", config.TeslaGateway.Address+endpoint, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, data)
	if err != nil {
		err = fmt.Errorf("%w; raw body %s", err, body)
		return err
	}

	return nil
}

func GetAll(config *config.Configuration, client *http.Client) (model.Teg, error) {

	teg := model.Teg{}

	endpoint := "/api/meters/aggregates"
	err := GetEndpoint(config, client, endpoint, &teg.Meters)
	if err != nil {
		return teg, fmt.Errorf("error when querying %s, %s", endpoint, err)
	}
	teg.Meters.Timestamp = time.Now()
	err = teg.Meters.ParseTime()
	if err != nil {
		return teg, fmt.Errorf("error when parsing time for endpoint %s, %s", endpoint, err)
	}
	endpoint = "/api/meters/status"
	err = GetEndpoint(config, client, endpoint, &teg.MetersStatus)
	if err != nil {
		return teg, fmt.Errorf("error when querying %s, %s", endpoint, err)
	}
	teg.MetersStatus.Timestamp = time.Now()

	endpoint = "/api/operation"
	err = GetEndpoint(config, client, endpoint, &teg.Operation)
	if err != nil {
		return teg, fmt.Errorf("error when querying %s, %s", endpoint, err)
	}
	teg.Operation.Timestamp = time.Now()

	endpoint = "/api/powerwalls"
	err = GetEndpoint(config, client, endpoint, &teg.Powerwalls)
	if err != nil {
		return teg, fmt.Errorf("error when querying %s, %s", endpoint, err)
	}
	teg.Powerwalls.Timestamp = time.Now()
	err = teg.Powerwalls.ParseTime()
	if err != nil {
		return teg, fmt.Errorf("error when parsing time for endpoint %s, %s", endpoint, err)
	}

	endpoint = "/api/site_info"
	err = GetEndpoint(config, client, endpoint, &teg.SiteInfo)
	if err != nil {
		return teg, fmt.Errorf("error when querying %s, %s", endpoint, err)
	}
	teg.SiteInfo.Timestamp = time.Now()

	endpoint = "/api/sitemaster"
	err = GetEndpoint(config, client, endpoint, &teg.Sitemaster)
	if err != nil {
		return teg, fmt.Errorf("error when querying %s, %s", endpoint, err)
	}
	teg.Sitemaster.Timestamp = time.Now()

	endpoint = "/api/solar_powerwall"
	err = GetEndpoint(config, client, endpoint, &teg.SolarPowerwall)
	if err != nil {
		return teg, fmt.Errorf("error when querying %s, %s", endpoint, err)
	}
	teg.SolarPowerwall.Timestamp = time.Now()
	err = teg.SolarPowerwall.ParseTime()
	if err != nil {
		return teg, fmt.Errorf("error when parsing time for endpoint %s, %s", endpoint, err)
	}

	endpoint = "/api/solars"
	err = GetEndpoint(config, client, endpoint, &teg.Solars)
	if err != nil {
		return teg, fmt.Errorf("error when querying %s, %s", endpoint, err)
	}
	ts := time.Now()
	for i := range teg.Solars {
		teg.Solars[i].Timestamp = ts
	}

	endpoint = "/api/system/networks/conn_tests"
	err = GetEndpoint(config, client, endpoint, &teg.NetworkConnectionTests)
	if err != nil {
		return teg, fmt.Errorf("error when querying %s, %s", endpoint, err)
	}
	teg.NetworkConnectionTests.Timestamp = time.Now()
	err = teg.NetworkConnectionTests.ParseTime()
	if err != nil {
		return teg, fmt.Errorf("error when parsing time for endpoint %s, %s", endpoint, err)
	}

	endpoint = "/api/status"
	err = GetEndpoint(config, client, endpoint, &teg.Status)
	if err != nil {
		return teg, fmt.Errorf("error when querying %s, %s", endpoint, err)
	}
	teg.Status.Timestamp = time.Now()
	err = teg.Status.ParseTime()
	if err != nil {
		return teg, fmt.Errorf("error when parsing time for endpoint %s, %s", endpoint, err)
	}

	endpoint = "/api/system/testing"
	err = GetEndpoint(config, client, endpoint, &teg.SystemTesting)
	if err != nil {
		return teg, fmt.Errorf("error when querying %s, %s", endpoint, err)
	}
	teg.SystemTesting.Timestamp = time.Now()

	endpoint = "/api/system/update/status"
	err = GetEndpoint(config, client, endpoint, &teg.UpdateStatus)
	if err != nil {
		return teg, fmt.Errorf("error when querying %s, %s", endpoint, err)
	}
	teg.UpdateStatus.Timestamp = time.Now()

	endpoint = "/api/system_status"
	err = GetEndpoint(config, client, endpoint, &teg.SystemStatus)
	if err != nil {
		return teg, fmt.Errorf("error when querying %s, %s", endpoint, err)
	}
	teg.SystemStatus.Timestamp = time.Now()
	err = teg.SystemStatus.ParseTime()
	if err != nil {
		return teg, fmt.Errorf("error when parsing time for endpoint %s, %s", endpoint, err)
	}
	err = teg.SystemStatus.ParseFaults()
	if err != nil {
		return teg, fmt.Errorf("error when parsing faults for endpoint %s, %s", endpoint, err)
	}

	endpoint = "/api/system_status/grid_status"
	err = GetEndpoint(config, client, endpoint, &teg.SystemGridStatus)
	if err != nil {
		return teg, fmt.Errorf("error when querying %s, %s", endpoint, err)
	}
	teg.SystemGridStatus.Timestamp = time.Now()

	endpoint = "/api/system_status/soe"
	err = GetEndpoint(config, client, endpoint, &teg.SystemStateOfEnergy)
	if err != nil {
		return teg, fmt.Errorf("error when querying %s, %s", endpoint, err)
	}
	teg.SystemStateOfEnergy.Timestamp = time.Now()

	return teg, nil
}
