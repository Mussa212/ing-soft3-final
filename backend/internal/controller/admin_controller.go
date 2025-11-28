package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	controllerdto "vesuvio/internal/dto/controller"
	servicedto "vesuvio/internal/dto/service"
	"vesuvio/internal/middleware"
	"vesuvio/internal/service"
)

type AdminController struct {
	reservationService *service.ReservationService
}

func NewAdminController(reservationService *service.ReservationService) *AdminController {
	return &AdminController{reservationService: reservationService}
}

func (ctl *AdminController) ListReservations(c *gin.Context) {
	date := c.Query("date")
	var status *string
	if s := c.Query("status"); s != "" {
		status = &s
	}

	res, err := ctl.reservationService.AdminListReservations(c.Request.Context(), servicedto.AdminListReservationsInput{
		Date:   date,
		Status: status,
	})
	if err != nil {
		switch err {
		case service.ErrInvalidInput, service.ErrInvalidStatus:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list reservations"})
		}
		return
	}

	resp := make([]controllerdto.AdminReservationResponse, 0, len(res))
	for _, r := range res {
		user := controllerdto.AdminUserInfo{
			ID:    r.User.ID,
			Name:  r.User.Name,
			Email: r.User.Email,
		}
		resp = append(resp, controllerdto.AdminReservationResponse{
			ID:        r.ID,
			User:      user,
			Date:      r.Date.Format("2006-01-02"),
			Time:      r.Time,
			People:    r.People,
			Comment:   r.Comment,
			Status:    r.Status,
			CreatedAt: r.CreatedAt.Format(time.RFC3339),
			UpdatedAt: r.UpdatedAt.Format(time.RFC3339),
		})
	}

	c.JSON(http.StatusOK, resp)
}

func (ctl *AdminController) ConfirmReservation(c *gin.Context) {
	currentUser := c.MustGet(middleware.ContextUserKey).(servicedto.User)
	reservationID, ok := parseIDParam(c.Param("id"))
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid reservation id"})
		return
	}

	res, err := ctl.reservationService.ConfirmReservation(c.Request.Context(), currentUser, reservationID)
	if err != nil {
		switch err {
		case service.ErrReservationNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "reservation not found"})
		case service.ErrInvalidInput, service.ErrInvalidStatus:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case service.ErrUnauthorized:
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to confirm reservation"})
		}
		return
	}

	c.JSON(http.StatusOK, toReservationResponse(*res))
}

func (ctl *AdminController) CancelReservation(c *gin.Context) {
	currentUser := c.MustGet(middleware.ContextUserKey).(servicedto.User)
	reservationID, ok := parseIDParam(c.Param("id"))
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid reservation id"})
		return
	}

	res, err := ctl.reservationService.AdminCancelReservation(c.Request.Context(), currentUser, reservationID)
	if err != nil {
		switch err {
		case service.ErrReservationNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "reservation not found"})
		case service.ErrInvalidInput, service.ErrInvalidStatus:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case service.ErrUnauthorized:
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to cancel reservation"})
		}
		return
	}

	c.JSON(http.StatusOK, toReservationResponse(*res))
}

func parseIDParam(param string) (uint, bool) {
	id, err := strconv.ParseUint(param, 10, 64)
	if err != nil || id == 0 {
		return 0, false
	}
	return uint(id), true
}
