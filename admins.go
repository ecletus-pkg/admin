package admin_plugin

import (
	"github.com/ecletus/admin"
	"github.com/ecletus/plug"
	"github.com/moisespsena-go/error-wrap"
)

type Admins struct {
	ByName map[string]*admin.Admin
}

func (a *Admins) Set(name string, Admin *admin.Admin) {
	if a.ByName == nil {
		a.ByName = make(map[string]*admin.Admin)
	}
	a.ByName[name] = Admin

	if Admin.SiteName == "" {
		Admin.SiteName = name
	}

	Admin.Name = name
	Admin.Init()
}

func (a *Admins) Names() (n []string) {
	for name := range a.ByName {
		n = append(n, name)
	}
	return
}

func (adms *Admins) Each(cb func(adminName string, Admin *admin.Admin) (err error)) (err error) {
	for adminName, Admin := range adms.ByName {
		err = cb(adminName, Admin)
		if err != nil {
			return errwrap.Wrap(err, "Admin %q", adminName)
		}
	}
	return nil
}

func (a *Admins) SetDefault(Admin *admin.Admin) {
	a.Set(DEFAULT_ADMIN, Admin)
}

func (a *Admins) GetDefault() *admin.Admin {
	return a.ByName[DEFAULT_ADMIN]
}

func (a *Admins) On(ef func(name string) string, dis plug.EventDispatcherInterface, cb interface{}) (err error) {
	for name := range a.ByName {
		err = dis.OnE(ef(name), cb)
		if err != nil {
			return
		}
	}
	return
}

func (a *Admins) OnAdmin(dis plug.EventDispatcherInterface, cb func(e *AdminEvent)) (err error) {
	return a.On(EAdmin, dis, func(e plug.PluginEventInterface) {
		cb(e.(*AdminEvent))
	})
}

func (a *Admins) OnAdminE(dis plug.EventDispatcherInterface, cb func(e *AdminEvent) error) (err error) {
	return a.On(EAdmin, dis, func(e plug.PluginEventInterface) error {
		return cb(e.(*AdminEvent))
	})
}

func (a *Admins) OnInitResources(dis plug.EventDispatcherInterface, cb func(e *AdminEvent)) (err error) {
	return a.On(EInitResources, dis, func(e plug.PluginEventInterface) {
		cb(e.(*AdminEvent))
	})
}

func (a *Admins) OnInitResourcesE(dis plug.EventDispatcherInterface, cb func(e *AdminEvent) error) (err error) {
	return a.On(EInitResources, dis, func(e plug.PluginEventInterface) error {
		return cb(e.(*AdminEvent))
	})
}

func (a *Admins) OnFuncMap(dis plug.EventDispatcherInterface, cb func(e *AdminFuncMapEvent)) (err error) {
	return a.On(EFuncMap, dis, func(e plug.PluginEventInterface) {
		cb(e.(*AdminFuncMapEvent))
	})
}

func (a *Admins) OnFuncMapE(dis plug.EventDispatcherInterface, cb func(e *AdminFuncMapEvent) error) (err error) {
	return a.On(EFuncMap, dis, func(e plug.PluginEventInterface) error {
		return cb(e.(*AdminFuncMapEvent))
	})
}

func (a *Admins) OnRoute(dis plug.EventDispatcherInterface, cb func(e *AdminRouterEvent)) (err error) {
	return a.On(ERoute, dis, func(e plug.PluginEventInterface) {
		cb(e.(*AdminRouterEvent))
	})
}

func (a *Admins) OnRouteE(dis plug.EventDispatcherInterface, cb func(e *AdminRouterEvent) error) (err error) {
	return a.On(ERoute, dis, func(e plug.PluginEventInterface) error {
		return cb(e.(*AdminRouterEvent))
	})
}
