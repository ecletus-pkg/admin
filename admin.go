package admin_plugin

import (
	"path/filepath"

	"github.com/ecletus/assets"
	"github.com/gobwas/glob"

	"github.com/moisespsena-go/aorm"

	"github.com/ecletus/db"

	"strings"

	"github.com/ecletus/admin"
	"github.com/ecletus/plug"
	"github.com/ecletus/router"
	"github.com/moisespsena-go/xroute"
	"github.com/moisespsena-go/error-wrap"
	"github.com/moisespsena-go/path-helpers"
)

const DEFAULT_ADMIN = "default"

var (
	PKG = path_helpers.GetCalledDir()
)

type Plugin struct {
	plug.EventDispatcher
	db.DBNames

	AdminsKey          string
	SystemDBDialectKey string
}

func (p *Plugin) RequireOptions() []string {
	return []string{p.AdminsKey, p.SystemDBDialectKey}
}

func (p *Plugin) NameSpace() string {
	return filepath.Join("github.com", "ecletus", "admin")
}

func (p *Plugin) AssetsRootPath() (pth string) {
	pth = filepath.Join(path_helpers.GetCalledDir(true), "..", "..", "ecletus", "admin")
	return
}

func (p *Plugin) OnRegister(options *plug.Options) {
	plug.OnFS(p, func(e *plug.FSEvent) {
		e.RegisterAssetPath(e.PathOf(&admin.Admin{}))
	})

	_ = plug.OnPostInit(p, func(e plug.PluginEventInterface) {
		admins := options.GetInterface(p.AdminsKey).(*Admins)
		Events(p, func() []string {
			return admins.Names()
		}).InitResources(initResources)

		db.Events(p).DBOnMigrate(migrate)
	})

	adminsCalled := map[string]bool{}

	log := p.Logger()

	p.On(E_ADMIN, func(e plug.PluginEventInterface) (err error) {
		adminEvent := e.(*AdminEvent)
		if _, ok := adminsCalled[adminEvent.AdminName]; ok {
			return nil
		}
		adminsCalled[adminEvent.AdminName] = true
		adminName, Admin := adminEvent.AdminName, adminEvent.Admin

		if systemDialect := options.GetString(p.SystemDBDialectKey, Admin.Config.FakeDBDialect); systemDialect != Admin.Config.FakeDBDialect {
			Admin.FakeDB = aorm.FakeDB(systemDialect)
		}
		// github.com/ecletus-pkg/admin
		//
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

	err := router.OnRouteE(p, func(e *router.RouterEvent) (err error) {
		admins := e.Options().GetInterface(p.AdminsKey).(*Admins)
		err = admins.Trigger(e.PluginDispatcher())
		if err != nil {
			return errwrap.Wrap(err, "Trigger Admins [%s]", strings.Join(admins.Names(), ", "))
		}
		return admins.Each(func(adminName string, Admin *admin.Admin) error {
			log.Debugf("[admin=%q] mounted on %v", adminName, Admin.Config.MountPath)
			var router xroute.Router = e.Router.Mux
			if Admin.Config.MountPath != "/" {
				router = xroute.NewMux(Admin.Name)
				Admin.InitRoutes(router)
				e.Router.Mux.Mount(Admin.Config.MountPath, router)
			} else {
				Admin.InitRoutes(router)
			}
			return errwrap.Wrap(e.PluginDispatcher().TriggerPlugins(&AdminRouterEvent{&AdminEvent{
				plug.NewPluginEvent(ERoute(adminName)), Admin,
				adminName, e}, router}),
				"AdminRouterEvent")
		})
	})
	if err != nil {
		panic(err)
	}

	assets.Dis(p).OnSyncConfig(func(e *assets.PreRepositorySyncEvent) {
		e.Repo.IgnorePath(
			glob.MustCompile("static/admin/javascripts/{app,qor}").Match,
			glob.MustCompile("*.map").Match,
		)
	})
}
