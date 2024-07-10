package config

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Dbconnect() *gorm.DB {
	destinesion := "host=aws-0-ap-southeast-1.pooler.supabase.com user=postgres.btjyjwufuhqrzhlqpoie password=@Mfatur14jr dbname=postgres port=5432 TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(destinesion), &gorm.Config{})
	if err != nil {
		panic("error connection database")
	}

	return db
}
