package models

import (
	"log"

	"github.com/antoniomralmeida/k2/initializers"
	"github.com/antoniomralmeida/k2/lib"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func (user *KBUser) Persist() error {
	ctx, collection := initializers.GetCollection("KBUser")
	if user.Id.IsZero() {
		user.Id = primitive.NewObjectID()
		_, err := collection.InsertOne(ctx, user)
		return err
	} else {
		_, err := collection.UpdateOne(ctx, bson.D{{Key: "_id", Value: user.Id}}, user)
		return err
	}
}

func (user *KBUser) FindOne(p bson.D) error {
	ctx, collection := initializers.GetCollection("KBUser")
	err := collection.FindOne(ctx, p).Decode(&user)
	return err
}

func NewUser(name, email, pwd, image string) error {
	copy, err := lib.LoadImage(image)
	if err != nil {
		return err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u := KBUser{Email: email, Name: name, Hash: hash, FaceImage: copy, Profile: Empty}
	err = u.Persist()
	return err
}

func InitSecurity() {
	user := KBUser{}
	CheckIndexs()
	if user.FindOne(bson.D{{Key: "profile", Value: Admin}}) != nil {
		pwd := lib.GeneratePassword(12, 1, 3, 2)
		hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
		initializers.Log(err, initializers.Fatal)
		user = KBUser{Name: "Default Admin", Email: "admin@k2.com", Hash: hash, Profile: Admin}
		initializers.Log(user.Persist(), initializers.Fatal)
		log.Println("Default Hash " + pwd)
	}
}

func CheckIndexs() {
	ctx, collection := initializers.GetCollection("KBUser")
	idx := collection.Indexes()
	ret, err := idx.List(ctx)
	initializers.Log(err, initializers.Fatal)
	var results []interface{}
	err = ret.All(ctx, &results)
	initializers.Log(err, initializers.Fatal)
	if len(results) == 1 {
		initializers.CreateUniqueIndex("KBUser", "email")
	}
}
