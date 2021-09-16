package UserDS

import (
	"encoding/json"
	"io"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	FirstName string             `json:"firstname,omitempty" bson:"firstname,omitempty"`
	LastName  string             `json:"lastname,omitempty" bson:"lastname,omitempty"`
	CreatedOn time.Time          `json:"createdon,omitempty" bson:"createdon,omitempty"`
	UpdatedOn time.Time          `json:"updatedon" bson:"updatedon"`
	Mobile    string             `json:"mobile,omitempty" bson:"mobile,omitempty"`
	Active    bool               `json:"active,omitempty" bson:"active,omitempty"`
	Age       ageDS              `json:"age,omitempty" bson:"age,omitempty"`
}

type ageDS struct {
	Value    int    `json:"age,omitempty" bson:"age,omitempty"`
	Interval string `json:"interval,omitempty" bson:"interval,omitempty"`
}

func (u *User) FromJson(r io.Reader) error {
	err := json.NewDecoder(r)
	return err.Decode(u)
}

func (u *User) ToJson(rw io.Writer) error {
	err := json.NewEncoder(rw)
	return err.Encode(u)
}
