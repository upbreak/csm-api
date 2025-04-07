package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type ResultRole string

const (
	Success ResultRole = "Success"
	Failure ResultRole = "Failure"
)

type ErrDetailsRole string

const (
	ResponseEncodeError ErrDetailsRole = "Response Encode Error"
	InvalidToken        ErrDetailsRole = "Invalid Token or Not Found Token"
	BodyDataParseError  ErrDetailsRole = "Body Data Parse Error"
	InvalidUser         ErrDetailsRole = "Invalid User"
	TokenCreatedFail    ErrDetailsRole = "Token Created Fail"
	NotFoundParam       ErrDetailsRole = "Not Found Parameter"
	ParsingError        ErrDetailsRole = "Parsing Error"
	DataAddFailed       ErrDetailsRole = "Data Add Failed"
	DataModifyFailed    ErrDetailsRole = "Data Modify Failed"
	DataRemoveFailed    ErrDetailsRole = "Data Remove Failed"
	DataMergeFailed     ErrDetailsRole = "Data Merge Failed"
	CallApiFailed       ErrDetailsRole = "Call Api Failed"
)

type ErrResponse struct {
	Result         ResultRole     `json:"result"`
	Message        string         `json:"message"`
	Details        ErrDetailsRole `json:"details"`
	HttpStatusCode int            `json:"http-status-code"`
}

type Response struct {
	Result ResultRole `json:"result"`
	Values any        `json:"values"`
}

func RespondJSON(ctx context.Context, w http.ResponseWriter, body any, status int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	bodyBytes, err := json.Marshal(body)

	if err != nil {
		fmt.Printf("encode response error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		rsp := ErrResponse{
			Result:  Success,
			Message: http.StatusText(http.StatusInternalServerError),
			Details: ResponseEncodeError,
		}
		if err := json.NewEncoder(w).Encode(rsp); err != nil {
			fmt.Printf("write error response error: %v", err)
		}
		return
	}

	w.WriteHeader(status)
	if _, err = fmt.Fprintf(w, "%s", bodyBytes); err != nil {
		fmt.Printf("write response error: %v", err)
	}
}
