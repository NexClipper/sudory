package v1

import "testing"

func TestAllSync(t *testing.T) {
	TestClientSync(t)          //1
	TestClusterSync(t)         //2
	TestServiceStepSync(t)     //3
	TestServiceSync(t)         //4
	TestTemplateCommandSync(t) //5
	TestTemplateSync(t)        //6
}
