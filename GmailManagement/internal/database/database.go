package database

import (
	"GmailManagement/internal/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service interface {
	Health() map[string]string
	StoreUser(users models.User) error
	GetUserByEmail(email string) (*models.User, error)
	UpdateUser(users models.User) error
	StoreEmail(email models.RawEmails) error
	GetAllRawEmails() ([]models.RawEmails, error)
	GetAllCategorizedEmails() ([]models.CategorizedEmail, error)
}

type service struct {
	db *mongo.Client
}

var (
	database         = os.Getenv("BLUEPRINT_DB_DATABASE")
	connectionString = os.Getenv("MONGO_CONNECTION_URL")
)

func New() Service {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(connectionString))

	if err != nil {
		log.Fatal(err)

	}
	return &service{
		db: client,
	}
}

func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := s.db.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("db down: %v", err)
	}

	return map[string]string{
		"message": "It's healthy",
	}
}

func (s *service) StoreUser(user models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	collection := s.db.Database(database).Collection("User")
	_, err := collection.InsertOne(ctx, user)
	if err != nil {
		// Print the error to see the exact message
		log.Printf("Insert error: %v", err)
		return err
	}
	return err
}

func (s *service) GetUserById(id string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	user := models.User{}
	err := s.db.Database(database).Collection("User").FindOne(ctx, bson.M{"id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *service) GetUserByEmail(email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	user := models.User{}
	err := s.db.Database(database).Collection("User").FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *service) UpdateUser(user models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	collection := s.db.Database(database).Collection("User")
	_, err := collection.ReplaceOne(ctx, bson.M{"email": user.Email}, user)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) StoreEmail(email models.RawEmails) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	collection := s.db.Database(database).Collection("RawEmails")
	_, err := collection.InsertOne(ctx, email)
	if err != nil {
		// Print the error to see the exact message
		log.Printf("Insert error: %v", err)
		return err
	}
	return err
}

func (s *service) GetAllRawEmails() ([]models.RawEmails, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	rawEmails := []models.RawEmails{}
	collection := s.db.Database(database).Collection("RawEmails")
	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	if err = cur.All(ctx, &rawEmails); err != nil {
		return nil, err
	}
	return rawEmails, nil
}

func (s *service) GetAllCategorizedEmails() ([]models.CategorizedEmail, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	rawEmails := []models.CategorizedEmail{}
	collection := s.db.Database(database).Collection("CategorizedEmails1")
	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	if err = cur.All(ctx, &rawEmails); err != nil {
		return nil, err
	}
	return rawEmails, nil
}
