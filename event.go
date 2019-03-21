package admin_plugin

import (
	"github.com/ecletus/admin"
	"github.com/ecletus/plug"
	"github.com/moisespsena-go/xroute"
	"github.com/moisespsena/go-error-wrap"
)

var (
	E_ADMIN                = PKG + ".admin"
	E_ADMIN_DONE           = E_ADMIN + ".done"
	E_ADMIN_FUNC_MAP       = E_ADMIN + ".funcMap"
	E_ADMIN_ROUTE          = E_ADMIN + ".route"
	E_ADMIN_INIT_RESOURCES = E_ADMIN + ".initResources"
)

type AdminEvent struct {
	plug.PluginEventInterface
	Admin       *admin.Admin
	AdminName   string
	PluginEvent plug.PluginEventInterface
}

type AdminFuncMapEvent struct {
	*AdminEvent
}

func (afm *AdminFuncMapEvent) Register(name string, value interface{}) {
	afm.Admin.RegisterFuncMap(name, value)
}

type AdminRouterEvent struct {
	*AdminEvent
	router xroute.Router
}

func (are *AdminRouterEvent) Router() xroute.Router {
	return are.router
}

func EAdmin(adminKey string) string {
	if adminKey == "" {
		panic("adminKey is blank")
	}
	return E_ADMIN + ":" + adminKey
}

func EDone(adminKey string) string {
	if adminKey == "" {
		panic("adminKey is blank")
	}
	return E_ADMIN_DONE + ":" + adminKey
}

func EFuncMap(adminKey string) string {
	if adminKey == "" {
		panic("adminKey is blank")
	}
	return E_ADMIN_FUNC_MAP + ":" + adminKey
}

func ERoute(adminName string) string {
	if adminName == "" {
		panic("AdminName is blank")
	}
	return E_ADMIN_ROUTE + ":" + adminName
}

func EInitResources(adminName string) string {
	if adminName == "" {
		panic("AdminName is blank")
	}
	return E_ADMIN_INIT_RESOURCES + ":" + adminName
}

func (admins *Admins) Trigger(d plug.PluginEventDispatcherInterface) error {
	return admins.Each(func(adminName string, Admin *admin.Admin) (err error) {
		e := &AdminEvent{plug.NewPluginEvent(E_ADMIN), Admin, adminName, nil}
		if err = d.TriggerPlugins(e); err != nil {
			return errwrap.Wrap(err, "Admin %q: event %q", adminName, e.Name())
		}
		return nil
	})
}
