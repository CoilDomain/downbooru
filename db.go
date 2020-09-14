package main

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
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

// Add logging for troubleshooting
var newLogger = logger.New(
	log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
	logger.Config{
		SlowThreshold: time.Second,   // Slow SQL threshold
		LogLevel:      logger.Silent, // Log level
		Colorful:      false,         // Disable color
	},
)

// Database path
var usr, _ = user.Current()
var path = usr.HomeDir
var filename = "downbooru.db"
var databasepath = filepath.Join(path, filename)

// Configure connection to database
var db, _ = gorm.Open(sqlite.Open(databasepath), &gorm.Config{
	Logger: logger.Default.LogMode(logger.Silent),
})

// Query table for undownloaded images
func dbquery() {
	var image []Image
	db.Where(&Image{Downloaded: false}).Find(&image)
	length := len(image)
	for n := 1; n < length; n++ {
		fmt.Println(image[n].URL)
	}
}

// Insert scraped images into database
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
