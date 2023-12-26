package main

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client              *mongo.Client
	collection          *mongo.Collection
	healthPostsCollection *mongo.Collection
)

type Patient struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	FullName       string             `bson:"full_name"`
	Email          string             `bson:"email"`
	CPF            string             `bson:"cpf"`
	PhoneNumber    string             `bson:"phone_number"`
	HealthPostIDs  []primitive.ObjectID `bson:"health_post_ids"`
}

type HealthPost struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	PatientID primitive.ObjectID `bson:"patient_id"`
	Content   string             `bson:"content"`
	Timestamp time.Time          `bson:"timestamp"`
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	mongodbURI := os.Getenv("MONGODB_URI")
	mongodbUsername := os.Getenv("MONGODB_USERNAME")
	mongodbPassword := os.Getenv("MONGODB_PASSWORD")

	router := gin.Default()

	initMongoDB(mongodbURI, mongodbUsername, mongodbPassword)

	router.LoadHTMLGlob("./templates/*")


    // Set the default route to index.html
    router.GET("/", func(c *gin.Context) {
        c.HTML(http.StatusOK, "index.html", gin.H{})
    })

    // Use /register to show the form and register the user
    router.GET("/register", showForm)
    router.POST("/register", registerUser)

    // Add other routes as needed
    router.GET("/profile/:id", displayProfile)
    router.POST("/add-health-post/:id", addHealthPost)

	router.Run(":8080")
}

func initMongoDB(uri, username, password string) {
	clientOptions := options.Client().ApplyURI(uri)
	clientOptions.Auth = &options.Credential{
		Username: username,
		Password: password,
	}
	var err error
	client, err = mongo.Connect(nil, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(nil, nil)
	if err != nil {
		log.Fatal(err)
	}

	collection = client.Database("petri-dish").Collection("patients")
	healthPostsCollection = client.Database("petri-dish").Collection("health_posts")
}

func getPatientByID(patientID string) (Patient, error) {
	var patient Patient
	objID, err := primitive.ObjectIDFromHex(patientID)
	if err != nil {
		return patient, err
	}

	err = collection.FindOne(nil, bson.M{"_id": objID}).Decode(&patient)
	return patient, err
}

func getHealthPostsByPatientID(patientID string) ([]HealthPost, error) {
	var healthPosts []HealthPost
	objID, err := primitive.ObjectIDFromHex(patientID)
	if err != nil {
		return healthPosts, err
	}

	cursor, err := healthPostsCollection.Find(nil, bson.M{"patient_id": objID})
	if err != nil {
		return healthPosts, err
	}
	defer cursor.Close(nil)

	for cursor.Next(nil) {
		var healthPost HealthPost
		if err := cursor.Decode(&healthPost); err != nil {
			return healthPosts, err
		}
		healthPosts = append(healthPosts, healthPost)
	}

	return healthPosts, nil
}

func showForm(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", gin.H{})
}

func registerUser(c *gin.Context) {
	fullName := c.PostForm("full_name")
	email := c.PostForm("email")
	cpf := c.PostForm("cpf")
	phoneNumber := c.PostForm("phone_number")

	result, err := collection.InsertOne(nil, bson.M{
		"full_name":      fullName,
		"email":          email,
		"cpf":            cpf,
		"phone_number":   phoneNumber,
		"health_post_ids": []primitive.ObjectID{},
	})
	if err != nil {
		log.Println(err)
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "Internal Server Error"})
		return
	}

	patientID := result.InsertedID.(primitive.ObjectID)
	c.Redirect(http.StatusSeeOther, "/profile/"+patientID.Hex())
}

func displayProfile(c *gin.Context) {
	patientID := c.Param("id")

	// Remove quotes from the patientID if present
	patientID = strings.Trim(patientID, "\"")

	patient, err := getPatientByID(patientID)
	if err != nil {
		log.Println(err)
		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": "Patient Not Found"})
		return
	}

	healthPosts, err := getHealthPostsByPatientID(patientID)
	if err != nil {
		log.Println(err)
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "Internal Server Error"})
		return
	}

	log.Printf("Health Posts for Patient %s: %+v", patientID, healthPosts)

	// Ensure that healthPosts is passed to the template
	c.HTML(http.StatusOK, "profile.html", gin.H{"patient": patient, "healthPosts": healthPosts})
}





func addHealthPost(c *gin.Context) {
	patientID := c.Param("id")

	content := c.PostForm("content")

	_, err := healthPostsCollection.InsertOne(nil, bson.M{
		"patient_id": patientID,
		"content":    content,
		"timestamp":  time.Now(),
	})
	if err != nil {
		log.Println(err)
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "Internal Server Error"})
		return
	}

	c.Redirect(http.StatusSeeOther, "/profile/"+url.QueryEscape(patientID))
}






