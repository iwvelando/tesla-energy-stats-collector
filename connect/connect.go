package connect

import (
	"bytes"
	"crypto/tls"
	b64 "encoding/base64"
	"encoding/json"
	"github.com/iwvelando/tesla-energy-stats-collector/config"
	"github.com/iwvelando/tesla-energy-stats-collector/model"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

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
	err = json.Unmarshal(body, bodyJson)
	if err != nil {
		return client, time.Now(), err
	}
	err = bodyJson.ParseTime()
	if err != nil {
		return client, time.Now(), err
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
		return err
	}

	return nil
}

func GetAll(config *config.Configuration, client *http.Client) (model.Teg, error) {

	teg := model.Teg{}

	err := GetEndpoint(config, client, "/api/meters/aggregates", &teg.Meters)
	if err != nil {
		return teg, err
	}
	teg.Meters.Timestamp = time.Now()
	err = teg.Meters.ParseTime()
	if err != nil {
		return teg, err
	}

	err = GetEndpoint(config, client, "/api/meters/status", &teg.MetersStatus)
	if err != nil {
		return teg, err
	}
	teg.MetersStatus.Timestamp = time.Now()

	err = GetEndpoint(config, client, "/api/operation", &teg.Operation)
	if err != nil {
		return teg, err
	}
	teg.Operation.Timestamp = time.Now()

	err = GetEndpoint(config, client, "/api/powerwalls", &teg.Powerwalls)
	if err != nil {
		return teg, err
	}
	teg.Powerwalls.Timestamp = time.Now()
	err = teg.Powerwalls.ParseTime()
	if err != nil {
		return teg, err
	}

	err = GetEndpoint(config, client, "/api/site_info", &teg.SiteInfo)
	if err != nil {
		return teg, err
	}
	teg.SiteInfo.Timestamp = time.Now()

	err = GetEndpoint(config, client, "/api/sitemaster", &teg.Sitemaster)
	if err != nil {
		return teg, err
	}
	teg.Sitemaster.Timestamp = time.Now()

	err = GetEndpoint(config, client, "/api/solar_powerwall", &teg.SolarPowerwall)
	if err != nil {
		return teg, err
	}
	teg.SolarPowerwall.Timestamp = time.Now()
	err = teg.SolarPowerwall.ParseTime()
	if err != nil {
		return teg, err
	}

	err = GetEndpoint(config, client, "/api/solars", &teg.Solars)
	if err != nil {
		return teg, err
	}
	ts := time.Now()
	for i := range teg.Solars {
		teg.Solars[i].Timestamp = ts
	}

	err = GetEndpoint(config, client, "/api/system/networks/conn_tests", &teg.NetworkConnectionTests)
	if err != nil {
		return teg, err
	}
	teg.NetworkConnectionTests.Timestamp = time.Now()
	err = teg.NetworkConnectionTests.ParseTime()
	if err != nil {
		return teg, err
	}

	err = GetEndpoint(config, client, "/api/status", &teg.Status)
	if err != nil {
		return teg, err
	}
	teg.Status.Timestamp = time.Now()
	err = teg.Status.ParseTime()
	if err != nil {
		return teg, err
	}

	err = GetEndpoint(config, client, "/api/system/testing", &teg.SystemTesting)
	if err != nil {
		return teg, err
	}
	teg.SystemTesting.Timestamp = time.Now()

	err = GetEndpoint(config, client, "/api/system/update/status", &teg.UpdateStatus)
	if err != nil {
		return teg, err
	}
	teg.UpdateStatus.Timestamp = time.Now()

	err = GetEndpoint(config, client, "/api/system_status", &teg.SystemStatus)
	if err != nil {
		return teg, err
	}
	teg.SystemStatus.Timestamp = time.Now()
	err = teg.SystemStatus.ParseFaults()
	if err != nil {
		return teg, err
	}

	err = GetEndpoint(config, client, "/api/system_status/grid_status", &teg.SystemGridStatus)
	if err != nil {
		return teg, err
	}
	teg.SystemGridStatus.Timestamp = time.Now()

	err = GetEndpoint(config, client, "/api/system_status/soe", &teg.SystemStateOfEnergy)
	if err != nil {
		return teg, err
	}
	teg.SystemStateOfEnergy.Timestamp = time.Now()

	return teg, nil
}
