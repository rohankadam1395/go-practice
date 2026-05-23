package memory

import (
	"fmt"
	"go-practice/internal/album"
	"sync"
)

type Store struct {
	albums map[string]*album.Album
	mu     sync.RWMutex
}

func NewStore(seedAlbums []album.Album) album.Store {
	albums := make(map[string]*album.Album)
	for _, album := range seedAlbums {
		albums[album.ID] = &album
	}
	return &Store{albums: albums}
}

func (s *Store) GetAlbumByID(id string) (*album.Album, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	album, ok := s.albums[id]
	if !ok {
		return nil, fmt.Errorf("album not found")
	}
	return album, nil
}

func (s *Store) CreateAlbum(album *album.Album) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.albums[album.ID]; ok {
		return fmt.Errorf("album already exists")
	}
	s.albums[album.ID] = album
	return nil
}

func (s *Store) DeleteAlbum(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.albums[id]
	if !ok {
		return fmt.Errorf("album not found")
	}
	delete(s.albums, id)
	return nil
}

func (s *Store) UpdateAlbum(album *album.Album) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.albums[album.ID]
	if !ok {
		return fmt.Errorf("album not found")
	}
	s.albums[album.ID] = album
	return nil
}

func (s *Store) ListAlbums() ([]*album.Album, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	albums := make([]*album.Album, 0, len(s.albums))
	for _, album := range s.albums {
		albums = append(albums, album)
	}
	return albums, nil
}

func (s *Store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.albums = nil
	return nil
}
