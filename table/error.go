package table

import "errors"

var (
	ErrTableExists       = errors.New("table already exists")
	ErrEntityExists      = errors.New("entity already exists")
	ErrMissingProperties = errors.New("values are not specified for all properties in the entity (likely RowKey or PartitionKey)")
	ErrBadRequest        = errors.New("bad request")
)

type OdataErrorResponse struct {
	OdataError struct {
		Code    string `json:"code"`
		Message struct {
			Lang  string `json:"lang"`
			Value string `json:"value"`
		} `json:"message"`
	} `json:"odata.error"`
}
