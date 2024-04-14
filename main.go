package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
)

type Banner struct {
	Content    BannerContent `json:content,omitempty`
	Id         int           `json:"banner_id,omitempty"`
	tag_id     int
	Feature_id int       `json:"feature_id,omitempty"`
	Is_active  bool      `json:"is_active,omitempty"`
	Created_at time.Time `json:"created_at,omitempty"`
	Updated_at time.Time `json:"updated_at,omitempty"`
	Tag_ids    []int     `json:"tag_ids,omitempty"`
}

type BannerContent struct {
	Title string `json:"title"`
	Text  string `json:"text"`
	Url   string `json:"url"`
}

type ErrorBody struct {
	ErrorBody string `json:"error"`
}

func getUserBannerHandler(w http.ResponseWriter, r *http.Request) {
	var tagId int
	var featureId int
	token := r.Header.Get("token")
	if token == "" {
		errorBody := ErrorBody{}
		errorBody.ErrorBody = "string"
		w.Header().Set("Content-Type", "text/JSON; charset=utf-8")
		w.WriteHeader(http.StatusUnauthorized)
		out, _ := json.Marshal(errorBody)
		w.Write(out)
		return
	} else if token != "admin_token" && token != "user_token" {
		errorBody := ErrorBody{}
		errorBody.ErrorBody = "string"

		w.Header().Set("Content-Type", "text/JSON; charset=utf-8")
		w.WriteHeader(http.StatusForbidden)
		out, _ := json.Marshal(errorBody)
		w.Write(out)
		return
	}

	_, tagOk := r.URL.Query()["tag_id"]
	_, featueOk := r.URL.Query()["feature_id"]
	if !tagOk || !featueOk {
		errorBody := ErrorBody{}
		errorBody.ErrorBody = "string"
		w.Header().Set("Content-Type", "text/JSON; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		out, _ := json.Marshal(errorBody)
		w.Write(out)
		return
	} else {
		var err error
		tagId, err = strconv.Atoi(r.URL.Query().Get("tag_id"))
		if err != nil {
			errorBody := ErrorBody{}
			errorBody.ErrorBody = "string"
			w.Header().Set("Content-Type", "text/JSON; charset=utf-8")
			w.WriteHeader(http.StatusBadRequest)
			out, _ := json.Marshal(errorBody)
			w.Write(out)
			return
		}

		featureId, err = strconv.Atoi(r.URL.Query().Get("feature_id"))
		if err != nil {
			errorBody := ErrorBody{}
			errorBody.ErrorBody = "string"
			w.Header().Set("Content-Type", "text/JSON; charset=utf-8")
			w.WriteHeader(http.StatusBadRequest)
			out, _ := json.Marshal(errorBody)
			w.Write(out)
			return
		}

	}

	useLastRevision := false
	_, ok := r.URL.Query()["use_last_revision"]
	if ok {
		useLastRevision, _ = strconv.ParseBool(r.URL.Query().Get("use_last_revision"))
	}

	connStr := "user=postgres password=mypass dbname=banners sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		w.Header().Set("Content-Type", "text/JSON; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		errorBanner := ErrorBody{}
		errorBanner.ErrorBody = "string"
		out, _ := json.Marshal(errorBanner)
		w.Write(out)
	}

	if useLastRevision {
		row, _ := db.Query("SELECT is_active FROM banner WHERE tag_id = $1 AND feature_id = $2", tagId, featureId)
		defer row.Close()
		row.Next()
		var isActive bool
		row.Scan(&isActive)
		if !isActive {
			w.Header().Set("Content-Type", "text/JSON; charset=utf-8")
			w.WriteHeader(http.StatusNotFound)
			errorBanner := ErrorBody{}
			errorBanner.ErrorBody = "string"
			out, _ := json.Marshal(errorBanner)
			w.Write(out)
			return
		}

	}

	defer db.Close()
	exists, _ := db.Query("SELECT EXISTS (SELECT id FROM Banner WHERE tag_id = $1 AND feature_id = $2)", tagId, featureId)
	exists.Next()
	var existsOk bool
	exists.Scan(&existsOk)
	defer exists.Close()
	if !existsOk {
		w.Header().Set("Content-Type", "text/JSON; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
		errorBanner := ErrorBody{}
		errorBanner.ErrorBody = "string"
		out, _ := json.Marshal(errorBanner)
		w.Write(out)
	}

	rows, err := db.Query(fmt.Sprintf("SELECT id FROM Banner WHERE tag_id = %d AND feature_id = %d", tagId, featureId))
	if err != nil {
		w.Header().Set("Content-Type", "text/JSON; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		errorBanner := ErrorBody{}
		errorBanner.ErrorBody = "string"
		out, _ := json.Marshal(errorBanner)
		w.Write(out)
	}

	defer rows.Close()
	for rows.Next() {
		banner := Banner{}
		err := rows.Scan(&banner.Id)
		if err != nil {
			fmt.Println(err)
			continue
		}

		row, _ := db.Query(fmt.Sprintf("SELECT title, banner_text, url FROM Banner_content WHERE id = %d", banner.Id))
		row.Next()
		err = row.Scan(&banner.Content.Title, &banner.Content.Text, &banner.Content.Url)
		if err != nil {
			fmt.Println(err)
			continue
		}

		out, _ := json.Marshal(banner.Content)
		w.Header().Set("Content-Type", "text/JSON; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write(out)
	}

}

func getAllBannersHandler(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("token")
	if token != "" {
		if token != "admin_token" && token != "user_token" {
			w.Header().Set("Content-Type", "text/JSON; charset=utf-8")
			w.WriteHeader(http.StatusForbidden)
			errorBanner := ErrorBody{}
			errorBanner.ErrorBody = "string"
			out, _ := json.Marshal(errorBanner)
			w.Write(out)
			return
		}

	}

	tag_id := 0
	_, tagOk := r.URL.Query()["tag_id"]
	if tagOk {
		tag_id, _ = strconv.Atoi(r.URL.Query().Get("tag_id"))
	}

	feature_id := 0
	_, featureOk := r.URL.Query()["feature_id"]
	if featureOk {
		feature_id, _ = strconv.Atoi(r.URL.Query().Get("feature_id"))
	}

	connStr := "user=postgres password=mypass dbname=banners sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		w.Header().Set("Content-Type", "text/JSON; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Внутренняя ошибка сервера"))
		return
	}

	defer db.Close()
	allrows, _ := db.Query("SELECT COUNT (*) FROM Banner")
	allrows.Next()
	var limit int
	allrows.Scan(&limit)
	_, limitOk := r.URL.Query()["limit"]
	if limitOk {
		limit, _ = strconv.Atoi(r.URL.Query().Get("limit"))
	}

	offset := 0
	_, offsetOk := r.URL.Query()["offset"]
	if offsetOk {
		offset, _ = strconv.Atoi(r.URL.Query().Get("offset"))
	}

	var rows *sql.Rows
	if tagOk && featureOk {
		rows, _ = db.Query("SELECT * FROM Banner WHERE tag_id = $1 AND feature_id = $2 LIMIT $3 OFFSET $4", tag_id, feature_id, limit, offset)
	} else if tagOk && !featureOk {
		rows, _ = db.Query("SELECT * FROM Banner WHERE tag_id = $1 LIMIT $2 OFFSET $3", tag_id, limit, offset)
	} else if !tagOk && featureOk {
		rows, _ = db.Query("SELECT * FROM Banner WHERE feature_id = %1 LIMIT $2 OFFSET $3", feature_id, limit, offset)
	} else if !tagOk && !featureOk {
		rows, _ = db.Query("SELECT * FROM Banner LIMIT $1 OFFSET $2", limit, offset)
	}

	defer rows.Close()
	for rows.Next() {
		banner := Banner{}
		err := rows.Scan(&banner.Id, &banner.tag_id, &banner.Feature_id, &banner.Created_at, &banner.Updated_at, &banner.Is_active)
		banner.Tag_ids = append(banner.Tag_ids, banner.tag_id)
		if err != nil {
			fmt.Println(err)
			continue
		}

		row, _ := db.Query(fmt.Sprintf("SELECT title, banner_text, url FROM Banner_content WHERE id = %d", banner.Id))
		row.Next()
		err = row.Scan(&banner.Content.Title, &banner.Content.Text, &banner.Content.Url)
		if err != nil {
			fmt.Println(err)
			continue
		}

		out, _ := json.Marshal(banner)
		w.Header().Set("Content-Type", "text/JSON; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write(out)
	}

}

func newBannerHandler(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("token")
	if token != "" {
		if token != "admin_token" {
			w.Header().Set("Content-Type", "text/JSON; charset=utf-8")
			w.WriteHeader(http.StatusForbidden)
			errorBanner := ErrorBody{}
			errorBanner.ErrorBody = "string"
			out, _ := json.Marshal(errorBanner)
			w.Write(out)
			return
		}

	}

	banner := Banner{}
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		w.Header().Set("Content-Type", "text/JSON; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		errorBanner := ErrorBody{}
		errorBanner.ErrorBody = "string"
		out, _ := json.Marshal(errorBanner)
		w.Write(out)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &banner); err != nil {
		w.Header().Set("Content-Type", "text/JSON; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		errorBanner := ErrorBody{}
		out, _ := json.Marshal(errorBanner)
		w.Write(out)
		return
	}

	connStr := "user=postgres password=mypass dbname=banners sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		w.Header().Set("Content-Type", "text/JSON; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		errorBanner := ErrorBody{}
		errorBanner.ErrorBody = "string"
		out, _ := json.Marshal(errorBanner.ErrorBody)
		w.Write(out)
	}

	defer db.Close()
	id := 0
	for _, tag_id := range banner.Tag_ids {

		db.Exec("INSERT INTO banner_content VALUES (DEFAULT, $1, $2, $3, $4)", banner.Content.Title, banner.Content.Text, banner.Content.Url, banner.Feature_id)
		row, _ := db.Query("SELECT (id) FROM banner_content ORDER BY id DESC LIMIT 1")
		row.Next()
		row.Scan(&id)
		db.Exec("INSERT INTO banner VALUES($1, $2, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, $4)", &id, tag_id, banner.Feature_id, banner.Is_active)
	}

	w.Header().Set("Content-Type", "text/JSON; charset=utf-8")
	out, _ := json.Marshal(id)
	w.WriteHeader(http.StatusCreated)
	w.Write(out)
}

func deleteBannerHandler(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("token")
	if token != "" {
		if token != "admin_token" {
			w.Header().Set("Content-Type", "text/JSON; charset=utf-8")
			w.WriteHeader(http.StatusForbidden)
			errorBanner := ErrorBody{}
			errorBanner.ErrorBody = "string"
			out, _ := json.Marshal(errorBanner)
			w.Write(out)
			return
		}

	}

	id, err := strconv.Atoi(strings.Split(r.URL.Path, "/")[2])
	if err != nil {
		w.Header().Set("Content-Type", "text/JSON; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		errorBanner := ErrorBody{}
		errorBanner.ErrorBody = "string"
		out, _ := json.Marshal(errorBanner)
		w.Write(out)
	}

	connStr := "user=postgres password=mypass dbname=banners sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		w.Header().Set("Content-Type", "text/JSON; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		errorBanner := ErrorBody{}
		errorBanner.ErrorBody = "string"
		out, _ := json.Marshal(errorBanner)
		w.Write(out)
	}

	db.Exec("DELETE FROM banner_content WHERE id = $1", id)
	db.Exec("DELETE FROM banner WHERE id = $1", id)
	w.WriteHeader(http.StatusNoContent)
}

func updateBannerHandler(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("token")
	if token != "" {
		if token != "admin_token" {
			w.Header().Set("Content-Type", "text/JSON; charset=utf-8")
			w.WriteHeader(http.StatusForbidden)
			errorBanner := ErrorBody{}
			errorBanner.ErrorBody = "string"
			out, _ := json.Marshal(errorBanner)
			w.Write(out)
			return
		}

	}

	id, err := strconv.Atoi(strings.Split(r.URL.Path, "/")[2])
	if err != nil {
		w.Header().Set("Content-Type", "text/JSON; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		errorBanner := ErrorBody{}
		errorBanner.ErrorBody = "string"
		out, _ := json.Marshal(errorBanner)
		w.Write(out)
	}

	connStr := "user=postgres password=mypass dbname=banners sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		w.Header().Set("Content-Type", "text/JSON; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		errorBanner := ErrorBody{}
		errorBanner.ErrorBody = "string"
		out, _ := json.Marshal(errorBanner)
		w.Write(out)
	}

	exists, _ := db.Query("SELECT EXISTS (SELECT id FROM Banner WHERE id = $1)", id)
	exists.Next()
	var existsOk bool
	exists.Scan(&existsOk)
	defer exists.Close()
	if !existsOk {
		w.Header().Set("Content-Type", "text/JSON; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
		errorBanner := ErrorBody{}
		errorBanner.ErrorBody = "string"
		out, _ := json.Marshal(errorBanner)
		w.Write(out)
	}

	banner := Banner{}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(r.Body)
	if err != nil {
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &banner); err != nil {
		return
	}

	db.Exec("UPDATE banner_content set title = $1, banner_text = $2, url = $3, feature_id = $4 WHERE id = $5", banner.Content.Title, banner.Content.Text, banner.Content.Url, banner.Feature_id, id)

	for _, tag_id := range banner.Tag_ids {

		db.Exec("UPDATE banner set tag_id = $1, feature_id = $2, updated_at = CURRENT_TIMESTAMP, is_active = $3 WHERE id = $4", tag_id, banner.Feature_id, banner.Is_active, id)
	}

}

func main() {
	r := chi.NewRouter()
	r.Get("/banner", getAllBannersHandler)
	r.Get("/user_banner", getUserBannerHandler)
	r.Post("/banner", newBannerHandler)
	r.Patch("/banner/{id}", updateBannerHandler)
	r.Delete("/banner/{id}", deleteBannerHandler)
	if err := http.ListenAndServe(":3000", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
	}

}
