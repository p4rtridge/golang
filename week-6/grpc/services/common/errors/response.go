package errors

import (
	"encoding/json"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func WriteError(e error) error {
	if err, ok := e.(CodeCarrier); ok {
		return status.Errorf(err.Code(), e.Error())
	}

	return status.Errorf(codes.Internal, "Internal Error")
}

func WriteJSONError(w http.ResponseWriter, status int, err error) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}
