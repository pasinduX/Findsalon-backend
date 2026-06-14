package apiHandlers

import (
	"findsalon-backend/api"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func Router(app *fiber.App) {
	app.Use(cors.New())

	app.Get("/health", api.HealthApi)

	v1 := app.Group("/api/v1")
	v1.Get("/health", api.HealthApi)

	registerPublicRoutes(v1)
	registerAuthenticatedRoutes(v1)
}

func registerPublicRoutes(router fiber.Router) {
	auth := router.Group("/auth")
	auth.Post("/sync", api.AuthSyncApi)
	auth.Post("/CreateUser", api.CreateUserApi)
	auth.Post("/login", api.LoginApi)
	auth.Post("/register", api.RegisterApi)
	auth.Get("/FindUser", api.FindUserApi)
	auth.Get("/FindallUser", api.FindAllUserApi)
	auth.Get("/google/login", api.GoogleLoginApi)
	auth.Get("/google/callback", api.GoogleCallbackApi)
	auth.Post("/refresh", api.RefreshTokenApi)

	salon := router.Group("/salon")
	salon.Get("/FindSalon", api.FindSalonApi)
	salon.Get("/FindallSalon", api.FindAllSalonApi)
	salon.Get("/FindBarber", api.FindBarberApi)
	salon.Get("/FindallBarber", api.FindAllBarberApi)
	salon.Get("/FindSalonService", api.FindSalonServiceApi)
	salon.Get("/FindallSalonService", api.FindAllSalonServiceApi)
	salon.Get("/FindallWorkingHours", api.FindAllWorkingHoursApi)
	salon.Get("/FindallGallery", api.FindAllGalleryApi)
	salon.Get("/FindDistrict", api.FindDistrictApi)
	salon.Get("/FindallDistrict", api.FindAllDistrictApi)
	salon.Get("/FindallSpecialty", api.FindAllSpecialtyApi)
	salon.Get("/FindallQuote", api.FindAllQuoteApi)

	booking := router.Group("/booking")
	booking.Get("/FindTimeSlot", api.FindTimeSlotApi)
	booking.Get("/FindallTimeSlot", api.FindAllTimeSlotApi)
	booking.Get("/GetWeeklySchedule", api.GetWeeklyScheduleApi)
	booking.Get("/FindallScheduleBlock", api.FindAllScheduleBlockApi)
	booking.Get("/Availability", api.GetAvailabilityApi)
	booking.Get("/FindReview", api.FindReviewApi)
	booking.Get("/FindallReview", api.FindAllReviewApi)

	review := router.Group("/review")
	review.Get("/FindReview", api.FindReviewApi)
	review.Get("/FindallReviewBySalon", api.FindAllReviewBySalonApi)
	review.Get("/FindallReviewByBarber", api.FindAllReviewByBarberApi)
	review.Get("/RatingSummarySalon", api.FindSalonRatingSummaryApi)
	review.Get("/RatingSummaryBarber", api.FindBarberRatingSummaryApi)
}

func registerAuthenticatedRoutes(router fiber.Router) {
	noAuth := func(c *fiber.Ctx) error { return c.Next() }
	jwtAuth := noAuth
	adminOnly := fiber.Handler(noAuth)
	adminOrMod := fiber.Handler(noAuth)

	auth := router.Group("/auth", jwtAuth)
	auth.Put("/UpdateUser", api.UpdateUserApi)
	auth.Delete("/DeleteUser", api.DeleteUserApi)
	auth.Post("/logout", api.LogoutApi)

	user := router.Group("/user", jwtAuth)
	user.Get("/FindProfile", api.FindProfileApi)
	user.Put("/UpdateProfile", api.UpdateProfileApi)
	user.Delete("/DeleteProfile", api.DeleteProfileApi)
	user.Get("/FindallProfile", adminOrMod, api.FindAllProfileApi)
	user.Post("/UploadAvatar", api.UploadAvatarApi)
	user.Get("/DownloadProfile", adminOrMod, api.DownloadProfileApi)

	user.Post("/CreateUserRole", adminOnly, api.CreateUserRoleApi)
	user.Put("/UpdateUserRole", adminOnly, api.UpdateUserRoleApi)
	user.Delete("/DeleteUserRole", adminOnly, api.DeleteUserRoleApi)
	user.Get("/FindUserRole", api.FindUserRoleApi)
	user.Get("/FindallUserRole", adminOrMod, api.FindAllUserRoleApi)

	user.Get("/Dashboard", api.FindDashboardApi)
	user.Get("/SalonOwnerDashboard", api.FindSalonOwnerDashboardApi)
	user.Get("/BarberDashboard", api.FindBarberDashboardApi)
	user.Get("/AdminDashboard", api.FindAdminDashboardApi)

	salon := router.Group("/salon", jwtAuth)
	salon.Post("/CreateSalon", api.CreateSalonApi)
	salon.Put("/UpdateSalon", api.UpdateSalonApi)
	salon.Delete("/DeleteSalon", api.DeleteSalonApi)
	salon.Post("/UploadSalon", api.UploadSalonApi)
	salon.Get("/DownloadSalon", api.DownloadSalonApi)
	salon.Post("/UploadImage", api.UploadImageApi)

	salon.Post("/CreateBarber", api.CreateBarberApi)
	salon.Put("/UpdateBarber", api.UpdateBarberApi)
	salon.Delete("/DeleteBarber", api.DeleteBarberApi)
	salon.Get("/FindBarberByUser", api.FindBarberByUserApi)

	salon.Post("/CreateSalonService", api.CreateSalonServiceApi)
	salon.Put("/UpdateSalonService", api.UpdateSalonServiceApi)
	salon.Delete("/DeleteSalonService", api.DeleteSalonServiceApi)

	salon.Post("/CreateWorkingHours", api.CreateWorkingHoursApi)
	salon.Put("/UpdateWorkingHours", api.UpdateWorkingHoursApi)
	salon.Delete("/DeleteWorkingHours", api.DeleteWorkingHoursApi)

	salon.Post("/CreateGallery", api.CreateGalleryApi)
	salon.Delete("/DeleteGallery", api.DeleteGalleryApi)

	salon.Post("/CreateDistrict", adminOnly, api.CreateDistrictApi)
	salon.Put("/UpdateDistrict", adminOnly, api.UpdateDistrictApi)
	salon.Delete("/DeleteDistrict", adminOnly, api.DeleteDistrictApi)

	salon.Post("/CreateQuote", api.CreateQuoteApi)
	salon.Put("/UpdateQuote", api.UpdateQuoteApi)
	salon.Delete("/DeleteQuote", api.DeleteQuoteApi)

	booking := router.Group("/booking", jwtAuth)
	booking.Post("/CreateTimeSlot", api.CreateTimeSlotApi)
	booking.Post("/BulkCreateTimeSlot", api.BulkCreateTimeSlotApi)
	booking.Put("/UpdateTimeSlot", api.UpdateTimeSlotApi)
	booking.Delete("/DeleteTimeSlot", api.DeleteTimeSlotApi)

	booking.Post("/CreateWeeklySchedule", api.CreateWeeklyScheduleApi)

	booking.Post("/CreateScheduleBlock", api.CreateScheduleBlockApi)
	booking.Put("/UpdateScheduleBlock", api.UpdateScheduleBlockApi)
	booking.Delete("/DeleteScheduleBlock", api.DeleteScheduleBlockApi)

	booking.Post("/CreateBooking", api.CreateBookingApi)
	booking.Post("/CreateWalkInBooking", api.CreateWalkInBookingApi)
	booking.Put("/UpdateBooking", api.UpdateBookingApi)
	booking.Delete("/DeleteBooking", api.DeleteBookingApi)
	booking.Get("/FindBooking", api.FindBookingApi)
	booking.Get("/FindallBooking", api.FindAllBookingApi)
	booking.Get("/FindallBookingByUser", api.FindAllBookingByUserApi)
	booking.Get("/DownloadBooking", api.DownloadBookingApi)
	booking.Post("/DirectBooking", api.DirectBookingApi)

	booking.Post("/CreateReview", api.CreateReviewApi)
	booking.Put("/UpdateReview", api.UpdateReviewApi)
	booking.Delete("/DeleteReview", api.DeleteReviewApi)

	review := router.Group("/review", jwtAuth)
	review.Post("/CreateReview", api.CreateReviewApi)
	review.Put("/UpdateReview", api.UpdateReviewApi)
	review.Delete("/DeleteReview", api.DeleteReviewApi)
	review.Get("/FindallReview", adminOrMod, api.FindAllReviewApi)
	review.Get("/FindallReviewByUser", api.FindAllReviewByUserApi)
	review.Get("/DownloadReview", adminOrMod, api.DownloadReviewApi)

	notify := router.Group("/notify", jwtAuth)
	notify.Get("/FindNotification", api.FindNotificationApi)
	notify.Get("/FindallNotification", api.FindAllNotificationApi)
	notify.Put("/MarkRead", api.MarkNotificationReadApi)
	notify.Put("/MarkAllRead", api.MarkAllReadApi)
	notify.Get("/CountUnread", api.CountUnreadApi)
	notify.Delete("/DeleteNotification", api.DeleteNotificationApi)

	notify.Post("/CreateTemplate", adminOnly, api.CreateTemplateApi)
	notify.Put("/UpdateTemplate", adminOnly, api.UpdateTemplateApi)
	notify.Delete("/DeleteTemplate", adminOnly, api.DeleteTemplateApi)
	notify.Get("/FindTemplate", adminOrMod, api.FindTemplateApi)
	notify.Get("/FindallTemplate", adminOrMod, api.FindAllTemplateApi)

	notify.Post("/booking", api.SendBookingNotificationApi)
	notify.Post("/SendCustom", adminOrMod, api.SendCustomNotificationApi)
	notify.Post("/SendBulk", adminOrMod, api.SendBulkNotificationApi)
}
