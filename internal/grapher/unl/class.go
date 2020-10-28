package uml

import (
	"fmt"
	"strings"

	"github.com/i582/phpstats/internal/stats"
	"github.com/i582/phpstats/internal/utils"
)

func GetUmlForClass(c *stats.Class) string {
	return GetUmlForClassWithFilter(c, func(*stats.Function) bool {
		return true
	}, func(*stats.Field) bool {
		return true
	})
}

func GetShortUmlForClass(c *stats.Class) string {
	return GetUmlForClassWithFilter(c, func(*stats.Function) bool {
		return false
	}, func(*stats.Field) bool {
		return false
	})
}

func GetUmlForClassWithFilter(c *stats.Class, predMethods func(m *stats.Function) bool, predFields func(f *stats.Field) bool) string {
	id := utils.ClassNameNormalize(c.Name)
	idParts := strings.Split(c.Name, `\`)
	shortName := idParts[len(idParts)-1]

	var fields string
	for _, f := range c.Fields.Fields {
		if !predFields(f) {
			continue
		}

		fields += fmt.Sprintf("+ %s\\n", f.Name)
	}

	if fields == "" && c.Fields.Len() == 0 {
		fields = "no fields"
	}

	if fields == "" && c.Fields.Len() != 0 {
		fields = "..."
	}

	var methods string
	for _, m := range c.Methods.Funcs {
		if !predMethods(m) {
			continue
		}

		methods += fmt.Sprintf("+ %s\\n", m.Name.Name)
	}

	if methods == "" && c.Methods.Len() == 0 {
		methods = "no methods"
	}

	if methods == "" && c.Methods.Len() != 0 {
		methods = "..."
	}

	var typ string

	if c.IsInterface {
		typ = "Interface\\n"
	} else if c.IsAbstract {
		typ = "Abstract\\n"
	} else {
		typ = ""
	}

	label := fmt.Sprintf("{%s%s|<fields>%s|<methods>%s|}", typ, shortName, fields, methods)

	return fmt.Sprintf("%s[label = \"%s\"]\n", id, label)
}
