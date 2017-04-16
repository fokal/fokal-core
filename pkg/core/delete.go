package core

import (
	"net/http"
	"reflect"

	"github.com/sprioc/composer/pkg/model"
	"github.com/sprioc/composer/pkg/sql"
	"github.com/sprioc/composer/pkg/rsp"
)

func DeleteImage(requestFrom model.Ref, image model.Ref) rsp.Response {
	// checking if the user has permission to delete the item
	valid, err := sql.Permissions(requestFrom, model.CanDelete, image)
	if err != nil {
		return rsp.Response{Code: http.StatusInternalServerError, Message: "Unable to retrieve user permissions."}
	}
	if !valid {
		return rsp.Response{Code: http.StatusForbidden, Message: "User does not have permission to delete item."}
	}

	err = sql.DeleteImage(image)
	if err != nil {
		return rsp.Response{Code: http.StatusInternalServerError,
			Message: "Unable to delete user."}
	}
	return rsp.Response{Code: http.StatusAccepted}

}

func DeleteUser(requestFrom model.Ref, user model.Ref) rsp.Response {
	// checking if the user has permission to delete the item
	valid, err := sql.Permissions(requestFrom, model.CanDelete, user)
	if err != nil {
		return rsp.Response{Code: http.StatusInternalServerError,
			Message: "Unable to retrieve user permissions."}
	}
	if !valid {
		return rsp.Response{Code: http.StatusForbidden,
			Message: "User does not have permission to delete item."}
	}

	err = sql.DeleteUser(user)
	if err != nil {
		return rsp.Response{Code: http.StatusInternalServerError,
			Message: "Unable to delete user."}
	}
	return rsp.Response{Code: http.StatusAccepted}

}

func inRef(item model.Ref, collection []model.Ref) bool {
	for _, x := range collection {
		if reflect.DeepEqual(x, item) {
			return true
		}
	}
	return false
}
