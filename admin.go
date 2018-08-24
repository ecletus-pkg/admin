package admin_plugin

import (
	"path/filepath"

	"strings"

	"github.com/aghape/admin"
	"github.com/aghape/plug"
	"github.com/aghape/router"
	"github.com/moisespsena/go-default-logger"
	"github.com/moisespsena/go-error-wrap"
	"github.com/moisespsena/go-path-helpers"
)

const DEFAULT_ADMIN = "default"

var (
	PKG = path_helpers.GetCalledDir()
	log = defaultlogger.NewLogger(PKG)
)

type Plugin struct {
	plug.EventDispatcher
	AdminsKey string
}

func (p *Plugin) RequireOptions() []string {
	return []string{p.AdminsKey}
}

func (p *Plugin) NameSpace() string {
	return filepath.Dir(path_helpers.GetCalledDir(false))
}

func (p *Plugin) AssetsRootPath() string {
	return filepath.Dir(path_helpers.GetCalledDir(true))
}

func (p *Plugin) OnRegister() {
	plug.OnAssetFS(p, func(e *plug.AssetFSEvent) {
		e.RegisterAssets(e.PathOf(&admin.Admin{}))
	})

	adminsCalled := map[string]bool{}

	p.On(E_ADMIN, func(e plug.PluginEventInterface) (err error) {
		adminEvent := e.(*AdminEvent)
		if _, ok := adminsCalled[adminEvent.AdminName]; ok {
			return nil
		}
		adminsCalled[adminEvent.AdminName] = true
		adminName, Admin := adminEvent.AdminName, adminEvent.Admin
		log.Debugf("trigger AdminEvent")
		if err = e.PluginDispatcher().TriggerPlugins(&AdminEvent{plug.NewPluginEvent(EAdmin(adminName)), Admin, adminName, e}); err != nil {
			return errwrap.Wrap(err, "AdminEvent")
		}
		log.Debugf("trigger AdminInitResourcesEvent")
		if err = e.PluginDispatcher().TriggerPlugins(&AdminEvent{plug.NewPluginEvent(EInitResources(adminName)), Admin, adminName, e}); err != nil {
			return errwrap.Wrap(err, "AdminInitResourcesEvent")
		}
		log.Debugf("trigger AdminFuncMapEvent")
		if err = e.PluginDispatcher().TriggerPlugins(&AdminFuncMapEvent{&AdminEvent{plug.NewPluginEvent(EFuncMap(adminName)), Admin, adminName, e}}); err != nil {
			return errwrap.Wrap(err, "AdminFuncMapEvent")
		}
		log.Debugf("trigger AdminDone")
		if err = Admin.TriggerDone(&admin.AdminEvent{plug.NewEvent(admin.E_DONE), Admin}); err != nil {
			return errwrap.Wrap(err, "Admin.Done")
		}
		if err = e.PluginDispatcher().TriggerPlugins(&AdminEvent{plug.NewPluginEvent(EDone(adminName)), Admin, adminName, e}); err != nil {
			return errwrap.Wrap(err, "AdminDone")
		}
		return nil
	})

	router.OnRouteE(p, func(e *router.RouterEvent) (err error) {
		admins := e.Options().GetInterface(p.AdminsKey).(*Admins)
		err = admins.Trigger(e.PluginDispatcher())
		if err != nil {
			return errwrap.Wrap(err, "Trigger Admins [%s]", strings.Join(admins.Names(), ", "))
		}
		return admins.Each(func(adminName string, Admin *admin.Admin) error {
			log.Debugf("[admin=%q] mounted on %v", adminName, Admin.Config.MountPath)
			mux := Admin.NewServeMux()
			e.Router.Mux.Mount(Admin.Config.MountPath, mux)
			return errwrap.Wrap(e.PluginDispatcher().TriggerPlugins(&AdminRouterEvent{&AdminEvent{
				plug.NewPluginEvent(ERoute(adminName)), Admin,
				adminName, e}, mux}),
				"AdminRouterEvent")
		})
	})
}
