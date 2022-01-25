package receipt

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"taphoa-iot-backend/internal"
	"taphoa-iot-backend/package/api-interface/receipt"

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

func Install(r *mux.Router, s receipt.API) {
	r.HandleFunc("/receipt/add", addReceipt(s)).Methods(http.MethodPost, http.MethodHead)
	r.HandleFunc("/receipt/update", updateReceipt(s)).Methods(http.MethodPut, http.MethodHead)
	r.HandleFunc("/receipt/list/{page}", listReceipt(s)).Methods(http.MethodGet, http.MethodHead)
	r.HandleFunc("/receipt/delete/{id}", deleteReceipt(s)).Methods(http.MethodDelete, http.MethodHead)

	r.HandleFunc("/receipt/report", reportReceipt(s)).Methods(http.MethodGet, http.MethodHead)
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

func parseReqDatas(r *http.Request) (receipt.Instance, error) {
	decoder := json.NewDecoder(r.Body)
	var instance receipt.Instance
	err := decoder.Decode(&instance)
	return instance, err
}

func addReceipt(s receipt.API) http.HandlerFunc {
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

		instance, err = s.AddReceipt(instance, user)
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

func updateReceipt(s receipt.API) http.HandlerFunc {
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

		if instance.ID == 0 {
			writeRsp(w, err, &errorRsp{
				Msg: fmt.Sprintf("Invalid recipe ID"),
			}, http.StatusNotFound)
			return
		}

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

		instance, err = s.UpdateReceipt(instance, user)
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

func listReceipt(s receipt.API) http.HandlerFunc {
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
		sStatus := q.Get(STATUS)
		status, err := strconv.Atoi(sStatus)
		if err != nil {
			log.Println("Get data failed: unable to parse status input")
			writeRsp(w, err, &errorRsp{
				Msg: fmt.Sprintf("%v", err),
			}, http.StatusBadRequest)
			return
		}

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

		ret, err := s.ListReceipts(page, user, search, status)
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

func deleteReceipt(s receipt.API) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %v, method: %v\n", r.URL.String(), r.Method)
		log.Printf("Details: %v\n", r.Body)

		sid := mux.Vars(r)[ID]
		// id, err := strconv.Atoi(sid)
		id, err := strconv.ParseInt(sid, 10, 64)
		if err != nil {
			log.Println("Get data failed: unable to parse id input")
			writeRsp(w, err, &errorRsp{
				Msg: fmt.Sprintf("Invalid token"),
			}, http.StatusBadRequest)
			return
		}
		if id <= 0 {
			writeRsp(w, err, &errorRsp{
				Msg: fmt.Sprintf("Invalid id (id should > 0)"),
			}, http.StatusBadRequest)
			return
		}

		q := r.URL.Query()
		token := q.Get(TOKEN)
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

		err = s.DeleteReceipt(id, user)
		if err == nil {
			writeRsp(w, nil, errorRsp{
				Msg: "Delete successfully",
			}, http.StatusOK)
		} else {
			log.Printf("Failed to get items with error: %v", err)
			writeRsp(w, err, &errorRsp{
				Msg: fmt.Sprintf("%v", err),
			}, http.StatusNotFound)
		}
	}
}

func reportReceipt(s receipt.API) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %v, method: %v\n", r.URL.String(), r.Method)
		log.Printf("Details: %v\n", r.Body)

		q := r.URL.Query()
		token := q.Get(TOKEN)
		from := q.Get(FROM)
		to := q.Get(TO)
		sStatus := q.Get(STATUS)
		status, err := strconv.Atoi(sStatus)
		if err != nil {
			log.Println("Get data failed: unable to parse status input")
			writeRsp(w, err, &errorRsp{
				Msg: fmt.Sprintf("%v", err),
			}, http.StatusBadRequest)
			return
		}

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

		ret, err := s.ReportReceipt(user, from, to, status)
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
