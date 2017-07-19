package retrieval

import (
	"net/http"

	"errors"

	"github.com/devinmcgloin/fokal/pkg/handler"
	"github.com/devinmcgloin/fokal/pkg/model"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

func UserHandler(store *handler.State, w http.ResponseWriter, r *http.Request) (handler.Response, error) {
	var rsp handler.Response
	username := mux.Vars(r)["ID"]

	ref, err := GetUserRef(store.DB, username)
	if err != nil {
		return rsp, err
	}

	user, err := GetUser(store, ref.Id)
	return handler.Response{
		Code: http.StatusOK,
		Data: user,
	}, nil
}

func UserImagesHandler(store *handler.State, w http.ResponseWriter, r *http.Request) (handler.Response, error) {
	var rsp handler.Response
	username := mux.Vars(r)["ID"]

	ref, err := GetUserRef(store.DB, username)
	if err != nil {
		return rsp, err
	}

	images, err := GetUserImages(store, ref.Id)
	return handler.Response{
		Code: http.StatusOK,
		Data: images,
	}, nil
}

func LoggedInUserHandler(store *handler.State, w http.ResponseWriter, r *http.Request) (handler.Response, error) {
	var rsp handler.Response

	val, ok := context.GetOk(r, "auth")
	if !ok {
		return rsp, handler.StatusError{
			Code: http.StatusUnauthorized,
			Err:  errors.New("Must be logged in to use this endpoint")}
	}

	usrRef := val.(model.Ref)
	user, err := GetUser(store, usrRef.Id)
	if err != nil {
		return rsp, err
	}

	return handler.Response{
		Code: http.StatusOK,
		Data: user,
	}, nil
}

func ImageHandler(store *handler.State, w http.ResponseWriter, r *http.Request) (handler.Response, error) {
	var rsp handler.Response

	id := mux.Vars(r)["ID"]

	ref, err := GetImageRef(store.DB, id)
	if err != nil {
		return rsp, err
	}

	img, err := GetImage(store, ref.Id)
	if err != nil {
		return rsp, err
	}

	return handler.Response{
		Code: http.StatusOK,
		Data: img,
	}, nil
}