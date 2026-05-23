package album

type Store interface {
	GetAlbumByID(id string) (*Album, error)
	CreateAlbum(album *Album) error
	UpdateAlbum(album *Album) error
	DeleteAlbum(id string) error
	ListAlbums() ([]*Album, error)
	Close() error
}
