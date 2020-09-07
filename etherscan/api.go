package etherscan

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

type APIClient struct {
	Key string
}

type ABIResponse struct {
	Status  string // "1"
	Message string // "OK"
	Result  string
}

// {"status":"1","message":"OK","result": ... }

func (api *APIClient) GetABI(address string) (abistr string, err error) {
	u, err := url.Parse("https://api.etherscan.io/api")
	if err != nil {
		return
	}

	// https://api.etherscan.io/api?module=contract&action=getabi&address=0xBB9bc244D798123fDe783fCc1C72d3Bb8C189413&apikey=YourApiKeyToken

	q := u.Query()
	q.Add("apikey", api.Key)
	q.Add("module", "contract")
	q.Add("action", "getabi")
	q.Add("address", address)
	u.RawQuery = q.Encode()

	// fmt.Println("url", u)

	res, err := http.Get(u.String())
	if err != nil {
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", errors.New("cannot get ABI data")
	}

	dec := json.NewDecoder(res.Body)

	var data ABIResponse
	err = dec.Decode(&data)
	if err != nil {
		return
	}

	return data.Result, nil
}
