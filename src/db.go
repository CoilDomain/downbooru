package db

import (
	"gorm.io/gorm"
)

type Image struct {
	gorm.Model
	URL        string
	Downloaded bool
}

func db() {

}
