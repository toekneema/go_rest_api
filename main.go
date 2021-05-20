package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/mux"
)

type Item struct { // essentially a class/object
	UID   string  `json:"uid"`
	Name  string  `json:"name"`
	Desc  string  `json:"desc"`
	Price float64 `json:"price"`
}

var inventory []Item //global var

func homePage(w http.ResponseWriter, r *http.Request) { // handles loading the home page
	fmt.Fprintf(w, "Endpoint called: homePage()")
}

func getInventory(w http.ResponseWriter, r *http.Request) { // handles the GET request
	w.Header().Set("Content-Type", "application/json") // just setting the header to use JSON format, w is used to print/write out to http response
	fmt.Println("Function called: getInventory()")

	json.NewEncoder(w).Encode(inventory) // this will write out the full inventory/DB to the http GET request
}

func createItem(w http.ResponseWriter, r *http.Request) { // handles the POST request
	w.Header().Set("Content-Type", "application/json")

	var item Item
	_ = json.NewDecoder(r.Body).Decode(&item) //decode/deserialize the passed in json object into that {item} created in the line above, dont care about the return value so we use underscore
	inventory = append(inventory, item)       //add to database

	json.NewEncoder(w).Encode(item) //then encode the item back to json to display to webpage
}

func deleteItem(w http.ResponseWriter, r *http.Request) { // handles the DELETE request
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	_deleteItemAtUid(params["uid"])
	json.NewEncoder(w).Encode(inventory) //then encode the item back to json to display to webpage
}
func _deleteItemAtUid(uid string) {
	for idx, item := range inventory {
		if item.UID == uid {
			// delete from slice
			inventory = append(inventory[:idx], inventory[idx+1:]...)
			break
		}
	}
}

func updateItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var item Item
	_ = json.NewDecoder(r.Body).Decode(&item)

	params := mux.Vars(r)

	_deleteItemAtUid(params["uid"])     // call my custom delete method
	inventory = append(inventory, item) // then add new item to list/arr

	json.NewEncoder(w).Encode(inventory)
}

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", homePage).Methods("GET")
	router.HandleFunc("/inventory", getInventory).Methods("GET") // be careful, there is no ending slash, so don't include it when making http calls, otherwise it breaks
	router.HandleFunc("/inventory", createItem).Methods("POST")
	router.HandleFunc("/inventory/{uid}", deleteItem).Methods("DELETE")
	router.HandleFunc("/inventory/{uid}", updateItem).Methods("PUT")

	log.Fatal(http.ListenAndServe(":8000", router))
}

func main() {
	fmt.Println("Starting....")
	inventory = append(inventory, Item{
		UID:   "0",
		Name:  "Shoes",
		Desc:  "Nike Dunk Low Black White",
		Price: 736,
	})
	inventory = append(inventory, Item{
		UID:   "1",
		Name:  "Fancy Shoes",
		Desc:  "Nike Air Dior High 1s",
		Price: 9995,
	})
	handleRequests()
	fmt.Println("Done!")
}
