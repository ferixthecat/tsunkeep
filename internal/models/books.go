package models

import (
	"database/sql"
	"errors"
	"time"
)

type Book struct {
	ID     int
	Title  string
	Author string
	Genres string
	Pages  int
	Added  time.Time
	UserID int
	Status string
}

type BookModel struct {
	DB *sql.DB
}

func (m *BookModel) Insert(title, author, genres string, pages, userID int) (int, error) {
	stmt := `INSERT INTO books (title, author, genres, pages, added, user_id)
	VALUES(?, ?, ?, ?, UTC_TIMESTAMP(), ?)`

	result, err := m.DB.Exec(stmt, title, author, genres, pages, userID)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *BookModel) Get(userID, bookID int) (Book, error) {
	stmt := `SELECT id, title, author, genres, pages, added, user_id, status FROM books WHERE user_id = ? AND id = ?`

	row := m.DB.QueryRow(stmt, userID, bookID)

	var b Book

	err := row.Scan(&b.ID, &b.Title, &b.Author, &b.Genres, &b.Pages, &b.Added, &b.UserID, &b.Status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Book{}, ErrNoRecord
		} else {
			return Book{}, err
		}
	}

	return b, nil
}

func (m *BookModel) Unfinished(userID int) ([]Book, error) {
	stmt := `SELECT id, title, author, genres, pages, added, user_id, status FROM books WHERE user_id = ? AND status != 'finished' ORDER BY id`

	rows, err := m.DB.Query(stmt, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var books []Book

	for rows.Next() {
		var b Book
		err = rows.Scan(&b.ID, &b.Title, &b.Author, &b.Genres, &b.Pages, &b.Added, &b.UserID, &b.Status)
		if err != nil {
			return nil, err
		}

		books = append(books, b)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}

func (m *BookModel) Finished(userID int) ([]Book, error) {
	stmt := `SELECT id, title, author, genres, pages, added, user_id, status FROM books WHERE user_id = ? AND status = 'finished' ORDER BY id`

	rows, err := m.DB.Query(stmt, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var books []Book

	for rows.Next() {
		var b Book
		err = rows.Scan(&b.ID, &b.Title, &b.Author, &b.Genres, &b.Pages, &b.Added, &b.UserID, &b.Status)
		if err != nil {
			return nil, err
		}

		books = append(books, b)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}

func (m *BookModel) UpdateStatus(userID, bookID int, newStatus string) error {
	stmt := `UPDATE books SET status = ? WHERE id = ? AND user_id = ?`

	_, err := m.DB.Exec(stmt, newStatus, bookID, userID)
	if err != nil {
		return err
	}

	return nil
}
