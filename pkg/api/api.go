// API приложения GoNews.
package api

import (
	"CommentApi/pkg/db"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type API struct {
	db *db.DB
	r  *mux.Router
}

type News struct {
	Id    int
	Title string
}

// Конструктор API.
func New(db *db.DB) *API {
	a := API{db: db, r: mux.NewRouter()}
	a.endpoints()
	return &a
}

// Router возвращает маршрутизатор для использования
// в качестве аргумента HTTP-сервера.
func (api *API) Router() *mux.Router {
	return api.r
}

// Регистрация методов API в маршрутизаторе запросов.
func (api *API) endpoints() {
	api.r.HandleFunc("/comment", api.comment).Methods(http.MethodPost, http.MethodOptions)
}

func (api *API) comment(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}

	var comment db.PostComment
	err = json.Unmarshal(b, &comment)
	if err != nil {
		log.Fatal(err)
	}

	err = api.db.AddComment(comment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
