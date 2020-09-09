package main

import (
	"os/user"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Image database structure
type Image struct {
	gorm.Model
	URL        string
	Downloaded bool
}

func query(ImageURL string) {
	// Create Database
	usr, _ := user.Current()
	path := usr.HomeDir
	filename := "downbooru.db"
	databasepath := filepath.Join(path, filename)
	// Configure connection to database
	db, _ := gorm.Open(sqlite.Open(databasepath), &gorm.Config{})
	// Create tables
	db.AutoMigrate(&Image{})
	// Input URLs into database
	db.Create(&Image{
		URL:        ImageURL,
		Downloaded: true,
	})
}
