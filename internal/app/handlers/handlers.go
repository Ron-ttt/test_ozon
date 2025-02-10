package handlers

import (
	"encoding/json"
	"net/http"
	"net/url"
	"testozon/internal/app/config"
	"testozon/internal/app/storage"
	"testozon/internal/app/utils"

	"github.com/gorilla/mux"
)

type URLRegistry struct {
	URL string `json:"url"`
}
type URLRegistryResult struct {
	Result string `json:"result"`
}

func Init() handlerWrapper {
	localhost, baseURL, dbAdress := config.Flags()

	//dbAdress := "postgresql://postgres:190603@localhost:5432/postgres"
	if dbAdress != "" {
		dBStorage, err := storage.NewDBStorage(dbAdress)
		if err == nil {
			return handlerWrapper{storageInterface: dBStorage, Localhost: localhost, baseURL: baseURL + "/"}
		}
	}
	return handlerWrapper{storageInterface: storage.NewMapStorage(), Localhost: localhost, baseURL: baseURL + "/"}
}

func MInit() handlerWrapper {
	return handlerWrapper{storageInterface: storage.NewMockStorage(), Localhost: "localhost:8080", baseURL: "http://localhost:8080/"}
}

type handlerWrapper struct {
	storageInterface storage.Storage
	Localhost        string
	baseURL          string
}

func (hw handlerWrapper) IndexPage(res http.ResponseWriter, req *http.Request) { // post
	var longURL URLRegistry
	if err := json.NewDecoder(req.Body).Decode(&longURL); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	_, err1 := url.ParseRequestURI(longURL.URL)
	if err1 != nil {
		http.Error(res, "invalid url", http.StatusBadRequest)
		return
	}
	var rez URLRegistryResult
	oldshort, err := hw.storageInterface.Find(string(longURL.URL))
	if err == nil {
		res.Header().Set("content-type", "application/json")
		res.WriteHeader(http.StatusConflict)
		rez.Result = hw.baseURL + oldshort
		if err := json.NewEncoder(res).Encode(rez); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	res.Header().Set("content-type", "application/json")
	res.WriteHeader(http.StatusCreated)
	length := 10 // Укажите длину строки
	randomString := utils.RandString(length)
	rez.Result = hw.baseURL + randomString
	err = hw.storageInterface.Add(randomString, string(longURL.URL))
	if err != nil {
		http.Error(res, "error adding to database", http.StatusBadRequest)
		return
	}
	if err := json.NewEncoder(res).Encode(rez); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (hw handlerWrapper) Redirect(res http.ResponseWriter, req *http.Request) { //get
	params := mux.Vars(req)
	id := params["id"]
	originalURL, ok := hw.storageInterface.Get(id)
	if ok != nil {
		http.Error(res, "not found", http.StatusBadRequest)
		return
	}
	res.Header().Set("Location", originalURL)
	res.WriteHeader(http.StatusTemporaryRedirect)
}
