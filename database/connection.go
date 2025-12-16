package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq" // Driver PostgreSQL
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	PostgresDB *sql.DB
	MongoDB    *mongo.Database
)

func ConnectPostgres() {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("❌ Error: DB_DSN tidak ditemukan di .env")
	}

	var err error
	PostgresDB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("❌ Gagal membuka driver Postgres:", err)
	}

	// Tes Ping
	err = PostgresDB.Ping()
	if err != nil {
		log.Fatal("❌ Gagal koneksi ke PostgreSQL:", err)
	}

	fmt.Println("✅ PostgreSQL Connected Successfully!")
}

func ConnectMongo() {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		log.Fatal("❌ Error: MONGO_URI tidak ditemukan di .env")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("❌ Gagal membuat client Mongo:", err)
	}

	// Tes Ping
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("❌ Gagal ping MongoDB:", err)
	}

	// Set database global
	// Nama database diambil dari URI atau default
	MongoDB = client.Database("pelaporan-prestasi")
	fmt.Println("✅ MongoDB Connected Successfully!")
}