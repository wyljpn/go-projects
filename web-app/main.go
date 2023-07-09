package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

// 创建全局变量
var db *sql.DB

// Movie 结构体表示电影
type Movie struct {
	ID          int
	Name        string
	ReleaseDate string
}

func main() {
	// 连接数据库
	connectDB()

	// 创建路由处理程序
	http.HandleFunc("/movies", moviesHandler)

	// 启动HTTP服务器
	log.Fatal(http.ListenAndServe(":8080", nil))
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

	fmt.Println("成功连接到数据库！")
}

// 处理/movies路由的请求
func moviesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getMovies(w)
	case http.MethodPost:
		createMovie(w, r)
	case http.MethodPut:
		updateMovie(w, r)
	case http.MethodDelete:
		deleteMovie(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// 获取所有电影
func getMovies(w http.ResponseWriter) {
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

	for _, movie := range movies {
		fmt.Fprintf(w, "ID: %d, Name: %s, Release Date: %s\n", movie.ID, movie.Name, movie.ReleaseDate)
	}
}

// 创建电影
func createMovie(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	releaseDate := r.FormValue("release_date")

	_, err := db.Exec("INSERT INTO movies (name, release_date) VALUES (?, ?)", name, releaseDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, "电影创建成功！")
}

// 更新电影
func updateMovie(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	name := r.FormValue("name")
	releaseDate := r.FormValue("release_date")

	_, err := db.Exec("UPDATE movies SET name = ?, release_date = ? WHERE id = ?", name, releaseDate, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, "电影更新成功！")
}

// 删除电影
func deleteMovie(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	_, err := db.Exec("DELETE FROM movies WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, "电影删除成功！")
}
