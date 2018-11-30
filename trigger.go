package flogoaztrigger

import (
	
	"context"
	"fmt"
	
	"github.com/project-flogo/core/trigger"
	"github.com/project-flogo/core/support/log"
)

var triggerMd = trigger.NewMetadata(&Output{}, &Reply{})

func init() {
	trigger.Register(&Trigger{}, &Factory{})
}

var singleton *Trigger

type Trigger struct {
	id      string
	config   *trigger.Config
	handlers []trigger.Handler
	logger   log.Logger
}

type Factory struct {
}

func (*Factory) New(config *trigger.Config) (trigger.Trigger, error) {
	singleton = &Trigger{config: config}
	return singleton, nil
}

func (f *Factory) Metadata() *trigger.Metadata {
	return triggerMd
}
func Invoke() (string, error) {
	
	return singleton.Invoke()
}

func (t *Trigger) Invoke() (string, error) {

	handler := t.handlers[0]

	inputData := map[string]interface{}{
		"body": "Abc",
	}

	result, err := handler.Handle(context.Background(), inputData)

	if err != nil {
		t.logger.Debug("Azure Error", err.Error())
		return "", err
	}
	reply := &Reply{}
	reply.FromMap(result)
	
	if reply.Data == nil{
		return "Something didn't work out right", nil
	}
	return reply.Data.(string), nil
}
func (t *Trigger) Initialize(ctx trigger.InitContext) error {
	fmt.Println("Initializing the Azure Trigger by getting handlers", len(ctx.GetHandlers()))
	if len(ctx.GetHandlers()) == 0 {
		return fmt.Errorf("no commands found for cli trigger ")
	}
	t.handlers = ctx.GetHandlers()

	t.logger = ctx.Logger()

	return nil
}

// Start implements util.Managed.Start
func (t *Trigger) Start() error {
	//start servers/services if necessary
	return nil
}

// Stop implements util.Managed.Stop
func (t *Trigger) Stop() error {
	//stop servers/services if necessary
	return nil
}
