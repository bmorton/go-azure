package table

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func (c *Client) Insert(table string, entity interface{}) error {
	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	enc.Encode(entity)

	url := fmt.Sprintf("%s/%s", c.BaseURL(), table)
	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return err
	}

	date := strings.Replace(time.Now().UTC().Add(-time.Minute).Format(time.RFC1123), "UTC", "GMT", -1)
	req.Header.Add("Date", date)
	req.Header.Add("x-ms-version", "2014-02-14")
	req.Header.Add("Authorization", c.authorizationHeader(date, table))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json;odata=nometadata")
	req.Header.Add("Prefer", "return-no-content")
	c.printDebug(req)

	resp, err := c.connection.Do(req)
	if err != nil {
		return err
	}
	c.printDebug(resp)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode == 409 {
		return ErrEntityExists
	}

	if resp.StatusCode == 400 {
		var decoded OdataErrorResponse
		err = json.Unmarshal(body, &decoded)
		if err != nil {
			return err
		}

		if decoded.OdataError.Code == "PropertiesNeedValue" {
			return ErrMissingProperties
		} else {
			return ErrBadRequest
		}
	}

	return nil
}
