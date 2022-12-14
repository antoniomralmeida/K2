package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Workspace struct {
	Workspace       string `json:"Workspace"`
	BackgroundImage string `json:"BackgroundImage"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type KBProfile byte

const (
	Empty KBProfile = iota
	User
	Manager
	Admin
)

type KBUser struct {
	Id         primitive.ObjectID   `bson:"_id" json:"id"`
	Name       string               `bson:"name"`
	Email      string               `bson:"email"`
	Hash       []byte               `bson:"hash" json:"-"`
	Profile    KBProfile            `bson:"profile"`
	Workspaces []primitive.ObjectID `bson:"workspaces"`
}

type Context struct {
	User       string
	Title      string
	DataInput  string
	Workspace  string
	Alerts     string
	ApiKernel  string
	Avatar     string
	Workspaces []Workspace
}
