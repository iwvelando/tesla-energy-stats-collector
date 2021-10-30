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
	fmt.Println(endpoint)
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

func GetAll(config *config.Configuration, client *http.Client) error {

	tegMeters := model.TegMeters{}
	err := GetEndpoint(config, client, "/api/meters/aggregates", &tegMeters)
	if err != nil {
		return err
	}
	err = tegMeters.ParseTime()
	if err != nil {
		return err
	}
	fmt.Println(tegMeters)

	tegMetersStatus := model.TegMetersStatus{}
	err = GetEndpoint(config, client, "/api/meters/status", &tegMetersStatus)
	if err != nil {
		return err
	}
	fmt.Println(tegMetersStatus)

	tegOperation := model.TegOperation{}
	err = GetEndpoint(config, client, "/api/operation", &tegOperation)
	if err != nil {
		return err
	}
	fmt.Println(tegOperation)

	tegPowerwalls := model.TegPowerwalls{}
	err = GetEndpoint(config, client, "/api/powerwalls", &tegPowerwalls)
	if err != nil {
		return err
	}
	err = tegPowerwalls.ParseTime()
	if err != nil {
		return err
	}
	fmt.Println(tegPowerwalls)

	tegSiteInfo := model.TegSiteInfo{}
	err = GetEndpoint(config, client, "/api/site_info", &tegSiteInfo)
	if err != nil {
		return err
	}
	fmt.Println(tegSiteInfo)

	tegSitemaster := model.TegSitemaster{}
	err = GetEndpoint(config, client, "/api/sitemaster", &tegSitemaster)
	if err != nil {
		return err
	}
	fmt.Println(tegSitemaster)

	tegSolarPowerwall := model.TegSolarPowerwall{}
	err = GetEndpoint(config, client, "/api/solar_powerwall", &tegSolarPowerwall)
	if err != nil {
		return err
	}
	err = tegSolarPowerwall.ParseTime()
	if err != nil {
		return err
	}
	fmt.Println(tegSolarPowerwall)

	tegSolars := []model.TegSolars{}
	err = GetEndpoint(config, client, "/api/solars", &tegSolars)
	if err != nil {
		return err
	}
	for _, i := range tegSolars {
		fmt.Println(i)
	}

	tegNetworkConnectionTests := model.TegNetworkConnectionTests{}
	err = GetEndpoint(config, client, "/api/system/networks/conn_tests", &tegNetworkConnectionTests)
	if err != nil {
		return err
	}
	err = tegNetworkConnectionTests.ParseTime()
	if err != nil {
		return err
	}
	fmt.Println(tegNetworkConnectionTests)

	tegStatus := model.TegStatus{}
	err = GetEndpoint(config, client, "/api/status", &tegStatus)
	if err != nil {
		return err
	}
	err = tegStatus.ParseTime()
	if err != nil {
		return err
	}
	fmt.Println(tegStatus)

	tegSystemTesting := model.TegSystemTesting{}
	err = GetEndpoint(config, client, "/api/system/testing", &tegSystemTesting)
	if err != nil {
		return err
	}
	fmt.Println(tegSystemTesting)

	tegUpdateStatus := model.TegUpdateStatus{}
	err = GetEndpoint(config, client, "/api/system/update/status", &tegUpdateStatus)
	if err != nil {
		return err
	}
	fmt.Println(tegUpdateStatus)

	tegSystemStatus := model.TegSystemStatus{}
	err = GetEndpoint(config, client, "/api/system_status", &tegSystemStatus)
	if err != nil {
		return err
	}
	err = tegSystemStatus.ParseFaults()
	if err != nil {
		return err
	}
	fmt.Println(tegSystemStatus)

	tegSystemGridStatus := model.TegSystemGridStatus{}
	err = GetEndpoint(config, client, "/api/system_status/grid_status", &tegSystemGridStatus)
	if err != nil {
		return err
	}
	fmt.Println(tegSystemGridStatus)

	tegSystemStateOfEnergy := model.TegSystemStateOfEnergy{}
	err = GetEndpoint(config, client, "/api/system_status/soe", &tegSystemStateOfEnergy)
	if err != nil {
		return err
	}
	fmt.Println(tegSystemStateOfEnergy)

	return nil
}
