package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	db        *gorm.DB
	err       error
	uploadDir = "./public/uploads"
)

// Model untuk objek Article
type Article struct {
	ID      uint   `json:"id" gorm:"primaryKey" `
	Title   string `json:"title" binding:"required" form:"title"`
	Content string `json:"content" binding:"required" form:"content"`
	Banner  string `json:"banner" binding:"" form:""`
}

func connectDB() {
	// Membuat koneksi ke SQLite
	db, err = gorm.Open(sqlite.Open("golangcrud.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Tidak bisa terhubung database")
	}
}

func init() {
	// Koneksi ke MySQL Database
	connectDB()
}

func main() {

	// Migrasi tabel
	db.AutoMigrate(&Article{})

	// Menginisialisasi router
	router := gin.Default()
	router.LoadHTMLGlob("views/*")

	router.Static("/public", "./public")
	router.Static("/assets", "./public/assets")
	router.StaticFS("/uploads", http.Dir("public/uploads"))
	router.StaticFile("/favicon.ico", "./public/favicon.ico")

	// router.Static("/", "./public")
	// Menambahkan middleware untuk penanganan kesalahan validasi
	router.Use(gin.ErrorLogger(), gin.Recovery())

	router.GET("/", pageHome)

	routerBlog := router.Group("/blogs")
	routerBlog.GET("/", pageBlogs)
	routerBlog.GET("/:id", pageBlog)

	// Rute API
	// API Version1 endpoints and routes
	api := router.Group("/api/v1")
	{
		apiuser := api.Group("/articles")

		apiuser.GET("/", getArticles)
		apiuser.GET("/:id", getArticle)
		apiuser.POST("/", createArticle)
		apiuser.PUT("/:id", updateArticle)
		apiuser.DELETE("/:id", deleteArticle)
	}

	// Menjalankan server pada port 8080
	router.Run(":8080")
}

// handle page home
func pageHome(c *gin.Context) {
	c.HTML(
		200,
		"index.html",
		gin.H{
			"title": "Golang CRUD",
			"nav":   "home",
		},
	)
}

// handle blogs page
func pageBlogs(c *gin.Context) {
	c.HTML(
		200,
		"blogs.html",
		gin.H{
			"title": "Golang CRUD",
			"nav":   "blog",
		},
	)
}

// handle blogs page
func pageBlog(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format", "id": id})
		return
	}

	c.HTML(
		200,
		"blog.html",
		gin.H{
			"title": "Golang CRUD",
			"nav":   "blog",
			"id":    id,
		},
	)
}

// Handler untuk mendapatkan semua article
func getArticles(c *gin.Context) {
	var articles []Article
	db.Order("id desc").Find(&articles)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    articles,
		"message": "Artikel berhasil didapatkan",
	})
}

// Handler untuk mendapatkan article berdasarkan ID
func getArticle(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var article Article
	if err := db.First(&article, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    article,
		"message": "Artikel ditemukan",
	})
}

func (a *Article) handleUpload(c *gin.Context) error {
	// Menerima file gambar dari form-data
	file, err := c.FormFile("banner_file")
	if err != nil {
		return err
	}

	// Membuat direktori penyimpanan gambar (jika belum ada)

	checkDir(uploadDir)

	// Menghasilkan nama file unik
	filename := GenerateUniqueFileName(file.Filename)

	// Menyimpan file di direktori uploads
	err = c.SaveUploadedFile(file, filepath.Join(uploadDir, filename))
	if err != nil {
		return err
	} else {
		/* if a.Banner != "" {
			if _, err := os.Stat(filepath.Join(uploadDir, a.Banner)); err == nil {
				// path/to/whatever exists
				err = os.Remove(filepath.Join(uploadDir, a.Banner))
				if err != nil {
					return err
				}
			} else if errors.Is(err, os.ErrNotExist) {
				// path/to/whatever does *not* exist
			}
		} */
		bannerPath := filepath.Join(uploadDir, a.Banner)
		err = a.deleteImage(bannerPath)
		a.Banner = filename
		return nil
	}
}

func (a *Article) deleteImage(pathUpload string) error {
	if a.Banner != "" {
		if _, err := os.Stat(pathUpload); err == nil {
			// path/to/whatever exists
			err = os.Remove(pathUpload)
			if err != nil {
				return err
			}
		} else if errors.Is(err, os.ErrNotExist) {
			// path/to/whatever does *not* exist
		}
	}
	return nil
}

// Handler untuk membuat article baru
func createArticle(c *gin.Context) {
	var input Article
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Menerima file gambar dari form-data
	/* file, err := c.FormFile("banner_file")
	if err == nil {
		// Membuat direktori penyimpanan gambar (jika belum ada)
		uploadDir := "./public/uploads"
		checkDir(uploadDir)

		// Menghasilkan nama file unik
		filename := GenerateUniqueFileName(file.Filename)

		// Menyimpan file di direktori uploads
		err = c.SaveUploadedFile(file, filepath.Join(uploadDir, filename))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
			return
		} else {
			input.Banner = filename
		}
	} */

	_ = input.handleUpload(c)

	db.Create(&input)
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    input,
		"message": "Artikel baru berhasil disimpan",
	})
}

// Handler untuk mengupdate article berdasarkan ID
func updateArticle(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var article Article
	if err := db.First(&article, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
		return
	}

	var input Article
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_ = article.handleUpload(c)

	/* // Menerima file gambar dari form-data
	file, err = c.FormFile("banner_file")
	if err == nil {

		// Membuat direktori penyimpanan gambar (jika belum ada)
		uploadDir := "./public/uploads"
		checkDir(uploadDir)

		// Menghasilkan nama file unik
		filename := GenerateUniqueFileName(file.Filename)

		// Menyimpan file di direktori uploads
		err = c.SaveUploadedFile(file, filepath.Join(uploadDir, filename))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
			return
		} else {
			input.Banner = filename
			article.Banner = input.Banner
		}
	}*/

	// Update nilai article
	article.Title = input.Title
	article.Content = input.Content
	db.Save(&article)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    article,
		"message": "Artikel berhasil diupdate",
	})
}

// Handler untuk menghapus article berdasarkan ID
func deleteArticle(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var article Article
	if err := db.First(&article, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
		return
	}

	bannerPath := filepath.Join(uploadDir, article.Banner)
	article.deleteImage(bannerPath)

	db.Delete(&article)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    [0]string{},
		"message": "Artikel berhasil dihapus",
	})
}

// GenerateUniqueFileName menghasilkan nama file unik dengan menambahkan timestamp
func GenerateUniqueFileName(originalFilename string) string {
	ext := filepath.Ext(originalFilename)
	filename := strings.TrimSuffix(originalFilename, ext)

	// Menambahkan timestamp sebagai bagian dari nama file
	timestamp := fmt.Sprintf("%d", NowUnixNano())
	return fmt.Sprintf("%s_%s%s", filename, timestamp, ext)
}

// NowUnixNano mengembalikan waktu sekarang dalam format UnixNano
func NowUnixNano() int64 {
	return time.Now().UnixNano()
}

// check directory
func checkDir(uploadDir string) {
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.Mkdir(uploadDir, os.ModePerm)
	}
}
