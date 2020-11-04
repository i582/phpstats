package symbols

import (
	"strings"
	"sync"
)

type Namespaces struct {
	m sync.Mutex

	Namespaces map[string]*Namespace
}

func NewNamespaces() *Namespaces {
	return &Namespaces{
		Namespaces: map[string]*Namespace{},
	}
}

func (n *Namespaces) Len() int {
	return len(n.Namespaces)
}

func (n *Namespaces) CreateNamespace(nsName string) *Namespace {
	splitted := splitNsName(nsName)
	ns, isNew := n.createNamespace(splitted, "")
	if isNew {
		n.Add(ns)
	}
	return ns
}

func (n *Namespaces) AddFileToNamespace(nsName string, file *File) {
	ns, ok := n.GetNamespace(nsName)
	if !ok {
		return
	}
	ns.Files.Add(file)
}

func (n *Namespaces) AddClassToNamespace(nsName string, class *Class) {
	ns, ok := n.GetNamespace(nsName)
	if !ok {
		return
	}
	ns.Classes.Add(class)
}

func (n *Namespaces) createNamespace(nsParts []string, fullName string) (*Namespace, bool) {
	if len(nsParts) == 0 {
		return nil, false
	}

	current := nsParts[0]
	fullName += `\` + current

	curNs, ok := n.Get(current)
	if ok {
		ns, _ := curNs.Childs.createNamespace(nsParts[1:], fullName)
		curNs.Childs.Add(ns)
		return curNs, false
	}

	if len(nsParts) == 1 {
		ns := NewNamespace(current, fullName)
		return ns, true
	}

	nsCur := NewNamespace(current, fullName)

	ns, _ := n.createNamespace(nsParts[1:], fullName)
	nsCur.Childs.Add(ns)

	return nsCur, true
}

func (n *Namespaces) GetNamespace(nsName string) (*Namespace, bool) {
	splitted := splitNsName(nsName)
	ns, found := n.getNamespace(splitted)
	return ns, found
}

func (n *Namespaces) getNamespace(nsParts []string) (*Namespace, bool) {
	current := nsParts[0]

	if len(nsParts) == 1 {
		curNs, ok := n.Get(current)
		if !ok {
			return nil, false
		}
		return curNs, true
	}

	curNs, ok := n.Get(current)
	if !ok {
		return nil, false
	}

	ns, found := curNs.Childs.getNamespace(nsParts[1:])

	return ns, found
}

func (n *Namespaces) Add(ns *Namespace) {
	if ns == nil {
		return
	}

	n.m.Lock()
	n.Namespaces[ns.Name] = ns
	n.m.Unlock()
}

func (n *Namespaces) Get(name string) (*Namespace, bool) {
	n.m.Lock()
	ns, ok := n.Namespaces[name]
	n.m.Unlock()
	return ns, ok
}

func (n *Namespaces) GetNamespacesWithSpecificLevel(level int64) []*Namespace {
	if level < 0 {
		return nil
	}

	return n.getNamespacesWithSpecificLevel(0, level)
}

func (n *Namespaces) getNamespacesWithSpecificLevel(curLevel int64, level int64) []*Namespace {
	if curLevel == level {
		res := make([]*Namespace, 0, len(n.Namespaces))

		for _, ns := range n.Namespaces {
			res = append(res, ns)
		}

		return res
	}

	res := make([]*Namespace, 0, len(n.Namespaces)*10)

	for _, ns := range n.Namespaces {
		nss := ns.Childs.getNamespacesWithSpecificLevel(curLevel+1, level)
		res = append(res, nss...)
	}

	return res
}

type Namespace struct {
	Name     string
	FullName string

	Files   *Files
	Classes *Classes

	Childs *Namespaces

	MetricsResolved bool
	Aff             float64
	Eff             float64
	Instab          float64
}

func NewNamespace(name string, fullName string) *Namespace {
	return &Namespace{
		Name:     name,
		FullName: fullName,
		Files:    NewFiles(),
		Classes:  NewClasses(),
		Childs:   NewNamespaces(),
	}
}

func joinNsName(parts []string) string {
	return `\` + strings.Join(parts, `\`)
}

func splitNsName(name string) []string {
	name = strings.TrimPrefix(name, `\`)
	res := strings.Split(name, `\`)
	if len(res) == 0 {
		return []string{name}
	}

	return res
}
