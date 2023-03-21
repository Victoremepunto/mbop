package handlers

import (
	"net/http"

	"github.com/redhatinsights/mbop/internal/config"
	"github.com/redhatinsights/mbop/internal/models"
)

func Status(w http.ResponseWriter, _ *http.Request) {
	status := models.Status{
		ConfiguredModules: models.ConfiguredModules{
			Users:  config.Get().UsersModule,
			Mailer: config.Get().MailerModule,
			JWT:    config.Get().JwtModule,
		},
	}

	sendJSON(w, status)
}
