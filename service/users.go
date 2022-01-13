package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type ICollection interface {
	InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error)
}

// User model
type User struct {
	ID        string    `json:"id,omitempty" bson:"_id,omitempty"`
	FirstName string    `json:"first_name,omitempty" bson:"first_name,omitempty"`
	LastName  string    `json:"last_name,omitempty" bson:"last_name,omitempty"`
	Nickname  string    `json:"nickname,omitempty" bson:"nickname,omitempty"`
	Password  string    `json:"password,omitempty" bson:"password,omitempty"`
	Email     string    `json:"email,omitempty" bson:"email,omitempty"`
	Country   string    `json:"country,omitempty" bson:"country,omitempty"`
	CreateAt  time.Time `json:"create_at,omitempty" bson:"create_at,omitempty"`
	UpdateAt  time.Time `json:"update_at,omitempty" bson:"update_at,omitempty"`
}

// PageRequest is a struct for pagination request used in ListUsers
type PageRequest struct {
	Page    int
	Size    int
	Country string
}

// SaveUser is a function for save user to database
// It will hash the password before save
func SaveUser(collection ICollection, user *User) error {
	user.ID = uuid.New().String()
	user.CreateAt = time.Now()
	user.UpdateAt = time.Now()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = collection.InsertOne(ctx, user)
	return err
}

// UpdateUser is a function for update user
func UpdateUser(collection ICollection, user *User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	user.UpdateAt = time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = collection.UpdateOne(ctx, bson.M{"_id": user.ID}, bson.M{"$set": user})
	return err
}

// RemoveUser is a function for remove user
func RemoveUser(collection ICollection, ID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := collection.DeleteOne(ctx, bson.M{"_id": ID})
	return err
}

// ListUsers is a function for list users based on pageRequest page, size and filter with country
func ListUsers(collection ICollection, pageRequest *PageRequest) ([]User, error) {

	var users []User
	var filter bson.M

	if len(pageRequest.Country) > 0 {
		filter = bson.M{"country": pageRequest.Country}
	}

	options := options.Find()
	options.SetSkip(int64(pageRequest.Size * (pageRequest.Page - 1)))
	options.SetLimit(int64(pageRequest.Size))
	options.SetSort(bson.D{primitive.E{Key: "create_at", Value: -1}})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cur, err := collection.Find(ctx, filter, options)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var user User
		err := cur.Decode(&user)
		if err != nil {
			continue
		}
		users = append(users, user)
	}
	return users, nil
}
