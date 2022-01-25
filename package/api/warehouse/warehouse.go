package warehouse

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"taphoa-iot-backend/internal"
	"taphoa-iot-backend/package/api-interface/warehouse"

	"github.com/gorilla/mux"
)

const (
	TOKEN  = "token"
	PAGE   = "page"
	USER   = "user"
	ROLE   = "role"
	SEARCH = "search"
)

func Install(r *mux.Router, s warehouse.API) {
	r.HandleFunc("/warehouse/item/create", createItem(s)).Methods(http.MethodPost, http.MethodHead)
	r.HandleFunc("/warehouse/item/update", updateItem(s)).Methods(http.MethodPut, http.MethodHead)
	r.HandleFunc("/warehouse/item/list/{page}", getItemList(s)).Methods(http.MethodGet, http.MethodHead)
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

func parseReqDatas(r *http.Request) (warehouse.Instance, error) {
	decoder := json.NewDecoder(r.Body)
	var instance warehouse.Instance
	err := decoder.Decode(&instance)
	return instance, err
}

func createItem(s warehouse.API) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %v, method: %v\n", r.URL.String(), r.Method)
		log.Printf("Details: %v\n", r.Body)

		instance, err := parseReqDatas(r)
		if err != nil {
			log.Println("Create failed: unable to parse user input: ", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		log.Printf("Quyen debug body: %v\n", instance)

		if !internal.IsTokenExist(instance.Token) {
			log.Println("Create failed: unable to get token")
			writeRsp(w, err, &errorRsp{
				Msg: fmt.Sprintf("Invalid token"),
			}, http.StatusNotFound)
			return
		}

		user, err := internal.GetTokenUser(instance.Token)
		if err != nil {
			log.Println("Create failed: unable to get user from token")
			writeRsp(w, err, &errorRsp{
				Msg: fmt.Sprintf("Invalid token"),
			}, http.StatusNotFound)
			return
		}

		instance, err = s.CreateItem(instance, user)
		if err == nil {
			writeRsp(w, nil, fmt.Sprintf("Create/Update sucessfully"), http.StatusOK)
		} else {
			log.Printf("Failed to create/update patient with data: %v - %v", instance, err)
			writeRsp(w, err, &errorRsp{
				Msg: fmt.Sprintf("%v", err),
			}, http.StatusNotFound)
		}
	}
}

func updateItem(s warehouse.API) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %v, method: %v\n", r.URL.String(), r.Method)
		log.Printf("Details: %v\n", r.Body)

		instance, err := parseReqDatas(r)
		if err != nil {
			log.Println("Create failed: unable to parse user input: ", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		log.Printf("Quyen debug body: %v\n", instance)

		if !internal.IsTokenExist(instance.Token) {
			log.Println("Update failed: unable to get token")
			writeRsp(w, err, &errorRsp{
				Msg: fmt.Sprintf("Invalid token"),
			}, http.StatusNotFound)
			return
		}

		user, err := internal.GetTokenUser(instance.Token)
		if err != nil {
			log.Println("Update failed: unable to get user from token")
			writeRsp(w, err, &errorRsp{
				Msg: fmt.Sprintf("Invalid token"),
			}, http.StatusNotFound)
			return
		}

		if instance.ID == 0 {
			log.Println("Update failed: invalid ID")
			writeRsp(w, err, &errorRsp{
				Msg: fmt.Sprintf("Invalid ID"),
			}, http.StatusNotFound)
			return
		}

		instance, err = s.UpdateItem(instance, user)
		if err == nil {
			writeRsp(w, nil, fmt.Sprintf("Update/Update sucessfully"), http.StatusOK)
		} else {
			log.Printf("Failed to Update/update patient with data: %v - %v", instance, err)
			writeRsp(w, err, &errorRsp{
				Msg: fmt.Sprintf("%v", err),
			}, http.StatusNotFound)
		}
	}
}

func getItemList(s warehouse.API) http.HandlerFunc {
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

		ret, err := s.DescribeItems(page, user, search)
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
