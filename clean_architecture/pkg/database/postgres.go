package database

import (
	"database/sql"
	"time"
	// Sürücüyü Go'nun SQL motoruna gizlice (Blank Identifier) kaydediyoruz.
	_ "github.com/jackc/pgx/v5/stdlib"
)

// NewPostgresConnection driverName : pxg olmalı
// "postgres://kullanici:sifre@host:port/veritabani_adi?sslmode=disable" -> connectionString
func NewPostgresConnection(driverName string, dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	//db.SetConnMaxIdleTime(5 * time.Minute) max iddle connection time hangi durumda ne kadar kullanılmalı ?
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}
