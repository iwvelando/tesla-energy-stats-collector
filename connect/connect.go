package connect

import (
	"bytes"
	"crypto/tls"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/iwvelando/tesla-energy-stats-collector/config"
	"github.com/iwvelando/tesla-energy-stats-collector/model"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"sync"
	"time"
)

const expectedHTTPStatus = 200

// Auth authenticates to the Tesla Energy Gateway and returns an HTTP client
func Auth(conf *config.Configuration) (*http.Client, time.Time, error) {

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: conf.TeslaGateway.SkipVerifySsl}
	client := &http.Client{}

	data := &model.AuthPayload{
		Username:   "customer",
		Password:   conf.TeslaGateway.Password,
		Email:      conf.TeslaGateway.Email,
		ForceSmOff: false,
	}
	dataJSON, err := json.Marshal(data)
	if err != nil {
		return client, time.Now(), err
	}

	req, err := http.NewRequest("POST", conf.TeslaGateway.Address+"/api/login/Basic", bytes.NewBuffer(dataJSON))
	if err != nil {
		return client, time.Now(), err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return client, time.Now(), err
	}
	defer resp.Body.Close()

	bodyJSON := &model.AuthResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return client, time.Now(), err
	}

	status := resp.StatusCode
	if status != expectedHTTPStatus {
		err = fmt.Errorf("expected %d HTTP status code but got %d; raw body %s", expectedHTTPStatus, resp.StatusCode, body)
		return client, time.Now(), err
	}

	err = json.Unmarshal(body, bodyJSON)
	if err != nil {
		return client, time.Now(), err
	}
	err = bodyJSON.ParseTime()
	if err != nil {
		return client, time.Now(), fmt.Errorf("error when parsing authentication time, %s", err)
	}

	jar, _ := cookiejar.New(nil)
	var cookies []*http.Cookie
	cookie := &http.Cookie{
		Name:    "AuthCookie",
		Value:   bodyJSON.Token,
		Path:    "/",
		Domain:  "",
		Expires: bodyJSON.LoginTime.Add(24 * time.Hour),
	}
	cookies = append(cookies, cookie)
	cookie = &http.Cookie{
		Name:    "UserRecord",
		Value:   b64.StdEncoding.EncodeToString(body),
		Path:    "/",
		Domain:  "",
		Expires: bodyJSON.LoginTime.Add(24 * time.Hour),
	}
	cookies = append(cookies, cookie)
	u, _ := url.Parse(conf.TeslaGateway.Address)
	jar.SetCookies(u, cookies)
	client.Jar = jar

	return client, bodyJSON.LoginTime.Add(23*time.Hour + 55*time.Minute), nil
}

// GetEndpoint queries an individual endpoint and stores the results in the provided data structure
func GetEndpoint(conf *config.Configuration, client *http.Client, endpoint string, data interface{}) error {
	req, err := http.NewRequest("GET", conf.TeslaGateway.Address+endpoint, nil)
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

	status := resp.StatusCode
	if status != expectedHTTPStatus {
		err = fmt.Errorf("expected %d HTTP status code but got %d; raw body %s", expectedHTTPStatus, resp.StatusCode, body)
		return err
	}

	if endpoint == "/api/devices/vitals" {
		err = proto.Unmarshal(body, data.(protoreflect.ProtoMessage))
	} else {
		err = json.Unmarshal(body, data)
	}
	if err != nil {
		err = fmt.Errorf("%w; raw body %s", err, body)
		return err
	}

	return nil
}

// GetAll calls GetEndpoint for all specified endpoints
func GetAll(conf *config.Configuration, client *http.Client) (model.Teg, error) {

	teg := model.Teg{}

	errChan := make(chan error)
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func(waitgroup *sync.WaitGroup) {
		defer waitgroup.Done()
		endpoint := "/api/meters/aggregates"
		err := GetEndpoint(conf, client, endpoint, &teg.Meters)
		if err != nil {
			errChan <- fmt.Errorf("error when querying %s, %s", endpoint, err)
			return
		}
		teg.Meters.Timestamp = time.Now()
		err = teg.Meters.ParseTime()
		if err != nil {
			errChan <- fmt.Errorf("error when parsing time for endpoint %s, %s", endpoint, err)
			return
		}
	}(&wg)

	wg.Add(1)
	go func(waitgroup *sync.WaitGroup) {
		defer waitgroup.Done()
		endpoint := "/api/meters/status"
		err := GetEndpoint(conf, client, endpoint, &teg.MetersStatus)
		if err != nil {
			errChan <- fmt.Errorf("error when querying %s, %s", endpoint, err)
			return
		}
		teg.MetersStatus.Timestamp = time.Now()
	}(&wg)

	wg.Add(1)
	go func(waitgroup *sync.WaitGroup) {
		defer waitgroup.Done()
		endpoint := "/api/operation"
		err := GetEndpoint(conf, client, endpoint, &teg.Operation)
		if err != nil {
			errChan <- fmt.Errorf("error when querying %s, %s", endpoint, err)
			return
		}
		teg.Operation.Timestamp = time.Now()
	}(&wg)

	wg.Add(1)
	go func(waitgroup *sync.WaitGroup) {
		defer waitgroup.Done()
		endpoint := "/api/powerwalls"
		err := GetEndpoint(conf, client, endpoint, &teg.Powerwalls)
		if err != nil {
			errChan <- fmt.Errorf("error when querying %s, %s", endpoint, err)
			return
		}
		teg.Powerwalls.Timestamp = time.Now()
		err = teg.Powerwalls.ParseTime()
		if err != nil {
			errChan <- fmt.Errorf("error when parsing time for endpoint %s, %s", endpoint, err)
			return
		}
	}(&wg)

	wg.Add(1)
	go func(waitgroup *sync.WaitGroup) {
		defer waitgroup.Done()
		endpoint := "/api/site_info"
		err := GetEndpoint(conf, client, endpoint, &teg.SiteInfo)
		if err != nil {
			errChan <- fmt.Errorf("error when querying %s, %s", endpoint, err)
			return
		}
		teg.SiteInfo.Timestamp = time.Now()
	}(&wg)

	wg.Add(1)
	go func(waitgroup *sync.WaitGroup) {
		defer waitgroup.Done()
		endpoint := "/api/sitemaster"
		err := GetEndpoint(conf, client, endpoint, &teg.Sitemaster)
		if err != nil {
			errChan <- fmt.Errorf("error when querying %s, %s", endpoint, err)
			return
		}
		teg.Sitemaster.Timestamp = time.Now()
	}(&wg)

	wg.Add(1)
	go func(waitgroup *sync.WaitGroup) {
		defer waitgroup.Done()
		endpoint := "/api/solars"
		err := GetEndpoint(conf, client, endpoint, &teg.Solars)
		if err != nil {
			errChan <- fmt.Errorf("error when querying %s, %s", endpoint, err)
			return
		}
		ts := time.Now()
		for i := range teg.Solars {
			teg.Solars[i].Timestamp = ts
		}
	}(&wg)

	wg.Add(1)
	go func(waitgroup *sync.WaitGroup) {
		defer waitgroup.Done()
		endpoint := "/api/system/networks/conn_tests"
		err := GetEndpoint(conf, client, endpoint, &teg.NetworkConnectionTests)
		if err != nil {
			errChan <- fmt.Errorf("error when querying %s, %s", endpoint, err)
			return
		}
		teg.NetworkConnectionTests.Timestamp = time.Now()
		err = teg.NetworkConnectionTests.ParseTime()
		if err != nil {
			errChan <- fmt.Errorf("error when parsing time for endpoint %s, %s", endpoint, err)
			return
		}
	}(&wg)

	wg.Add(1)
	go func(waitgroup *sync.WaitGroup) {
		defer waitgroup.Done()
		endpoint := "/api/status"
		err := GetEndpoint(conf, client, endpoint, &teg.Status)
		if err != nil {
			errChan <- fmt.Errorf("error when querying %s, %s", endpoint, err)
			return
		}
		teg.Status.Timestamp = time.Now()
		err = teg.Status.ParseTime()
		if err != nil {
			errChan <- fmt.Errorf("error when parsing time for endpoint %s, %s", endpoint, err)
			return
		}
	}(&wg)

	wg.Add(1)
	go func(waitgroup *sync.WaitGroup) {
		defer waitgroup.Done()
		endpoint := "/api/system/testing"
		err := GetEndpoint(conf, client, endpoint, &teg.SystemTesting)
		if err != nil {
			errChan <- fmt.Errorf("error when querying %s, %s", endpoint, err)
			return
		}
		teg.SystemTesting.Timestamp = time.Now()
	}(&wg)

	wg.Add(1)
	go func(waitgroup *sync.WaitGroup) {
		defer waitgroup.Done()
		endpoint := "/api/system/update/status"
		err := GetEndpoint(conf, client, endpoint, &teg.UpdateStatus)
		if err != nil {
			errChan <- fmt.Errorf("error when querying %s, %s", endpoint, err)
			return
		}
		teg.UpdateStatus.Timestamp = time.Now()
	}(&wg)

	wg.Add(1)
	go func(waitgroup *sync.WaitGroup) {
		defer waitgroup.Done()
		endpoint := "/api/system_status"
		err := GetEndpoint(conf, client, endpoint, &teg.SystemStatus)
		if err != nil {
			errChan <- fmt.Errorf("error when querying %s, %s", endpoint, err)
			return
		}
		teg.SystemStatus.Timestamp = time.Now()
		err = teg.SystemStatus.ParseTime()
		if err != nil {
			errChan <- fmt.Errorf("error when parsing time for endpoint %s, %s", endpoint, err)
			return
		}
		err = teg.SystemStatus.ParseFaults()
		if err != nil {
			errChan <- fmt.Errorf("error when parsing faults for endpoint %s, %s", endpoint, err)
			return
		}
	}(&wg)

	wg.Add(1)
	go func(waitgroup *sync.WaitGroup) {
		defer waitgroup.Done()
		endpoint := "/api/system_status/grid_status"
		err := GetEndpoint(conf, client, endpoint, &teg.SystemGridStatus)
		if err != nil {
			errChan <- fmt.Errorf("error when querying %s, %s", endpoint, err)
			return
		}
		teg.SystemGridStatus.Timestamp = time.Now()
	}(&wg)

	wg.Add(1)
	go func(waitgroup *sync.WaitGroup) {
		defer waitgroup.Done()
		endpoint := "/api/system_status/soe"
		err := GetEndpoint(conf, client, endpoint, &teg.SystemStateOfEnergy)
		if err != nil {
			errChan <- fmt.Errorf("error when querying %s, %s", endpoint, err)
			return
		}
		teg.SystemStateOfEnergy.Timestamp = time.Now()
	}(&wg)

	go func() {
		wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		if err != nil {
			return teg, err
		}
	}

	return teg, nil
}
