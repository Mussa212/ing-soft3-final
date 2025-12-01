package client

import (
	"testing"
	"time"

	servicedto "vesuvio/internal/dto/service"
	"vesuvio/internal/model"
)

func TestToServiceUser(t *testing.T) {
	now := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)
	user := model.UserModel{
		ID:        1,
		Name:      "Alice",
		Email:     "alice@example.com",
		IsAdmin:   true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	result := toServiceUser(&user)
	if result.ID != user.ID || result.Name != user.Name || result.Email != user.Email || !result.CreatedAt.Equal(now) || !result.UpdatedAt.Equal(now) {
		t.Fatalf("unexpected user mapping: %+v", result)
	}
	if !result.IsAdmin {
		t.Fatalf("expected IsAdmin to be true")
	}
}

func TestMapReservationsWithoutUser(t *testing.T) {
	comment := "Window"
	date := time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)
	created := time.Date(2024, 2, 1, 10, 0, 0, 0, time.UTC)
	models := []model.ReservationModel{{
		ID:        10,
		UserID:    5,
		Date:      date,
		Time:      "19:30",
		People:    2,
		Comment:   &comment,
		Status:    "pending",
		CreatedAt: created,
		UpdatedAt: created,
	}}

	res := mapReservations(models, nil)
	if len(res) != 1 {
		t.Fatalf("expected 1 reservation, got %d", len(res))
	}
	if res[0].User != nil {
		t.Fatalf("expected nil user when mapper is nil")
	}
	if res[0].ID != models[0].ID || res[0].UserID != models[0].UserID || res[0].Time != "19:30" || res[0].Status != "pending" {
		t.Fatalf("unexpected reservation mapping: %+v", res[0])
	}
	if !res[0].Date.Equal(date) || !res[0].CreatedAt.Equal(created) || !res[0].UpdatedAt.Equal(created) {
		t.Fatalf("unexpected timestamps in mapping: %+v", res[0])
	}
	if res[0].Comment == nil || *res[0].Comment != comment {
		t.Fatalf("expected comment to be preserved")
	}
}

func TestMapReservationsWithUserMapper(t *testing.T) {
	date := time.Date(2024, 3, 5, 0, 0, 0, 0, time.UTC)
	models := []model.ReservationModel{{
		ID:     11,
		UserID: 7,
		User: model.UserModel{
			ID:    7,
			Name:  "Bob",
			Email: "bob@example.com",
		},
		Date:   date,
		Time:   "20:00",
		People: 4,
		Status: "confirmed",
	}}

	res := mapReservations(models, func(m model.ReservationModel) *servicedto.User {
		return toServiceUser(&m.User)
	})

	if len(res) != 1 {
		t.Fatalf("expected 1 reservation, got %d", len(res))
	}
	if res[0].User == nil || res[0].User.Email != "bob@example.com" {
		t.Fatalf("expected embedded user to be mapped, got %+v", res[0].User)
	}
	if res[0].Status != "confirmed" || res[0].Time != "20:00" || res[0].People != 4 {
		t.Fatalf("unexpected reservation data: %+v", res[0])
	}
	if !res[0].Date.Equal(date) {
		t.Fatalf("date not copied correctly")
	}
}
