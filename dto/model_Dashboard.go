package dto

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	PageSize   int         `json:"pageSize"`
	TotalCount int         `json:"totalCount"`
	TotalPages int         `json:"totalPages"`
}

type BookingSummary struct {
	BookingId    string `json:"BookingId" bson:"BookingId"`
	SalonId      string `json:"SalonId" bson:"SalonId"`
	BarberId     string `json:"BarberId" bson:"BarberId"`
	Status       string `json:"Status" bson:"Status"`
	BookingType  string `json:"BookingType" bson:"BookingType"`
	CustomerName string `json:"CustomerName" bson:"CustomerName"`
}

type SalonSummary struct {
	SalonId  string `json:"SalonId" bson:"SalonId"`
	OwnerId  string `json:"OwnerId" bson:"OwnerId"`
	Name     string `json:"Name" bson:"Name"`
	IsActive bool   `json:"IsActive" bson:"IsActive"`
}

type BarberProfile struct {
	BarberId string `json:"BarberId" bson:"BarberId"`
	SalonId  string `json:"SalonId" bson:"SalonId"`
	UserId   string `json:"UserId" bson:"UserId"`
	Name     string `json:"Name" bson:"Name"`
}

type SalonMembership struct {
	SalonId   string `json:"SalonId"`
	SalonName string `json:"SalonName"`
	Role      string `json:"Role"`
	BarberId  string `json:"BarberId"`
}

type ProfileSummary struct {
	UserId   string `json:"UserId"`
	FullName string `json:"FullName"`
	Email    string `json:"Email"`
	Role     string `json:"Role"`
}

type UserDashboard struct {
	UserId          string           `json:"UserId"`
	FullName        string           `json:"FullName"`
	Email           string           `json:"Email"`
	AvatarUrl       string           `json:"AvatarUrl"`
	Role            string           `json:"Role"`
	TotalBookings   int              `json:"TotalBookings"`
	UpcomingSlots   int              `json:"UpcomingSlots"`
	CompletedVisits int              `json:"CompletedVisits"`
	CancelledCount  int              `json:"CancelledCount"`
	RecentBookings  []BookingSummary `json:"RecentBookings"`
}

type SalonOwnerDashboard struct {
	TotalSalons   int            `json:"TotalSalons"`
	TotalBarbers  int            `json:"TotalBarbers"`
	TotalBookings int            `json:"TotalBookings"`
	TodayBookings int            `json:"TodayBookings"`
	PendingSlots  int            `json:"PendingSlots"`
	TotalRevenue  float64        `json:"TotalRevenue"`
	Salons        []SalonSummary `json:"Salons"`
}

type BarberDashboard struct {
	SalonId        string            `json:"SalonId"`
	SalonName      string            `json:"SalonName"`
	Memberships    []SalonMembership `json:"Memberships"`
	TodaySlots     int               `json:"TodaySlots"`
	TodayBookings  int               `json:"TodayBookings"`
	TotalCompleted int               `json:"TotalCompleted"`
	AverageRating  float64           `json:"AverageRating"`
}

type AdminDashboard struct {
	TotalUsers    int `json:"TotalUsers"`
	TotalSalons   int `json:"TotalSalons"`
	ActiveSalons  int `json:"ActiveSalons"`
	TotalBarbers  int `json:"TotalBarbers"`
	TotalBookings int `json:"TotalBookings"`
}
