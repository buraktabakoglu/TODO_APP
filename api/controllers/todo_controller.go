package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/buraktabakoglu/GOLANGAPPX/api/auth"
	"github.com/buraktabakoglu/GOLANGAPPX/api/models"
	formaterror "github.com/buraktabakoglu/GOLANGAPPX/api/utils"
	"github.com/gin-gonic/gin"
	//"github.com/gorilla/mux"
)

func (server *Server) CreateTodo(c *gin.Context) {

	errList = map[string]string{}

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		errList["Invalid_body"] = "unable to get request"
		c.JSON(http.StatusUnprocessableEntity,gin.H{
			"status":http.StatusUnprocessableEntity,
			"error": errList,
		})
		return
	}
	todo := models.Todo{}

	err = json.Unmarshal(body, &todo)
	if err != nil {
		errList["unmarshal_error"] = "cannow unmarshal body"
		c.JSON(http.StatusUnprocessableEntity , gin.H{

			"status":http.StatusUnprocessableEntity,
			"error":errList,
		})
		return
	}

	uid, err := auth.ExtractTokenID(c.Request)
	if err != nil{
		errList["Unauhthorizeid"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized,gin.H{
			"status":http.StatusUnauthorized,
			"error":errList,
		})
		return
	}

	user := models.User{}

	err = server.DB.Debug().Model(models.User{}).Where("id = ?" ,uid).Take(&user).Error
	if err != nil {

		errList["unathorized"] = "Unathorized"
		c.JSON(http.StatusUnauthorized,gin.H{
			"status":http.StatusUnauthorized,
			"error":errList,
		})
		return
	}

	todo.AuthorID=uid

	todo.Prepare()
	errorMessages := todo.Validate()

	if len(errorMessages) > 0 {

		errList = errorMessages 
		c.JSON(http.StatusUnprocessableEntity,gin.H{
			"status":http.StatusUnprocessableEntity,
			"error":errList,
		})
		return
	}

	todoCreated, err := todo.CreateTodo(server.DB)

	if err != nil{

		errList := formaterror.FormatError(err.Error())
		c.JSON(http.StatusInternalServerError , gin.H{
			"status":http.StatusInternalServerError,
			"error":errList,
		})
		return
	}

	c.JSON(http.StatusCreated,gin.H{
		"status":http.StatusCreated,
		"response":todoCreated,
	})

}


	

	func (server *Server) GetTodos(c *gin.Context) {

		todo := models.Todo{}
	
		todos, err := todo.FindAllTodos(server.DB)
		if err != nil {
			errList["No_post"] = "No Post Found"
			c.JSON(http.StatusNotFound, gin.H{
				"status": http.StatusNotFound,
				"error":  errList,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status":   http.StatusOK,
			"response": todos,
		})
	}

func (server *Server) GetTodo(c *gin.Context) {

	todoID := c.Param("id")

	pid, err := strconv.ParseUint(todoID, 10 , 64)
	if err != nil {
		errList["Invalid_request"] = "Invalid Request"

		c.JSON(http.StatusBadRequest, gin.H{
			"status":http.StatusBadRequest,
			"error":errList,
		})
		return
	}

	todo := models.Todo{}


	postRec, err := todo.FindTodoByID(server.DB,pid)
	if err != nil {
		errList["no_todo"] = "Not todo found"
		c.JSON(http.StatusNotFound,gin.H{
			"status":http.StatusNotFound,
			"error":errList,
		})
		return
	}

	c.JSON(http.StatusOK,gin.H{
		"status":http.StatusOK,
		"response":postRec,
	})

	
}

func (server *Server) UpdateATodo(c *gin.Context) {

	errList = map[string]string{}

	postID := c.Param("id")

	pid, err := strconv.ParseUint(postID, 10, 64)
	if err != nil {
		errList["Invalied_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{

			"status":http.StatusBadRequest,
			"error":errList,
		})
		return
	}

	uid, err := auth.ExtractTokenID(c.Request)
	if err != nil {
		errList["unathorized"] = "unauthorrized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":http.StatusUnauthorized,
			"error":errList,
		})
		return
	}

	todoCont := models.Todo{}
	err = server.DB.Debug().Model(models.Todo{}).Where("id = ? ", pid).Take(&todoCont).Error
	if err != nil {
		errList["no_todo"] = "No todo found"
		c.JSON(http.StatusNotFound,gin.H{
			"status":http.StatusNotFound,
			"error":errList,
		})
		return
	}

	if uid != todoCont.AuthorID {
		errList["Unauthroized"] = "unauthorized"
		c.JSON(http.StatusUnauthorized,gin.H{
			"status":http.StatusUnauthorized,
			"error":errList,
		})
		return
	}

	
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		errList["Invalid_body"] = "Unable to get request"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	
	todo := models.Todo{}
	err = json.Unmarshal(body, &todo)
	if err != nil {
		errList["Unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	todo.ID = todoCont.ID 
	todo.AuthorID = todoCont.AuthorID

	todo.Prepare()
	errorMessages := todo.Validate()
	if len(errorMessages) > 0 {
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	todoAUpdated, err := todo.UpdateATodo(server.DB)
	if err != nil {
		errList := formaterror.FormatError(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  errList,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": todoAUpdated,
	})
}

func (server *Server) DeleteATodo(c *gin.Context) {

	todoID := c.Param("id")

	pid, err := strconv.ParseUint(todoID,10, 64)
	if err != nil {
		errList["Invalid_request"] = "Invalid requestr"
		c.JSON(http.StatusBadRequest,gin.H{

			"status":http.StatusBadRequest,
			"error":errList,

		})
		return
	}

	uid , err := auth.ExtractTokenID(c.Request)
	if err != nil {
		errList["unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":http.StatusUnauthorized,
			"error":errList,
		})
		return
	}
	
	todo := models.Todo{}
	err = server.DB.Debug().Model(models.Todo{}).Where("id = ? ",pid).Take(&todo).Error
	if err != nil {
		errList["no todo"] = "not todo found"
		c.JSON(http.StatusNotFound,gin.H{
			"status":http.StatusNotFound,
			"error":errList,
		})
		return
	}
	if uid != todo.AuthorID {
		errList["Unauthorized"] = "unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":http.StatusUnauthorized,
			"error":errList,
		})
		return
		
	}

	_, err = todo.DeleteATodo(server.DB)
	if err != nil {
		errList["other_error"] = "Try again"
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":http.StatusInternalServerError,
			"error":errList,
		})
		return
	}

}
