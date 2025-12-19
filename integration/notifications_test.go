package integration

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"SpaceBookProject/internal/auth"
	"SpaceBookProject/internal/config"
	"SpaceBookProject/internal/db"
	"SpaceBookProject/internal/domain"
	"SpaceBookProject/internal/handlers"
	"SpaceBookProject/internal/repository"
	"SpaceBookProject/internal/services"
	"SpaceBookProject/middleware"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

)

func setupTestServer(t *testing.T) (*gin.Engine, *sql.DB, string) {
	t.Helper()

	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatal(err)
	}

	database, err := db.InitDB(&cfg.Database)
	if err != nil {
		t.Fatal(err)
	}

	jwtManager := auth.NewJWTManager(cfg.JWT.SecretKey)

	notificationRepo := repository.NewNotificationRepository(database)
	notificationService := services.NewNotificationService(notificationRepo)
	notificationHandler := handlers.NewNotificationHandler(notificationService)

	r := gin.New()
	r.Use(gin.Recovery())

	api := r.Group(cfg.API.Prefix + "/" + cfg.API.Version)
	api.GET(
		"/notifications",
		middleware.AuthMiddleware(jwtManager),
		notificationHandler.List,
	)

	// создаём тестовый JWT
	token, err := jwtManager.GenerateAccessToken(1, "test@example.com", string(domain.RoleTenant))
	if err != nil {
    	t.Fatal(err)
	}

	if err != nil {
		t.Fatal(err)
	}

	return r, database, token
}

func TestGetNotifications_Integration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r, database, token := setupTestServer(t)
	defer database.Close()

	// clean up данных
	_, err := database.Exec(`DELETE FROM notifications`)
	if err != nil {
		t.Fatal(err)
	}

	// seed данных
	_, err = database.Exec(`
		INSERT INTO notifications (user_id, type, message)
		VALUES (1, 'booking_created', 'Your booking was created')
	`)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/notifications",
		bytes.NewBuffer(nil),
	)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", w.Code, w.Body.String())
	}

	var resp struct {
		Items []map[string]interface{} `json:"items"`
	}

	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}

	if len(resp.Items) != 1 {
		t.Fatalf("expected 1 notification, got %d", len(resp.Items))
	}
}

