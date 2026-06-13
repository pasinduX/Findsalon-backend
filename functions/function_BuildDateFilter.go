package functions

import "go.mongodb.org/mongo-driver/bson"

func BuildDateFilter(date, fromDate, toDate string) bson.M {
	if date != "" {
		return bson.M{"Date": date}
	}
	if fromDate != "" && toDate != "" {
		return bson.M{"Date": bson.M{"$gte": fromDate, "$lte": toDate}}
	}
	if fromDate != "" {
		return bson.M{"Date": bson.M{"$gte": fromDate}}
	}
	if toDate != "" {
		return bson.M{"Date": bson.M{"$lte": toDate}}
	}
	return bson.M{}
}
