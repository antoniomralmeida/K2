package kb

import (
	"encoding/json"

	"github.com/antoniomralmeida/k2/initializers"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (obj *KBWorkspace) Persist() error {
	return initializers.Persist(obj)

}

func (obj *KBWorkspace) GetPrimitiveUpdateAt() primitive.DateTime {
	return primitive.NewDateTimeFromTime(obj.UpdatedAt)
}

func (w *KBWorkspace) String() string {
	j, err := json.MarshalIndent(*w, "", "\t")
	initializers.Log(err, initializers.Error)
	return string(j)
}
