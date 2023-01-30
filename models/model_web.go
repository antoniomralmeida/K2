package models

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Workspace struct {
	Workspace       string `json:"Workspace"`
	BackgroundImage string `json:"BackgroundImage"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SigupRequest struct {
	Name      string `form:"name"`
	Email     string `form:"email"`
	Password  string `form:"password"`
	Password2 string `form:"password2"`
	//	FaceImage string `form:"faceimage"`
}

type KBProfile byte

const (
	Empty KBProfile = iota
	User
	Manager
	Admin
)

type KBUser struct {
	mgm.DefaultModel `json:",inline" bson:",inline"`
	Name             string               `bson:"name"`
	Email            string               `bson:"email"`
	Hash             []byte               `bson:"hash" json:"-"`
	Profile          KBProfile            `bson:"profile"`
	Workspaces       []primitive.ObjectID `bson:"workspaces"`
	FaceImage        string               `bson:"faceimage,omitempty"`
}

type KBAlert struct {
	mgm.DefaultModel `json:",inline" bson:",inline"`
	Message          string               `bson:"message"`
	User             primitive.ObjectID   `bson:"user"`
	Views            []primitive.ObjectID `bson:"views"`
}

type Context struct {
	I18n              map[string]string
	Locales           map[string]string
	UserId            string
	User              string
	ApiKernel         string
	Avatar            string
	SpeechSynthesisId int
	Workspaces        []Workspace
}
