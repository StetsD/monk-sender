package app

type SenderStrategy interface {
	Send(msg *onSendMsg)
	Init(conf map[string]string) error
}

type Context struct {
	channel *SenderStrategy
}

func (c *Context) Init(strategy *SenderStrategy, conf map[string]string) error {
	c.channel = strategy
	err := (*c.channel).Init(conf)
	if err != nil {
		return nil
	}

	return nil
}

func (c *Context) Send(msg *onSendMsg) {
	(*c.channel).Send(msg)
}
