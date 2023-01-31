package models

type Context struct {
	I18n              map[string]string
	Locales           map[string]string
	UserId            string
	User              string
	ApiKernel         string
	Avatar            string
	SpeechSynthesisId int
	Workspaces        []WorkspaceInfo
}
