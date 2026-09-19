package timetable

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/openschool-org/openschool/internal/apierror"
)

type StatusHistoryItem struct {
	ID            uuid.UUID `json:"id"`
	TimetableID   uuid.UUID `json:"timetable_id"`
	FromStatus    *string   `json:"from_status"`
	ToStatus      string    `json:"to_status"`
	ChangedBy     uuid.UUID `json:"changed_by"`
	Comment       *string   `json:"comment"`
	ChangedAt     time.Time `json:"changed_at"`
	ChangedByName string    `json:"changed_by_name"`
}

type statusHistoryStore interface {
	listStatusHistory(context.Context, uuid.UUID) ([]StatusHistoryItem, error)
}

type statusHistoryService struct{ store statusHistoryStore }

func (s *statusHistoryService) list(ctx context.Context, timetableID uuid.UUID) ([]StatusHistoryItem, error) {
	return s.store.listStatusHistory(ctx, timetableID)
}

type statusHistoryHandler struct{ service *statusHistoryService }

func newStatusHistoryHandler(store statusHistoryStore) *statusHistoryHandler {
	return &statusHistoryHandler{service: &statusHistoryService{store: store}}
}

func RegisterTimetableStatusHistoryRoute(teacherOrAdmin *gin.RouterGroup, pool *pgxpool.Pool) {
	handler := newStatusHistoryHandler(newTimetableRepository(pool))
	teacherOrAdmin.GET("/timetables/:id/status-history", handler.list)
}

func (h *statusHistoryHandler) list(c *gin.Context) {
	id, ok := parseTimetableID(c)
	if !ok {
		return
	}
	history, err := h.service.list(c.Request.Context(), id)
	if err != nil {
		apierror.RespondInternal(c, err)
		return
	}
	c.JSON(http.StatusOK, history)
}
