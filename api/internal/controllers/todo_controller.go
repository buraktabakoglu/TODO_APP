package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"


	"github.com/buraktabakoglu/GOLANGAPPX/api/pkg/auth"
	
	"github.com/buraktabakoglu/GOLANGAPPX/api/pkg/models"
	formaterror "github.com/buraktabakoglu/GOLANGAPPX/api/pkg/utils"

	
	"github.com/gin-gonic/gin"
	//"github.com/gorilla/mux"
)


// @Summary Create a Todo
// @Description Creates a Todo item and assigns it to the authenticated user
// @Produce json
// @Security ApiKeyAuth
// @Param body body models.Todo true "Todo content and status"
// @Success 201 json models.Todo
// @Router /api/todos [post]
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


	err = server.DB.Debug().Model(models.User{}).Where("id = ?", uid).Take(&user).Error
	if err != nil {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}

	todo.AuthorID = uid //the authenticated user is the one creating the post

	todo.Prepare()

	errorMessages := todo.Validate()
	if len(errorMessages) > 1 {
		errList["StatusUnprocessableEntity"] = "StatusUnprocessableEntity"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	postCreated, err := todo.CreateTodo(server.DB)
	if err != nil {
		errList["InternalServerError"] = "InternalServerError"
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  errList,
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"status":   http.StatusCreated,
		"response": postCreated,
	})

}


	
// @Summary Get Todos
// @Description Retrieves all todos created by the authenticated user
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} models.Todo
// @Router /api/todos [get]
func (server *Server) GetTodos(c *gin.Context) {
    userID, err := auth.ExtractTokenID(c.Request)
    if err != nil {
        errList["Unauthorized"] = "Unauthorized"
        c.JSON(http.StatusUnauthorized, gin.H{
            "status": http.StatusUnauthorized,
            "error":  errList,
        })
        return
    }
    todo := models.Todo{}
    todos, err := todo.FindTodosByUserID(server.DB, userID)
    if err != nil {
        errList["No_todo"] = "No Todo Found"
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

// @Summary Get Todo by ID
// @Description Retrieves a todo by ID
// @Produce json
// @Param id path uint64 true "Todo ID"
// @Success 200 {object} models.Todo
// @Router /api/todos/{id} [get]
// @Security ApiKeyAuth
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
// UpdateATodo godoc
// @Summary Update a Todo by ID
// @Description Update a Todo by ID
// @Tags Todos
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Todo ID"
// @Param todoUpdatee body models.Todo true "Update Todo"
// @Success 200 {object} models.Todo
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/todos/{id} [put]
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

	todoA := models.Todo{}
	err = server.DB.Debug().Model(models.Todo{}).Where("id = ? ", pid).Take(&todoA).Error
	if err != nil {
		errList["no_todo"] = "No todo found"
		c.JSON(http.StatusNotFound,gin.H{
			"status":http.StatusNotFound,
			"error":errList,
		})
		return
	}

	if uid != todoA.AuthorID {
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
	
	todoUpdatee := models.Todo{}
	err = json.Unmarshal(body, &todoUpdatee)
	if err != nil {
		errList["Unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	todoUpdatee.ID = todoA.ID
	todoUpdatee.AuthorID = todoA.AuthorID

	todoUpdatee.Prepare()
	errorMessages := todoUpdatee.Validate()
	if len (errorMessages) > 0 {
		errList = errorMessages
		
		errList["Unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}


	todoUpdated, err := todoUpdatee.UpdateATodo(server.DB)
	if err != nil {
		formaterror := formaterror.FormatError(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":http.StatusInternalServerError,
			"error":formaterror,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": todoUpdated,
	})





	
}
// DeleteATodo godoc
// @Summary Delete a Todo by ID
// @Description Delete a Todo by ID
// @Tags Todos
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Todo ID"
// @Success 200
// @Failure 400
// @Failure 401 
// @Failure 404 
// @Failure 500 
// @Router /api/todos/{id} [delete]
func (server *Server) DeleteATodo(c *gin.Context) {

	todoID := c.Param("id")

	pid, err := strconv.ParseUint(todoID,10, 64)
	if err != nil {
		errList["Invalid_request"] = "Invalid request"
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

	_, err = todo.DeleteATodo(server.DB,pid,uid)
	if err != nil {
		errList["other_error"] = "Try again"
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":http.StatusInternalServerError,
			"error":errList,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": "Todo deleted",
	})
	

}
