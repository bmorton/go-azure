package table

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type CreateParams struct {
	TableName string
}

func (c *Client) Create(name string) error {
	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	enc.Encode(CreateParams{TableName: name})

	url := fmt.Sprintf("%s/Tables", c.BaseURL())
	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return err
	}

	date := strings.Replace(time.Now().UTC().Add(-time.Minute).Format(time.RFC1123), "UTC", "GMT", -1)
	req.Header.Add("Date", date)
	req.Header.Add("x-ms-version", "2014-02-14")
	req.Header.Add("Authorization", c.authorizationHeader(date, "Tables"))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json;odata=nometadata")
	c.printDebug(req)

	resp, err := c.connection.Do(req)
	if err != nil {
		return err
	}
	c.printDebug(resp)

	if resp.StatusCode == 409 {
		return ErrTableExists
	}

	return nil
}
