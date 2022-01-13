package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type IMockCollection struct {
	mock.Mock
}

func (m *IMockCollection) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	args := m.Called(ctx, document, opts)
	return args.Get(0).(*mongo.InsertOneResult), args.Error(1)
}
func (m *IMockCollection) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	args := m.Called(ctx, filter, update, opts)
	return args.Get(0).(*mongo.UpdateResult), args.Error(1)
}
func (m *IMockCollection) DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	args := m.Called(ctx, filter, opts)
	return args.Get(0).(*mongo.DeleteResult), args.Error(1)
}
func (m *IMockCollection) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	args := m.Called(ctx, filter, opts)
	return args.Get(0).(*mongo.Cursor), args.Error(1)
}

func TestSaveUser(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		assert := assert.New(t)
		mockColl := new(IMockCollection)
		mockColl.On("InsertOne", mock.Anything, mock.Anything, mock.Anything).Return(&mongo.InsertOneResult{}, nil)
		user := &User{
			Password: "password",
		}

		err := SaveUser(mockColl, user)
		assert.Nil(err)
		assert.NotNil(user.ID)

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("password"))
		assert.Nil(err)
	})

	t.Run("Error", func(t *testing.T) {
		assert := assert.New(t)
		mockColl := new(IMockCollection)
		mockColl.On("InsertOne", mock.Anything, mock.Anything, mock.Anything).Return(&mongo.InsertOneResult{}, errors.New("cannot insert!"))
		user := &User{
			Password: "password",
		}

		err := SaveUser(mockColl, user)
		assert.NotNil(err)
		assert.Equal("cannot insert!", err.Error())
	})
}

func TestUpdateUser(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		assert := assert.New(t)
		mockColl := new(IMockCollection)
		mockColl.On("UpdateOne", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&mongo.UpdateResult{}, nil)
		user := &User{
			Password: "password",
		}

		err := UpdateUser(mockColl, user)
		assert.Nil(err)

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("password"))
		assert.Nil(err)
	})

	t.Run("Fail", func(t *testing.T) {
		assert := assert.New(t)
		mockColl := new(IMockCollection)
		mockColl.On("UpdateOne", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&mongo.UpdateResult{}, errors.New("cannot update!"))
		user := &User{
			Password: "password",
		}

		err := UpdateUser(mockColl, user)
		assert.NotNil(err)
		assert.Equal("cannot update!", err.Error())
	})
}

func TestRemoveUser(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		assert := assert.New(t)
		mockColl := new(IMockCollection)
		mockColl.On("DeleteOne", mock.Anything, mock.Anything, mock.Anything).Return(&mongo.DeleteResult{}, nil)

		err := RemoveUser(mockColl, "1")
		assert.Nil(err)
	})

	t.Run("Fail", func(t *testing.T) {
		assert := assert.New(t)
		mockColl := new(IMockCollection)
		mockColl.On("DeleteOne", mock.Anything, mock.Anything, mock.Anything).Return(&mongo.DeleteResult{}, errors.New("cannot remove!"))

		err := RemoveUser(mockColl, "1")
		assert.NotNil(err)
		assert.Equal("cannot remove!", err.Error())
	})
}
