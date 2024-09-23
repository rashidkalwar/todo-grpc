package main

import (
	"context"
	"io"
	"log"
	"time"

	pb "github.com/rashidkalwar/todo-grpc/protos/todo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:9001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("Failed to connect to gRPC server")
	}

	defer conn.Close()

	c := pb.NewTodoServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Create two new Todos
	newTodo1, err := c.CreateTodo(ctx, &pb.CreateTodoRequest{
		Text: "Some task 1",
	})
	log.Printf("%v", newTodo1)

	newTodo2, err := c.CreateTodo(ctx, &pb.CreateTodoRequest{
		Text: "Some task 2",
	})
	log.Printf("%v", newTodo2)

	// Read one Todo
	getTodo, err := c.ReadTodo(ctx, &pb.ReadTodoRequest{Id: "1"})
	log.Printf("%v", getTodo)

	// Read all Todos
	todoStream, err := c.ReadAllTodos(ctx, &pb.NullRequest{})
	if err != nil {
		log.Fatal("Error calling function ReadAllTodos")
	}
	done := make(chan bool)
	// var todos []*pb.Todo
	go func() {
		for {
			resp, err := todoStream.Recv()
			if err == io.EOF {
				done <- true
				return
			}
			if err != nil {
				log.Fatalf("can not recieve %v", err)
			}
			// todos = resp.Todos
			log.Printf("Response recieved: %v", resp.Todos)
		}
	}()
	<-done

	// Update one Todo
	updateTodo, err := c.UpdateTodo(ctx, &pb.UpdateTodoRequest{
		Id:        "1",
		Text:      "Some task 1 Updated",
		Completed: true,
	})
	log.Printf("%v", updateTodo)

	// Delete one Todo
	deleteTodo, err := c.DeleteTodo(ctx, &pb.DeleteTodoRequest{
		Id: "1",
	})
	log.Printf("%v", deleteTodo)

	// Delete the already deleted Todo to see error handling
	deleteTodo2, err := c.DeleteTodo(ctx, &pb.DeleteTodoRequest{
		Id: "1",
	})
	log.Printf("%v", deleteTodo2)
}
