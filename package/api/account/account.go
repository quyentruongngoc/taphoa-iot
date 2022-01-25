package account

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"taphoa-iot-backend/package/api-interface/account"

	"github.com/gorilla/mux"
)

func Install(r *mux.Router, s account.API) {
	r.HandleFunc("/auth", authticateHandler(s)).Methods(http.MethodPost, http.MethodHead)

	r.HandleFunc("/user/create", createUser(s)).Methods(http.MethodPost, http.MethodHead)
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

func parseReqDatas(r *http.Request) (account.Instance, error) {
	decoder := json.NewDecoder(r.Body)
	var instance account.Instance
	err := decoder.Decode(&instance)
	return instance, err
}

func authticateHandler(s account.API) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %v, method: %v\n", r.URL.String(), r.Method)
		log.Println("Details:", r.Body)

		instance, err := parseReqDatas(r)
		if err != nil {
			log.Println("Authenticate failed: unable to parse user input")
			writeRsp(w, err, &errorRsp{
				Msg: fmt.Sprintf("%v", err),
			}, http.StatusBadRequest)
			return
		}

		instance, err = s.Authenticate(instance)
		if err == nil {
			instance.Passwd = ""
			writeRsp(w, nil, instance, http.StatusOK)
		} else {
			log.Printf("Failed to authenticate with data: %v - %v", instance, err)
			writeRsp(w, err, &errorRsp{
				Msg: fmt.Sprintf("%v", err),
			}, http.StatusNotFound)
		}

	}
}

func createUser(s account.API) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %+v, method: %+v\n", r.URL.String(), r.Method)
		log.Println("Details:", r.Body)

		instance, err := parseReqDatas(r)
		if err != nil {
			log.Println("Create failed: unable to parse user input: ", err)
			writeRsp(w, err, &errorRsp{
				Msg: fmt.Sprintf("%v", err),
			}, http.StatusBadRequest)
			return
		}

		log.Printf("Quyen debug body: %+v\n", instance)

		// if !internal.IsTokenExist(instance.Token) {
		// 	log.Println("Create failed: unable to get token")
		// 	writeRsp(w, err, &errorRsp{
		// 		Msg: fmt.Sprintf("Invalid token"),
		// 	}, http.StatusBadRequest)
		// 	return
		// }

		// role, err := internal.GetTokenRole(instance.Token)
		// if err != nil {
		// 	writeRsp(w, err, &errorRsp{
		// 		Msg: fmt.Sprintf("Invalid role"),
		// 	}, http.StatusBadRequest)
		// 	return
		// }

		// if (role == constant.SystemRole) || (role == constant.AdminRole) || (role == constant.ClinicRole) || (role == constant.AdministrativeRole) {
		// 	instance, err = s.Create(instance, true)
		// 	if err == nil {
		// 		writeRsp(w, nil, fmt.Sprintf("Create sucessfully"), http.StatusOK)
		// 	} else {
		// 		instance.Passwd = ""
		// 		log.Printf("Failed to authenticate with data: %v - %v", instance, err)
		// 		writeRsp(w, err, &errorRsp{
		// 			Msg: fmt.Sprintf("%v", err),
		// 		}, http.StatusNotFound)
		// 	}
		// } else {
		// 	writeRsp(w, err, &errorRsp{
		// 		Msg: fmt.Sprintf("Invalid role"),
		// 	}, http.StatusBadRequest)
		// }

		instance, err = s.Create(instance)
		if err == nil {
			writeRsp(w, nil, fmt.Sprintf("Create sucessfully"), http.StatusOK)
		} else {
			instance.Passwd = ""
			log.Printf("Failed to create with data: %+v - %+v", instance, err)
			writeRsp(w, err, &errorRsp{
				Msg: fmt.Sprintf("%v", err),
			}, http.StatusNotFound)
		}
	}
}
