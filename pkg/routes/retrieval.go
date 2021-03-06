package routes

import (
	"github.com/fokal/fokal-core/pkg/cache"
	"github.com/fokal/fokal-core/pkg/handler"
	"github.com/fokal/fokal-core/pkg/model"
	"github.com/fokal/fokal-core/pkg/retrieval"
	"github.com/fokal/fokal-core/pkg/security"
	"github.com/fokal/fokal-core/pkg/security/permissions"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

func RegisterRetrievalRoutes(state *handler.State, api *mux.Router, chain alice.Chain) {
	get := api.Methods("GET").Subrouter()
	opts := api.Methods("OPTIONS").Subrouter()

	c := chain.Append(alice.Constructor(handler.Middleware{State: state, M: cache.Handler}.Handler))
	get.Handle("/images/{ID:[a-zA-Z]{12}}",
		chain.Append(
			handler.Middleware{
				State: state,
				M:     security.SetAuthenticatedUser,
			}.Handler,
			permissions.Middleware{State: state,
				T:          permissions.CanView,
				TargetType: model.Images,
				M:          permissions.PermissionMiddle,
			}.Handler).Then(handler.Handler{State: state, H: retrieval.ImageHandler}))
	opts.Handle("/images/{ID:[a-zA-Z]{12}}", chain.Then(handler.Options("GET")))

	get.Handle("/images/featured",
		c.Append(
			handler.Middleware{
				State: state,
				M:     security.SetAuthenticatedUser,
			}.Handler).Then(handler.Handler{State: state, H: retrieval.FeaturedImageHandler}))
	opts.Handle("/images/featured", chain.Then(handler.Options("GET")))

	get.Handle("/images/recent",
		c.Append(
			handler.Middleware{
				State: state,
				M:     security.SetAuthenticatedUser,
			}.Handler).Then(handler.Handler{State: state, H: retrieval.RecentImageHandler}))
	opts.Handle("/images/recent", chain.Then(handler.Options("GET")))

	get.Handle("/images/trending",
		c.Append(
			handler.Middleware{
				State: state,
				M:     security.SetAuthenticatedUser,
			}.Handler).Then(handler.Handler{State: state, H: retrieval.TrendingImagesHander}))
	opts.Handle("/images/trending", chain.Then(handler.Options("GET")))

	get.Handle("/users/me", chain.Append(
		handler.Middleware{
			State: state,
			M:     security.Authenticate,
		}.Handler).Then(handler.Handler{State: state, H: retrieval.LoggedInUserHandler}))
	opts.Handle("/users/me", chain.Then(handler.Options("GET")))

	get.Handle("/users/me/images", chain.Append(
		handler.Middleware{
			State: state,
			M:     security.Authenticate,
		}.Handler).Then(handler.Handler{State: state, H: retrieval.LoggedInUserImagesHandler}))
	opts.Handle("/users/me/images", chain.Then(handler.Options("GET")))

	get.Handle("/users/{ID}", chain.Then(handler.Handler{State: state, H: retrieval.UserHandler}))
	opts.Handle("/users/{ID}", chain.Then(handler.Options("GET")))

	get.Handle("/users/{ID}/images", c.Then(handler.Handler{State: state, H: retrieval.UserImagesHandler}))
	opts.Handle("/users/{ID}/images", chain.Then(handler.Options("GET")))

	get.Handle("/users/{ID}/favorites", c.Then(handler.Handler{State: state, H: retrieval.UserFavoritesHandler}))
	opts.Handle("/users/{ID}/favorites", chain.Then(handler.Options("GET")))

	get.Handle("/tags/{ID}", c.Then(handler.Handler{State: state, H: retrieval.TagHandler}))
	opts.Handle("/tags/{ID}", chain.Then(handler.Options("GET")))
}
