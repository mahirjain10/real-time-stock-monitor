package alert

import (
	"context"
	"net/http"

	"github/mahirjain_10/sse-backend/backend/internal/helpers"
	"github/mahirjain_10/sse-backend/backend/internal/types"
	"github/mahirjain_10/sse-backend/backend/internal/utils"
	"github/mahirjain_10/sse-backend/backend/internal/websocket"

	"github.com/gin-gonic/gin"
)

func SendAlertNotification(c *gin.Context, r *gin.Engine, app *types.App,hub *websocket.Hub) {
	ctx := context.Background()
	
	var UpdateActiveStatus types.UpdateActiveStatus
	if !helpers.BindAndValidateJSON(c, &UpdateActiveStatus) {
		return
	}
	
	// Unregister the WebSocket client before updating the status
	// if hub := hub.UnregisterClientByAlertID(UpdateActiveStatus; hub != nil {
		// hub.UnregisterClientByAlertID(UpdateActiveStatus.ID)
	// }
	
	utils.UpdateActiveStatusUtil(c, ctx, UpdateActiveStatus.UserID, UpdateActiveStatus.ID, UpdateActiveStatus.Active, app)
	
	data := map[string]interface{}{
		"user_id": UpdateActiveStatus.UserID,
		"alert_id": UpdateActiveStatus.ID,
	}
	
	
	helpers.SendResponse(c, http.StatusOK, "Alert condition met", data, nil, true)
}


