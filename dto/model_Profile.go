package dto

type Profile struct {
	UserId          string `json:"UserId" bson:"UserId"`
	FullName        string `json:"FullName" bson:"FullName"`
	Email           string `json:"Email" bson:"Email"`
	Phone           string `json:"Phone" bson:"Phone"`
	AvatarUrl       string `json:"AvatarUrl" bson:"AvatarUrl"`
	GoogleAvatarUrl string `json:"GoogleAvatarUrl" bson:"GoogleAvatarUrl"`
	Provider        string `json:"Provider" bson:"Provider"`
	Role            string `json:"Role" bson:"Role"`
	IsActive        bool   `json:"IsActive" bson:"IsActive"`
}
