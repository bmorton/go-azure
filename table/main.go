package table

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/bmorton/go-azure/azure"
)

type Client struct {
	credentials azure.StorageCredentials
	connection  *http.Client
	Debug       bool
}

func New(creds azure.StorageCredentials) *Client {
	return &Client{
		credentials: creds,
		connection:  &http.Client{},
	}
}

func (c *Client) BaseURL() string {
	return fmt.Sprintf("https://%s.table.core.windows.net", c.credentials.AccountName)
}

func (c *Client) authorizationHeader(date string, resource string) string {
	signable := fmt.Sprintf("%s\n/%s/%s", date, c.credentials.AccountName, resource)

	key, _ := base64.StdEncoding.DecodeString(c.credentials.AccessKey)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(signable))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	return fmt.Sprintf("SharedKeyLite %s:%s", c.credentials.AccountName, signature)
}

func (c *Client) printDebug(reqOrResp interface{}) {
	if !c.Debug {
		return
	}

	switch r := reqOrResp.(type) {
	default:
		fmt.Printf("Can't debug type %T", r)
	case *http.Request:
		dump, _ := httputil.DumpRequest(r, true)
		fmt.Printf("Request:\n%s\n", dump)
	case *http.Response:
		dump, _ := httputil.DumpResponse(r, true)
		fmt.Printf("Response:\n%s\n", dump)
	}
}
