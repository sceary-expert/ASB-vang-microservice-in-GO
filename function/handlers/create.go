package handlers

import (
	"fmt"
	"net/http"

	"example.com/core/pkg/log"
	"example.com/core/types"
	"example.com/core/utils"
	"example.com/function/database"
	domain "example.com/function/dto"
	socialModels "example.com/function/models"
	service "example.com/function/services"
	"github.com/gofiber/fiber/v2"
)

// CreateUserRelHandle handle create a new userRel
func CreateUserRelHandle(c *fiber.Ctx) error {

	// Create the model object
	model := new(domain.UserRel)
	if err := c.BodyParser(model); err != nil {
		errorMessage := fmt.Sprintf("Parse UserRel Error %s", err.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/parseModel", "Error happened while parsing model!"))
	}
	// Create service
	userRelService, serviceErr := service.NewUserRelService(database.Db)
	if serviceErr != nil {
		log.Error("NewUserRelService %s", serviceErr.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/userRelService", "Error happened while creating userRelService!"))
	}

	if err := userRelService.SaveUserRel(model); err != nil {
		errorMessage := fmt.Sprintf("Save UserRel Error %s", err.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/saveUserRel", "Error happened while saving UserRel!"))
	}

	return c.JSON(fiber.Map{
		"objectId": model.ObjectId.String(),
	})

}

// FollowHandle handle create a new userRel
func FollowHandle(c *fiber.Ctx) error {
	// a -> b
	fmt.Println("at follow handle function")
	// Create the model object
	model := new(socialModels.FollowModel)
	fmt.Println("model has been assigned")
	if err := c.BodyParser(model); err != nil {
		errorMessage := fmt.Sprintf("Parse UserRel Error %s", err.Error())
		log.Error(errorMessage)
		fmt.Println("error ", errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/parseModel", "Error happened while parsing model!"))
	}

	// Create service
	fmt.Println("creating service")
	userRelService, serviceErr := service.NewUserRelService(database.Db)
	if serviceErr != nil {
		log.Error("NewUserRelService %s", serviceErr.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/userRelService", "Error happened while creating userRelService!"))
	}
	fmt.Println("creating current user") //the code is coming till line no. 66 as we have added a dummy rel meta model to pass the body parser
	// currentUser, ok := c.Locals(types.UserCtxName).(types.UserContext)
	// if !ok {
	// 	log.Error("[FollowHandle] Can not get current user")
	// 	return c.Status(http.StatusBadRequest).JSON(utils.Error("invalidCurrentUser",
	// 		"Can not get current user"))
	// }
	//creating a demo user current user

	currentUser := types.UserContext{

		DisplayName: "Current User",
		Avatar:      "Current User Avatar",
	}

	// Left User Meta
	fmt.Println("creating left user meta")
	leftUserMeta := domain.UserRelMeta{
		// UserId:   currentUser.UserID,
		FullName: currentUser.DisplayName,
		Avatar:   currentUser.Avatar,
	}

	// Right User Meta
	fmt.Println("creating right user meta")
	rightUserMeta := domain.UserRelMeta{
		// UserId:   model.RightUser.UserId,
		FullName: model.RightUser.FullName,
		Avatar:   model.RightUser.Avatar,
	}

	// Store the relation
	fmt.Println("storing the relation")
	if err := userRelService.FollowUser(leftUserMeta, rightUserMeta, model.CircleIds, []string{"status:follow"}); err != nil {
		errorMessage := fmt.Sprintf("Save UserRel Error %s", err.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/saveUserRel", "Error happened while saving UserRel!"))
	}

	// Create notification
	fmt.Println("creating notification")
	go sendFollowNotification(model, getUserInfoReq(c))
	// Increase user follow count
	fmt.Println("increasing follow count")
	go increaseUserFollowCount(currentUser.UserID, 1, getUserInfoReq(c))
	// Increase user follower count
	fmt.Println("increasing follower count")
	go increaseUserFollowerCount(model.RightUser.UserId, 1, getUserInfoReq(c))

	return c.SendStatus(http.StatusOK)
}
