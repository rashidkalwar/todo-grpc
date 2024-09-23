package main

import (
	"context"
	"log"
	"net"
	"strconv"

	pb "github.com/rashidkalwar/todo-grpc/protos/todo"
	"google.golang.org/grpc"
)

// Using this array as a temporary data storage solution instead of a database!
var todos []pb.Todo

type server struct {
	pb.UnimplementedTodoServiceServer
}

func (s *server) CreateTodo(context context.Context, createTodoRequest *pb.CreateTodoRequest) (*pb.Todo, error) {
	newTodo := &pb.Todo{
		Id:        strconv.Itoa(len(todos) + 1),
		Text:      createTodoRequest.Text,
		Completed: false,
	}
	todos = append(todos, pb.Todo{
		Id:        newTodo.Id,
		Text:      newTodo.Text,
		Completed: newTodo.Completed,
	})
	return newTodo, nil
}

func (s *server) ReadTodo(context context.Context, readTodoRequest *pb.ReadTodoRequest) (*pb.Todo, error) {
	for _, todo := range todos {
		if todo.Id == readTodoRequest.Id {
			return &pb.Todo{
				Id:        todo.Id,
				Text:      todo.Text,
				Completed: todo.Completed,
			}, nil
		}
	}
	return nil, nil
}

func (s *server) ReadAllTodos(nullRequest *pb.NullRequest, srv pb.TodoService_ReadAllTodosServer) error {
	// for i, _ := range todos {
	// 	srv.Send(&pb.Todos{
	// 		// Todos: todos[0 : i+1],
	// 	})
	// }
	for i := range todos {
		todosPtr := make([]*pb.Todo, i+1)
		for j := 0; j <= i; j++ {
			todosPtr[j] = &todos[j]
		}
		srv.Send(&pb.Todos{
			Todos: todosPtr,
		})
	}
	return nil
}

func (s *server) UpdateTodo(context context.Context, updateTodoRequest *pb.UpdateTodoRequest) (*pb.Todo, error) {
	for i, todo := range todos {
		if todo.Id == updateTodoRequest.Id {
			todos[i].Text = updateTodoRequest.Text
			todos[i].Completed = updateTodoRequest.Completed
			return &pb.Todo{
				Id:        todo.Id,
				Text:      todos[i].Text,
				Completed: todos[i].Completed,
			}, nil
		}
	}
	return nil, nil
}

func (s *server) DeleteTodo(context context.Context, deleteTodoRequest *pb.DeleteTodoRequest) (*pb.DeleteTodoResponse, error) {
	for i, todo := range todos {
		if todo.Id == deleteTodoRequest.Id {
			todos = append(todos[:i], todos[i+1:]...)
			return &pb.DeleteTodoResponse{
				Message: "Todo deleted Successfully!",
			}, nil
		}
	}
	return &pb.DeleteTodoResponse{
		Message: "Todo not found!",
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":9001")
	if err != nil {
		log.Fatalf("Failed to listen %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterTodoServiceServer(grpcServer, &server{})
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %s", err)
	}
}
