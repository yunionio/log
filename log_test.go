package log

import (
	"testing"
)

func TestDebug(t *testing.T) {
	//Fatalf("debug test")
	V(1).Infof("verbose printf")
	V(0).Debugf("verbose debugf, level %d", 0)
	Printf("printf test")
	Infof("info test %d", 1)
	Warningf("warningf test %f", 2.0)
	Errorf("errorf test %q", "string")
}
