package postgres

import (
	"context"
	"fmt"
	"github.com/IKostarev/yandex-go-dev/internal/logger"
	"github.com/IKostarev/yandex-go-dev/internal/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	id "github.com/vgarvardt/pgx-google-uuid/v5"
	"strconv"
	"strings"
	"time"
)

type DB struct {
	db *pgxpool.Pool
}

func NewPostgresDB(addrConn string) (*DB, error) {
	cfg, err := pgxpool.ParseConfig(addrConn)
	if err != nil {
		return nil, fmt.Errorf("error parse config: %w", err)
	}

	cfg.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		id.Register(conn.TypeMap())
		return nil
	}

	db, err := pgxpool.NewWithConfig(context.Background(), cfg)
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

	return psql, nil
}

func (psql *DB) Save(longURL, corrID string, cookie string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	var count string

	shortURL := utils.RandomString()

	s := psql.db.QueryRow(ctx, `SELECT COUNT(*) FROM yandex`)

	err := s.Scan(&count)
	if err != nil {
		logger.Errorf("error in Scan count in SELECT query: %s", err)
	}

	if corrID == "" {
		corrID = shortURL
	}

	_, err = psql.db.Exec(ctx, `INSERT INTO yandex (id, longurl, shorturl, correlation, cookie, deleted) VALUES ($1, $2, $3, $4, $5, $6);`, count, longURL, shortURL, corrID, cookie, false)
	if err != nil {
		return "", fmt.Errorf("error is INSERT data in database: %w", err)
	}

	return shortURL, nil
}

func (psql *DB) Get(shortURL, corrID string, _ string) (string, string) {
	var longURL string

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	row := psql.db.QueryRow(ctx, `SELECT longurl FROM yandex WHERE shorturl = $1`, shortURL)

	err := row.Scan(&longURL)
	if err != nil {
		logger.Errorf("error in Scan longURL in SELECT query: %s", err)
	}

	fmt.Println("longURL = ", &longURL)

	return longURL, corrID
}

func (psql *DB) IsDel(shortURL string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	var isDel bool

	row := psql.db.QueryRow(ctx, `SELECT deleted FROM yandex WHERE shorturl = $1`, shortURL)
	_ = row.Scan(&isDel)

	return isDel
}

func (psql *DB) DeleteURL(shortURLs []string, _ string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	values := make([]interface{}, len(shortURLs)*2)
	for i, shortURL := range shortURLs {
		values[i*2] = true
		values[i*2+1] = shortURL
	}

	query := `UPDATE yandex SET deleted = $1 WHERE shorturl IN (` + placeholders(len(shortURLs), 2) + `);`
	_, err := psql.db.Exec(ctx, query, values...)

	return err == nil
}

func placeholders(count, offset int) string {
	ph := make([]string, count)
	for i := range ph {
		ph[i] = "$" + strconv.Itoa(i+offset)
	}
	return strings.Join(ph, ", ")
}

func (psql *DB) Close() error {
	psql.db.Close()
	return nil
}

func (psql *DB) createTable() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	_, err := psql.db.Exec(ctx,
		`CREATE TABLE IF NOT EXISTS yandex (
    		id SERIAL PRIMARY KEY,
   			longurl VARCHAR(255) NOT NULL,
    		shorturl VARCHAR(255) NOT NULL,
    		cookie VARCHAR(255) NOT NULL,
    		deleted bool,
   			correlation VARCHAR(255) NOT NULL);`)

	return err
}

func (psql *DB) checkIsTablesExists() (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	row := psql.db.QueryRow(ctx, `SELECT shorturl FROM yandex WHERE longurl = $1`, longURL)

	var res string

	err := row.Scan(&res)
	if err != nil {
		return "", fmt.Errorf("error in Scan res in SELECT query: %w", err)
	}

	return res, nil
}

func (psql *DB) GetAllURLs(cookie string) ([]string, string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	rows, err := psql.db.Query(ctx, `SELECT longurl, shorturl FROM yandex WHERE cookie = $1`, cookie)
	if err != nil {
		// Обработка ошибки запроса
		return nil, ""
	}
	defer rows.Close()

	var res []string

	for rows.Next() {
		var longURL, shortURL string
		err := rows.Scan(&longURL, &shortURL)
		if err != nil {
			// Обработка ошибки сканирования результатов запроса
			return nil, ""
		}
		res = append(res, longURL)
	}

	if err := rows.Err(); err != nil {
		// Обработка ошибки после итерации по результатам запроса
		return nil, ""
	}

	return res, ""
}

func (psql *DB) GetAllShortURLs(cookie string) ([]string, string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	rows, err := psql.db.Query(ctx, `SELECT shorturl, longurl FROM yandex WHERE cookie = $1`, cookie)
	if err != nil {
		// Обработка ошибки запроса
		return nil, ""
	}
	defer rows.Close()

	var res []string

	for rows.Next() {
		var longURL, shortURL string
		err := rows.Scan(&longURL, &shortURL)
		if err != nil {
			// Обработка ошибки сканирования результатов запроса
			return nil, ""
		}
		res = append(res, longURL)
	}

	if err := rows.Err(); err != nil {
		// Обработка ошибки после итерации по результатам запроса
		return nil, ""
	}

	return res, ""
}

func (psql *DB) Ping() bool {
	return psql.db.Ping(context.Background()) == nil
}
