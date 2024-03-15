package webService

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/DoubleH7/presenceHoursLog/database"
	"github.com/labstack/echo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetUserbyname(client *mongo.Client) func(c echo.Context) error {
	return func(c echo.Context) error {
		logger := getWeblogger()

		name := c.Param("name")
		if name == "" {
			c.String(http.StatusBadRequest, "name not specified correctly")
		}

		user, err := database.FindByName(client, name)

		if err == mongo.ErrNoDocuments {
			return c.String(http.StatusBadRequest, "user not found")
		}
		if err != nil {
			logger.Println("unknown error: ", err.Error())
			return c.String(http.StatusInternalServerError, "Something went wrong on our side!")
		}

		return c.JSON(http.StatusOK, user)
	}
}

func GetUserbyid(client *mongo.Client) func(c echo.Context) error {
	return func(c echo.Context) error {
		logger := getWeblogger()

		id := c.Param("id")
		if id == "" {
			c.String(http.StatusBadRequest, "id not specified correctly")
		}
		user, err := database.FindByID(client, id)

		if err == primitive.ErrInvalidHex {
			return c.String(http.StatusBadRequest, "invalid id")
		}

		if err == mongo.ErrNoDocuments {
			return c.String(http.StatusBadRequest, "user not found")
		}

		if err != nil {
			logger.Println("unknown error: ", err.Error())
			return c.String(http.StatusInternalServerError, "Something went wrong on our side!")
		}

		return c.JSON(http.StatusOK, user)
	}
}

func CreateSitting(client *mongo.Client) func(c echo.Context) error {
	return func(c echo.Context) error {
		logger := getWeblogger()

		id := c.Param("id")
		if id == "" {
			c.String(http.StatusBadRequest, "id not specified correctly")
		}

		active, err := database.SittingActive(client, id)

		if err == primitive.ErrInvalidHex {
			return c.String(http.StatusBadRequest, "invalid id")
		}

		if err == mongo.ErrNoDocuments {
			return c.String(http.StatusBadRequest, "user not found")
		}

		if err != nil {
			logger.Println("unknown error: ", err.Error())
			return c.String(http.StatusInternalServerError, "Something went wrong on our side!")
		}

		if active {
			return c.String(http.StatusNotAcceptable, "User already has an active Session")
		}

		name, err := database.AddSittingGiveName(client, id)

		if err == mongo.ErrNoDocuments {
			return c.String(http.StatusBadRequest, "User not found")
		}
		if err != nil {
			logger.Println("new sitting not added for: ", id)
			return c.String(http.StatusInternalServerError, "something has gone wrong on our side!")
		}
		return c.String(http.StatusOK, fmt.Sprintf("A sitting was started for %s", name))
	}
}

func StopSitting(client *mongo.Client) func(c echo.Context) error {
	return func(c echo.Context) error {
		logger := getWeblogger()
		id := c.Param("id")

		active, err := database.SittingActive(client, id)

		if err == primitive.ErrInvalidHex {
			return c.String(http.StatusBadRequest, "invalid id")
		}

		if err == mongo.ErrNoDocuments {
			return c.String(http.StatusBadRequest, "user not found")
		}

		if err != nil {
			logger.Println("unknown error: ", err.Error())
			return c.String(http.StatusInternalServerError, "Something went wrong on our side!")
		}

		if !active {
			return c.String(http.StatusNotAcceptable, "User already has no active sessions")
		}

		name, dur, err := database.AddEndGiveInfo(client, id)

		if err == mongo.ErrNoDocuments {
			return c.String(http.StatusBadRequest, "User not found")
		}
		if err != nil {
			logger.Println("Sitting not stopped for: ", id)
			return c.String(http.StatusInternalServerError, "something has gone wrong on our side!")
		}
		return c.String(http.StatusOK, fmt.Sprintf("%s's sitting was stopped after %v", name, dur))
	}
}

func GetUsers(client *mongo.Client) func(c echo.Context) error {
	return func(c echo.Context) error {
		logger := getWeblogger()

		opts := options.Find().SetProjection(bson.D{
			{"_id", 1},
			{"name", 1},
			{"age", 1},
		})

		cur, err := client.Database("presenceLog").Collection("users").Find(
			context.TODO(),
			bson.D{{}},
			opts,
		)

		if err != nil {
			logger.Println("error when fetching users from database: ", err.Error())
			return c.String(http.StatusInternalServerError, "something has gone wrong on our side!")
		}

		var results []struct {
			ID       primitive.ObjectID `json:"id" bson:"_id"`
			FullName string             `json:"name" bson:"name"`
		}
		cur.All(context.TODO(), &results)

		return c.JSON(http.StatusOK, results)
	}
}

func ServerAlive(client *mongo.Client) func(c echo.Context) error {
	return func(c echo.Context) error {
		return c.String(http.StatusOK, "Server is responsive and listening!")
	}
}

func CreateUser(client *mongo.Client) func(c echo.Context) error {
	return func(c echo.Context) error {
		datajs, err := io.ReadAll(c.Request().Body)
		defer c.Request().Body.Close()

		if err != nil {
			return c.String(http.StatusBadRequest, "failed reading request data")
		}

		user := database.User{}
		err = json.Unmarshal(datajs, &user)
		if err != nil {
			return c.String(http.StatusBadRequest, "failed unmarshalling request data")
		}
		user.Sittings = make([]database.Sitting, 0)

		result, err := client.Database("presenceLog").Collection("users").InsertOne(context.TODO(), user)
		if err != nil {
			fmt.Println(err.Error())
			return c.String(http.StatusInternalServerError, "user creation failed due to DB issues")
		}

		newId := result.InsertedID
		return c.String(http.StatusOK, fmt.Sprintf("New user successfully created with the ID : %s", newId))
	}
}
