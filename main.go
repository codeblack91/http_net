package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

var CountryCapital = map[string]string{
	"Russia":         "Moscow",
	"USA":            "Washington, D.C.",
	"Germany":        "Berlin",
	"France":         "Paris",
	"China":          "Beijing",
	"Japan":          "Tokyo",
	"Brazil":         "Brasília",
	"Canada":         "Ottawa",
	"India":          "New Delhi",
	"South Africa":   "Pretoria",
	"Australia":      "Canberra",
	"United Kingdom": "London",
	"Italy":          "Rome",
	"Mexico":         "Mexico City",
	"Turkey":         "Ankara",
}

func main() {
	port := ":" + "8080"

	http.HandleFunc("/", welcomToTheServer)
	http.HandleFunc("/info", info)
	http.HandleFunc("/user/", userId)
	http.HandleFunc("/items/", items)
	http.HandleFunc("/search", search)
	http.HandleFunc("/error-test", errorTest)

	fmt.Printf("запускаем серверна на: %s", port)
	http.ListenAndServe(port, nil)
}

func welcomToTheServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcom to the server!")
	//fmt.Println(r.URL.Query().Get(""))
}

func info(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	name := queryParams.Get("name")

	w.Header().Set("Content-Type", "application/json")
	response := `{"message": "Hello,` + name + `"}`

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func userId(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/user/"):]

	w.Header().Set("Content-Type", "application/json")
	response := `{"userID": ` + `"` + id + `"` + `}`
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func items(w http.ResponseWriter, r *http.Request) {

	//curl -X GET http://localhost:8080/items/

	//curl -X POST http://localhost:8080/items/ \
	//-H "Content-Type: application/json" \
	//-d '{"key1": "value1", "key2": "value2"}'

	//curl -X PUT http://localhost:8080/items/Russia \
	//-H "Content-Type: application/json" \
	//-d '{"value": "St. Petersburg"}'

	//curl -X DELETE http://localhost:8080/items/USA

	// GET: выводит все данные из CountryCapital
	if r.Method == http.MethodGet {
		for country, capital := range CountryCapital {
			fmt.Printf("Country: %s, Capital: %s\n", country, capital)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Data printed to console\n"))
		return
	}

	// POST: выводит полученные данные в консоль
	if r.Method == http.MethodPost {
		var data map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		fmt.Println("Received POST data:", data)
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("POST data printed to console\n"))
		return
	}

	// PUT: проверяет наличие ключа и выводит в консоль
	if r.Method == http.MethodPut {
		key := strings.TrimPrefix(r.URL.Path, "/items/")
		if key == "" {
			http.Error(w, "Key not provided", http.StatusBadRequest)
			return
		}

		var data struct {
			Value string `json:"value"`
		}
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		fmt.Printf("Key: %s, Updated Value: %s\n", key, data.Value)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("PUT data printed to console\n"))
		return
	}

	// DELETE: проверяет наличие ключа и выводит удаление в консоль
	if r.Method == http.MethodDelete {
		key := strings.TrimPrefix(r.URL.Path, "/items/")

		fmt.Printf("Key %s will be deleted\n", key)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("DELETE operation printed to console\n"))
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func search(w http.ResponseWriter, r *http.Request) {
	//curl "http://localhost:8080/search?key=Rus&value=Mos"

	// Получаем query-параметры key и value
	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")

	if key == "" || value == "" {
		http.Error(w, "Both key and value query parameters are required", http.StatusBadRequest)
		return
	}

	// Ищем все элементы в map, которые соответствуют фильтру
	var result []map[string]string
	for k, v := range CountryCapital {
		if strings.Contains(k, key) && strings.Contains(v, value) {
			result = append(result, map[string]string{
				"country": k,
				"capital": v,
			})
		}
	}

	// Если результат пустой, возвращаем ошибку 404
	if len(result) == 0 {
		http.Error(w, "No matching results found", http.StatusNotFound)
		return
	}

	// Возвращаем результат в формате JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func errorTest(w http.ResponseWriter, r *http.Request) {
	xTest := r.Header.Get("X-Test")
	if xTest == "" {
		http.Error(w, "X-Test does not exist", http.StatusForbidden)
	} else if xTest == "fail" {
		http.Error(w, "X-test does not exist fail", http.StatusInternalServerError)
	}
}
