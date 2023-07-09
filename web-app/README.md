# Go-test

## MySQL
### Connect to MySQL
```go
db, err = sql.Open("mysql", "root:wangyulong6@tcp(localhost:3306)/yulong_test")
```

### Create tables in MySQL
```sql
CREATE TABLE yulong_test.movies (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    release_date DATE
);
```

### Insert sample data
```sql
INSERT INTO movies (name, release_date) VALUES
    ('Movie 1', '2022-01-01'),
    ('Movie 2', '2022-02-15'),
    ('Movie 3', '2022-03-10');
```


## Create go project
### Init
```shell
go mod init yulong_test
```
### Install libraries
```shell
go mod tidy
```

### Build
```shell
go build
```

### Start 
```shell
./yulong_test
```

## Test
### Search
```shell
curl http://localhost:8080/movies
ID: 1, Name: Movie 1, Release Date: 2022-01-01
ID: 2, Name: Movie 2, Release Date: 2022-02-15
ID: 3, Name: Movie 3, Release Date: 2022-03-10
```

### Insert
```shell
curl -X POST -d "name=Movie 4&release_date=2022-04-20" http://localhost:8080/movies
电影创建成功！%
```

### Delete
```shell
curl -X DELETE "http://localhost:8080/movies?id=2"
电影删除成功！%
```

### Update
```shell
curl -X PUT -d "id=1&name=Updated Movie&release_date=2022-05-01" http://localhost:8080/movies
电影更新成功！%
```
