package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var pool *pgxpool.Pool

func Connect(dbURL string) *pgxpool.Pool {
	var err error
	pool, err = pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	err = pool.Ping(context.Background())
	if err != nil {
		panic("Could not connect to the database")
	}

	log.Println("Connected to database")


	return pool
}

func ConnectAdmin(dbURL string) *pgxpool.Pool {
	var err error
	p, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	err = pool.Ping(context.Background())
	if err != nil {
		panic("Could not connect to the database")
	}

	log.Println("Connected to database")


	return p
}

func GetConnection() *pgxpool.Pool {
	if pool != nil {
		return pool
	}

	panic("Database connection not established")
}
