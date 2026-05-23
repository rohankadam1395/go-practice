package mysql

import (
	"database/sql"
	"go-practice/internal/album"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) album.Store {
	return &Store{db: db}
}

func (s *Store) GetAlbumByID(id string) (albumRes *album.Album, err error) {
	row := s.db.QueryRow("SELECT id, title, artist, price FROM album WHERE id=?", id)
	album := album.Album{}
	if err := row.Scan(&album.ID, &album.Title, &album.Artist, &album.Price); err != nil {
		return nil, err
	}
	return &album, nil
}

func (s *Store) CreateAlbum(album *album.Album) error {
	_, err := s.db.Exec("INSERT INTO album (id, title, artist, price) VALUES (?, ?, ?, ?)", album.ID, album.Title, album.Artist, album.Price)
	return err
}

func (s *Store) UpdateAlbum(album *album.Album) error {
	_, err := s.db.Exec("UPDATE album SET title=?, artist=?, price=? WHERE id=?", album.Title, album.Artist, album.Price, album.ID)
	return err
}

func (s *Store) DeleteAlbum(id string) error {
	_, err := s.db.Exec("DELETE FROM album WHERE id=?", id)
	return err
}

func (s *Store) ListAlbums() ([]*album.Album, error) {
	rows, err := s.db.Query("SELECT id, title, artist, price FROM album")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	albums := make([]*album.Album, 0)
	for rows.Next() {
		album := album.Album{}
		if err := rows.Scan(&album.ID, &album.Title, &album.Artist, &album.Price); err != nil {
			return nil, err
		}
		albums = append(albums, &album)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return albums, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}
