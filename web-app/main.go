package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// 创建全局变量
var db *sql.DB

// Movie 结构体表示电影
type Movie struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	ReleaseDate string `json:"release_date"`
}

func main() {
	// 连接数据库
	connectDB()

	// 创建路由
	router := mux.NewRouter()

	// 定义路由处理程序
	router.HandleFunc("/movies", getMovies).Methods("GET")
	router.HandleFunc("/movies", createMovie).Methods("POST")
	router.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	router.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	router.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	// 启动HTTP服务器
	log.Fatal(http.ListenAndServe(":8080", router))
}

// 连接数据库
func connectDB() {
	var err error
	db, err = sql.Open("mysql", "root:wangyulong6@tcp(localhost:3306)/yulong_test")
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("成功连接到数据库！")
}

// 获取所有电影
func getMovies(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT * FROM movies")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var movies []Movie
	for rows.Next() {
		var movie Movie
		err = rows.Scan(&movie.ID, &movie.Name, &movie.ReleaseDate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		movies = append(movies, movie)
	}

	jsonResponse(w, movies)
}

// 创建电影
func createMovie(w http.ResponseWriter, r *http.Request) {
	var movie Movie
	err := json.NewDecoder(r.Body).Decode(&movie)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := db.Exec("INSERT INTO movies (name, release_date) VALUES (?, ?)", movie.Name, movie.ReleaseDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	lastInsertID, _ := result.LastInsertId()
	movie.ID = int(lastInsertID)
	jsonResponse(w, movie)
}

// 获取单个电影
func getMovie(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var movie Movie
	err := db.QueryRow("SELECT * FROM movies WHERE id = ?", id).Scan(&movie.ID, &movie.Name, &movie.ReleaseDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse(w, movie)
}

// 更新电影
func updateMovie(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var movie Movie
	err := json.NewDecoder(r.Body).Decode(&movie)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = db.Exec("UPDATE movies SET name = ?, release_date = ? WHERE id = ?", movie.Name, movie.ReleaseDate, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	movie.ID, _ = strconv.Atoi(id)
	jsonResponse(w, movie)
}

// 删除电影
func deleteMovie(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	_, err := db.Exec("DELETE FROM movies WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// 返回JSON响应
func jsonResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
