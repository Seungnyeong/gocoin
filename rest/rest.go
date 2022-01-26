package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/snkim/sncoin/blockchain"
	"github.com/snkim/sncoin/utils"
)
var port string

type url string

func (u url) MarshalText() ([]byte, error) {
	url := fmt.Sprintf("http://localhost%s%s",port, u)
	return []byte(url), nil
}
 
type urlDescription struct {
	URL url `json:"url"`
	Method string `json:"method"`
	Description string `json:"description"`
	Payload string `json:"payload,omitempty"`
}

type addBlockBody struct {
	Message string
}

type errorResponse struct {
	ErrorMessage string `json:"errorMessage"`
}


func documentation(rw http.ResponseWriter, r *http.Request) {
	data := []urlDescription{
		{
			 URL: url("/"),
			 Method: "GET",
			 Description: "See Documentation",
		},
		{
			URL : url("/status"),
			Method: "GET",
			Description: "Get blockchain status",
		},
		{
			 URL: url("/blocks"),
			 Method: "GET",
			 Description: "See All Blocks",
		},
		{
			 URL: url("/blocks"),
			 Method: "POST",
			 Description: "Add a block",
			 Payload: "data:string",
		},
		{
			 URL: url("/blocks/{hash}"),
			 Method: "GET",
			 Description: "See a block",
		},
	}
	json.NewEncoder(rw).Encode(data)
}

func blocks(rw http.ResponseWriter, r *http.Request) {
	switch r.Method{
	case  "GET":
		json.NewEncoder(rw).Encode(blockchain.Blockchain().Blocks())
	case "POST":
		var addBlockBody addBlockBody
		utils.HandleErr(json.NewDecoder(r.Body).Decode(&addBlockBody))
		blockchain.Blockchain().AddBlock(addBlockBody.Message)
		rw.WriteHeader(http.StatusCreated)
	default:
		rw.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func block(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]
	
	block, err := blockchain.FindBlock(hash)
	encoder := json.NewEncoder(rw)
	
	if err == blockchain.ErrNotFound {
		encoder.Encode(errorResponse{fmt.Sprint(err)})
	} else {
		encoder.Encode(block)
	}
	
}

func jsonContentTypeMiddleWare(next http.Handler) http.Handler {
	// adapter pattern 6.8 chapter
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(rw, r)
	})
}

func status(rw http.ResponseWriter, r *http.Request) {
	json.NewEncoder(rw).Encode(blockchain.Blockchain())
}

func Start(aPort int) {
	// ServeMux는 url(/block) 와 url(blocks)를 연결해주는 역할을 한다.
	router := mux.NewRouter()
	port = fmt.Sprintf(":%d", aPort)
	router.Use(jsonContentTypeMiddleWare)
	router.HandleFunc("/", documentation).Methods("GET")
	router.HandleFunc("/status", status)
	router.HandleFunc("/blocks", blocks).Methods("GET", "POST")
	router.HandleFunc("/blocks/{hash:[a-f0-9]+}", block).Methods("GET")
	fmt.Printf("Listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}