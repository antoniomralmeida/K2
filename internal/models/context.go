package models

type Context struct {
	I18n        map[string]string
	Locales     map[string]string
	UserId      string
	User        string
	Avatar      string
	WellcomeMsg string
	Workspaces  []WorkspaceInfo
}
