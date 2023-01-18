package models

import (
	"github.com/antoniomralmeida/k2/initializers"
	"github.com/antoniomralmeida/k2/lib"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func (obj *KBUser) Persist() error {
	return initializers.Persist(obj)

}

func (obj *KBUser) GetPrimitiveUpdateAt() primitive.DateTime {
	return primitive.NewDateTimeFromTime(obj.UpdatedAt)
}

func (user *KBUser) FindOne(p bson.D) error {
	err := mgm.Coll(user).First(p, user)
	return err
}

func NewUser(name, email, pwd, image string) (err error) {
	var copy string
	if image != "" {
		copy, err = lib.LoadImage(image)
		if err != nil {
			return
		}
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return
	}
	u := KBUser{Email: email, Name: name, Hash: hash, FaceImage: copy, Profile: Empty}
	return u.Persist()
}

func InitSecurity() {
	user := KBUser{}
	CheckIndexs()
	err := user.FindOne(bson.D{{Key: "profile", Value: Admin}})
	if err == mongo.ErrNoDocuments {
		pwd := lib.GeneratePassword(12, 1, 3, 2)
		hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
		initializers.Log(err, initializers.Fatal)
		user = KBUser{Name: "Default Admin", Email: "admin@k2.com", Hash: hash, Profile: Admin}
		initializers.Log(user.Persist(), initializers.Fatal)
		initializers.Log("Default Hash "+pwd, initializers.Info)
	} else {
		initializers.Log(err, initializers.Fatal)
	}
}

func CheckIndexs() {
	coll := mgm.Coll(&KBUser{})
	idx := coll.Indexes()
	ret, err := idx.List(mgm.Ctx())
	initializers.Log(err, initializers.Fatal)
	var results []interface{}
	err = ret.All(mgm.Ctx(), &results)
	initializers.Log(err, initializers.Fatal)
	if len(results) == 1 {
		initializers.CreateUniqueIndex(coll, "email")
	}
}
