package object

type Context struct {
	store map[string]Object
	outer *Context
}

func NewContext() *Context {
	s := make(map[string]Object)
	return &Context{store: s, outer: nil}
}

func (c *Context) Get(name string) (Object, bool) {
	obj, ok := c.store[name]
	if !ok && c.outer != nil {
		obj, ok = c.outer.Get(name)
	}
	return obj, ok
}

func (c *Context) Set(name string, value Object) Object {
	c.store[name] = value
	return value
}

func NewEnclosedContext(outer *Context) *Context {
	ctx := NewContext()
	ctx.outer = outer
	return ctx
}
