package descriptor

import (
	"fmt"

	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
	gw_descriptor "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway/descriptor"
	acopts "github.com/grpc-ecosystem/grpc-gateway/third_party/appscodeapis/appscode/api"
)

// loadServices registers services and their methods from "targetFile" to "r".
// It must be called after loadFile is called for all files so that loadServices
// can resolve names of message types and their fields.
func (r *Registry) loadServices(file *gw_descriptor.File) error {
	glog.V(1).Infof("Loading services from %s", file.GetName())
	var svcs []*gw_descriptor.Service
	for _, sd := range file.GetService() {
		glog.V(2).Infof("Registering %s", sd.GetName())
		svc := &gw_descriptor.Service{
			File: file,
			ServiceDescriptorProto: sd,
		}
		for _, md := range sd.GetMethod() {
			glog.V(2).Infof("Processing %s.%s", sd.GetName(), md.GetName())
			corsOpts, err := extractCorsOptions(md)
			if err != nil {
				glog.Errorf("Failed to extract ApiMethodOptions from %s.%s: %v", svc.GetName(), md.GetName(), err)
				return err
			}
			if corsOpts == nil {
				glog.V(1).Infof("Skip non-target method: %s.%s", svc.GetName(), md.GetName())
				continue
			}
			if corsOpts.Enable {
				opts, err := gw_descriptor.ExtractAPIOptions(md)
				if err != nil {
					glog.Errorf("Failed to extract ApiMethodOptions from %s.%s: %v", svc.GetName(), md.GetName(), err)
					return err
				}
				if opts == nil {
					glog.V(1).Infof("Skip non-target method: %s.%s", svc.GetName(), md.GetName())
					continue
				}
				meth, err := r.NewMethod(svc, md, opts)
				if err != nil {
					return err
				}
				svc.Methods = append(svc.Methods, meth)
			}
		}
		glog.V(2).Infof("Registered %s with %d method(s)", svc.GetName(), len(svc.Methods))
		svcs = append(svcs, svc)
	}
	file.Services = svcs
	return nil
}

func extractCorsOptions(meth *descriptor.MethodDescriptorProto) (*acopts.CorsRule, error) {
	if meth.Options == nil {
		return nil, nil
	}
	if !proto.HasExtension(meth.Options, acopts.E_Cors) {
		return nil, nil
	}
	ext, err := proto.GetExtension(meth.Options, acopts.E_Cors)
	if err != nil {
		return nil, err
	}
	opts, ok := ext.(*acopts.CorsRule)
	if !ok {
		return nil, fmt.Errorf("extension is %T; want an CorsRule", ext)
	}
	return opts, nil
}
