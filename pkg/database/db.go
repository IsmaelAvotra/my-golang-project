package database

import (
	"context"
	"errors"

	"github.com/IsmaelAvotra/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var DB *mongo.Database

// Connect database
func ConnectDatabase() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}

	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}

	DB = client.Database("my-project")
}

// For users
func GetUserByEmail(email string) (*models.User, error) {
	user := models.User{}
	err := DB.Collection("users").FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func GetUserByUsername(username string) (*models.User, error) {
	user := models.User{}
	err := DB.Collection("users").FindOne(context.TODO(), bson.M{"username": username}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func GetAllUsers() ([]models.User, error) {
	users := []models.User{}

	cursor, err := DB.Collection("users").Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		user := models.User{}
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func GetUserByID(id string) (*models.User, error) {
	objID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, err
	}

	user := models.User{}

	err = DB.Collection("users").FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&user)

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func DeleteUser(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return err
	}
	_, err = DB.Collection("users").DeleteOne(context.TODO(), bson.M{"_id": objID})
	if err != nil {
		return err
	}
	return nil
}

func UpdateUser(id string, update bson.M) error {
	objID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return err
	}
	_, err = DB.Collection("users").UpdateOne(context.TODO(), bson.M{"_id": objID}, bson.M{"$set": update})
	if err != nil {
		return err
	}
	return nil
}

// for university
func GetUnivByName(univName string) (*models.University, error) {
	university := models.University{}
	err := DB.Collection("universities").FindOne(context.TODO(), bson.M{"name": univName}).Decode(&university)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &university, nil
}

func GetUnivById(id string) (*models.University, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	university := models.University{}

	err = DB.Collection("universities").FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&university)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("university not found")
		}
		return nil, err
	}
	return &university, nil
}

func GetAllUniversities() ([]models.University, error) {
	universities := []models.University{}
	cursor, err := DB.Collection("universities").Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		university := models.University{}
		err := cursor.Decode(&university)

		if err != nil {
			return nil, err
		}
		universities = append(universities, university)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return universities, nil
}

func DeleteUniversity(id string) error {
	objId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return err
	}

	result, err := DB.Collection("universities").DeleteOne(context.TODO(), bson.M{"_id": objId})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("university not found")
	}
	return nil
}

func UpdateUniversity(id string, update bson.M) error {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	result, err := DB.Collection("universities").UpdateOne(context.TODO(), bson.M{"_id": objId}, bson.M{"$set": update})
	if err != nil {
		return err
	}
	if result.ModifiedCount == 0 {
		return errors.New("no changes made")
	}
	return nil
}
