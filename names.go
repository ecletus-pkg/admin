package admin_plugin

type AdminNames struct {
	NamesFunc func() []string
	Names     []string
}

func (dn *AdminNames) EachOrDefaultE(cb func(adminName string) error) (err error) {
	for _, name := range dn.GetNames() {
		err = cb(name)
		if err != nil {
			return
		}
	}
	return
}

func (dn *AdminNames) EachOrDefault(cb func(adminName string)) {
	dn.EachOrDefaultE(func(adminName string) error {
		cb(adminName)
		return nil
	})
}

func (a *AdminNames) GetNames() []string {
	if a.NamesFunc != nil {
		names := a.NamesFunc()
		if len(names) == 0 {
			names = []string{DEFAULT_ADMIN}
		}
		return names
	}
	if len(a.Names) == 0 {
		a.Names = []string{DEFAULT_ADMIN}
	}
	return a.Names
}
