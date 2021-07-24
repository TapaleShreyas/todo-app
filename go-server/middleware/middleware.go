package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"go-server/models"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection

func init() {
	loadEnv()
	createDatabaseInstance()
}

func loadEnv() {
	err := godotenv.Load("local.env")
	if err != nil {
		log.Fatalf("Error in loading .env file: " + err.Error())
	}
}

func createDatabaseInstance() {
	connectionString := os.Getenv("DB_URI")
	dbName := os.Getenv("DB_NAME")
	collName := os.Getenv("DB_COLLECTION_NAME")
	clientOptions := options.Client().ApplyURI(connectionString)

	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB...")

	collection = client.Database(dbName).Collection(collName)

	fmt.Println("Collection instance created...")
}

//GetTask route to get task
func GetTask(resWriter http.ResponseWriter, req *http.Request) {
	resWriter.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	resWriter.Header().Set("Access-Control-Allow-Origin", "*")
	params := mux.Vars(req)
	payload, err := getTask(params["id"])
	if err != "" {
		var response models.Response
		response.Message = err
		json.NewEncoder((resWriter)).Encode(response)
	} else {
		json.NewEncoder(resWriter).Encode(payload)
	}

}

//GetAllTask route to get all tasks
func GetAllTask(resWriter http.ResponseWriter, req *http.Request) {
	resWriter.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	resWriter.Header().Set("Access-Control-Allow-Origin", "*")
	payload, errString := getAllTask()
	if payload == nil {
		var response models.Response
		response.Message = errString
		json.NewEncoder((resWriter)).Encode(response)
	} else {
		json.NewEncoder(resWriter).Encode(payload)
	}

}

//CreateTask route to create task
func CreateTask(resWriter http.ResponseWriter, req *http.Request) {
	resWriter.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	resWriter.Header().Set("Access-Control-Allow-Origin", "*")
	resWriter.Header().Set("Access-Control-Allow-Methods", "POST")
	resWriter.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	var task models.ToDoList
	_ = json.NewDecoder(req.Body).Decode(&task)
	createTask(&task)
	json.NewEncoder((resWriter)).Encode(task)
}

//DeleteTask route to delete task
func DeleteTask(resWriter http.ResponseWriter, req *http.Request) {
	resWriter.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	resWriter.Header().Set("Access-Control-Allow-Origin", "*")
	resWriter.Header().Set("Access-Control-Allow-Methods", "POST")
	resWriter.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	params := mux.Vars(req)
	result := deleteTask(params["id"])
	var response models.Response
	response.Message = result
	json.NewEncoder((resWriter)).Encode(response)
}

//DeleteAllTask route to delete all task's
func DeleteAllTask(resWriter http.ResponseWriter, req *http.Request) {
	resWriter.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	resWriter.Header().Set("Access-Control-Allow-Origin", "*")
	result := deleteAllTask()
	var response models.Response
	response.Message = result
	json.NewEncoder(resWriter).Encode(response)
}

//CompleteTask route to complete task
func CompleteTask(resWriter http.ResponseWriter, req *http.Request) {
	resWriter.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	resWriter.Header().Set("Access-Control-Allow-Origin", "*")
	params := mux.Vars(req)
	result := completeTask(params["id"])

	var response models.Response
	response.Message = result
	json.NewEncoder(resWriter).Encode(response)
}

//UndoTask route to undo task
func UndoTask(resWriter http.ResponseWriter, req *http.Request) {
	resWriter.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	resWriter.Header().Set("Access-Control-Allow-Origin", "*")
	params := mux.Vars(req)
	result := undoTask(params["id"])
	var response models.Response
	response.Message = result
	json.NewEncoder(resWriter).Encode(response)
}

//getTask get single task
func getTask(id string) (task bson.M, errString string) {
	ID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": ID}

	err := collection.FindOne(context.Background(), filter).Decode(&task)
	if err != nil {
		errString = "task not found"
	}
	return task, errString
}

//getAllTask get all tasks from database and return
func getAllTask() (tasks []primitive.M, errString string) {
	cursor, err := collection.Find(context.Background(), bson.D{{}})
	if err != nil || cursor.Current == nil {
		errString = "task not found"
	}

	for cursor.Next(context.Background()) {
		var task bson.M
		er := cursor.Decode(&task)
		if er != nil {
			log.Fatal(er)
		}
		tasks = append(tasks, task)
	}

	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}
	cursor.Close(context.Background())
	return tasks, errString
}

//createTask create task in database
func createTask(task *models.ToDoList) {
	insertedRecord, err := collection.InsertOne(context.Background(), task)
	if err != nil {
		log.Fatal(err)
	}
	newID := insertedRecord.InsertedID
	task.ID = newID.(primitive.ObjectID)
}

//delteTask delete task from database
func deleteTask(id string) string {
	ID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": ID}
	deleteCount, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}
	if deleteCount.DeletedCount > 0 {
		return "task deleted"
	}
	return "task not found"
}

//deleteAllTask delete all the task from database
func deleteAllTask() string {
	cnt, err := collection.DeleteMany(context.Background(), bson.D{{}}, nil)
	if err != nil {
		log.Fatal(err)
	}
	if cnt.DeletedCount > 0 {
		return "task deleted"
	}
	return "task not found"
}

//completeTask mark task as complete
func completeTask(id string) string {
	ID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": ID}
	update := bson.M{"$set": bson.M{"status": true}}
	updateCount, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	if updateCount.ModifiedCount > 0 {
		return "task marked as completed"
	}
	return "task not found"
}

//undoTask mark task as incomplete
func undoTask(id string) string {
	ID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": ID}
	update := bson.M{"$set": bson.M{"status": false}}
	updateCount, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	if updateCount.ModifiedCount > 0 {
		return "task undo successfully"
	}
	return "task not found"
}
