package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)




type Todo struct
{
	ID 			primitive.ObjectID 	`json:"id,omitempty" bson:"_id,omitempty"` //omitempty tells  not to use duplicate Key values
	Completed 	bool  	`json:"completed"`
	Body 		string 	`json:"body"`
}

var collection *mongo.Collection

func main() {
	
	
	fmt.Println(" Welcome to Go World")


err  := godotenv.Load(".env")
if err != nil {
	log.Fatal("Unable to load .env file",err)
}

MONGODB_URI := os.Getenv("MONGODB_URI")
fmt.Println(MONGODB_URI)
serverAPI := options.ServerAPI(options.ServerAPIVersion1)
clientOptions := options.Client().ApplyURI(MONGODB_URI).SetServerAPIOptions(serverAPI)
client, err := mongo.Connect(context.Background(),clientOptions)

defer client.Disconnect(context.Background())

if err != nil{
	log.Fatal("Error in connection part ",err)
}

err  = client.Ping(context.Background(),nil)
if err != nil{
log.Fatal(err)
}

fmt.Println("Connected to MongoDB ATLAS Successfully --> KPS")


collection = client.Database("golang_db").Collection("todos")

app := fiber.New()



PORT := os.Getenv("PORT")

if PORT == ""{
	PORT = "5000"
}

app.Get("/api/todos", getTodos)
app.Post("/api/todos", createToo)
app.Patch("/api/todos/:id", updateTodo)
app.Delete("/api/todos/:id", deleteTodo)

log.Fatal(app.Listen(":" + PORT))

}

func getTodos(c *fiber.Ctx) error{
	var todos []Todo
   
	cursor,err :=	collection.Find(context.Background(),bson.M{})

	if err != nil {
	return err
	}

	defer cursor.Close(context.Background())

		for cursor.Next(context.Background()){
			var todo Todo
			if err := cursor.Decode(&todo); err != nil{
				return err
			}
			todos = append(todos, todo)
		}
   return c.JSON(todos)
}

func createToo(c *fiber.Ctx) error{
	 todo := new(Todo)
	if err :=  c.BodyParser(todo); err != nil{
		return err
	}

	if todo.Body == "" {
		return c.Status(400).JSON(fiber.Map{"error"  :  "Todo body cannot be empty"})
	}

inserResults , err  := collection.InsertOne(context.Background(),todo)
 
if err != nil{
	return err
}

todo.ID = inserResults.InsertedID.(primitive.ObjectID)
return c.Status(200).JSON(todo)
}

func updateTodo(c *fiber.Ctx) error{
	todo := new(Todo)
	id := c.Params("id")
	if errdata  :=  c.BodyParser(todo); errdata != nil {
		return errdata
	}
	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil{
		return c.Status(400).JSON(fiber.Map{"error":"Invalid todo ID"})
	}

	filter := bson.M{"_id":objectID}
	update := bson.M{"$set":bson.M{"completed":true,"body":todo.Body}}

	_,err = collection.UpdateOne(context.Background(),filter,update)
    if err != nil{
		return err
	}

	return c.Status(200).JSON(fiber.Map{"Success":true})
}


func deleteTodo(c *fiber.Ctx) error{
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error":"Invalid todo ID"})
	}

	filter := bson.M{"_id" : objectID}
	_, err = collection.DeleteOne(context.Background(),filter)

	if err != nil{
		return err
	}

	return c.Status(200).JSON(fiber.Map{"Success": true})
}




