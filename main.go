package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

type Album struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Image       any    `json:"image"`
	Description any    `json:"description"`
}

type Singer struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Profile any    `json:"profile"`
}

type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description any    `json:"description"`
}

type Song struct {
	ID               int    `json:"id"`
	Title            string `json:"title"`
	AlbumID          any    `json:"album_id"`
	SingerID         any    `json:"singer_id"`
	CategoryID       any    `json:"category_id"`
	SongsDescription any    `json:"songs_description"`
	IsFavorite       any    `json:"is_favorite"`
}

func main() {

	err := godotenv.Load()
	// if err != nil {
	// 	log.Fatalf("Error loading .env file: %v", err)
	// }

	db, err := sql.Open("mysql", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	defer db.Close()

	// Make sure connection is available
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Get("/api/data", func(c *fiber.Ctx) error {
		return c.JSON(map[string]string{
			"message": "Hello From go server",
		})
	})

	app.Get("/albums", func(c *fiber.Ctx) error {
		rows, err := db.Query("SELECT * FROM album_table")
		if err != nil {
			return err
		}

		var albums []Album // Assuming Album is your struct for albums

		for rows.Next() {
			var album Album
			err := rows.Scan(&album.ID, &album.Title, &album.Image, &album.Description)
			if err != nil {
				return err
			}

			albums = append(albums, album)
		}

		return c.JSON(albums)

	})

	app.Get("/albums/:id", func(c *fiber.Ctx) error {

		id := c.Params("id")

		row := db.QueryRow("SELECT * FROM album_table WHERE album_id = ?", id)
		var album Album

		err := row.Scan(&album.ID, &album.Title, &album.Image, &album.Description)

		if err != nil {
			return err
		}

		return c.JSON(album)

	})

	app.Post("/albums", func(c *fiber.Ctx) error {

		body := &Album{}

		if err := c.BodyParser(body); err != nil {
			return err
		}

		res, err := db.Exec("INSERT INTO album_table (album_title, album_image, album_description) VALUES (?, ?, ?)", body.Title, body.Image, body.Description)
		if err != nil {
			return err
		}

		id, err := res.LastInsertId()

		if err != nil {
			return err
		}

		body.ID = int(id)

		return c.JSON(body)
	})

	app.Put("/albums/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		body := &Album{}

		if err := c.BodyParser(body); err != nil {
			return err
		}

		_, err := db.Exec("UPDATE album_table SET album_title = ?, album_image = ?, album_description = ? WHERE album_id = ?", body.Title, body.Image, body.Description, id)

		if err != nil {
			return err
		}

		return c.SendStatus(fiber.StatusOK)

	})

	app.Delete("/albums/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		_, err := db.Exec("DELETE FROM album_table WHERE album_id = ?", id)

		if err != nil {
			return err
		}

		return c.SendStatus(fiber.StatusOK)

	})

	app.Get("/singers", func(c *fiber.Ctx) error {

		rows, err := db.Query("SELECT * FROM singers_table")
		if err != nil {
			return err
		}

		var singers []Singer // Assuming Singer is your struct for singers

		for rows.Next() {
			var singer Singer
			err := rows.Scan(&singer.ID, &singer.Name, &singer.Profile)
			if err != nil {
				return err
			}

			singers = append(singers, singer)
		}

		return c.JSON(singers)

	})

	app.Get("/singers/:id", func(c *fiber.Ctx) error {

		id := c.Params("id")

		row := db.QueryRow("SELECT * FROM singers_table WHERE singer_id = ?", id)
		var singer Singer

		err := row.Scan(&singer.ID, &singer.Name, &singer.Profile)

		if err != nil {
			return err
		}

		return c.JSON(singer)

	})

	app.Post("/singers", func(c *fiber.Ctx) error {

		body := &Singer{}

		if err := c.BodyParser(body); err != nil {
			return err
		}

		res, err := db.Exec("INSERT INTO singers_table (singer_name, singer_profile) VALUES (?, ?)", body.Name, body.Profile)
		if err != nil {
			return err
		}

		id, err := res.LastInsertId()

		if err != nil {
			return err
		}

		body.ID = int(id)

		return c.JSON(body)

	})

	app.Put("/singers/:id", func(c *fiber.Ctx) error {

		id := c.Params("id")

		body := &Singer{}

		if err := c.BodyParser(body); err != nil {
			return err
		}

		_, err := db.Exec("UPDATE singers_table SET singer_name = ?, singer_profile = ? WHERE singer_id = ?", body.Name, body.Profile, id)

		if err != nil {

			return err
		}

		return c.SendStatus(fiber.StatusOK)

	})

	app.Delete("/singers/:id", func(c *fiber.Ctx) error {

		id := c.Params("id")

		_, err := db.Exec("DELETE FROM singers_table WHERE singer_id = ?", id)

		if err != nil {

			return err
		}

		return c.SendStatus(fiber.StatusOK)

	})

	app.Get("/categories", func(c *fiber.Ctx) error {

		rows, err := db.Query("SELECT * FROM song_category_table")
		if err != nil {
			return err
		}

		var categories []Category // Assuming Category is your struct for categories

		for rows.Next() {
			var category Category
			err := rows.Scan(&category.ID, &category.Name, &category.Description)
			if err != nil {
				return err
			}

			categories = append(categories, category)
		}

		return c.JSON(categories)

	})

	app.Get("/categories/:id", func(c *fiber.Ctx) error {

		id := c.Params("id")

		row := db.QueryRow("SELECT * FROM song_category_table WHERE category_id = ?", id)

		var category Category

		err := row.Scan(&category.ID, &category.Name, &category.Description)

		if err != nil {

			return err

		}

		return c.JSON(category)

	})

	app.Post("/categories", func(c *fiber.Ctx) error {

		body := &Category{}

		if err := c.BodyParser(body); err != nil {

			return err

		}

		res, err := db.Exec("INSERT INTO song_category_table (category_name, category_description) VALUES (?, ?)", body.Name, body.Description)

		if err != nil {

			return err

		}

		id, err := res.LastInsertId()

		if err != nil {

			return err

		}

		body.ID = int(id)

		return c.JSON(body)

	})

	app.Put("/categories/:id", func(c *fiber.Ctx) error {

		id := c.Params("id")

		body := &Category{}

		if err := c.BodyParser(body); err != nil {

			return err

		}

		_, err := db.Exec("UPDATE song_category_table SET category_name = ?, category_description = ? WHERE category_id = ?", body.Name, body.Description, id)

		if err != nil {

			return err

		}

		return c.SendStatus(fiber.StatusOK)

	})

	app.Delete("/categories/:id", func(c *fiber.Ctx) error {

		id := c.Params("id")

		_, err := db.Exec("DELETE FROM song_category_table WHERE category_id = ?", id)

		if err != nil {

			return err

		}

		return c.SendStatus(fiber.StatusOK)

	})

	app.Get("/songs", func(c *fiber.Ctx) error {

		rows, err := db.Query("SELECT * FROM songs_table")
		if err != nil {
			return err
		}

		var songs []Song // Assuming Song is your struct for songs

		for rows.Next() {
			var song Song
			err := rows.Scan(&song.ID, &song.Title, &song.AlbumID, &song.SingerID, &song.CategoryID, &song.SongsDescription, &song.IsFavorite)
			if err != nil {
				return err
			}

			songs = append(songs, song)
		}

		return c.JSON(songs)

	})

	app.Get("/songs/:id", func(c *fiber.Ctx) error {

		id := c.Params("id")

		row := db.QueryRow("SELECT * FROM songs_table WHERE song_id = ?", id)

		var song Song

		err := row.Scan(&song.ID, &song.Title, &song.AlbumID, &song.SingerID, &song.CategoryID, &song.SongsDescription, &song.IsFavorite)

		if err != nil {

			return err

		}

		return c.JSON(song)

	})

	app.Post("/songs", func(c *fiber.Ctx) error {

		body := &Song{}

		if err := c.BodyParser(body); err != nil {

			return err

		}

		res, err := db.Exec("INSERT INTO songs_table (song_title, is_favorite, song_description, album_id, singer_id, category_id) VALUES (?, ?)", body.Title, body.IsFavorite, body.SongsDescription, body.AlbumID, body.SingerID, body.CategoryID)

		if err != nil {

			return err

		}

		id, err := res.LastInsertId()

		if err != nil {

			return err

		}

		body.ID = int(id)

		return c.JSON(body)
	})

	app.Put("/songs/:id", func(c *fiber.Ctx) error {

		id := c.Params("id")

		body := &Song{}

		if err := c.BodyParser(body); err != nil {

			return err

		}

		_, err := db.Exec("UPDATE songs_table SET song_title = ?, album_id = ?, singer_id = ?, category_id = ?, song_description = ?, is_favorite = ? WHERE songs_id = ?", body.Title, body.AlbumID, body.SingerID, body.CategoryID, body.SongsDescription, body.IsFavorite, id)

		if err != nil {

			return err

		}

		return c.SendStatus(fiber.StatusOK)

	})

	app.Delete("/songs/:id", func(c *fiber.Ctx) error {

		id := c.Params("id")

		_, err := db.Exec("DELETE FROM songs_table WHERE songs_id = ?", id)

		if err != nil {

			return err

		}

		return c.SendStatus(fiber.StatusOK)

	})

	fmt.Println("Listening on :" + os.Getenv("PORT"))

	app.Listen(":" + os.Getenv("PORT"))
}
