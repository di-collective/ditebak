package global

// Globalvar
var (
	Context = &context{
		email: key("email"),
	}
)

type key string

func (k key) String() string {
	return "ditebak." + string(k)
}

type context struct {
	email key
}

func (ctx *context) Email() key {
	return ctx.email
}
