package dbConfig

import "go.mongodb.org/mongo-driver/mongo"

var DATABASE_URL string
var DATABASE_NAME string
var DATABASE *mongo.Database

// Auth / User
const USERS_COLLECTION = "Users"
const USERROLES_COLLECTION = "UserRoles"

// Salon
const SALONS_COLLECTION = "Salons"
const BARBERS_COLLECTION = "Barbers"
const SALONSERVICES_COLLECTION = "SalonServices"
const HOURS_COLLECTION = "WorkingHours"
const GALLERY_COLLECTION = "Gallery"
const QUOTES_COLLECTION = "Quotes"
const DISTRICTS_COLLECTION = "Districts"
const SPECIALTIES_COLLECTION = "Specialties"

// Booking
const TIMESLOTS_COLLECTION = "TimeSlots"
const BOOKINGS_COLLECTION = "Bookings"
const WEEKLY_SCHEDULES_COLLECTION = "WeeklySchedules"
const SCHEDULE_BLOCKS_COLLECTION = "ScheduleBlocks"

// Review
const REVIEWS_COLLECTION = "Reviews"

// Notification
const NOTIFICATIONS_COLLECTION = "Notifications"
const TEMPLATES_COLLECTION = "Templates"
