package postgres

import (
	"context"
	"fmt"
	"github.com/IKostarev/yandex-go-dev/internal/logger"
	"github.com/IKostarev/yandex-go-dev/internal/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	uuid "github.com/vgarvardt/pgx-google-uuid/v5"
	"time"
)

type DB struct {
	db *pgxpool.Pool
}

var Count uint64

func NewPostgresDB(addrConn string) (*DB, error) {
	conn, err := pgxpool.ParseConfig(addrConn)
	if err != nil {
		return nil, fmt.Errorf("error parse config: %w", err)
	}

	conn.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		uuid.Register(conn.TypeMap())
		return nil
	}

	db, err := pgxpool.NewWithConfig(context.Background(), conn)
	if err != nil {
		return nil, fmt.Errorf("error new config: %w", err)
	}

	psql := &DB{db: db}

	exists, err := psql.checkIsTablesExists()
	if err != nil {
		return nil, fmt.Errorf("error check is table exists: %w", err)
	}

	if !exists {
		err = psql.createTable()
		if err != nil {
			return nil, fmt.Errorf("error create table: %w", err)
		}
	}

	Count++

	return psql, nil
}

func (psql *DB) Save(longURL, corrID string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var count string

	shortURL := utils.RandomString()

	sh, err := psql.CheckIsURLExists(longURL)
	if err != nil {
		logger.Errorf("error in Check Is URL Exists: %s", err)
	}

	if sh != "" {
		err := psql.UpdateDB(sh, longURL, corrID)
		if err != nil {
			logger.Errorf("error update DB: %s", err)
		}

		return sh, nil
	}

	s := psql.db.QueryRow(ctx, `SELECT COUNT(*) FROM yandex`)

	_ = s.Scan(&count)

	if corrID == "" {
		corrID = shortURL
	}

	_, err = psql.db.Exec(ctx, `INSERT INTO yandex (id, longurl, shorturl, correlation) VALUES ($1, $2, $3, $4);`, count, longURL, shortURL, corrID)
	if err != nil {
		return "", fmt.Errorf("error is INSERT data in database: %w", err)
	}

	return shortURL, nil
}

func (psql *DB) Get(shortURL, corrID string) (string, string) {
	var longURL string

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if corrID == "" {
		row := psql.db.QueryRow(ctx, `SELECT longurl FROM yandex WHERE shorturl = $1`, shortURL)

		_ = row.Scan(&longURL)

		return longURL, corrID
	} else {
		row := psql.db.QueryRow(ctx, `SELECT longurl FROM yandex WHERE shorturl = $1 AND correlation = $2`, shortURL, corrID)

		_ = row.Scan(&longURL)

		return longURL, corrID
	}
}

func (psql *DB) Close() error {
	//if err := psql.db.Close; err != nil {
	//	logger.Errorf("error close db: %s", err)
	//}

	return nil //TODO заглушка на будущее, кажется что в бд этот метод вообще не нужен
}
func (psql *DB) createTable() error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, err := psql.db.Exec(ctx,
		`CREATE TABLE IF NOT EXISTS yandex (
    		id SERIAL PRIMARY KEY,
   			longurl VARCHAR(255) NOT NULL,
    		shorturl VARCHAR(255) NOT NULL,
   			correlation VARCHAR(255) NOT NULL);`)

	return err
}

func (psql *DB) checkIsTablesExists() (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	row := psql.db.QueryRow(ctx, `SELECT EXISTS (SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'yandex')`)

	var res bool

	err := row.Scan(&res)
	if err != nil {
		return false, fmt.Errorf("error scan: %w", err)
	}

	return res, nil
}

func (psql *DB) CheckIsURLExists(longURL string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	row := psql.db.QueryRow(ctx, `SELECT shorturl FROM yandex WHERE longurl = $1`, longURL)

	var res string

	_ = row.Scan(&res)

	return res, nil
}

func (psql *DB) UpdateDB(shortURL, longURL, corrID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, err := psql.db.Exec(ctx, `UPDATE yandex SET shorturl = $1 WHERE longurl = $2 OR correlation = $3;`, shortURL, longURL, corrID)
	if err != nil {
		return fmt.Errorf("error is INSERT data in database: %w", err)
	}

	return nil
}

func (psql *DB) Ping() bool {
	return psql.db.Ping(context.Background()) == nil
}
