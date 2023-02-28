package models

type Context struct {
	Locale     string
	I18n       map[string]string
	Locales    map[string]string
	UserId     string
	User       string
	Avatar     string
	Voice      string
	Workspaces []WorkspaceInfo
}
