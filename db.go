package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Image database structure
type Image struct {
	URL        string `gorm:"primaryKey"`
	Downloaded bool
}

// File exists function
func fileExists(file string) bool {
	info, err := os.Stat(file)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// Database path
var usr, _ = user.Current()
var path = usr.HomeDir
var filename = "downbooru.db"
var databasepath = filepath.Join(path, filename)

// Configure connection to database
var db, _ = gorm.Open(sqlite.Open(databasepath), &gorm.Config{})

func dbinsert(ImageURL string) {
	// Create tables
	db.AutoMigrate(&Image{})
	// Test if database file exists, if not create
	if fileExists(databasepath) {
	} else {
		fmt.Println("Database does not exist, creating:")
		dbfile, _ := os.Create(databasepath)
		defer dbfile.Close()
		fmt.Println("Done")
	}
	// Input URLs into database
	db.Clauses(clause.OnConflict{DoNothing: true}).Create(&Image{
		URL:        ImageURL,
		Downloaded: false,
	})
}
