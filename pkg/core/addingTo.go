package core

import (
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/sprioc/composer/pkg/model"
	"github.com/sprioc/composer/pkg/mongo"
	"github.com/sprioc/composer/pkg/redis"
	"github.com/sprioc/composer/pkg/refs"
	"github.com/sprioc/composer/pkg/rsp"
)

func ModifyImagesInCollection(requestFrom model.Ref, col model.Ref, additions map[string][]string) rsp.Response {
	if col.Valid(model.Collections) {
		return rsp.Response{Message: "Invalid reference", Code: http.StatusBadRequest}
	}

	var addLinks []string
	var remLinks []string
	var ok bool

	if addLinks, ok = additions["add"]; !ok {
		return rsp.Response{Message: "Invalid body", Code: http.StatusBadRequest}
	}

	if remLinks, ok = additions["remove"]; !ok {
		return rsp.Response{Message: "Invalid body", Code: http.StatusBadRequest}
	}

	valid, err := redis.Permissions(requestFrom, model.CanEdit, col)
	if err != nil {
		return rsp.Response{Code: http.StatusInternalServerError}
	}
	if !valid {
		return rsp.Response{Message: "User cannot modify collection.", Code: http.StatusUnauthorized}
	}

	exists, err := redis.Exists(col)
	if err != nil {
		return rsp.Response{Code: http.StatusInternalServerError}
	}
	if !exists {
		return rsp.Response{Message: "Collection does not exist.", Code: http.StatusNotFound}
	}

	refs := refs.GetRefs(addLinks)
	for _, ref := range refs {
		err := redis.LinkItems(col, redis.Collection, ref, false)
		if err != nil {
			return rsp.Response{Code: http.StatusInternalServerError}
		}
	}

	refs = refs.GetRefs(remLinks)
	for _, ref := range refs {
		err := redis.LinkItems(col, redis.Collection, ref, true)
		if err != nil {
			return rsp.Response{Code: http.StatusInternalServerError}
		}
	}

	return rsp.Response{Code: http.StatusAccepted}
}

func AddTagsToImage(requestFrom model.Ref, imageRef model.Ref, additions map[string][]string) rsp.Response {
	if imageRef.Collection != model.Images {
		return rsp.Response{Message: "Invalid reference", Code: http.StatusBadRequest}
	}

	var tags []string
	var ok bool

	if tags, ok = additions["tags"]; !ok {
		return rsp.Response{Message: "Invalid body", Code: http.StatusBadRequest}
	}

	valid, err := redis.Permissions(requestFrom, model.CanEdit, imageRef)
	if err != nil {
		return rsp.Response{Code: http.StatusInternalServerError}
	}
	if !valid {
		return rsp.Response{Message: "User cannot modify collection.", Code: http.StatusUnauthorized}
	}

	exists, err := redis.Exists(imageRef)
	if err != nil {
		return rsp.Response{Code: http.StatusInternalServerError}
	}
	if !exists {
		return rsp.Response{Message: "Image does not exist.", Code: http.StatusNotFound}
	}

	err = mongo.Modify(imageRef, bson.M{"$addToSet": bson.M{"tags": bson.M{"$each": tags}}})
	if err != nil {
		return rsp.Response{Code: http.StatusInternalServerError}
	}

	return rsp.Response{Code: http.StatusAccepted}
}

func RemoveTagsFromImage(requestFrom model.Ref, imageRef model.Ref, deletions map[string][]string) rsp.Response {
	if imageRef.Collection != model.Images {
		return rsp.Response{Message: "Invalid reference", Code: http.StatusBadRequest}
	}

	var tags []string
	var ok bool

	if tags, ok = deletions["tags"]; !ok {
		return rsp.Response{Message: "Invalid body", Code: http.StatusBadRequest}
	}

	valid, err := redis.Permissions(requestFrom, model.CanEdit, imageRef)
	if err != nil {
		return rsp.Response{Code: http.StatusInternalServerError}
	}
	if !valid {
		return rsp.Response{Message: "User cannot modify collection.", Code: http.StatusUnauthorized}
	}

	exists, err := redis.Exists(imageRef)
	if err != nil {
		return rsp.Response{Code: http.StatusInternalServerError}
	}
	if !exists {
		return rsp.Response{Message: "Image does not exist.", Code: http.StatusNotFound}
	}

	err = mongo.Modify(imageRef, bson.M{"$pull": bson.M{"tags": bson.M{"$each": tags}}})
	if err != nil {
		return rsp.Response{Code: http.StatusInternalServerError}
	}

	return rsp.Response{Code: http.StatusAccepted}
}
