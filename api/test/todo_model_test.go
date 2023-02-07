package test

import (
	"log"
	"testing"

	
	"github.com/buraktabakoglu/GOLANGAPPX/api/pkg/models"
	"gopkg.in/go-playground/assert.v1"
	
)





func TestFindAllTodos(t *testing.T) {
	

	err := refreshUserAndTodoTable()
	if err != nil {
		log.Fatalf("Error refreshing user and todo table %v\n", err)
	}
	_, _, err = seedUsersAndTodos()
	if err != nil {
		log.Fatalf("Error seeding user and todo  table %v\n", err)
	}
	todos, err := todoInstance.FindTodosByUserID(server.DB,todoInstance.AuthorID)
	if err != nil {
		t.Errorf("this is the error getting the todos: %v\n", err)
		return
	}
	assert.Equal(t, len(*todos), 0)
}

func TestSaveTodo(t *testing.T) {

	err := refreshUserAndTodoTable()
	if err != nil {
		log.Fatalf("Error user and todo refreshing table %v\n", err)
	}

	user, err := seedOneUser()
	if err != nil {
		log.Fatalf("Cannot seed user %v\n", err)
	}

	newTodo := models.Todo{
		Status:    "This is the status",
		Description:  "This is the description",
		AuthorID: uint32(user.ID),
	}
	savedTodo, err := newTodo.CreateTodo(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the Todo: %v\n", err)
		return
	}
	
	assert.Equal(t, newTodo.Status, savedTodo.Status)
	assert.Equal(t, newTodo.Description, savedTodo.Description)
	assert.Equal(t, newTodo.AuthorID, savedTodo.AuthorID)

}

func TestGetTodoByID(t *testing.T) {

	err := refreshUserAndTodoTable()
	if err != nil {
		log.Fatalf("Error refreshing user and Todo table: %v\n", err)
	}
	_,todo, err := seedOneUserAndOneTodo()
	if err != nil {
		log.Fatalf("Error Seeding table")
	}
	foundTodo, err := todoInstance.FindTodoByID(server.DB, todo.ID)
	if err != nil {
		t.Errorf("this is the error getting one user: %v\n", err)
		return
	}
	
	assert.Equal(t, foundTodo.Status, todo.Status)
	assert.Equal(t, foundTodo.Description, todo.Description)
}

func TestUpdateATodo(t *testing.T) {

	err := refreshUserAndTodoTable()
	if err != nil {
		log.Fatalf("Error refreshing user and todo table: %v\n", err)
	}
	_,todo, err := seedOneUserAndOneTodo()
	if err != nil {
		log.Fatalf("Error Seeding table")
	}
	todoUpdate := models.Todo{
		
		Status:  "modiUpdate",
		Description: "modiupdate@gmail.com",
		AuthorID: todo.AuthorID,
	}
	updatedTodo, err := todoUpdate.UpdateATodo(server.DB)
	if err != nil {
		t.Errorf("this is the error updating the user: %v\n", err)
		return
	}
	
	assert.Equal(t, updatedTodo.ID, todoUpdate.ID)
	
	
	
}

func TestDeleteATodo(t *testing.T) {

	err := refreshUserAndTodoTable()
	if err != nil {
		log.Fatalf("Error refreshing user and todo table: %v\n", err)
	}
	_,todo, err := seedOneUserAndOneTodo()
	if err != nil {
		log.Fatalf("Error Seeding tables")
	}
	isDeleted, err := todo.DeleteATodo(server.DB, todo.ID, todo.AuthorID)
	if err != nil {
		t.Errorf("this is the error updating the user: %v\n", err)
		return
	}
	
	assert.Equal(t, isDeleted, int64(1))
}

func TestDeleteUserTodos(t *testing.T) {

	err := refreshUserAndTodoTable()
	if err != nil {
		log.Fatalf("Error refreshing user and todo table: %v\n", err)
	}
	user, _, err := seedOneUserAndOneTodo()
	if err != nil {
		log.Fatalf("Error Seeding tables")
	}

	numberDeleted, err := todoInstance.DeleteUserTodos(server.DB, user.ID)
	if err != nil {
		t.Errorf("this is the error deleting the post: %v\n", err)
		return
	}
	assert.Equal(t, numberDeleted, int64(1))
}