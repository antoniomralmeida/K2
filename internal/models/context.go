package models

type Context struct {
	Locale     string
	I18n       map[string]string
	Locales    map[string]string
	UserId     string
	User       string
	ApiKernel  string
	Avatar     string
	Voice      string
	Workspaces []WorkspaceInfo
}
