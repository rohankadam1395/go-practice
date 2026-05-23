package httpserver

import (
	"go-practice/internal/album"
	"net/http"

	"github.com/gin-gonic/gin"
)

func getAlbums(c *gin.Context, store album.Store) {
	albums, err := store.ListAlbums()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, albums)
}

func postAlbums(c *gin.Context, store album.Store) {
	var newAlbum album.Album
	if err := c.BindJSON(&newAlbum); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := store.CreateAlbum(&newAlbum)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

func getAlbumByID(c *gin.Context, store album.Store) {
	var id = c.Param("id")
	alb, err := store.GetAlbumByID(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, alb)
}

func updateAlbum(c *gin.Context, store album.Store) {
	var album album.Album
	if err := c.BindJSON(&album); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := store.UpdateAlbum(&album)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, album)
}

func deleteAlbum(c *gin.Context, store album.Store) {
	id := c.Param("id")
	err := store.DeleteAlbum(id)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Album deleted"})
}
