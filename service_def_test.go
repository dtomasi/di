package di_test

import (
	"github.com/dtomasi/di"
	"testing"
)

func TestNewServiceDef(t *testing.T) {
	sd := di.NewServiceDef(di.StringRef("foo")).
		Opts().
		Provider(func() {}).
		Args(
			di.ContextArg(),
			di.ContainerArg(),
			di.ParamProviderArg(),
			di.InterfaceArg(""),
			di.ServiceArg(di.StringRef("bar")),
			di.ParamArg(""),
		)

	if sd == nil {
		t.Error("NewServiceDef returns nil value")
	}
}
