package resources

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

type AzureError struct {
	Response ResponseError `json:"error"`
}

func (err *AzureError) Error() string {
	return err.Response.Message
}

func asAzureError(err error) *AzureError {
	e, ok := err.(*azcore.ResponseError)
	if !ok {
		panic(fmt.Sprintf("error is not of type *azcore.ResponseError - %t", err))
	}

	azerr := AzureError{}

	out, _ := ioutil.ReadAll(e.RawResponse.Body)
	if err := json.Unmarshal(out, &azerr); err != nil {
		panic(err)
	}

	if azerr.Response.Code == "" {
		panic("error parsing error message from Azure")
	}

	return &azerr
}

type ResponseError struct {
	Code           string                   `json:"code"`
	Target         string                   `json:"target"`
	Message        string                   `json:"message"`
	AdditionalInfo []ResponseAdditionalInfo `json:"additionalInfo"`
}

type ResponseAdditionalInfo struct {
	Type string                     `json:"type"`
	Info ResponseAdditionalInfoInfo `json:"info"`
}

type ResponseAdditionalInfoInfo struct {
	Type                        string `json:"Type"`
	PolicyDefinitionEffect      string `json:"policyDefinitionEffect"`
	PolicyAssignmentName        string `json:"policyAssignmentName"`
	PolicyAssignmentDisplayName string `json:"policyAssignmentDisplayName"`
}
