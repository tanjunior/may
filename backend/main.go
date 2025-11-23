package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

// respondSuccess sends a consistent success JSON envelope.
func respondSuccess(c *gin.Context, status int, data interface{}, meta map[string]interface{}) {
	if meta == nil {
		c.JSON(status, gin.H{
			"success": true,
			"status":  status,
			"data":    data,
		})
		return
	}
	c.JSON(status, gin.H{
		"success": true,
		"status":  status,
		"data":    data,
		"meta":    meta,
	})
}

// respondError sends a consistent error JSON envelope.
func respondError(c *gin.Context, status int, code string, message string, details interface{}) {
	c.JSON(status, gin.H{
		"success": false,
		"status":  status,
		"error": gin.H{
			"code":    code,
			"message": message,
			"details": details,
		},
	})
}

func main() {
	// Try to load environment variables from the repo-level `.env` (one dir up),
	// then fall back to a `.env` in the current working dir.
	if err := godotenv.Load("../.env"); err != nil {
		if err2 := godotenv.Load(); err2 != nil {
			log.Printf(".env not found (tried ../.env and .env): %v; %v", err, err2)
		}
	}

	// Run the frontend generator in non-release (development) mode so
	// TypeScript types stay in sync during development. In release mode
	// we skip generation to avoid requiring a Go toolchain at runtime.
	if gin.Mode() != gin.ReleaseMode {
		if err := runGenerator(); err != nil {
			log.Printf("warning: generator failed: %v", err)
		}
	}

	// Initialize DB after loading env.
	database = db()

	// Create router and start server
	r := newRouter()
	r.Run()
}

// newRouter sets up and returns the Gin engine with routes (useful for tests).
func newRouter() *gin.Engine {
	r := gin.Default()
	r.Use(cors.Default())

	// maxPerPage can be configured via env var `MAX_PER_PAGE` (defaults to 100)
	maxPerPage := 100
	if v := os.Getenv("MAX_PER_PAGE"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			maxPerPage = n
		}
	}

	r.GET("/ping", func(c *gin.Context) {
		respondSuccess(c, http.StatusOK, gin.H{"message": "pong"}, nil)
	})

	r.GET("/products", func(c *gin.Context) {
		// Pagination parameters
		pageStr := c.DefaultQuery("page", "1")
		perPageStr := c.DefaultQuery("per_page", "20")
		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			page = 1
		}
		perPage, err := strconv.Atoi(perPageStr)
		if err != nil || perPage < 1 {
			perPage = 20
		}

		if perPage > maxPerPage {
			RespondBadRequest(c, CodePerPageTooLarge, map[string]interface{}{"requested": perPage, "max_per_page": maxPerPage})
			return
		}

		products, total, err := getAllProducts(page, perPage)
		if err != nil {
			RespondInternal(c, CodeInternalError, err.Error())
			return
		}

		totalPages := 0
		if total > 0 {
			totalPages = int((total + int64(perPage) - 1) / int64(perPage))
		}

		meta := map[string]interface{}{
			"page":        page,
			"per_page":    perPage,
			"total":       total,
			"total_pages": totalPages,
		}

		respondSuccess(c, http.StatusOK, products, meta)
	})

	r.GET("/product/latest", func(c *gin.Context) {
		var product, err = getLatestProduct()
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				RespondNotFound(c, CodeProductNotFound, nil)
				return
			}
			RespondInternal(c, CodeInternalError, err.Error())
			return
		}
		respondSuccess(c, http.StatusOK, product, nil)
	})

	r.GET("/product/:id", func(c *gin.Context) {
		idParam := c.Param("id")
		product, err := getProductByID(idParam)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				RespondNotFound(c, CodeProductNotFound, nil)
				return
			}
			RespondInternal(c, CodeInternalError, err.Error())
			return
		}
		respondSuccess(c, http.StatusOK, product, nil)
	})

	r.POST("/product", func(c *gin.Context) {
		var json struct {
			Code  string `json:"code" binding:"required"`
			Price uint   `json:"price" binding:"required"`
		}

		if err := c.ShouldBindJSON(&json); err != nil {
			RespondBadRequest(c, CodeInvalidRequest, err.Error())
			return
		}

		created, err := addProduct(json.Code, json.Price)
		if err != nil {
			RespondInternal(c, CodeInternalError, err.Error())
			return
		}

		// Set Location header for the created resource
		c.Header("Location", fmt.Sprintf("/product/%d", created.ID))
		respondSuccess(c, http.StatusCreated, created, nil)
	})

	r.PUT("/product/:id", func(c *gin.Context) {
		idParam := c.Param("id")
		var json struct {
			Code  string `json:"code" binding:"required"`
			Price uint   `json:"price" binding:"required"`
		}

		if err := c.ShouldBindJSON(&json); err != nil {
			RespondBadRequest(c, CodeInvalidRequest, err.Error())
			return
		}

		var id uint
		_, err := fmt.Sscanf(idParam, "%d", &id)
		if err != nil {
			RespondBadRequest(c, CodeInvalidID, nil)
			return
		}

		updated, err := updateProduct(id, json.Code, json.Price)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				RespondNotFound(c, CodeProductNotFound, nil)
				return
			}
			RespondInternal(c, CodeInternalError, err.Error())
			return
		}

		respondSuccess(c, http.StatusOK, updated, nil)
	})

	r.DELETE("/product/:id", func(c *gin.Context) {
		idParam := c.Param("id")

		var id uint
		_, err := fmt.Sscanf(idParam, "%d", &id)
		if err != nil {
			RespondBadRequest(c, CodeInvalidID, nil)
			return
		}

		if err := deleteProduct(id); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				RespondNotFound(c, CodeProductNotFound, nil)
				return
			}
			RespondInternal(c, CodeInternalError, err.Error())
			return
		}

		respondSuccess(c, http.StatusOK, gin.H{"message": "product deleted"}, nil)
	})

	return r
}

// runGenerator runs the small Go CLI that emits TypeScript types into the
// frontend source tree. It intentionally logs output and returns an error
// if the generator fails; callers can decide how to handle the error.
func runGenerator() error {
	// Run `go run ./cmd/genfrontend -out ../frontend/src` from backend dir
	cmd := exec.Command("go", "run", "./cmd/genfrontend", "-out", "../frontend/src")
	cmd.Env = os.Environ()
	// keep working dir as backend (where main.go lives)
	out, err := cmd.CombinedOutput()
	if len(out) > 0 {
		log.Printf("generator output: %s", string(out))
	}
	if err != nil {
		return fmt.Errorf("generator failed: %w", err)
	}
	return nil
}
