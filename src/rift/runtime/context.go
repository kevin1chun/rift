package runtime

import (
	"rift/lang"
	"rift/support/collections"
)

type Context struct{
	environment collections.PersistentMap	
}

func BuildContext(nodes []*lang.Node) *Context {
	ctx := Context{collections.NewPersistentMap()}
	for _, node := range nodes {
		head := node.Values[0].(*lang.Node)
		switch {
		default:
		case head.Type == lang.REF:
			// ctx.Inject(lang.Ref{head}.String(), node)
		}
	}
	return &ctx	
}

func (c *Context) Inject(ref string, node *lang.Node) {
	c.environment.Set(ref, node)
}

func (c *Context) Exists(ref string) bool {
	return c.environment.Contains(ref)
}

func (c *Context) Dereference(ref string) *lang.Node {
	referent := c.environment.GetOrNil(ref)
	if referent != nil {
		return referent.(*lang.Node)
	} else {
		return nil
	}
}