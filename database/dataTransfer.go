package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func FindByID(client *mongo.Client, id string) (User, error) {

	var user User

	idDB, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return user, err
	}

	err = client.Database("presenceLog").Collection("users").FindOne(
		context.TODO(),
		bson.D{{"_id", idDB}},
	).Decode(&user)

	return user, err
}

func FindByName(client *mongo.Client, name string) (UserwID, error) {

	var user UserwID

	err := client.Database("presenceLog").Collection("users").FindOne(
		context.TODO(),
		bson.D{{"name", name}},
	).Decode(&user)

	return user, err
}

func SittingActive(client *mongo.Client, id string) (bool, error) {
	var user User

	idDB, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, err
	}

	opt := options.FindOne().SetProjection(bson.D{
		{"sittings", bson.D{{"$slice", -1}}},
	})
	err = client.Database("presenceLog").Collection("users").FindOne(
		context.TODO(),
		bson.D{{"_id", idDB}},
		opt,
	).Decode(&user)
	if err != nil {
		return false, err
	}

	if len(user.Sittings) == 0 {
		return false, nil
	}

	if user.Sittings[0].Start.After(user.Sittings[0].End) {
		return true, nil
	}

	return false, nil
}

func AddSittingGiveName(client *mongo.Client, id string) (string, error) {

	idDB, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return "", err
	}
	var user User
	var sitting Sitting
	sitting.Start = time.Now()

	opt := options.FindOneAndUpdate().SetProjection(bson.D{
		{"name", 1},
	})

	err = client.Database("presenceLog").Collection("users").FindOneAndUpdate(
		context.TODO(),
		bson.M{"_id": idDB},
		bson.D{
			{"$push", bson.D{{"sittings", sitting}}},
		},
		opt,
	).Decode(&user)

	return user.FullName, err
}

func AddEndGiveInfo(client *mongo.Client, id string) (string, time.Duration, error) {
	idDB, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return "", 0, err
	}
	var user User

	opt := options.FindOneAndUpdate()

	opt.SetProjection(bson.D{
		{"name", 1},
		{"sittings", bson.D{{"$slice", -1}}},
	})

	opt.SetArrayFilters(options.ArrayFilters{
		Filters: []interface{}{bson.D{{"i.end", new(time.Time)}}},
	})

	err = client.Database("presenceLog").Collection("users").FindOneAndUpdate(
		context.TODO(),
		bson.M{"_id": idDB},
		bson.D{
			{"$set", bson.D{
				{"sittings.$[i].end", time.Now()},
			}},
		},
		opt,
	).Decode(&user)

	return user.FullName, time.Since(user.Sittings[0].Start), err
}
