package mongoDBRepos

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"ticken-ticket-service/models"
)

const UserCollectionName = "users"

type UserMongoDBRepository struct {
	baseRepository
}

func NewUserRepository(dbClient *mongo.Client, dbName string) *UserMongoDBRepository {
	return &UserMongoDBRepository{
		baseRepository{
			dbClient:       dbClient,
			dbName:         dbName,
			collectionName: UserCollectionName,
		},
	}
}

func (r *UserMongoDBRepository) AddUser(user *models.User) error {
	storeContext, cancel := r.generateOpSubcontext()
	defer cancel()

	users := r.getCollection()
	_, err := users.InsertOne(storeContext, user)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserMongoDBRepository) FindUser(userUUID string) *models.User {
	findContext, cancel := r.generateOpSubcontext()
	defer cancel()

	users := r.getCollection()
	result := users.FindOne(findContext, bson.M{"uuid": userUUID})

	var foundUser models.User
	err := result.Decode(&foundUser)

	if err != nil {
		return nil
	}

	return &foundUser
}
