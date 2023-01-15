package database

import (
	"fmt"
	"server/models"
	"server/pkg/mysql"
)

func RunMigration() {
	err := mysql.DB.AutoMigrate(
		&models.User{},
		&models.Donation{},
		&models.Funding{},
	)

	if err != nil {
		fmt.Println(err)
		panic("Migration Failed")
	}

	fmt.Println("Migration Successfully !!")
}
