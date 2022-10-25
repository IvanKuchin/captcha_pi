package apihandlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	configreader "github.com/ivankuchin/timecard.ru-api/config-reader"
	"github.com/ivankuchin/timecard.ru-api/logs"
)

var config configreader.Config

func SetConfig(c configreader.Config) {
	config = c
}

func generateTraceID() string {
	const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	const transIDLen = 10

	transID := make([]byte, transIDLen)

	for i := range transID {
		transID[i] = alphabet[rand.Intn(len(alphabet))]
	}

	return string(transID)
}

func getClientRequestBody(tID string, r *http.Request) ([]byte, error) {
	body, err := io.ReadAll(r.Body)

	if err != nil {
		logs.Sugar.Errorw(err.Error(),
			"traceID", tID,
		)
		return []byte{}, err
	}

	if len(body) == 0 {
		error_message := "not enough parameters\n"
		logs.Sugar.Errorw(error_message,
			"traceID", tID,
		)
		return []byte{}, fmt.Errorf("%s", error_message)
	}

	logs.Sugar.Debugw("request url: "+r.RequestURI,
		"traceID", tID,
	)
	logs.Sugar.Debugw("request body: "+string(body),
		"traceID", tID,
	)

	return body, nil
}

func sendReqToServer(tID, url, sessid string) ([]byte, error) {
	buf := new(bytes.Buffer)

	req, err := http.NewRequest(http.MethodGet, url, buf)
	if err != nil {
		logs.Sugar.Errorw(err.Error(),
			"traceID", tID,
		)
		return []byte{}, fmt.Errorf("incorrect http request")
	}
	req.AddCookie(&http.Cookie{Name: "lng", Value: "us"})

	if len(sessid) > 0 {
		req.AddCookie(&http.Cookie{Name: "sessid", Value: sessid})
	}

	// tr := &http.Transport{
	// 	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	// }
	c := &http.Client{
		Timeout: 15 * time.Second,
		// Transport: tr,
	}

	resp, err := c.Do(req)
	if err != nil {
		logs.Sugar.Errorw(err.Error(),
			"traceID", tID,
		)
		return []byte{}, fmt.Errorf("error returned by server")
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		error_message := fmt.Sprintf("server reply http.code %d", resp.StatusCode)
		logs.Sugar.Errorw(error_message,
			"traceID", tID,
		)
		return []byte{}, fmt.Errorf("%s", error_message)
	}

	text, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.Sugar.Errorw(err.Error(),
			"traceID", tID,
		)
		return []byte{}, fmt.Errorf("can't read from server")
	}

	return text, nil
}

func sendReqToServerNoAuth(tID string, url string) ([]byte, error) {
	return sendReqToServer(tID, url, "")
}

func parseServerResponse(tID string, sr []byte) error {

	var server_response login_response
	err := json.Unmarshal(sr, &server_response)
	if err != nil {
		error_message := "incorrect json format"
		logs.Sugar.Errorw(error_message+" (unmarshal error: "+err.Error()+")",
			"traceID", tID,
		)
		return fmt.Errorf("%s", error_message)
	}

	if server_response.Result == "error" {
		logs.Sugar.Debugw("server returned error: "+server_response.Description,
			"traceID", tID,
		)
		return fmt.Errorf("%s", server_response.Description)
	}

	return nil
}
