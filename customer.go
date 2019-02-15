package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	_ "strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Infrom struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Status string `json:"status"`
}

var Recs []Infrom

func Setup() *gin.Engine {
	r := gin.Default()
	r.Use(logon)
	//api := r.Group("/api")
	r.POST("/customers", AddCusHandler)
	r.GET("/customers", GetCusAllHandler)
	r.GET("/customers/:id", GetCusByIdHandler)
	r.DELETE("/customers/:id", DelCusHandler)
	r.PUT("/customers/:id", UpdateCusHandler)

	return r
}

//func ADD//
func AddCusHandler(c *gin.Context) {
	//	id, _ := strconv.Atoi(c.Param("id"))
	var newRec Infrom
	err := c.ShouldBindJSON(&newRec)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	row := db.QueryRow("INSERT INTO customers(name,email,status) values ($1,$2,$3) RETURNING id", newRec.Name, newRec.Email, newRec.Status)
	var id int
	err = row.Scan(&id)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	newRec.ID = id
	c.JSON(http.StatusCreated, newRec)
}

//fun getall
func GetCusAllHandler(c *gin.Context) {
	stmt, err := db.Prepare("SELECT id,name,email,status FROM customers")
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	rows, err := stmt.Query()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	var items []Infrom
	for rows.Next() {
		var item Infrom
		err := rows.Scan(&item.ID, &item.Name, &item.Email, &item.Status)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		items = append(items, item)
	}
	c.JSON(http.StatusOK, items)
}

//func get by id
func GetCusByIdHandler(c *gin.Context) {
	id := c.Param("id")
	stmt, err := db.Prepare("SELECT id,name,email,status FROM customers WHERE id = $1")
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	row := stmt.QueryRow(id)
	var item Infrom
	err = row.Scan(&item.ID, &item.Name, &item.Email, &item.Status)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, item)
}

//func updata//
func UpdateCusHandler(c *gin.Context) {

	id := c.Param("id")
	var newRec Infrom
	err := c.ShouldBindJSON(&newRec)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	stmt, err := db.Prepare("UPDATE customers SET name=$2,email=$3,status=$4 where id=$1")
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	_, err = stmt.Exec(id, newRec.Name, newRec.Email, newRec.Status)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	stmt1, err := db.Prepare("SELECT id, name,email,status FROM customers where id =$1")

	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	row := stmt1.QueryRow(id)
	var update Infrom
	err = row.Scan(&update.ID, &update.Name, &update.Email, &update.Status)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, update)

}

//func delete//
func DelCusHandler(c *gin.Context) {
	id := c.Param("id")
	stmt, err := db.Prepare("DELETE FROM customers WHERE id=$1")
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	_, err = stmt.Exec(id)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(200, gin.H{"message": "customer deleted"})
}

//func create tabel//
func CreateTable() {
	cre := `
	CREATE TABLE IF NOT EXISTS customers (
		id SERIAL PRIMARY KEY,
		name TEXT,
		email TEXT,
		status TEXT
	);`

	_, err := db.Exec(cre)
	if err != nil {
		log.Fatal("Cannot create table", err)
	}
	fmt.Println("create sucess")
}

//login

func logon(c *gin.Context) {
	authKey := c.GetHeader("Authorization")
	if authKey != "token2019" {
		c.JSON(http.StatusUnauthorized, "Unauthorization")
		c.Abort()

		return
	}
	c.Next()

}

//DATABASE//
var db *sql.DB

func main() {

	var err error
	//url := "postgres://vasxfitw:vug40XStORLuFX6ouwCDUvc-VFFet97P@echo.db.elephantsql.com:5432/vasxfitw"
	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("can't connect to database", err)
	}
	defer db.Close()
	CreateTable()
	r := Setup()

	r.Run(":2019")
}
