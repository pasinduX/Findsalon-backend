package dto

type BookingNotificationPayload struct {
	BookingId     string `json:"BookingId"`
	UserId        string `json:"UserId"`
	SalonId       string `json:"SalonId"`
	BarberId      string `json:"BarberId"`
	CustomerName  string `json:"CustomerName"`
	CustomerEmail string `json:"CustomerEmail"`
	CustomerPhone string `json:"CustomerPhone"`
	Date          string `json:"Date"`
	StartTime     string `json:"StartTime"`
	EndTime       string `json:"EndTime"`
	EventType     string `json:"EventType"`
}

type CustomNotificationPayload struct {
	UserId    string `json:"UserId" validate:"required"`
	Title     string `json:"Title" validate:"required"`
	Body      string `json:"Body" validate:"required"`
	EventType string `json:"EventType"`
	RefId     string `json:"RefId"`
}

type BulkNotificationPayload struct {
	UserIds   []string `json:"UserIds" validate:"required"`
	Title     string   `json:"Title" validate:"required"`
	Body      string   `json:"Body" validate:"required"`
	EventType string   `json:"EventType"`
}

type TemplateData struct {
	CustomerName string
	SalonName    string
	BarberName   string
	Date         string
	StartTime    string
	EndTime      string
	BookingId    string
	Extra        map[string]string
}
