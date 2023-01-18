package kb

import (
	"encoding/json"

	"github.com/antoniomralmeida/k2/initializers"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (obj *KBObject) Persist() error {
	return initializers.Persist(obj)

}

func (obj *KBObject) GetPrimitiveUpdateAt() primitive.DateTime {
	return primitive.NewDateTimeFromTime(obj.UpdatedAt)
}

func (o *KBObject) String() string {
	j, err := json.MarshalIndent(*o, "", "\t")
	initializers.Log(err, initializers.Error)
	return string(j)
}

func (o *KBObject) Delete() error {

	mgm.Coll(o).Delete(o)

	// Restart KB
	Stop()
	Init()
	return nil
}

func (o *KBObject) GetWorkspaces() (ret []*KBWorkspace) {
	for i := range GKB.Workspaces {
		for j := range GKB.Workspaces[i].Objects {
			if GKB.Workspaces[i].Objects[j].KBObject == o {
				ret = append(ret, &GKB.Workspaces[i])
			}
		}
	}
	return
}
