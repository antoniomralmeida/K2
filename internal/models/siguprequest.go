package models

type SigupRequest struct {
	Name      string `form:"name"`
	Email     string `form:"email"`
	Password  string `form:"password"`
	Password2 string `form:"password2"`
	//	FaceImage string `form:"faceimage"`
}
