package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Books struct {
	Id                  int    `gorm:"primaryKey"`
	Name                string `gorm:"type:text"`
	Author              string `gorm:"type:text"`
	Quantity_page       int    `gorm:"type:int"`
	Quantity_of_readers int    `gorm:"type:int"`
}

type ManagementBooksHttp struct {
	db *gorm.DB
}

func (m ManagementBooksHttp) create(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		body, _ := ioutil.ReadAll(r.Body)
		fmt.Println(string(body))
		var body_map map[string]string
		json.Unmarshal(body, &body_map)
		valid_keys := []string{"name", "author", "quantity_page", "quantity_of_readers"}
		score := 0
		for key := range body_map {
			if valid_keys[score] == key {
				fmt.Println("valid_keys", valid_keys[score])
				fmt.Println("key", key)
				fmt.Println("Успешно\n")
				score++
				continue
			} else {
				fmt.Println("valid_keys", valid_keys[score])
				fmt.Println("key", key)
				fmt.Fprint(w, "Error")
				break
			}
		}
		quantity_page_int, _ := strconv.Atoi(body_map["quantity_page"])
		quantity_of_readers_int, _ := strconv.Atoi(body_map["quantity_page"])
		book := Books{
			Name:                body_map["name"],
			Author:              body_map["author"],
			Quantity_page:       quantity_page_int,
			Quantity_of_readers: quantity_of_readers_int,
		}
		m.db.Create(&book)
		result := map[string]string{"result": "created"}
		fmt.Fprint(w, result)
	} else if r.Method == "GET" {
		fmt.Fprint(w, "<h1>Hello user!</h1>")
	}
}

func (m ManagementBooksHttp) delete(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		body, _ := ioutil.ReadAll(r.Body)
		var data map[string]string
		info_of_delete := Books{}
		json.Unmarshal(body, &data)
		// fmt.Println(data["id"])
		err := m.db.Delete(&info_of_delete, data["id"])
		if err.Error != nil {
			fmt.Fprint(w, err.Error)
			return
		}
		fmt.Fprint(w, "Book deleted successfully")
	}
}

func (m ManagementBooksHttp) get(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		body, _ := ioutil.ReadAll(r.Body)
		var data map[string]string
		var book Books
		json.Unmarshal(body, &data)
		info := m.db.First(&book, data["id"])
		if info.Error != nil {
			fmt.Fprint(w, info.Error)
			return
		}
		if data["id"] == "" {
			fmt.Fprint(w, "Неверный словарь")
		} else {
			fmt.Fprint(w, book)
		}
	}
}

func (m ManagementBooksHttp) update(w http.ResponseWriter, r *http.Request) {
	if r.Method == "PATCH" {
		id := r.URL.Query().Get("id")
		if id == "" {
			fmt.Fprint(w, "Not id in url")
			return
		}
		var book Books
		info := m.db.First(&book, id)
		fmt.Println(info.Error)
		fmt.Println(book)
		if info.Error != nil {
			fmt.Fprint(w, info.Error)
			return
		}
		var data map[string]string
		body, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(body, &data)
		if name, name_is := data["name"]; name_is {
			book.Name = name
		}
		if author, author_is := data["author"]; author_is {
			book.Author = author
		}
		if quantity_page, page_is := data["quantity_page"]; page_is {
			value, _ := strconv.Atoi(quantity_page)
			book.Quantity_page = value
		}
		if quantity_of_readers, reades_is := data["quantity_of_readers"]; reades_is {
			value, _ := strconv.Atoi(quantity_of_readers)
			book.Quantity_of_readers = value
		}
		m.db.Save(&book)
		fmt.Fprint(w, "Book is successfully updated")
		// fmt.Println(book)
		// val, in :=
	}
}

func main() {
	password, _ := os.LookupEnv("PASSWORD_POSTGRES")
	dsn := fmt.Sprintf("host=localhost user=postgres dbname=book_accouting password=%s sslmode=disable", password)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&Books{})
	// db.Create(&Books{})
	base_path := "/api/v1/book/"
	management_books_http := ManagementBooksHttp{db: db}
	http.HandleFunc(base_path+"append", management_books_http.create)
	http.HandleFunc(base_path+"delete", management_books_http.delete)
	http.HandleFunc(base_path+"get", management_books_http.get)
	http.HandleFunc(base_path+"patch/", management_books_http.update)
	fmt.Println("Сервер запущен")
	http.ListenAndServe(":8000", nil)
}
