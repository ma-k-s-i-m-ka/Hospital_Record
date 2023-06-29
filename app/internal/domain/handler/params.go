package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func ReadIdParam64(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, fmt.Errorf("id must have type int64")
	}

	return id, nil
}

func ReadIdParam32(r *http.Request) (int32, error) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.ParseInt(params.ByName("id"), 10, 32)
	if err != nil || id < 1 {
		return 0, fmt.Errorf("id must have type int64")
	}

	return int32(id), nil
}

func ReadIdParam16(r *http.Request) (int16, error) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.ParseInt(params.ByName("id"), 10, 16)
	if err != nil || id < 1 {
		return 0, fmt.Errorf("id must have type int64")
	}

	return int16(id), nil
}
