package kb

import (
	"encoding/json"
	"errors"

	"github.com/antoniomralmeida/k2/initializers"
	"github.com/kamva/mgm/v3"
)

func (obj *KBWorkspace) Persist() error {
	if obj.ID.IsZero() {
		err := mgm.Coll(obj).Create(obj)
		return err
	} else {

		db_doc := new(KBWorkspace)
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

func (w *KBWorkspace) String() string {
	j, err := json.MarshalIndent(*w, "", "\t")
	initializers.Log(err, initializers.Error)
	return string(j)
}
