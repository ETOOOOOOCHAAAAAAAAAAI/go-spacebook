package services
//Unit tests for BookingService
//ТУТ Я ПРОВОЖУ ТЕСТЫ ДЛЯ BOOKING SERVICE :3
import (
	"testing"
	"time"
	"SpaceBookProject/internal/domain"
)

//fakeBookingRepo фейк бд, Чтобы не обращаться каждый раз к бд

type fakeBookingRepo struct {
	booking *domain.Booking
	overlap bool
}

func (f *fakeBookingRepo) Create(b *domain.Booking) error {
	b.ID = 1
	f.booking = b
	return nil
}

func (f *fakeBookingRepo) GetByID(id int) (*domain.Booking, error) {
	return f.booking, nil
}

func (f *fakeBookingRepo) ListByTenant(int) ([]domain.Booking, error) {
	return nil, nil
}

func (f *fakeBookingRepo) ListByOwner(int) ([]domain.Booking, error) {
	return nil, nil
}

func (f *fakeBookingRepo) UpdateStatus(id int, status domain.BookingStatus) error {
	f.booking.Status = status
	return nil
}

func (f *fakeBookingRepo) HasApprovedOverlap(int, time.Time, time.Time, *int) (bool, error) {
	return f.overlap, nil
}

type fakeSpaceRepo struct {
	ownerID int
}

func (f *fakeSpaceRepo) GetByID(int) (*domain.Space, error) {
	return &domain.Space{
		ID:      1,
		OwnerID: f.ownerID,
	}, nil
}




//1 FIRST TESTING CREATE BOOKING
// 1. Создание бронирования при валидных датах и отсутствии конфликтов
func TestCreateBooking_OK(t *testing.T) {
	bookings := &fakeBookingRepo{}
	spaces := &fakeSpaceRepo{}
	events := make(chan domain.BookingEvent, 1)

	service := NewBookingService(bookings, spaces, events)

	req := &domain.CreateBookingRequest{
		SpaceID:  1,
		DateFrom: "2025-01-01",
		DateTo:   "2025-01-10",
	}

	b, err := service.CreateBooking(42, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if b.ID == 0 {
		t.Fatal("booking ID not set")
	}

	if b.Status != domain.BookingStatusPending {
		t.Fatal("invalid booking status")
	}

	select {
	case e := <-events:
		if e.Type != domain.BookingEventCreated {
			t.Fatal("wrong event type")
		}
	default:
		t.Fatal("event not sent")
	}
}
// Проверка успешного создание, статус PENDING, генерацию события
//валидирует даты
//создает booking
//ставит статус PENDING
//отправляет событие BookingEventCreated



//2 SECOND TESTING CANCEL(FORBIDDEN)
//Попытка отмены бронирования не своим арендатором
func TestCancelBooking_Forbidden(t *testing.T) {
	bookings := &fakeBookingRepo{
		booking: &domain.Booking{
			ID:        1,
			TenantID:  100,
			Status:    domain.BookingStatusPending,
			DateFrom:  time.Now().Add(24 * time.Hour),
			DateTo:    time.Now().Add(48 * time.Hour),
		},
	}

	service := NewBookingService(bookings, &fakeSpaceRepo{}, nil)

	err := service.CancelBooking(1, 999)
	if err != ErrForbidden {
		t.Fatal("expected ErrForbidden")
	}
}
//Проверка запрет отмены → ErrForbidden
//отменять может только арендатор owner



//3 THIRD TESTING APPROVE (OWNER OK)
//Подтверждение бронирования владельцем пространства
func TestApproveBooking_OK(t *testing.T) {
	bookings := &fakeBookingRepo{
		booking: &domain.Booking{
			ID:       1,
			SpaceID:  1,
			TenantID: 50,
			Status:   domain.BookingStatusPending,
			DateFrom: time.Now(),
			DateTo:   time.Now().Add(24 * time.Hour),
		},
	}

	spaces := &fakeSpaceRepo{ownerID: 10}

	service := NewBookingService(bookings, spaces, nil)

	err := service.ApproveBooking(1, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if bookings.booking.Status != domain.BookingStatusApproved {
		t.Fatal("booking not approved")
	}
}
//Проверка владелец может approve, статус меняется на APPROVED
//только владелец может дать добро, точнее aprrove



//4 CREATE BOOKING — WRONG DATE
//Создание бронирования с некорректным диапазоном дат
func TestCreateBooking_InvalidDates(t *testing.T) {
	service := NewBookingService(
		&fakeBookingRepo{},
		&fakeSpaceRepo{},
		nil,
	)

	req := &domain.CreateBookingRequest{
		SpaceID:  1,
		DateFrom: "2025-01-10",
		DateTo:   "2025-01-01",
	}

	_, err := service.CreateBooking(1, req)
	if err == nil {
		t.Fatal("expected error for invalid dates")
	}
}
//Проверка валидацию дат -- ошибка



// 5 CREATE BOOKING — OVERLAP
//Создание бронирования при пересечении с существующим approved
func TestCreateBooking_Overlap(t *testing.T) {
	bookings := &fakeBookingRepo{overlap: true}

	service := NewBookingService(
		bookings,
		&fakeSpaceRepo{},
		nil,
	)

	req := &domain.CreateBookingRequest{
		SpaceID:  1,
		DateFrom: "2025-01-01",
		DateTo:   "2025-01-05",
	}

	_, err := service.CreateBooking(1, req)
	if err == nil {
		t.Fatal("expected overlap error")
	}
}
//Проверка запрет создания → ошибка overlap
//нельзя создать booking, если есть пересечение с approved
//сервис спрашивает репозиторий
//блокирует операцию



// 6 CANCEL BOOKING — ALREADY STARTED
//Попытка отмены бронирования, которое уже началось
func TestCancelBooking_AlreadyStarted(t *testing.T) {
	bookings := &fakeBookingRepo{
		booking: &domain.Booking{
			ID:        1,
			TenantID:  1,
			Status:    domain.BookingStatusPending,
			DateFrom:  time.Now().Add(-time.Hour),
			DateTo:    time.Now().Add(time.Hour),
		},
	}

	service := NewBookingService(bookings, &fakeSpaceRepo{}, nil)

	err := service.CancelBooking(1, 1)
	if err != ErrAlreadyStarted {
		t.Fatal("expected ErrAlreadyStarted")
	}
}
//Проверка. запрет отмены → ErrAlreadyStarted
//если бронирование уже началось, отмена запрещена


// 7 APPROVE BOOKING — NOT OWNER
//Подтверждение бронирования пользователем, не являющимся владельцем
func TestApproveBooking_Forbidden(t *testing.T) {
	bookings := &fakeBookingRepo{
		booking: &domain.Booking{
			ID:       1,
			SpaceID:  1,
			Status:   domain.BookingStatusPending,
			DateFrom: time.Now(),
			DateTo:   time.Now().Add(time.Hour),
		},
	}

	spaces := &fakeSpaceRepo{ownerID: 999}

	service := NewBookingService(bookings, spaces, nil)

	err := service.ApproveBooking(1, 1)
	if err != ErrForbidden {
		t.Fatal("expected ErrForbidden")
	}
}
//Проверка. контроль прав → ErrForbidden
//approve доступен только owner’у
//чужой -- ErrForbidden


// 8 APPROVE BOOKING — OVERLAP EXISTS
//Подтверждение бронирования при наличии пересекающегося approved
func TestApproveBooking_Overlap(t *testing.T) {
	bookings := &fakeBookingRepo{
		overlap: true,
		booking: &domain.Booking{
			ID:       1,
			SpaceID:  1,
			Status:   domain.BookingStatusPending,
			DateFrom: time.Now(),
			DateTo:   time.Now().Add(time.Hour),
		},
	}

	spaces := &fakeSpaceRepo{ownerID: 1}

	service := NewBookingService(bookings, spaces, nil)

	err := service.ApproveBooking(1, 1)
	if err != ErrOverlappingBooking {
		t.Fatal("expected ErrOverlappingBooking")
	}
}
//Проверка. Бизнес-правило эксклюзивности → ErrOverlappingBooking
//Занятность айди
//нельзя approved, если есть пересечение
//даже если ты владелец
//эксклюзивность пространства гарантируется



// 9 REJECT BOOKING — WRONG STATUS
//Отклонение бронирования в недопустимом статусе
func TestRejectBooking_WrongStatus(t *testing.T) {
	bookings := &fakeBookingRepo{
		booking: &domain.Booking{
			ID:      1,
			SpaceID: 1,
			Status:  domain.BookingStatusApproved,
		},
	}

	spaces := &fakeSpaceRepo{ownerID: 1}

	service := NewBookingService(bookings, spaces, nil)

	err := service.RejectBooking(1, 1)
	if err != ErrWrongStatus {
		t.Fatal("expected ErrWrongStatus")
	}
}
//Проверка. контроль переходов статуса → ErrWrongStatus
//reject можно только из PENDING

