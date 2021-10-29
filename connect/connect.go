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
)

func Auth(config *config.Configuration) (*http.Client, error) {

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
		return client, err
	}

	req, err := http.NewRequest("POST", config.TeslaGateway.Address+"/api/login/Basic", bytes.NewBuffer(dataJson))
	if err != nil {
		return client, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return client, err
	}
	defer resp.Body.Close()

	bodyJson := &model.AuthResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return client, err
	}
	err = json.Unmarshal(body, bodyJson)
	if err != nil {
		return client, err
	}

	jar, _ := cookiejar.New(nil)
	var cookies []*http.Cookie
	cookie := &http.Cookie{
		Name:   "AuthCookie",
		Value:  bodyJson.Token,
		Path:   "/",
		Domain: "",
	}
	cookies = append(cookies, cookie)
	cookie = &http.Cookie{
		Name:   "UserRecord",
		Value:  b64.StdEncoding.EncodeToString(body),
		Path:   "/",
		Domain: "",
	}
	cookies = append(cookies, cookie)
	u, _ := url.Parse(config.TeslaGateway.Address)
	jar.SetCookies(u, cookies)

	client.Jar = jar

	return client, nil
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

//func GetAll(config *config.Configuration) (data Metrics, error) {
//
//}
