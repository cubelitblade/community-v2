package transport

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/cubelitblade/community-v2/backend/pkg/common/httperr"
	"github.com/gin-gonic/gin"
)

// Checker is the interface for performing authorization checks.
type Checker interface {
	Check(ctx context.Context, user, relation, object string) (bool, error)
}

// Deps holds the dependencies required for the authz Handler.
type Deps struct {
	Checker Checker
	Logger  *slog.Logger
}

// Handler handles HTTP requests for authorization checks.
type Handler struct {
	checker Checker
	logger  *slog.Logger
}

// NewHandler creates and returns a new Handler.
func NewHandler(deps Deps) *Handler {
	if deps.Checker == nil {
		panic("nil checker")
	}

	if deps.Logger == nil {
		panic("nil logger")
	}

	return &Handler{
		checker: deps.Checker,
		logger:  deps.Logger,
	}
}

// RegisterRoutes registers the authz routes on the given router.
func (h *Handler) RegisterRoutes(router gin.IRouter) {
	g := router.Group("/authz")

	g.POST("/check", h.check)
}

type checkRequest struct {
	User     string `binding:"required" json:"user"`
	Relation string `binding:"required" json:"relation"`
	Object   string `binding:"required" json:"object"`
}

type checkResponse struct {
	Allowed bool `json:"allowed"`
}

func (h *Handler) check(c *gin.Context) {
	var req checkRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		httperr.WriteBadRequest(c, detailInvalidCheckBody())

		return
	}

	allowed, err := h.checker.Check(
		c.Request.Context(),
		req.User,
		req.Relation,
		req.Object,
	)
	if err != nil {
		h.logger.Debug(
			"authorization check failed",
			"user", req.User,
			"relation", req.Relation,
			"object", req.Object,
			"error", err,
		)
		httperr.WriteMappedError(c, err, authzProblem)

		return
	}

	c.JSON(http.StatusOK, checkResponse{Allowed: allowed})
}
