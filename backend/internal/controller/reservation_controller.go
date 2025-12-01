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

type ReservationController struct {
	reservationService *service.ReservationService
}

func NewReservationController(reservationService *service.ReservationService) *ReservationController {
	return &ReservationController{reservationService: reservationService}
}

func (ctl *ReservationController) CreateReservation(c *gin.Context) {
	currentUser := c.MustGet(middleware.ContextUserKey).(servicedto.User)

	var req controllerdto.CreateReservationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	out, err := ctl.reservationService.CreateReservation(c.Request.Context(), servicedto.CreateReservationInput{
		UserID:  currentUser.ID,
		Date:    req.Date,
		Time:    req.Time,
		People:  req.People,
		Comment: req.Comment,
	})
	if err != nil {
		switch err {
		case service.ErrInvalidInput:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create reservation"})
		}
		return
	}

	c.JSON(http.StatusCreated, toReservationResponse(out.Reservation))
}

func (ctl *ReservationController) ListMyReservations(c *gin.Context) {
	currentUser := c.MustGet(middleware.ContextUserKey).(servicedto.User)
	var status *string
	if s := c.Query("status"); s != "" {
		status = &s
	}

	res, err := ctl.reservationService.ListUserReservations(c.Request.Context(), servicedto.ListUserReservationsInput{
		UserID: currentUser.ID,
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

	resp := make([]controllerdto.ReservationResponse, 0, len(res))
	for _, r := range res {
		resp = append(resp, toReservationResponse(r))
	}
	c.JSON(http.StatusOK, resp)
}

func (ctl *ReservationController) CancelReservation(c *gin.Context) {
	currentUser := c.MustGet(middleware.ContextUserKey).(servicedto.User)
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid reservation id"})
		return
	}

	res, err := ctl.reservationService.CancelReservation(c.Request.Context(), servicedto.CancelReservationInput{
		UserID:        currentUser.ID,
		ReservationID: uint(id),
	})
	if err != nil {
		switch err {
		case service.ErrInvalidInput:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case service.ErrForbiddenReservation:
			c.JSON(http.StatusForbidden, gin.H{"error": "not allowed to cancel this reservation"})
		case service.ErrReservationNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "reservation not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to cancel reservation"})
		}
		return
	}

	c.JSON(http.StatusOK, toReservationResponse(*res))
}

func toReservationResponse(res servicedto.Reservation) controllerdto.ReservationResponse {
	return controllerdto.ReservationResponse{
		ID:        res.ID,
		UserID:    res.UserID,
		Date:      res.Date.Format("2006-01-02"),
		Time:      res.Time,
		People:    res.People,
		Comment:   res.Comment,
		Status:    res.Status,
		CreatedAt: res.CreatedAt.Format(time.RFC3339),
		UpdatedAt: res.UpdatedAt.Format(time.RFC3339),
	}
}
