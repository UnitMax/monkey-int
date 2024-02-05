package object

type Context struct {
	store map[string]Object
}

func NewContext() *Context {
	s := make(map[string]Object)
	return &Context{store: s}
}

func (c *Context) Get(name string) (Object, bool) {
	obj, ok := c.store[name]
	return obj, ok
}

func (c *Context) Set(name string, value Object) Object {
	c.store[name] = value
	return value
}
