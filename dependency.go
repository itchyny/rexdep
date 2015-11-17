package main

// Dependency relation between modules
type Dependency struct {
	relation map[string]map[string]bool
	modules  []string
}

func newDependency() *Dependency {
	dependency := new(Dependency)
	dependency.relation = make(map[string]map[string]bool)
	return dependency
}

func (d *Dependency) add(module string, to string) {
	if _, ok := d.relation[module]; !ok {
		d.relation[module] = make(map[string]bool)
		d.modules = append(d.modules, module)
	}
	d.relation[module][to] = true
}

func (d *Dependency) concat(e *Dependency) {
	for _, module := range e.modules {
		for to := range e.relation[module] {
			d.add(module, to)
		}
	}
}
