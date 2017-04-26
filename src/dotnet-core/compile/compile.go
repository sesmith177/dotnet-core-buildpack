package compile

import "github.com/cloudfoundry/libbuildpack"

type Compiler struct {
	Stager *libbuildpack.Stager
}

func (dc *Compiler) RestoreCache() error {
	return nil
}
