package models

type Workspace struct {
	Workspace       string `json:"Workspace"`
	BackgroundImage string `json:"BackgroundImage"`
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
