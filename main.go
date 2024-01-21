package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Структура артист
type Artist struct {
	ID    string   `json:"id"`    //id коллектива
	Name  string   `json:"name"`  //название группы
	Born  string   `json:"born"`  //год основания группы
	Genre string   `json:"genre"` //жанр
	Songs []string `json:'songs"` //популярные песни, это слайс строк, т.е. песен много
}

// Создаем мапу
var artists = map[string]Artist{
	"1": {
		ID:    "1",
		Name:  "30 Seconds To Mars",
		Born:  "1998",
		Genre: "alternative",
		Songs: []string{
			"The Kill",
			"A Beautiful Lie",
			"Attack",
			"Live Like A Dream",
		},
	},
	"2": {
		ID:    "2",
		Name:  "Garbage",
		Born:  "1994",
		Genre: "alternative",
		Songs: []string{
			"Queer",
			"Shut Your Mouth",
			"Cup of Coffee",
			"Til the Day I Die",
		},
	},
}

// Еще одна мапа Русский рок
var rusArtists = map[string]Artist{
	"1": {
		ID:    "1",
		Name:  "Король и Шут",
		Born:  "1992",
		Genre: "Панк",
		Songs: []string{
			"Ели мясо мужики",
			"Лесник",
			"Кузьма и барин",
			"Ром",
		},
	},
	"2": {
		ID:    "2",
		Name:  "Nautilus Pompilius",
		Born:  "1982",
		Genre: "Русский рок",
		Songs: []string{
			"Утро полины",
			"Титаник",
			"Скованные одной цепью",
			"Крылья",
		},
	},
	"3": {
		ID:    "3",
		Name:  "Жуки",
		Born:  "1997",
		Genre: "Рок`н`Ролл",
		Songs: []string{
			"Батарейка",
			"Танкист",
			"Властелин колец",
			"Влечение",
		},
	},
}

func getArtists(w http.ResponseWriter, r *http.Request) {
	//сериализуем данные из слайса artists
	resp, err := json.Marshal(artists)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// в заголовок записываем тип контента, у нас это данные в формате JSON
	w.Header().Set("Content-Type", "application/json")
	// так как все успешно, то статус OK
	w.WriteHeader(http.StatusOK)
	// записываем сериализованные в JSON данные в тело ответа
	w.Write(resp)
}

func postArtist(w http.ResponseWriter, r *http.Request) {
	var artist Artist
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &artist); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	artists[artist.ID] = artist

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func getArtist(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	artist, ok := artists[id]
	if !ok {
		http.Error(w, "Артист не найден", http.StatusNoContent)
		return
	}

	resp, err := json.Marshal(artist)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func getRusArtists(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	artist, ok := rusArtists[id]
	if !ok {
		http.Error(w, "Артист не найден", http.StatusNoContent)
		return
	}

	resp, err := json.Marshal(artist)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func main() {
	// создаем новый роутер
	r := chi.NewRouter()

	// регистрируем в роутере эндпоинт `/artists` с методом GET, для которого используется обработчик `getArtists`
	r.Get("/artists", getArtists)
	// регистрируем в роутере эндпоинт `/artists` с методом POST, для которого используется обработчик `postArtist`
	r.Post("/artists", postArtist)
	// регестрируем в роутере эндпоинт `/artist/{id}`  с методом GETб, для которого используется обработчик `getArtist`
	r.Get("/artist/{id}", getArtist)

	// Регистрируем роутер Русский рок
	r.Get("/russian_rock/{id}", getRusArtists)

	// запускаем сервер
	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}

// ```
// Чтобы посмотреть на работу нового обработчика, нужно сначала изучить имеющиеся данные. Для этого в браузере перейдите по адресу `localhost:8080/artists`. Браузер покажет те же строки, которые были в прошлый раз:

// ```JSON
// {"1":{"id":"1","name":"30 Seconds To Mars","born":"1998","genre":"alternative","songs":["The Kill","A Beautiful Lie","Attack","Live Like A Dream"]},"2":{"id":"2","name":"Garbage","born":"1994","genre":"alternative","songs":["Queer","Shut Your Mouth","Cup of Coffee","Til the Day I Die"]}}
