package stats

import (
	"strings"
	"sync"
)

type Namespaces struct {
	sync.Mutex

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
	splitted := splitNsName(nsName)
	ns, _ := n.createNamespace(splitted, "")
	ns.Files.Add(file)
}

func (n *Namespaces) createNamespace(nsParts []string, fullName string) (*Namespace, bool) {
	if len(nsParts) == 1 {
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

	n.Lock()
	n.Namespaces[ns.Name] = ns
	n.Unlock()
}

func (n *Namespaces) Get(name string) (*Namespace, bool) {
	n.Lock()
	ns, ok := n.Namespaces[name]
	n.Unlock()
	return ns, ok
}

type Namespace struct {
	Name     string
	FullName string

	Files *Files

	Childs *Namespaces
}

func NewNamespace(name string, fullName string) *Namespace {
	return &Namespace{
		Name:     name,
		FullName: fullName,
		Files:    NewFiles(),
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
