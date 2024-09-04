package extract

import (
	"net/http"
	"strconv"
)

func ExtractTaskID(req *http.Request, param string) (int, error) {
	return strconv.Atoi(req.PathValue(param))
}
