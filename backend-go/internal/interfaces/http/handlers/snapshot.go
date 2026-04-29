package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	apsnapshot "github.com/fintech/cbpi/backend-go/internal/application/snapshot"
	"github.com/fintech/cbpi/backend-go/internal/interfaces/http/dto"
)

type SnapshotHandler struct {
	history *apsnapshot.GetHistory
}

func NewSnapshotHandler(h *apsnapshot.GetHistory) *SnapshotHandler {
	return &SnapshotHandler{history: h}
}

func (h *SnapshotHandler) History(c *gin.Context) {
	items, err := h.history.Execute(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.NewSnapshotResponses(items))
}
