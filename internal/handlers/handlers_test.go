package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/yj-matmul/bookings/internal/driver"
	"github.com/yj-matmul/bookings/internal/models"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	expectedStatusCode int
}{
	{"home", "/", "GET", http.StatusOK},
	{"about", "/about", "GET", http.StatusOK},
	{"generals-quarters", "/generals-quarters", "GET", http.StatusOK},
	{"majors-suite", "/majors-suite", "GET", http.StatusOK},
	{"search-availability", "/search-availability", "GET", http.StatusOK},
	{"contact", "/contact", "GET", http.StatusOK},
	{"non-existent", "/green/ham", "GET", http.StatusNotFound},

	// new routes
	{"login", "/user/login", "GET", http.StatusOK},
	{"log out", "/user/logout", "GET", http.StatusOK},

	{"dashboard", "/admin/dashboard", "GET", http.StatusOK},
	{"reservation new", "/admin/reservations-new", "GET", http.StatusOK},
	{"reservation all", "/admin/reservations-all", "GET", http.StatusOK},
	{"show new res", "/admin/reservations/new/1/show", "GET", http.StatusOK},

	// {"post-search-avail", "/search-availability", "POST", []postData{
	// 	{key: "start", value: "2022-01-01"},
	// 	{key: "end", value: "2022-01-02"},
	// }, http.StatusOK},
	// {"post-search-avail-json", "/search-availability-json", "POST", []postData{
	// 	{key: "start", value: "2022-01-01"},
	// 	{key: "end", value: "2022-01-02"},
	// }, http.StatusOK},
	// {"make reservation post", "/make-reservation", "POST", []postData{
	// 	{key: "first_name", value: "John"},
	// 	{key: "last_name", value: "Smith"},
	// 	{key: "email", value: "me@here.com"},
	// 	{key: "phone", value: "555-555-5555"},
	// }, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range theTests {
		resp, err := ts.Client().Get(ts.URL + e.url)
		if err != nil {
			t.Log(err)
			log.Fatal(err)
		}

		if resp.StatusCode != e.expectedStatusCode {
			t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
		}
	}
}

// test data for the TestRepository_Reservation handler
var reservationTests = []struct {
	name               string
	reservation        models.Reservation
	expectedStatusCode int
	expectedHTML       string
	expectedLocation   string
}{
	{
		name:               "reservation-in-session",
		reservation:        models.Reservation{RoomID: 1, Room: models.Room{ID: 1, RoomName: "General's Quarters"}},
		expectedStatusCode: http.StatusOK, expectedHTML: `action="/make-reservation"`,
	},
	{
		name:               "reservation-not-in-session",
		reservation:        models.Reservation{},
		expectedStatusCode: http.StatusSeeOther, expectedLocation: "/",
	},
	{
		name:               "non-existent-room",
		reservation:        models.Reservation{RoomID: 100000, Room: models.Room{ID: 100000, RoomName: ""}},
		expectedStatusCode: http.StatusSeeOther, expectedLocation: "/",
	},
}

func TestRepository_Reservation(t *testing.T) {
	for _, e := range reservationTests {
		req, _ := http.NewRequest("GET", "/make-reservation", nil)
		ctx := getCtx(req)
		req = req.WithContext(ctx)
		rr := httptest.NewRecorder()
		if e.reservation.RoomID > 0 {
			session.Put(ctx, "reservation", e.reservation)
		}

		handler := http.HandlerFunc(Repo.Reservation)
		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedStatusCode {
			t.Errorf("%s returned wrong response code: got %d, wanted %d", e.name, rr.Code, e.expectedStatusCode)
		}

		if e.expectedLocation != "" {
			actualLoc, _ := rr.Result().Location()
			if actualLoc.String() != e.expectedLocation {
				t.Errorf("failed %s: expected loaction %s, but got %s", e.name, e.expectedLocation, actualLoc.String())
			}
		}

		if e.expectedHTML != "" {
			html := rr.Body.String()
			if !strings.Contains(html, e.expectedHTML) {
				t.Errorf("failed %s: expected to find %s but did not", e.name, e.expectedHTML)
			}
		}
	}
}

var postReservationTests = []struct {
	name               string
	postedData         url.Values
	expectedStatusCode int
	expectedHTML       string
	expectedLocation   string
}{
	{
		name: "valid-post-reservation",
		postedData: url.Values{
			"start_date": {"2050-01-02"}, "end_date": {"2050-01-03"}, "first_name": {"John"},
			"last_name": {"Smith"}, "email": {"john@smith.com"}, "phone": {"555-111"}, "room_id": {"1"},
		},
		expectedStatusCode: http.StatusSeeOther,
		expectedLocation:   "/reservation-summary",
	},
	{
		name:               "non-existent-post-reservation",
		postedData:         nil,
		expectedStatusCode: http.StatusSeeOther,
		expectedLocation:   "/",
	},
	{
		name: "invalid-start-date-post-reservation",
		postedData: url.Values{
			"start_date": {"invalid"}, "end_date": {"2050-01-03"}, "first_name": {"John"},
			"last_name": {"Smith"}, "email": {"john@smith.com"}, "phone": {"555-111"}, "room_id": {"1"},
		},
		expectedStatusCode: http.StatusSeeOther,
		expectedLocation:   "/",
	},
	{
		name: "invalid-end-date-post-reservation",
		postedData: url.Values{
			"start_date": {"2050-01-02"}, "end_date": {"invalid"}, "first_name": {"John"},
			"last_name": {"Smith"}, "email": {"john@smith.com"}, "phone": {"555-111"}, "room_id": {"1"},
		},
		expectedStatusCode: http.StatusSeeOther,
		expectedLocation:   "/",
	},
	{
		name: "invalid-room-id-post-reservation",
		postedData: url.Values{
			"start_date": {"2050-01-02"}, "end_date": {"2050-01-03"}, "first_name": {"John"},
			"last_name": {"Smith"}, "email": {"john@smith.com"}, "phone": {"555-111"}, "room_id": {"invalid"},
		},
		expectedStatusCode: http.StatusSeeOther,
		expectedLocation:   "/",
	},
	{
		name: "invalid-room-data-post-reservation",
		postedData: url.Values{
			"start_date": {"2050-01-02"}, "end_date": {"2050-01-03"}, "first_name": {"John"},
			"last_name": {"Smith"}, "email": {"john@smith.com"}, "phone": {"555-111"}, "room_id": {"100000"},
		},
		expectedStatusCode: http.StatusSeeOther,
		expectedLocation:   "/",
	},
	{
		name: "missing-first-name-post-reservation",
		postedData: url.Values{
			"start_date": {"2050-01-02"}, "end_date": {"2050-01-03"}, "first_name": {""},
			"last_name": {"Smith"}, "email": {"john@smith.com"}, "phone": {"555-111"}, "room_id": {"1"},
		},
		expectedStatusCode: http.StatusOK,
		expectedHTML:       `action="/make-reservation"`,
	},
	{
		name: "missing-last-name-post-reservation",
		postedData: url.Values{
			"start_date": {"2050-01-02"}, "end_date": {"2050-01-03"}, "first_name": {"John"},
			"last_name": {""}, "email": {"john@smith.com"}, "phone": {"555-111"}, "room_id": {"1"},
		},
		expectedStatusCode: http.StatusOK,
		expectedHTML:       `action="/make-reservation"`,
	},
	{
		name: "missing-email-post-reservation",
		postedData: url.Values{
			"start_date": {"2050-01-02"}, "end_date": {"2050-01-03"}, "first_name": {"John"},
			"last_name": {"Smith"}, "email": {""}, "phone": {"555-111"}, "room_id": {"1"},
		},
		expectedStatusCode: http.StatusOK,
		expectedHTML:       `action="/make-reservation"`,
	},
	{
		name: "short-first-name-post-reservation",
		postedData: url.Values{
			"start_date": {"2050-01-02"}, "end_date": {"2050-01-03"}, "first_name": {"J"},
			"last_name": {"Smith"}, "email": {"john@smith.com"}, "phone": {"555-111"}, "room_id": {"1"},
		},
		expectedStatusCode: http.StatusOK,
		expectedHTML:       `action="/make-reservation"`,
	},
	{
		name: "invalid-email-post-reservation",
		postedData: url.Values{
			"start_date": {"2050-01-02"}, "end_date": {"2050-01-03"}, "first_name": {"John"},
			"last_name": {"Smith"}, "email": {"invalid"}, "phone": {"555-111"}, "room_id": {"1"},
		},
		expectedStatusCode: http.StatusOK,
		expectedHTML:       `action="/make-reservation"`,
	},
	{
		name: "database-insert--fail-reservation-post-reservation",
		postedData: url.Values{
			"start_date": {"2050-01-02"}, "end_date": {"2050-01-03"}, "first_name": {"John"},
			"last_name": {"Smith"}, "email": {"john@smith.com"}, "phone": {"555-111"}, "room_id": {"10000"},
		},
		expectedStatusCode: http.StatusSeeOther,
		expectedLocation:   "/",
	},
	{
		name: "database-insert--fail-room-restriction-post-reservation",
		postedData: url.Values{
			"start_date": {"2050-01-02"}, "end_date": {"2050-01-03"}, "first_name": {"John"},
			"last_name": {"Smith"}, "email": {"john@smith.com"}, "phone": {"555-111"}, "room_id": {"10001"},
		},
		expectedStatusCode: http.StatusSeeOther,
		expectedLocation:   "/",
	},
}

func TestRepository_PostReservation(t *testing.T) {
	for _, e := range postReservationTests {
		var req *http.Request
		if e.postedData != nil {
			req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(e.postedData.Encode()))
		} else {
			req, _ = http.NewRequest("POST", "/make-reservation", nil)
		}
		ctx := getCtx(req)
		req = req.WithContext(ctx)
		req.Header.Set("Content-type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(Repo.PostReservation)
		handler.ServeHTTP(rr, req)

		if e.expectedStatusCode != rr.Code {
			t.Errorf("%s returned wrong response code: got %d, wanted %d", e.name, rr.Code, e.expectedStatusCode)
		}

		if e.expectedLocation != "" {
			actualLoc, _ := rr.Result().Location()
			if actualLoc.String() != e.expectedLocation {
				t.Errorf("failed %s: expected loaction %s, but got %s", e.name, e.expectedLocation, actualLoc.String())
			}
		}

		if e.expectedHTML != "" {
			html := rr.Body.String()
			if !strings.Contains(html, e.expectedHTML) {
				t.Errorf("failed %s: expected to find %s but did not", e.name, e.expectedHTML)
			}
		}
	}
}

func TestNewRepo(t *testing.T) {
	var db driver.DB
	testRepo := NewRepo(&app, &db)

	if reflect.TypeOf(testRepo).String() != "*handlers.Repository" {
		t.Errorf("Did not get correct type from NewRepo: got %s, wanted *handlers.Repository",
			reflect.TypeOf(testRepo).String())
	}
}

var reservationSummaryTests = []struct {
	name               string
	reservation        models.Reservation
	expectedStatusCode int
	expectedHTML       string
	expectedLocation   string
}{
	{
		name:               "reservation-in-session-reservation-summary",
		reservation:        models.Reservation{RoomID: 1, Room: models.Room{ID: 1, RoomName: "General's Quarters"}},
		expectedStatusCode: http.StatusOK, expectedHTML: `Reservation Summary`,
	},
	{
		name:               "reservation-not-in-session-reservation-summary",
		reservation:        models.Reservation{},
		expectedStatusCode: http.StatusSeeOther, expectedLocation: "/",
	},
}

func TestRepository_ReservationSummary(t *testing.T) {
	for _, e := range reservationSummaryTests {
		req, _ := http.NewRequest("GET", "/reservation-summary", nil)
		ctx := getCtx(req)
		req = req.WithContext(ctx)
		req.Header.Set("Content-type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()

		if e.reservation.RoomID > 0 {
			session.Put(ctx, "reservation", e.reservation)
		}

		handler := http.HandlerFunc(Repo.ReservationSummary)
		handler.ServeHTTP(rr, req)

		if e.expectedStatusCode != rr.Code {
			t.Errorf("%s returned wrong response code: got %d, wanted %d", e.name, rr.Code, e.expectedStatusCode)
		}

		if e.expectedLocation != "" {
			actualLoc, _ := rr.Result().Location()
			if actualLoc.String() != e.expectedLocation {
				t.Errorf("failed %s: expected loaction %s, but got %s", e.name, e.expectedLocation, actualLoc.String())
			}
		}

		if e.expectedHTML != "" {
			html := rr.Body.String()
			if !strings.Contains(html, e.expectedHTML) {
				t.Errorf("failed %s: expected to find %s but did not", e.name, e.expectedHTML)
			}
		}
	}
}

var postAvailbilityTests = []struct {
	name               string
	postedData         url.Values
	expectedStatusCode int
	expectedHTML       string
	expectedLocation   string
}{
	{
		name: "valid-date-post-availability", postedData: url.Values{"start": {"2050-01-02"}, "end": {"2050-01-03"}},
		expectedStatusCode: http.StatusOK, expectedHTML: `Choose a Room`,
	},
	{
		name: "empty-postdata-post-availability", postedData: nil,
		expectedStatusCode: http.StatusSeeOther, expectedLocation: "/",
	},
	{
		name: "invalid-start-date-post-availability", postedData: url.Values{"start": {"invalid"}, "end": {"2050-01-03"}},
		expectedStatusCode: http.StatusSeeOther, expectedLocation: "/",
	},
	{
		name: "invalid-end-date-post-availability", postedData: url.Values{"start": {"2050-01-02"}, "end": {"invalid"}},
		expectedStatusCode: http.StatusSeeOther, expectedLocation: "/",
	},
	{
		name: "start-after-end-post-availability", postedData: url.Values{"start": {"2050-01-04"}, "end": {"2050-01-03"}},
		expectedStatusCode: http.StatusSeeOther, expectedLocation: "/",
	},
	{
		name: "no-left-room-post-availability", postedData: url.Values{"start": {"2021-08-11"}, "end": {"2021-08-12"}},
		expectedStatusCode: http.StatusSeeOther, expectedLocation: "/search-availability",
	},
}

func TestRepository_PostAvailability(t *testing.T) {
	for _, e := range postAvailbilityTests {
		var req *http.Request
		if e.postedData != nil {
			req, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(e.postedData.Encode()))
		} else {
			req, _ = http.NewRequest("POST", "/search-availability", nil)
		}
		ctx := getCtx(req)
		req = req.WithContext(ctx)
		req.Header.Set("Content-type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(Repo.PostAvailability)
		handler.ServeHTTP(rr, req)

		if e.expectedStatusCode != rr.Code {
			t.Errorf("%s returned wrong response code: got %d, wanted %d", e.name, rr.Code, e.expectedStatusCode)
		}

		if e.expectedLocation != "" {
			actualLoc, _ := rr.Result().Location()
			if actualLoc.String() != e.expectedLocation {
				t.Errorf("failed %s: expected loaction %s, but got %s", e.name, e.expectedLocation, actualLoc.String())
			}
		}

		if e.expectedHTML != "" {
			html := rr.Body.String()
			if !strings.Contains(html, e.expectedHTML) {
				t.Errorf("failed %s: expected to find %s but did not", e.name, e.expectedHTML)
			}
		}
	}
}

var availabilityJsonTests = []struct {
	name       string
	postedData url.Values
	expectedOK bool
}{
	{
		name:       "valid-availability-json",
		postedData: url.Values{"start": {"2050-01-02"}, "end": {"2050-01-03"}, "room_id": {"1"}},
		expectedOK: true,
	},
	{
		name:       "empty-availability-json",
		postedData: nil,
		expectedOK: false,
	},
	{
		name:       "no-left-room-availability-json",
		postedData: url.Values{"start": {"2050-01-02"}, "end": {"2050-01-03"}, "room_id": {"10000"}},
		expectedOK: false,
	},
}

func TestRepository_AvailabilityJson(t *testing.T) {
	for _, e := range availabilityJsonTests {
		var req *http.Request
		if e.postedData != nil {
			req, _ = http.NewRequest("POST", "/search-availability-json", strings.NewReader(e.postedData.Encode()))
		} else {
			req, _ = http.NewRequest("POST", "/search-availability-json", nil)
		}
		ctx := getCtx(req)
		req = req.WithContext(ctx)
		req.Header.Set("Content-type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(Repo.AvailabilityJSON)
		handler.ServeHTTP(rr, req)

		var j jsonResponse
		err := json.Unmarshal([]byte(rr.Body.Bytes()), &j)
		if err != nil {
			t.Error("failed to parse json")
		}

		if j.OK != e.expectedOK {
			t.Errorf("%s: expected %v but got %v", e.name, e.expectedOK, j.OK)
		}
	}
}

var chooseRoomTests = []struct {
	name               string
	reservation        models.Reservation
	url                string
	expectedStatusCode int
	expectedLocation   string
}{
	{
		name:        "reservation-in-session-choose-room",
		reservation: models.Reservation{RoomID: 1, Room: models.Room{ID: 1, RoomName: "General's Quarters"}},
		url:         "/choose-room/1", expectedStatusCode: http.StatusSeeOther, expectedLocation: "/make-reservation",
	},
	{
		name:        "reservation-not-in-session-choose-room",
		reservation: models.Reservation{},
		url:         "/choose-room/1", expectedStatusCode: http.StatusSeeOther, expectedLocation: "/",
	},
	{
		name:        "invalid-url-choose-room",
		reservation: models.Reservation{RoomID: 1, Room: models.Room{ID: 1, RoomName: "General's Quarters"}},
		url:         "/choose-room/invalid", expectedStatusCode: http.StatusSeeOther, expectedLocation: "/",
	},
}

func TestRepository_ChooseRoom(t *testing.T) {
	for _, e := range chooseRoomTests {
		req, _ := http.NewRequest("GET", e.url, nil)
		ctx := getCtx(req)
		req = req.WithContext(ctx)
		req.RequestURI = e.url

		rr := httptest.NewRecorder()
		if e.reservation.RoomID > 0 {
			session.Put(ctx, "reservation", e.reservation)
		}

		handler := http.HandlerFunc(Repo.ChooseRoom)
		handler.ServeHTTP(rr, req)

		if e.expectedStatusCode != rr.Code {
			t.Errorf("%s returned wrong response code: got %d, wanted %d", e.name, rr.Code, e.expectedStatusCode)
		}

		if e.expectedLocation != "" {
			actualLoc, _ := rr.Result().Location()
			if actualLoc.String() != e.expectedLocation {
				t.Errorf("failed %s: expected loaction %s, but got %s", e.name, e.expectedLocation, actualLoc.String())
			}
		}
	}
}

var bookRoomTests = []struct {
	name               string
	url                string
	expectedStatusCode int
	expectedLocation   string
}{
	{
		name:               "valid-url-book-room",
		url:                "/book-room?s=2050-01-02&e=2050-01-03&id=1",
		expectedStatusCode: http.StatusSeeOther, expectedLocation: "/make-reservation",
	},
	{
		name:               "invalid-room-id-book-room",
		url:                "/book-room?s=2050-01-02&e=2050-01-03&id=100000",
		expectedStatusCode: http.StatusSeeOther, expectedLocation: "/",
	},
}

func TestRepository_BookRoom(t *testing.T) {
	for _, e := range bookRoomTests {
		req, _ := http.NewRequest("GET", e.url, nil)
		ctx := getCtx(req)
		req = req.WithContext(ctx)
		req.RequestURI = e.url

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(Repo.BookRoom)
		handler.ServeHTTP(rr, req)

		if e.expectedStatusCode != rr.Code {
			t.Errorf("%s returned wrong response code: got %d, wanted %d", e.name, rr.Code, e.expectedStatusCode)
		}

		if e.expectedLocation != "" {
			actualLoc, _ := rr.Result().Location()
			if actualLoc.String() != e.expectedLocation {
				t.Errorf("failed %s: expected loaction %s, but got %s", e.name, e.expectedLocation, actualLoc.String())
			}
		}
	}
}

var loginTests = []struct {
	name               string
	email              string
	expectedStatusCode int
	expectedHTML       string
	expectedLocation   string
}{
	{"normal-credentials", "adm@adm.com", http.StatusSeeOther, "", "/"},
	{"invaid-credentials", "ad@adm.co.kr", http.StatusSeeOther, "", "/user/login"},
	{"invalid-data", "j", http.StatusOK, `action="/user/login"`, ""},
}

func TestLogin(t *testing.T) {
	for _, e := range loginTests {
		postedData := url.Values{}
		postedData.Add("email", e.email)
		postedData.Add("password", "book")

		req, _ := http.NewRequest("POST", "/user/login", strings.NewReader(postedData.Encode()))
		ctx := getCtx(req)
		req = req.WithContext(ctx)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(Repo.PostShowLogin)
		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedStatusCode {
			t.Errorf("failed %s: expected code %d, but got %d", e.name, e.expectedStatusCode, rr.Code)
		}

		if e.expectedLocation != "" {
			// get the URL from the test
			actualLoc, _ := rr.Result().Location()
			if actualLoc.String() != e.expectedLocation {
				t.Errorf("failed %s: expected location %s, but got %s", e.name, e.expectedLocation, actualLoc.String())
			}
		}

		if e.expectedHTML != "" {
			// read the response body into a string
			html := rr.Body.String()
			if !strings.Contains(html, e.expectedHTML) {
				t.Errorf("failed %s: expected to find %s, but did not", e.name, e.expectedHTML)
			}
		}
	}
}

var adminPostShowReservationTests = []struct {
	name               string
	url                string
	postedData         url.Values
	expectedStatusCode int
	expectedHTML       string
	expectedLocation   string
}{
	{
		name: "valid-data-from-new-admin-post-show-res", url: "/admin/reservations/new/1",
		postedData:         url.Values{"first_name": {"John"}, "last_name": {"Smith"}, "email": {"john@smi.com"}, "phone": {"555"}},
		expectedStatusCode: http.StatusSeeOther, expectedLocation: "/admin/reservations-new",
	},
	{
		name: "valid-data-from-all-admin-post-show-res", url: "/admin/reservations/all/1",
		postedData:         url.Values{"first_name": {"John"}, "last_name": {"Smith"}, "email": {"john@smi.com"}, "phone": {"555"}},
		expectedStatusCode: http.StatusSeeOther, expectedLocation: "/admin/reservations-all",
	},
	{
		name: "valid-data-from-cal-admin-post-show-res", url: "/admin/reservations/cal/1",
		postedData: url.Values{
			"first_name": {"John"}, "last_name": {"Smith"}, "email": {"john@smi.com"}, "phone": {"555"},
			"year": {"2050"}, "month": {"01"}},
		expectedStatusCode: http.StatusSeeOther, expectedLocation: "/admin/reservations-calendar?y=2050&m=01",
	},
}

func TestAdminPostShowReservation(t *testing.T) {
	for _, e := range adminPostShowReservationTests {
		var req *http.Request
		if e.postedData != nil {
			req, _ = http.NewRequest("POST", e.url, strings.NewReader(e.postedData.Encode()))
		} else {
			req, _ = http.NewRequest("POST", e.url, nil)
		}
		ctx := getCtx(req)
		req = req.WithContext(ctx)
		req.RequestURI = e.url
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(Repo.AdminPostShowReservation)
		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedStatusCode {
			t.Errorf("failed %s: expected code %d, but got %d", e.name, e.expectedStatusCode, rr.Code)
		}

		if e.expectedLocation != "" {
			actualLoc, _ := rr.Result().Location()
			if actualLoc.String() != e.expectedLocation {
				t.Errorf("failed %s: expected location %s, but got %s", e.name, e.expectedLocation, actualLoc.String())
			}
		}

		if e.expectedHTML != "" {
			html := rr.Body.String()
			if !strings.Contains(html, e.expectedHTML) {
				t.Errorf("failed %s: expected to find %s, but did not", e.name, e.expectedHTML)
			}
		}
	}
}

var adminReservationsCalendarTests = []struct {
	name               string
	url                string
	expectedStatusCode int
	expectedHTML       string
}{
	{
		name:               "valid-no-year-info-admin-res-cal",
		url:                "/admin/reservations-calendar",
		expectedStatusCode: http.StatusOK, expectedHTML: `action="/admin/reservations-calendar"`,
	},
	{
		name:               "valid-existent-year-info-admin-res-cal",
		url:                "/admin/reservations-calendar?y=2050&m=01",
		expectedStatusCode: http.StatusOK, expectedHTML: `action="/admin/reservations-calendar"`,
	},
}

func TestAdminReservationsCalendar(t *testing.T) {
	for _, e := range adminReservationsCalendarTests {
		req, _ := http.NewRequest("GET", e.url, nil)
		ctx := getCtx(req)
		req = req.WithContext(ctx)
		req.RequestURI = e.url
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(Repo.AdminReservationsCalendar)
		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedStatusCode {
			t.Errorf("failed %s: expected code %d, but got %d", e.name, e.expectedStatusCode, rr.Code)
		}

		if e.expectedHTML != "" {
			html := rr.Body.String()
			if !strings.Contains(html, e.expectedHTML) {
				t.Errorf("failed %s: expected to find %s, but did not", e.name, e.expectedHTML)
			}
		}
	}
}

var adminProcessReservationTests = []struct {
	name               string
	url                string
	expectedStatusCode int
	expectedLocation   string
}{
	{
		name:               "res-new-without-year-admin-process-res",
		url:                "/admin/process-reservation/new/1/do",
		expectedStatusCode: http.StatusSeeOther, expectedLocation: "/admin/reservations-new",
	},
	{
		name:               "res-all-without-year-admin-process-res",
		url:                "/admin/process-reservation/all/1/do",
		expectedStatusCode: http.StatusSeeOther, expectedLocation: "/admin/reservations-all",
	},
	{
		name:               "res-cal-without-year-admin-process-res",
		url:                "/admin/process-reservation/cal/1/do",
		expectedStatusCode: http.StatusSeeOther, expectedLocation: "/admin/reservations-cal",
	},
	{
		name:               "res-cal-with-year-admin-process-res",
		url:                "/admin/process-reservation/cal/1/do?y=2050&m=01",
		expectedStatusCode: http.StatusSeeOther, expectedLocation: "/admin/reservations-calendar?y=2050&m=01",
	},
}

func TestAdminProcessReservation(t *testing.T) {
	for _, e := range adminProcessReservationTests {
		req, _ := http.NewRequest("GET", e.url, nil)
		ctx := getCtx(req)
		req = req.WithContext(ctx)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.RequestURI = e.url
		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(Repo.AdminProcessReservation)
		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedStatusCode {
			t.Errorf("failed %s: expected code %d, but got %d", e.name, e.expectedStatusCode, rr.Code)
		}

		if e.expectedLocation != "" {
			actualLoc, _ := rr.Result().Location()
			if actualLoc.String() != e.expectedLocation {
				t.Errorf("failed %s: expected location %s, but got %s", e.name, e.expectedLocation, actualLoc.String())
			}
		}
	}
}

var adminDeleteReservationTests = []struct {
	name               string
	url                string
	expectedStatusCode int
	expectedLocation   string
}{
	{
		name:               "res-new-without-year-admin-process-res",
		url:                "/admin/delete-reservation/new/1/do",
		expectedStatusCode: http.StatusSeeOther, expectedLocation: "/admin/reservations-new",
	},
	{
		name:               "res-all-without-year-admin-process-res",
		url:                "/admin/delete-reservation/all/1/do",
		expectedStatusCode: http.StatusSeeOther, expectedLocation: "/admin/reservations-all",
	},
	{
		name:               "res-cal-without-year-admin-process-res",
		url:                "/admin/delete-reservation/cal/1/do",
		expectedStatusCode: http.StatusSeeOther, expectedLocation: "/admin/reservations-cal",
	},
	{
		name:               "res-cal-with-year-admin-process-res",
		url:                "/admin/delete-reservation/cal/1/do?y=2050&m=01",
		expectedStatusCode: http.StatusSeeOther, expectedLocation: "/admin/reservations-calendar?y=2050&m=01",
	},
}

func TestAdminDeleteReservation(t *testing.T) {
	for _, e := range adminDeleteReservationTests {
		req, _ := http.NewRequest("GET", e.url, nil)
		ctx := getCtx(req)
		req = req.WithContext(ctx)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.RequestURI = e.url
		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(Repo.AdminDeleteReservation)
		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedStatusCode {
			t.Errorf("failed %s: expected code %d, but got %d", e.name, e.expectedStatusCode, rr.Code)
		}

		if e.expectedLocation != "" {
			actualLoc, _ := rr.Result().Location()
			if actualLoc.String() != e.expectedLocation {
				t.Errorf("failed %s: expected location %s, but got %s", e.name, e.expectedLocation, actualLoc.String())
			}
		}
	}
}

func getCtx(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}

	return ctx
}
