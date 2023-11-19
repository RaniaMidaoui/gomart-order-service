package db

import (
	"log"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/RaniaMidaoui/goMart-order-service/pkg/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Handler struct {
	DB *gorm.DB
}

func Init(url string) Handler {
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	db.AutoMigrate(&models.Order{})

	return Handler{db}
}

func Mock() Handler {
	mockDb, mock, _ := sqlmock.New()

	dialector := postgres.New(postgres.Config{
		DriverName: "postgres",
		Conn:       mockDb,
	})

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	mock.ExpectQuery(`SELECT \* FROM "products" WHERE "products"."id" = \$1 ORDER BY "products"."id" LIMIT 1`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "stock", "price"}).
			AddRow(1, "Prod A", 20, 15))

	mock.ExpectQuery(`SELECT \* FROM "products" WHERE "products"."id" = \$1 ORDER BY "products"."id" LIMIT 1`).
		WithArgs(2).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "stock", "price"}).
			AddRow(2, "Prod B", 10, 5))

	return Handler{db}
}
