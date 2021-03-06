package request

import (
	"net/http"

	"github.com/mholt/binding"
)

type CreateUserRequest struct {
	Username string
	Email    string
	Password string
}

func (cf *CreateUserRequest) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&cf.Username: binding.Field{
			Form:     "username",
			Required: true,
		},
		&cf.Email: binding.Field{
			Form:     "email",
			Required: true,
		},
		&cf.Password: binding.Field{
			Form:     "password",
			Required: true,
		},
	}
}

type PatchUserRequest struct {
	Username  string `structs:"username,omitempty"`
	Bio       string `structs:"bio,omitempty"`
	URL       string `structs:"url,omitempty"`
	Name      string `structs:"name,omitempty"`
	Location  string `structs:"location,omitempty"`
	Instagram string `structs:"instagram,omitempty"`
	Twitter   string `structs:"twitter,omitempty"`
}

func (cf *PatchUserRequest) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&cf.Bio:       "bio",
		&cf.URL:       "email",
		&cf.Name:      "password",
		&cf.Location:  "location",
		&cf.Twitter:   "twitter",
		&cf.Instagram: "instagram",
		&cf.Username:  "username",
	}
}
