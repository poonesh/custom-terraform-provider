package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	_ "github.com/lib/pq"
)

type Food struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Origin string `json:"origin"`
}

func main() {

	db, err := sql.Open("postgres", "user=[your PSQL USERNAME] password=[your PSQL PASSWORD] dbname=terraformproviderexample sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	} else {
		fmt.Println("DB Connected...")
	}

	e := echo.New()

	// Create food
	e.POST("/food", func(c echo.Context) error {
		u := new(Food)
		if err := c.Bind(u); err != nil {
			fmt.Printf("failed to unmarshall the food struct %v", err)
			return err
		}

		sqlStatement := "INSERT INTO foods (name, origin) VALUES ($1, $2) RETURNING id"
		rows, err := db.Query(sqlStatement, u.Name, u.Origin)
		if err != nil {
			fmt.Println(err)
		} else {
			defer rows.Close()
			for rows.Next() {
				if err := rows.Scan(&u.Id); err != nil {
					log.Fatal(err)
				}
				log.Printf("name is %s origin is %s and ID %d\n", u.Name, u.Origin, u.Id)
			}
			c.JSON(http.StatusCreated, u)
			return nil
		}

		return c.String(http.StatusBadRequest, "bad request")
	})

	// Get Food
	e.GET("/food/:id", func(c echo.Context) error {
		var (
			name   string
			origin string
		)

		foodId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.String(http.StatusBadRequest, "bad request, unable to fetch food ID")
		}

		sqlStatement := "SELECT name, origin FROM foods WHERE id=$1"
		rows, err := db.Query(sqlStatement, foodId)
		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusBadRequest, "bad request")
		}

		defer rows.Close()
		for rows.Next() {
			if err := rows.Scan(&name, &origin); err != nil {
				log.Fatal(err)
			}
			log.Printf("name is %s origin is %s\n", name, origin)
		}

		food := &Food{
			Id:     foodId,
			Name:   name,
			Origin: origin,
		}

		c.JSON(http.StatusOK, food)
		return nil
	})

	// Update Food
	e.PUT("/food/:id", func(c echo.Context) error {
		u := new(Food)
		foodId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.String(http.StatusBadRequest, "bad request, unable to fetch food ID")
		}

		if err := c.Bind(&u); err != nil {
			fmt.Printf("failed to unmarshall the food struct %v", err)
			return err
		}

		sqlStatement := "UPDATE foods SET name = $1, origin = $2 WHERE id = $3"
		_, err = db.Query(sqlStatement, u.Name, u.Origin, foodId)
		if err != nil {
			fmt.Println(err)
		} else {
			c.JSON(http.StatusOK, u)
			return nil
		}

		return c.String(http.StatusBadRequest, "bad request")
	})

	// DELETE Food
	e.DELETE("/food/:id", func(c echo.Context) error {
		var u struct{}
		foodId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.String(http.StatusBadRequest, "bad request, unable to fetch food ID")
		}

		sqlStatement := "DELETE FROM foods WHERE id = $1"
		_, err = db.Query(sqlStatement, foodId)
		if err != nil {
			fmt.Println(err)
		} else {
			c.JSON(http.StatusOK, u)
			return nil
		}

		return c.String(http.StatusBadRequest, "bad request")
	})

	if err := e.Start(":8080"); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
