package mongoDBRepos

import (
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"ticken-ticket-service/models"
)

const UserCollectionName = "attendants"

type UserMongoDBRepository struct {
	baseRepository
}

func NewUserRepository(dbClient *mongo.Client, dbName string) *UserMongoDBRepository {
	userRepo := &UserMongoDBRepository{
		baseRepository{
			dbClient:       dbClient,
			dbName:         dbName,
			collectionName: UserCollectionName,
			primKeyField:   "attendant_id",
		},
	}
	users := userRepo.getCollection()
	storeContext, cancel := userRepo.generateOpSubcontext()
	defer cancel()
	_, err := users.Indexes().CreateOne(storeContext, mongo.IndexModel{
		Keys: bson.D{
			{"uuid", 1},
		},
		Options: &options.IndexOptions{
			Unique: &[]bool{true}[0],
		},
	})
	if err != nil {
		panic("Error creating user repository: " + err.Error())
	}

	return userRepo
}

func (r *UserMongoDBRepository) FindUser(userUUID uuid.UUID) *models.Attendant {
	findContext, cancel := r.generateOpSubcontext()
	defer cancel()

	users := r.getCollection()
	result := users.FindOne(findContext, bson.M{"uuid": userUUID})

	var foundUser models.Attendant
	err := result.Decode(&foundUser)

	if err != nil {
		return nil
	}

	return &foundUser
}

func (r *UserMongoDBRepository) FindAll() []*models.Attendant {
	findContext, cancel := r.generateOpSubcontext()
	defer cancel()

	attendants := r.getCollection()
	result, err := attendants.Find(findContext, bson.M{})
	if err != nil {
		return nil
	}

	var foundAttendants []*models.Attendant
	if err := result.Decode(&foundAttendants); err != nil {
		return nil
	}

	return foundAttendants
}
