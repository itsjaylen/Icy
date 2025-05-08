package pastebin

import (
	"net/http"
	"time"

	"github.com/itsjaylen/IcyAPI/internal/api/middleware"
	"github.com/itsjaylen/IcyAPI/internal/api/services/pastebin"
	"github.com/itsjaylen/IcyAPI/internal/appinit"
)

func RegisterRoutes(mux *http.ServeMux, app *appinit.App) {
	controller := pastebin.NewPasteBinController(
		app.PostgresClient,
		app.MinioClient,
	)

	// Register the pastebin routes
	mux.HandleFunc("/paste", middleware.RateLimitMiddleware(controller.CreatePaste, 5*time.Second, 3))               // POST /paste
	mux.HandleFunc("/paste/", middleware.RateLimitMiddleware(controller.GetPaste, 5*time.Second, 3))                 // GET /paste/{id}
	mux.HandleFunc("/paste/update/", middleware.RateLimitMiddleware(controller.UpdatePaste, 5*time.Second, 3))       // PUT/PATCH /paste/update/{id}
	mux.HandleFunc("/paste/delete/", middleware.RateLimitMiddleware(controller.DeletePaste, 5*time.Second, 3))       // DELETE /paste/delete/{id}?token={token}
	mux.HandleFunc("/paste/upload", middleware.RateLimitMiddleware(controller.UploadImageHandler, 5*time.Second, 3)) // POST /paste/upload
}
