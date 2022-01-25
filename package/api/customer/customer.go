package customer

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"taphoa-iot-backend/internal"
	"taphoa-iot-backend/package/api-interface/customer"

	"github.com/gorilla/mux"
)

const (
	TOKEN  = "token"
	PAGE   = "page"
	USER   = "user"
	ROLE   = "role"
	SEARCH = "search"
	STATUS = "status"
	ID     = "id"
	FROM   = "from"
	TO     = "to"
)

func Install(r *mux.Router, s customer.API) {
	r.HandleFunc("/customer/list/{page}", listCustomers(s)).Methods(http.MethodGet, http.MethodHead)
}

type errorRsp struct {
	Msg string
}

func writeRsp(w http.ResponseWriter, err error, rsp interface{}, code int) {
	w.Header().Set("Content-Type", "application/json")
	if err == nil {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}

	if encError := json.NewEncoder(w).Encode(rsp); encError != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func listCustomers(s customer.API) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %v, method: %v\n", r.URL.String(), r.Method)
		log.Printf("Details: %v\n", r.Body)

		spage := mux.Vars(r)[PAGE]
		page, err := strconv.Atoi(spage)
		if err != nil {
			log.Println("Get data failed: unable to parse page input")
			writeRsp(w, err, &errorRsp{
				Msg: fmt.Sprintf("Invalid token"),
			}, http.StatusBadRequest)
			return
		}
		if page <= 0 {
			writeRsp(w, err, &errorRsp{
				Msg: fmt.Sprintf("Invalid page (page should > 0)"),
			}, http.StatusBadRequest)
			return
		}
		page = page - 1
		log.Printf("Process data for page: %v\n", page)

		q := r.URL.Query()
		token := q.Get(TOKEN)
		search := q.Get(SEARCH)

		if !internal.IsTokenExist(token) {
			log.Println("Get failed: unable to get token")
			writeRsp(w, err, &errorRsp{
				Msg: fmt.Sprintf("Invalid token"),
			}, http.StatusNotFound)
			return
		}

		user, err := internal.GetTokenUser(token)
		if err != nil {
			log.Println("Get failed: unable to get user from token")
			writeRsp(w, err, &errorRsp{
				Msg: fmt.Sprintf("Invalid token"),
			}, http.StatusNotFound)
			return
		}

		ret, err := s.ListCustomers(page, user, search)
		if err == nil {
			writeRsp(w, nil, ret, http.StatusOK)
		} else {
			log.Printf("Failed to get items with error: %v", err)
			writeRsp(w, err, &errorRsp{
				Msg: fmt.Sprintf("%v", err),
			}, http.StatusNotFound)
		}
	}
}
