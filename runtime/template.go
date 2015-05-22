package runtime

import (
	"bytes"
	"strings"
	"text/template"
)

type TemplateValues map[string]interface{}

type TemplatePool map[string]*template.Template

// names accessible from text templates: ex. "{{paragraph}}"
var coreFuncs template.FuncMap = template.FuncMap{
	"paragraph": func() string {
		return "\n\n"
	},
	"br": func() string {
		return "\n"
	},
	"action": func() map[string]TemplateValues {
		return templateValueStack.top()
	},
}

type TVS struct {
	stack []map[string]TemplateValues
}

var templateValueStack TVS = TVS{
	[]map[string]TemplateValues{make(map[string]TemplateValues)},
}

func (this *TVS) pushValues(top map[string]TemplateValues) {
	this.stack = append(this.stack, top)
}
func (this *TVS) pop() {
	this.stack = this.stack[0 : len(this.stack)-1]
}
func (this *TVS) top() map[string]TemplateValues {
	i := len(this.stack) - 1
	return this.stack[i]
}

func (this TemplatePool) New(name string, text string) (err error) {
	if !strings.Contains(text, "{{") {
		delete(this, name)
	} else if temp, e := template.New(name).Funcs(coreFuncs).Parse(text); e == nil {
		this[name] = temp
	} else {
		err = e
	}
	return err
}

func runTemplate(temp *template.Template, data interface{}) (ret string, err error) {
	buffer := new(bytes.Buffer)
	if e := temp.Execute(buffer, data); e != nil {
		err = e
	} else {
		ret = buffer.String()
	}
	return ret, err
}

func reallySlow(text string, data interface{}) (ret string, err error) {
	temp, e := template.New(text).Funcs(coreFuncs).Parse(text)
	if e != nil {
		err = e
	} else {
		ret, err = runTemplate(temp, data)
	}
	return ret, err
}
