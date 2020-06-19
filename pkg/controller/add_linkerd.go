package controller

import (
	"github.com/spaghettifunk/linkerd2-operator/pkg/controller/linkerd"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, linkerd.Add)
}
