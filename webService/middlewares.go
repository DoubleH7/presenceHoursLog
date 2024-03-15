package webService

import (
	"context"

	"github.com/DoubleH7/presenceHoursLog/database"
	"github.com/labstack/echo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func UserpassCheck(username, password string, c echo.Context) (bool, error) {
	logger := getWeblogger()

	client, err := database.ConnectDB()
	if err != nil {
		logger.Fatal("database connection failed: ", err)
	}
	cred := database.AdminCredentials{}
	err = client.Database("presenceLog").Collection("admins").FindOne(context.TODO(), bson.D{
		{Key: "username", Value: username},
	}).Decode(&cred)
	if err == mongo.ErrNoDocuments {
		logger.Printf("Failed login attempt with \n username: %s\npassword: %s\nfrom: %s", username, password, c.RealIP())
		return false, nil
	}
	if err != nil {
		logger.Printf("login error with \n username: %s\npassword: %s\nfrom: %s\nerror:%s\n",
			username, password, c.RealIP(), err.Error())
		return false, err
	}
	if cred.Password == password {
		logger.Printf("access granted for \nusername: %s\npassword: %s\nfrom: %s\n",
			username, password, c.RealIP())
		return true, nil
	}
	return false, nil
}
