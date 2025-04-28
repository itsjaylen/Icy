package pastebin

import (
	"net/http"

	"github.com/itsjaylen/IcyAPI/internal/api/services/pastebin"
	"github.com/itsjaylen/IcyAPI/internal/appinit"
)

func RegisterRoutes(mux *http.ServeMux, app *appinit.App) {
	controller := pastebin.NewPasteBinController(
		app.PostgresClient,
		app.MinioClient,
	)

	// Register the pastebin routes
	mux.HandleFunc("/paste", controller.CreatePaste)               // POST /paste
	mux.HandleFunc("/paste/", controller.GetPaste)                 // GET /paste/{id}
	mux.HandleFunc("/paste/update/", controller.UpdatePaste)       // PUT/PATCH /paste/update/{id}
	mux.HandleFunc("/paste/delete/", controller.DeletePaste)       // DELETE /paste/delete/{id}?token={token}
	mux.HandleFunc("/paste/upload", controller.UploadImageHandler) // POST /paste/upload
}
