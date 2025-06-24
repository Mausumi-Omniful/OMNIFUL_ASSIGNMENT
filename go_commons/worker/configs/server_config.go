package configs

import (
	"strings"

	"github.com/omniful/go_commons/util"
)

type ServerConfig struct {
	IncludeGroupsArg string
	ExcludeGroupsArg string
	ListenerNamesArg string
}

func (sc ServerConfig) GetIncludeGroups() []string {
	return util.RemoveEmptyStrings(strings.Split(sc.IncludeGroupsArg, ","))
}

func (sc ServerConfig) GetExcludeGroups() []string {
	return util.RemoveEmptyStrings(strings.Split(sc.ExcludeGroupsArg, ","))
}

func (sc ServerConfig) GetListenerNames() []string {
	return util.RemoveEmptyStrings(strings.Split(sc.ListenerNamesArg, ","))
}
