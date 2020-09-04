package admin_plugin

import (
	"github.com/ecletus/plug"
	"github.com/moisespsena-go/edis"

	"github.com/ecletus/admin"
)

type events struct {
	Dis       plug.EventDispatcherInterface
	Names     *AdminNames
	AdminsKey string
}

func Events(args ...interface{}) (d *events) {
	d = &events{}
	for _, arg := range args {
		switch at := arg.(type) {
		case func() []string:
			panic("deprecated")
		case *AdminNames:
			panic("deprecated")
		case edis.EventDispatcherInterface:
			d.Dis = at
		case string:
			panic("deprecated")
		default:
			panic("deprecated")
		}
	}
	return
}

func (e *events) SetNamesFunc(f func() []string) *events {
	panic("deprecated")
	return e
}

func (dn *events) EachOrAll(e plug.PluginEventInterface, cb func(adminName string, Admin *admin.Admin) (err error)) (err error) {
	Admins := e.Options().GetInterface(dn.AdminsKey).(*Admins)
	return cb(DEFAULT_ADMIN, Admins.GetDefault())
}

func (dn *events) AdminE(cb func(e *AdminEvent) error) *events {
	dn.Dis.On(EAdmin(DEFAULT_ADMIN), func(e plug.PluginEventInterface) error {
		return cb(e.(*AdminEvent))
	})
	return dn
}

func (dn *events) Admin(cb func(e *AdminEvent)) *events {
	dn.Dis.On(EAdmin(DEFAULT_ADMIN), func(e plug.PluginEventInterface) {
		cb(e.(*AdminEvent))
	})
	return dn
}

func (dn *events) DoneE(cb func(e *AdminEvent) error) *events {
	dn.Dis.On(EDone(DEFAULT_ADMIN), func(e plug.PluginEventInterface) error {
		return cb(e.(*AdminEvent))
	})
	return dn
}

func (dn *events) Done(cb func(e *AdminEvent)) *events {
	dn.Dis.On(EDone(DEFAULT_ADMIN), func(e plug.PluginEventInterface) {
		cb(e.(*AdminEvent))
	})
	return dn
}

func (dn *events) InitResourcesE(cb func(e *AdminEvent) error) *events {
	dn.Dis.On(EInitResources(DEFAULT_ADMIN), func(e plug.PluginEventInterface) error {
		return cb(e.(*AdminEvent))
	})
	return dn
}

func (dn *events) InitResources(cb func(e *AdminEvent)) *events {
	dn.Dis.On(EInitResources(DEFAULT_ADMIN), func(e plug.PluginEventInterface) {
		cb(e.(*AdminEvent))
	})
	return dn
}

func (dn *events) FuncMapE(cb func(e *AdminFuncMapEvent) error) *events {
	dn.Dis.On(EFuncMap(DEFAULT_ADMIN), func(e plug.PluginEventInterface) error {
		return cb(e.(*AdminFuncMapEvent))
	})
	return dn
}

func (dn *events) FuncMap(cb func(e *AdminFuncMapEvent)) *events {
	dn.Dis.On(EFuncMap(DEFAULT_ADMIN), func(e plug.PluginEventInterface) {
		cb(e.(*AdminFuncMapEvent))
	})
	return dn
}

func (dn *events) RouterE(cb func(e *AdminRouterEvent) error) *events {
	dn.Dis.On(ERoute(DEFAULT_ADMIN), func(e plug.PluginEventInterface) error {
		return cb(e.(*AdminRouterEvent))
	})
	return dn
}

func (dn *events) Router(cb func(e *AdminRouterEvent)) *events {
	dn.Dis.On(ERoute(DEFAULT_ADMIN), func(e plug.PluginEventInterface) {
		cb(e.(*AdminRouterEvent))
	})
	return dn
}
