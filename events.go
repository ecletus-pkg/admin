package admin_plugin

import (
	"reflect"

	"github.com/ecletus/admin"
	"github.com/ecletus/plug"
	"github.com/moisespsena/go-edis"
)

type events struct {
	Dis       plug.EventDispatcherInterface
	Names     *AdminNames
	AdminsKey string
}

func Events(args ...interface{}) (d *events) {
	d = &events{}
	var doArg = func(arg interface{}) {
		value := reflect.Indirect(reflect.ValueOf(args[0]))
		if value.Kind() != reflect.Struct {
			return
		}
		if f := value.FieldByName("AdminNames"); f.IsValid() {
			if f.Kind() == reflect.Struct {
				f = f.Addr()
			}
			d.Names = f.Interface().(*AdminNames)
		}
		if f := value.FieldByName("EventDispatcher"); f.IsValid() {
			if f.Kind() == reflect.Struct {
				f = f.Addr()
			}
			d.Dis = f.Interface().(edis.EventDispatcherInterface)
		}
		if f := value.FieldByName("EventDispatcherInterface"); f.IsValid() {
			if f.Kind() == reflect.Struct {
				f = f.Addr()
			}
			d.Dis = f.Interface().(edis.EventDispatcherInterface)
		}
		return
	}
	for _, arg := range args {
		switch at := arg.(type) {
		case func() []string:
			d.SetNamesFunc(at)
		case *AdminNames:
			d.Names = at
		case edis.EventDispatcherInterface:
			d.Dis = at
			doArg(arg)
		case string:
			d.AdminsKey = at
		default:
			doArg(arg)
		}
	}
	return
}

func (e *events) SetNamesFunc(f func() []string) *events {
	e.Names = &AdminNames{NamesFunc: f}
	return e
}

func (dn *events) EachOrAll(e plug.PluginEventInterface, cb func(adminName string, Admin *admin.Admin) (err error)) (err error) {
	Admins := e.Options().GetInterface(dn.AdminsKey).(*Admins)
	if len(dn.Names.Names) == 0 {
		for adminName, Admin := range Admins.ByName {
			err = cb(adminName, Admin)
			if err != nil {
				return
			}
		}
	} else {
		for _, adminName := range dn.Names.Names {
			err = cb(adminName, Admins.ByName[adminName])
			if err != nil {
				return
			}
		}
	}
	return
}

func (dn *events) AdminE(cb func(e *AdminEvent) error) *events {
	dn.Names.EachOrDefault(func(adminName string) {
		dn.Dis.On(EAdmin(adminName), func(e plug.PluginEventInterface) error {
			return cb(e.(*AdminEvent))
		})
	})
	return dn
}

func (dn *events) Admin(cb func(e *AdminEvent)) *events {
	dn.Names.EachOrDefault(func(adminName string) {
		dn.Dis.On(EAdmin(adminName), func(e plug.PluginEventInterface) {
			cb(e.(*AdminEvent))
		})
	})
	return dn
}

func (dn *events) DoneE(cb func(e *AdminEvent) error) *events {
	dn.Names.EachOrDefault(func(adminName string) {
		dn.Dis.On(EDone(adminName), func(e plug.PluginEventInterface) error {
			return cb(e.(*AdminEvent))
		})
	})
	return dn
}

func (dn *events) Done(cb func(e *AdminEvent)) *events {
	dn.Names.EachOrDefault(func(adminName string) {
		dn.Dis.On(EDone(adminName), func(e plug.PluginEventInterface) {
			cb(e.(*AdminEvent))
		})
	})
	return dn
}

func (dn *events) InitResourcesE(cb func(e *AdminEvent) error) *events {
	dn.Names.EachOrDefault(func(adminName string) {
		dn.Dis.On(EInitResources(adminName), func(e plug.PluginEventInterface) error {
			return cb(e.(*AdminEvent))
		})
	})
	return dn
}

func (dn *events) InitResources(cb func(e *AdminEvent)) *events {
	dn.Names.EachOrDefault(func(adminName string) {
		dn.Dis.On(EInitResources(adminName), func(e plug.PluginEventInterface) {
			cb(e.(*AdminEvent))
		})
	})
	return dn
}

func (dn *events) FuncMapE(cb func(e *AdminFuncMapEvent) error) *events {
	dn.Names.EachOrDefault(func(adminName string) {
		dn.Dis.On(EFuncMap(adminName), func(e plug.PluginEventInterface) error {
			return cb(e.(*AdminFuncMapEvent))
		})
	})
	return dn
}

func (dn *events) FuncMap(cb func(e *AdminFuncMapEvent)) *events {
	dn.Names.EachOrDefault(func(adminName string) {
		dn.Dis.On(EFuncMap(adminName), func(e plug.PluginEventInterface) {
			cb(e.(*AdminFuncMapEvent))
		})
	})
	return dn
}

func (dn *events) RouterE(cb func(e *AdminRouterEvent) error) *events {
	dn.Names.EachOrDefault(func(adminName string) {
		dn.Dis.On(ERoute(adminName), func(e plug.PluginEventInterface) error {
			return cb(e.(*AdminRouterEvent))
		})
	})
	return dn
}

func (dn *events) Router(cb func(e *AdminRouterEvent)) *events {
	dn.Names.EachOrDefault(func(adminName string) {
		dn.Dis.On(ERoute(adminName), func(e plug.PluginEventInterface) {
			cb(e.(*AdminRouterEvent))
		})
	})
	return dn
}
