package country

import (
	"errors"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"net/http"
	"whiteListApi/db"
)

// routes for this branch
func Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Post("/whitelistedbyip", CheckWhiteList)              // core functionality
	router.Get("/namebyip/{ipAddress}", GetCountryNameByIP)      // almost free
	router.Get("/latlngbyip/{ipAddress}", GetLatLngNameByIP)     // as long as we are here
	router.Get("/timezonebyip/{ipAddress}", GetTimeZoneNameByIP) // YAGNI?
	return router
}

// call DB package to get country, apply biz logic, return
// once could argue that the routing code/layer and the biz logic layer should be separated
// and this should be in a services package, but kept simple for this project
func CheckWhiteList(w http.ResponseWriter, r *http.Request) {

	// parse the posted JSON which is of type WhiteListRequest
	data := &WhiteListRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	ip := data.IPAddress
	country := db.CountryProvider.GetCountry(ip)

	whiteListed := contains(data.CountryCodes, country)

	render.JSON(w, r, WhitelistResponse{
		WhiteListed: whiteListed,
		IP:          ip,
		CountryName: country,
	})
}

// basic validation on above request
func (a *WhiteListRequest) Bind(r *http.Request) error {
	if a.IPAddress == "" {
		return errors.New("missing required IP field")
	}
	if a.CountryCodes == nil {
		return errors.New("missing required countrycodes field")
	}
	if len(a.CountryCodes) == 0 {
		return errors.New("no countrycodes supplied")
	}
	return nil
}

// helper function array contains
func contains(s []WhiteListCountryInfo, e string) bool {
	for _, a := range s {
		if a.IsoCode == e {
			return true
		}
	}
	return false
}

func GetCountryNameByIP(w http.ResponseWriter, r *http.Request) {
	ip := chi.URLParam(r, "ipAddress")
	country := db.CountryProvider.GetCountry(ip)
	retVal := Country{CountryName: country}
	render.JSON(w, r, retVal)
}

func GetLatLngNameByIP(w http.ResponseWriter, r *http.Request) {
	ip := chi.URLParam(r, "ipAddress")
	Lat, Lng := db.CountryProvider.GetLatLng(ip)
	retVal := LatLng{Latitude: Lat, Longitude: Lng}
	render.JSON(w, r, retVal)
}

func GetTimeZoneNameByIP(w http.ResponseWriter, r *http.Request) {
	ip := chi.URLParam(r, "ipAddress")
	tz := db.CountryProvider.GetTimeZone(ip)
	retVal := TimeZone{TimeZone: tz}
	render.JSON(w, r, retVal)
}

// Request and Response payloads for the REST api.
// normally these would be in another file or package

// return JSON instead of just bools or strings
type Country struct {
	CountryName string `json:"CountryName"`
}

type LatLng struct {
	Latitude  float64 `json:"Latitude"`
	Longitude float64 `json:"Longitude"`
}

type TimeZone struct {
	TimeZone string `json:"TimeZone"`
}

// WhiteListRequest is the request payload for whitelist POST method
type WhiteListRequest struct {
	IPAddress    string                 `json:"ip"`
	CountryCodes []WhiteListCountryInfo `json:"countrycodes"`
}

type WhiteListCountryInfo struct {
	IsoCode string `json:"iso_code"`
}

type WhitelistResponse struct {
	WhiteListed bool   `json:"WhiteListed"`
	IP          string `json:"IP"`
	CountryName string `json:"CountryName"`
}

// basic error handling
type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}
