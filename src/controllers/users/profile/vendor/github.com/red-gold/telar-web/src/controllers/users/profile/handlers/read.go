package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	handler "github.com/openfaas-incubator/go-function-sdk"
	server "github.com/red-gold/telar-core/server"
	utils "github.com/red-gold/telar-core/utils"
	models "github.com/red-gold/telar-web/src/models/users"
	service "github.com/red-gold/telar-web/src/services/users"
	uuid "github.com/satori/go.uuid"
)

// ReadProfileHandle a function invocation
func ReadProfileHandle(db interface{}) func(server.Request) (handler.Response, error) {

	return func(req server.Request) (handler.Response, error) {
		userId := req.GetParamByName("userId")
		userUUID, uuidErr := uuid.FromString(userId)
		if uuidErr != nil {
			return handler.Response{StatusCode: http.StatusBadRequest,
					Body: utils.MarshalError("parseUUIDError", "Can not parse user id!")},
				nil
		}
		// Create service
		userProfileService, serviceErr := service.NewUserProfileService(db)
		if serviceErr != nil {
			return handler.Response{StatusCode: http.StatusInternalServerError}, serviceErr
		}

		foundUser, err := userProfileService.FindByUserId(userUUID)
		if err != nil {
			return handler.Response{StatusCode: http.StatusInternalServerError}, err
		}

		profileModel := models.MyProfileModel{
			ObjectId:       foundUser.ObjectId,
			FullName:       foundUser.FullName,
			Avatar:         foundUser.Avatar,
			Banner:         foundUser.Banner,
			TagLine:        foundUser.TagLine,
			Birthday:       foundUser.Birthday,
			WebUrl:         foundUser.WebUrl,
			CompanyName:    foundUser.CompanyName,
			FacebookId:     foundUser.FacebookId,
			InstagramId:    foundUser.InstagramId,
			TwitterId:      foundUser.TwitterId,
			AccessUserList: foundUser.AccessUserList,
			Permission:     foundUser.Permission,
		}

		body, marshalErr := json.Marshal(profileModel)
		if marshalErr != nil {
			errorMessage := fmt.Sprintf("{error: 'Error while marshaling userProfile: %s'}",
				marshalErr.Error())
			return handler.Response{StatusCode: http.StatusBadRequest, Body: []byte(errorMessage)},
				marshalErr

		}
		return handler.Response{
			Body:       []byte(body),
			StatusCode: http.StatusOK,
		}, nil
	}
}

// ReadMyProfileHandle a function invocation to read authed user profile
func ReadMyProfileHandle(db interface{}) func(server.Request) (handler.Response, error) {

	return func(req server.Request) (handler.Response, error) {
		// Create service
		userProfileService, serviceErr := service.NewUserProfileService(db)
		if serviceErr != nil {
			return handler.Response{StatusCode: http.StatusInternalServerError}, serviceErr
		}

		foundUser, err := userProfileService.FindByUserId(req.UserID)
		if err != nil {
			return handler.Response{StatusCode: http.StatusInternalServerError}, err
		}

		profileModel := models.MyProfileModel{
			ObjectId:       foundUser.ObjectId,
			FullName:       foundUser.FullName,
			Avatar:         foundUser.Avatar,
			Banner:         foundUser.Banner,
			TagLine:        foundUser.TagLine,
			Birthday:       foundUser.Birthday,
			WebUrl:         foundUser.WebUrl,
			CompanyName:    foundUser.CompanyName,
			FacebookId:     foundUser.FacebookId,
			InstagramId:    foundUser.InstagramId,
			TwitterId:      foundUser.TwitterId,
			AccessUserList: foundUser.AccessUserList,
			Permission:     foundUser.Permission,
		}
		body, marshalErr := json.Marshal(profileModel)
		if marshalErr != nil {
			errorMessage := fmt.Sprintf("Error while marshaling userProfile: %s",
				marshalErr.Error())
			return handler.Response{StatusCode: http.StatusBadRequest, Body: utils.MarshalError("marshalUserProfileError", errorMessage)},
				marshalErr

		}
		return handler.Response{
			Body:       []byte(body),
			StatusCode: http.StatusOK,
		}, nil
	}
}
