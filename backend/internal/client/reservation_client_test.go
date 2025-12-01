package client

import (
	"context"
	"fmt"
	"testing"
	"time"

	"gorm.io/driver/sqlite"

	servicedto "vesuvio/internal/dto/service"
)

func newReservationTestClient(t *testing.T) *GormReservationClient {
	t.Helper()
	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())
	db, err := NewDBWithDialector(sqlite.Open(dsn))
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	return NewReservationClient(db)
}

func TestReservationClient_CRUD(t *testing.T) {
	ctx := context.Background()
	client := newReservationTestClient(t)

	date := time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC)
	created, err := client.CreateReservation(ctx, servicedto.CreateReservationParams{
		UserID:  1,
		Date:    date,
		Time:    "20:00",
		People:  2,
		Status:  servicedto.StatusPending,
		Comment: nil,
	})
	if err != nil {
		t.Fatalf("create reservation: %v", err)
	}
	if created.ID == 0 || created.Status != servicedto.StatusPending {
		t.Fatalf("unexpected created reservation: %+v", created)
	}

	byID, err := client.GetReservationByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("get by id: %v", err)
	}
	if byID == nil || byID.ID != created.ID {
		t.Fatalf("unexpected reservation by id: %+v", byID)
	}

	updated, err := client.UpdateReservationStatus(ctx, created.ID, servicedto.StatusConfirmed)
	if err != nil {
		t.Fatalf("update status: %v", err)
	}
	if updated.Status != servicedto.StatusConfirmed {
		t.Fatalf("expected status confirmed, got %s", updated.Status)
	}

	none, err := client.GetReservationByID(ctx, 999)
	if err != nil {
		t.Fatalf("unexpected error for missing reservation: %v", err)
	}
	if none != nil {
		t.Fatalf("expected nil for missing reservation, got %+v", none)
	}
}

func TestReservationClient_Listing(t *testing.T) {
	ctx := context.Background()
	client := newReservationTestClient(t)
	date := time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC)
	otherDate := time.Date(2025, 12, 2, 0, 0, 0, 0, time.UTC)

	_, _ = client.CreateReservation(ctx, servicedto.CreateReservationParams{
		UserID: 1, Date: date, Time: "20:00", People: 2, Status: servicedto.StatusPending,
	})
	_, _ = client.CreateReservation(ctx, servicedto.CreateReservationParams{
		UserID: 1, Date: date, Time: "21:00", People: 3, Status: servicedto.StatusConfirmed,
	})
	_, _ = client.CreateReservation(ctx, servicedto.CreateReservationParams{
		UserID: 2, Date: otherDate, Time: "22:00", People: 2, Status: servicedto.StatusCancelled,
	})

	status := servicedto.StatusConfirmed
	userList, err := client.ListReservationsByUser(ctx, 1, &status)
	if err != nil {
		t.Fatalf("list by user: %v", err)
	}
	if len(userList) != 1 || userList[0].Status != servicedto.StatusConfirmed {
		t.Fatalf("unexpected user list: %+v", userList)
	}

	dateStatus := servicedto.StatusPending
	dateList, err := client.ListReservationsByDate(ctx, date, &dateStatus)
	if err != nil {
		t.Fatalf("list by date: %v", err)
	}
	if len(dateList) != 1 || dateList[0].Status != servicedto.StatusPending {
		t.Fatalf("unexpected date list: %+v", dateList)
	}
}
