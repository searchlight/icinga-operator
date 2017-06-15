package descriptor

import (
	"fmt"

	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	gw_descriptor "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway/descriptor"
)

// Registry is a registry of information extracted from plugin.CodeGeneratorRequest.
type Registry struct {
	*gw_descriptor.Registry
}

// NewRegistry returns a new Registry.
func NewRegistry() *Registry {
	return &Registry{
		Registry: &gw_descriptor.Registry{
			Msgs:       make(map[string]*gw_descriptor.Message),
			Enums:      make(map[string]*gw_descriptor.Enum),
			Files:      make(map[string]*gw_descriptor.File),
			PkgMap:     make(map[string]string),
			PkgAliases: make(map[string]string),
		},
	}
}

// Load loads definitions of services, methods, messages, enumerations and fields from "req".
func (r *Registry) Load(req *plugin.CodeGeneratorRequest) error {
	for _, file := range req.GetProtoFile() {
		r.LoadFile(file)
	}

	var targetPkg string
	for _, name := range req.FileToGenerate {
		target := r.Files[name]
		if target == nil {
			return fmt.Errorf("no such file: %s", name)
		}
		name := gw_descriptor.PackageIdentityName(target.FileDescriptorProto)
		if targetPkg == "" {
			targetPkg = name
		} else {
			if targetPkg != name {
				return fmt.Errorf("inconsistent package names: %s %s", targetPkg, name)
			}
		}

		if err := r.loadServices(target); err != nil {
			return err
		}
	}
	return nil
}
