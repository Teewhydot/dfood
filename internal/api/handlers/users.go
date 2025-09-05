package handlers

// import (
// 	"github.com/gin-gonic/gin"
// 	"sirteefyapps.com.ng/goroutines/internal/service"
// 	"sirteefyapps.com.ng/goroutines/pkg/errors"
// )

// type UserHandler struct {
// 	userService service.UserService
// }

// func NewUserHandler(userService service.UserService) *UserHandler {
// 	return &UserHandler{
// 		userService: userService,
// 	}
// }

// func (h *UserHandler) GetByID(c *gin.Context) {
// 	userID := c.Param("id")
	
// 	result := errors.HandleError(
// 		func() (interface{}, error) {
// 			return h.userService.GetByID(userID)
// 		},
// 		"fetching user by ID",
// 	)
// 	result.RespondWithJSON(c)
// }