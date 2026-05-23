package media

import (
	"bytes"
	"database/sql"
	"log/slog"
	"net/http"

	db_pkg "github.com/maddsua/flippercardapp/db"
)

func NewHandler(dbconn *sql.DB) http.Handler {

	db := db_pkg.NewWrapper(dbconn)

	mux := http.NewServeMux()

	mux.Handle("GET /images/{id}", http.HandlerFunc(func(wrt http.ResponseWriter, req *http.Request) {

		image, err := db.GetImageById(req.Context(), req.PathValue("id"))
		if db_pkg.IsNull(err) {
			wrt.WriteHeader(http.StatusNotFound)
			wrt.Write([]byte("Image not found"))
			return
		} else if err != nil {
			slog.Error("MEDIA Unable to retreive an image",
				slog.String("op", "sqlc.GetImageById"),
				slog.String("err", err.Error()))
			wrt.WriteHeader(http.StatusInternalServerError)
			wrt.Write([]byte("Unable to retreive an image"))
			return
		}

		http.ServeContent(wrt, req, req.URL.Path, image.CreatedAt.Time, bytes.NewReader(image.Data))
	}))

	return mux
}
