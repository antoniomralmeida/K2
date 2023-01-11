package kb

import (
	"encoding/json"
	"errors"

	"github.com/antoniomralmeida/k2/initializers"
	"github.com/kamva/mgm/v3"
)

func (obj *KBObject) Persist() error {
	if obj.ID.IsZero() {
		err := mgm.Coll(obj).Create(obj)
		return err
	} else {

		db_doc := new(KBObject)
		err := mgm.Coll(obj).FindByID(obj.ID, db_doc)
		if err != nil {
			return err
		}
		if obj.UpdatedAt != db_doc.UpdatedAt {
			return errors.New("Old document!")
		}
		err = mgm.Coll(obj).Update(obj)
		return err
	}
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
