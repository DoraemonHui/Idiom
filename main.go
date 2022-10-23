package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type Idiom struct {
	Derivation   string `json:"derivation"`
	Example      string `json:"example"`
	Explanation  string `json:"explanation"`
	Pinyin       string `json:"pinyin"`
	Word         string `json:"word"`
	Abbreviation string `json:"abbreviation"`
}

type GetIdiomListRequest struct {
	StartWord string `json:"start_word"`
}

type GetIdiomDetail struct {
	Word string `json:"word"`
}

var idiomMap map[string]Idiom
var idiomWordList []string

func main() {
	if idiomMap == nil {
		loadData()
	}
	startHttpServer()
}

func loadData() {
	jsonFile, err := os.Open("idiom.json")
	if err != nil {
		fmt.Println(err)
	}

	// 要记得关闭
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	idiomList := make([]Idiom, 0)
	json.Unmarshal([]byte(byteValue), &idiomList)

	idiomMap = make(map[string]Idiom)
	for _, idiom := range idiomList {
		idiomMap[idiom.Word] = idiom
	}

	idiomWordList = getMapKeys(idiomMap)
}

func findValidIdiom(startWord string) []string {
	result := make([]string, 0)

	for _, idiomWord := range idiomWordList {
		if idiomWord[0:3] == startWord {
			result = append(result, idiomWord)
		}
	}

	return result
}

func fundIdiomDetail(word string) Idiom {
	return idiomMap[word]
}

func getMapKeys(m map[string]Idiom) []string {
	// 数组默认长度为map长度,后面append时,不需要重新申请内存和拷贝,效率很高
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func getIdiomList(w http.ResponseWriter, r *http.Request) {
	var request GetIdiomListRequest
	_ = json.NewDecoder(r.Body).Decode(&request)
	idiom := findValidIdiom(request.StartWord)
	enc := json.NewEncoder(w)
	enc.Encode(idiom)
}

func getIdiomDetail(w http.ResponseWriter, r *http.Request) {
	var request GetIdiomDetail
	_ = json.NewDecoder(r.Body).Decode(&request)
	idiom := fundIdiomDetail(request.Word)
	enc := json.NewEncoder(w)
	enc.Encode(idiom)
}

func startHttpServer() {
	router := mux.NewRouter()

	//通过完整的path来匹配
	router.HandleFunc("/api/getIdiomList", getIdiomList)
	router.HandleFunc("/api/getIdiomDetail", getIdiomDetail)
	router.Methods("POST")

	// 初始化
	srv := &http.Server{
		Handler:      router,
		Addr:         ":80",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
