package table

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type QueryResult struct {
	Value []interface{}
}

type RowQuery struct {
	Table        string
	PartitionKey string
	RowKey       string
	Fields       []string
}

func (q RowQuery) Path() string {
	return fmt.Sprintf("%s()", q.Table)
}

func (q RowQuery) Filter() string {
	return fmt.Sprintf("PartitionKey eq '%s' and RowKey eq '%s'", q.PartitionKey, q.RowKey)
}

func (q RowQuery) QueryString() string {
	return fmt.Sprintf("$filter=%s&$select=%s", url.QueryEscape(q.Filter()), strings.Join(q.Fields, ","))
}

func (c *Client) GetEntity(query RowQuery, result interface{}) error {
	reqURL := fmt.Sprintf("%s/%s?%s", c.BaseURL(), query.Path(), query.QueryString())
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return err
	}

	date := strings.Replace(time.Now().UTC().Add(-time.Minute).Format(time.RFC1123), "UTC", "GMT", -1)
	req.Header.Add("Date", date)
	req.Header.Add("x-ms-version", "2014-02-14")
	req.Header.Add("Authorization", c.authorizationHeader(date, url.QueryEscape(query.Path())))
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

	var parsed QueryResult
	err = json.Unmarshal(body, &parsed)
	if err != nil {
		return err
	}
	value, err := json.Marshal(parsed.Value[0])
	if err != nil {
		return err
	}
	err = json.Unmarshal(value, &result)
	if err != nil {
		return err
	}

	return nil
}
