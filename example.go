package main

import (
	"fmt"

	"github.com/bmorton/go-azure/azure"
	"github.com/bmorton/go-azure/table"
	"github.com/k0kubun/pp"
)

type Message struct {
	Body         string
	RowKey       string
	PartitionKey string
}

func main() {
	c := table.New(azure.StorageCredentials{
		AccountName: "my-storage-account",
		AccessKey:   "my-access-key",
	})
	c.Debug = true

	err := c.Create("messages")
	if err != nil {
		fmt.Println(err)
	}

	err = c.Insert("messages", Message{
		Body:         "My message body!",
		RowKey:       "1",
		PartitionKey: "network-1",
	})
	if err != nil {
		fmt.Println(err)
	}

	m := &Message{}
	err = c.GetEntity(table.RowQuery{
		Table:        "messages",
		PartitionKey: "network-1",
		RowKey:       "1",
		Fields:       []string{"Body", "RowKey", "PartitionKey"},
	}, m)

	if err != nil {
		fmt.Println(err)
	}

	pp.Println(m)
}
