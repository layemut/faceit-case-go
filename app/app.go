package app

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/layemut/faceit-case-go/handlers"
	"github.com/layemut/faceit-case-go/notify"
	"github.com/layemut/faceit-case-go/service"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// App is the core of the application holding all the necessary components
type App struct {
	Router      *gin.Engine
	MRouter     *gin.Engine
	MongoClient *mongo.Client
	Config      *appConfig
	PubSub      *notify.Pubsub
}

// appConfig is the configuration for the application
type appConfig struct {
	MongoURI   string `mapstructure:"MONGO_URI"`
	DBName     string `mapstructure:"MONGO_DB_NAME"`
	ServerPort string `mapstructure:"SERVER_PORT"`
	MPort      string `mapstructure:"MANAGEMENT_PORT"`
}

// Initialize is the function to initialize the application, loads config connects to database etc.
func (a *App) Initialize() {
	a.Config = loadConfig()
	a.PubSub = notify.New()

	a.Router = gin.Default()
	a.MongoClient = connectToMongo(a.Config.MongoURI)
	a.setRouters()
}

// Run is the function to start the server
func (a *App) Run() {
	go a.MRouter.Run(a.Config.MPort)
	_ = a.Router.Run(a.Config.ServerPort)
}

// StartNotificationService starts notification service to recieve user events and send notification
func (a *App) StartNotificationService() {
	notificationService := &service.NotificationService{
		PubSub: a.PubSub,
	}

	notificationService.SubscribeUserCreateEvent()
	notificationService.SubscribeUserUpdateEvent()
}

// loadConfig is a function to load configuration from config file app.env
func loadConfig() (config *appConfig) {
	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatal(err)
	}
	return
}

// connectToMongo is a function to connect to mongo database
func connectToMongo(mongoURI string) (mongoClient *mongo.Client) {
	clientOptions := options.Client().ApplyURI(mongoURI)

	mongoClient, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = mongoClient.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB!")
	return
}

// setRouters is a function to set all the paths and methods to routers
func (a *App) setRouters() {
	database := a.MongoClient.Database(a.Config.DBName)
	userCollection := database.Collection("users")

	a.Router = gin.Default()
	a.Router.POST("/user", handlers.SaveUser(userCollection, a.PubSub))
	a.Router.PUT("/user", handlers.UpdateUser(userCollection, a.PubSub))
	a.Router.GET("/user", handlers.ListUsers(userCollection))
	a.Router.DELETE("/user/:id", handlers.RemoveUser(userCollection))

	a.MRouter = gin.New()
	a.MRouter.GET("/health", healthCheck(database))
}

// healthCheck is a function to check the health of the application
func healthCheck(m *mongo.Database) func(c *gin.Context) {
	return func(c *gin.Context) {
		mongoHealth := gin.H{
			"status": "UP",
		}
		appHealth := gin.H{
			"status": "UP",
		}

		if err := m.Client().Ping(context.TODO(), nil); err != nil {
			mongoHealth["status"] = "DOWN"
		}

		c.JSON(200, gin.H{
			"mongo": mongoHealth,
			"app":   appHealth,
		})
	}
}
