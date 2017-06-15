#!/usr/bin/env python

import fnmatch
import json
import os
import re
import subprocess
import sys
from itertools import chain
from os.path import expandvars

# binary -> GOOS -> GOARCH
# windows	386
# windows	amd64

# Debian package
# https://gist.github.com/rcrowley/3728417

REPO_ROOT = expandvars("$GOPATH/src/github.com/appscode/api")
VALID_FORMATS = ['date-time',
                 'email',
                 'hostname',
                 'ipv4',
                 'ipv6',
                 'uri']
PROTO_PKG_REGEX = ur'^package\s*(?P<pkg>[^;]*);$'
GO_PKG_REGEX = ur'^option go_package\s*=\s*"(?P<pkg>[^;]*)"\s*;$'
REQ_REGEX = ur'^\s*rpc\s*\w+\s*\(\s*(?P<req>\w+)\s*\)'
RESP_REGEX = ur'^\s*rpc\s*\w+\s*\(\s*[^)]+\)\s*returns\s*\(\s*(?P<resp>\w+)\s*\)'
URL_REGREX = ur'^\s*(?P<verb>get|post|put|patch|delete)\s*:\s*"(?P<url>[^"]+)"'


def call(cmd, stdin=None, cwd=REPO_ROOT):
    print(cmd)
    subprocess.call([expandvars(cmd)], shell=True, stdin=stdin, cwd=cwd)


def write_file(name, content):
    dir = os.path.dirname(name)
    if not os.path.exists(dir):
        os.makedirs(dir)
    with open(name, 'w') as f:
        return f.write(content)


def append_file(name, content):
    with open(name, 'a') as f:
        return f.write(content)


# TODO: use unicode encoding
def read_json(name):
    try:
        with open(name, 'r') as f:
            return json.load(f)
    except IOError:
        return {}


def write_json(obj, name):
    with open(name, 'w') as f:
        return json.dump(obj, f, sort_keys=True, indent=2, separators=(',', ': '))


def gen_assets():
    call('go get github.com/jteeuwen/go-bindata/...')
    call('go-bindata -ignore=\\.go -ignore=\\.DS_Store -mode=0644 -modtime=1453795200 -o meta/data.go -pkg meta meta/...')


def fix_swagger_schema():
    for root, dirnames, filenames in os.walk(REPO_ROOT):
        for filename in fnmatch.filter(filenames, '*.swagger.json'):

            rel_path = root[len(REPO_ROOT) + 1:]
            parts = rel_path.split('/', 2)
            if len(parts) != 2:
                continue
            # api_name = parts[0]
            # api_version = parts[1]

            swagger = os.path.join(root, filename)
            spec = read_json(swagger)
            spec["basePath"] = "/"
            spec["host"] = "api.appscode.com"
            spec["info"]["version"] = parts[1][1:]
            spec["schemes"] = ["https"]
            write_json(spec, swagger)


def swagger_defs(defs):
    stack = []
    result = {
        'requests': {},
        'responses': {}
    }
    for key in defs.keys():
        if key.endswith("Request"):
            stack.append(key)
        elif key.endswith("Response"):
            result['responses'][key] = defs[key]
    while stack:
        name = stack.pop()
        schema = defs[name]
        result['requests'][name] = schema
        if 'properties' in schema:
            for p, v in schema['properties'].items():
                if '$ref' in v:
                    nw_obj = v['$ref'][len('#/definitions/'):]
                    if nw_obj not in result['requests']:
                        stack.append(nw_obj)
                if 'format' in v and not v['format'] in VALID_FORMATS:
                    v.pop('format', None)
                if 'additionalProperties' in v:
                    vprops = v['additionalProperties']
                    if 'format' in vprops and not vprops['format'] in VALID_FORMATS:
                        vprops.pop('format', None)
                    if '$ref' in vprops:
                        nw_obj = vprops['$ref'][len('#/definitions/'):]
                        if nw_obj not in result['requests']:
                            stack.append(nw_obj)
                if 'items' in v:
                    if '$ref' in v['items']:
                        nw_obj = v['items']['$ref'][len('#/definitions/'):]
                        if nw_obj not in result['requests']:
                            stack.append(nw_obj)
                    if 'format' in v['items'] and not v['items']['format'] in VALID_FORMATS:
                        v['items'].pop('format', None)
    return result


def generate_json_schema():
    call("find . | grep schema.json | xargs rename 's/schema.json/schema.json.ext/' {}")
    for root, dirnames, filenames in os.walk(REPO_ROOT):
        for filename in fnmatch.filter(filenames, '*.swagger.json'):
            swagger = os.path.join(root, filename)
            schema = os.path.join(root, filename.replace('.swagger.', '.schema.'))
            ext_schema = os.path.join(root, filename.replace('.swagger.json', '.schema.json.ext'))
            gen_defs = swagger_defs(read_json(swagger)['definitions'])['requests']
            if os.path.exists(ext_schema):
                # merge
                ext_defs = read_json(ext_schema)['definitions']
                for m, mspec in gen_defs.items():
                    if 'properties' in mspec.keys():
                        for f, fspec in mspec['properties'].items():
                            if f in [
                                'auth_secret_name',
                                'bucket_name',
                                'cloud_credential',
                                'cluster_name',
                                'name',
                                'namespace',
                                'secret_name',
                                'service_name',
                                'snapshot_name'
                            ]:
                                if 'maxLength' not in fspec:
                                    fspec['maxLength'] = 63
                                if 'pattern' not in fspec:
                                    fspec['pattern'] = "^[a-z0-9](?:[a-z0-9\-]{0,61}[a-z0-9])?$"

                            if m in ext_defs \
                                    and 'properties' in ext_defs[m] \
                                    and f in ext_defs[m]['properties'] \
                                    and set(fspec.keys()) != set(ext_defs[m]['properties'][f].keys()):
                                print(mspec['properties'][f])
                                print(ext_defs[m]['properties'][f])
                                mspec['properties'][f] = ext_defs[m]['properties'][f]
            write_json({'definitions': gen_defs}, schema)
    call("(find . | grep schema.json.ext | xargs rm) || true")


def schema_go(reqs, resps, defs):
    result = {
        'requests': {},
        'responses': {}
    }
    for key in defs['requests'].keys():
        if key in reqs and key.endswith("Request"):
            schema = defs['requests'][key]
            result['requests'][reqs[key]] = schema
            dep_defs = {}
            stack = []
            stack.append(schema)
            while stack:
                sch = stack.pop()
                if 'properties' in sch:
                    for p, v in sch['properties'].items():
                        if '$ref' in v:
                            nw_obj = v['$ref'][len('#/definitions/'):]
                            if nw_obj not in dep_defs:
                                dep_defs[nw_obj] = defs['requests'][nw_obj]
                                stack.append(defs['requests'][nw_obj])
                        if 'items' in v and '$ref' in v['items']:
                            nw_obj = v['items']['$ref'][len('#/definitions/'):]
                            if nw_obj not in dep_defs:
                                dep_defs[nw_obj] = defs['requests'][nw_obj]
                                stack.append(defs['requests'][nw_obj])
                        if 'additionalProperties' in v and '$ref' in v['additionalProperties']:
                                nw_obj = v['additionalProperties']['$ref'][len('#/definitions/'):]
                                if nw_obj not in dep_defs:
                                    dep_defs[nw_obj] = defs['requests'][nw_obj]
                                    stack.append(defs['requests'][nw_obj])
            if dep_defs:
                schema['definitions'] = dep_defs
            schema['$schema'] = 'http://json-schema.org/draft-04/schema#'
    for key in defs['responses'].keys():
        # print key
        if key in resps and key.endswith("Response"):
            schema = defs['responses'][key]
            result['responses'][resps[key]] = schema
    return result


def render_schema_go(pkg, schemas):
    contents = """package {}

// Auto-generated. DO NOT EDIT.
""".format(pkg)

    imports = []
    if schemas['requests']:
        imports.append("github.com/xeipuuv/gojsonschema")
        imports.append("github.com/golang/glog")
    imports.sort()
    if imports:
        contents += 'import (\n'
        for pkg in imports:
            contents += '	"{}"\n'.format(pkg)
        contents += ')\n\n'

    for key in schemas['requests'].keys():
        contents += 'var {0}Schema *gojsonschema.Schema\n'.format(key[0:1].lower() + key[1:])
    contents += '\n'

    if schemas['requests']:
        contents += """func init() {
	var err error
"""
        for key, sch in schemas['requests'].items():
            var_name = key[0:1].lower() + key[1:]
            sch_str = json.dumps(sch, sort_keys=True, indent=2, separators=(',', ': '))
            sch_str = sch_str.replace('`', '` + "`" + `')
            contents += '	{0}Schema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{1}`))\n'.format(
                var_name, sch_str)
            contents += """	if err != nil {
		glog.Fatal(err)
	}
"""
        contents += '}\n\n'

    for key in schemas['requests'].keys():
        contents += 'func (m *' + key + ') IsValid() (*gojsonschema.Result, error) {\n'
        contents += '	return {}Schema.Validate(gojsonschema.NewGoLoader(m))\n'.format(key[0:1].lower() + key[1:])
        contents += '}\n'
        contents += 'func (m *' + key + ') IsRequest() {}\n\n'

    return contents


def detect_objs(swagger):
    proto = swagger[:-len('.swagger.json')] + '.proto'
    proto_reqs = []
    proto_resps = []
    reqs = {}
    resps = {}
    with open(proto) as f:
        content = f.read()
        proto_reqs = re.findall(REQ_REGEX, content, re.MULTILINE)
        proto_resps = re.findall(RESP_REGEX, content, re.MULTILINE)
        defs = read_json(swagger)['definitions']
        m = re.search(PROTO_PKG_REGEX, content, re.MULTILINE)
        if m:
            pkg = m.group('pkg')
            parts = pkg.split(".")
            for x in reversed(range(0, len(parts))):
                prefix = str.join("", parts[x:])
                found = True
                for obj in chain(proto_reqs, proto_resps):
                    found &= ((prefix + obj) in defs)
                if found:
                    for r in proto_reqs:
                        reqs[prefix + r] = r
                    for r in proto_resps:
                        resps[prefix + r] = r
        else:
            print("Failed to detect package.")
            sys.exit(1)
    return reqs, resps


def detect_schema_pkg(swagger):
    proto = swagger[:-len('.swagger.json')] + '.proto'
    with open(proto) as f:
        content = f.read()
        m = re.search(PROTO_PKG_REGEX, content, re.MULTILINE)
        if m:
            pkg = m.group('pkg')
            parts = pkg.split(".")
            for x in range(0, len(parts)):
                prefix = str.join("", parts[x:])
                for key, defs in read_json(swagger)['definitions'].items():
                    if key.startswith(prefix) and key.endswith("Request"):
                        return prefix
        else:
            print("Failed to detect package.")
            sys.exit(1)


def detect_go_pkg(swagger):
    proto = swagger[:-len('.swagger.json')] + '.proto'
    with open(proto) as f:
        content = f.read()
        m = re.search(GO_PKG_REGEX, content, re.MULTILINE)
        if m:
            return m.group('pkg')
        else:
            m = re.search(PROTO_PKG_REGEX, content, re.MULTILINE)
            if m:
                return m.group('pkg')
            else:
                print("Failed to detect package.")
                sys.exit(1)


def generate_go_schema():
    for root, dirnames, filenames in os.walk(REPO_ROOT):
        for filename in fnmatch.filter(filenames, '*.swagger.json'):
            swagger = os.path.join(root, filename)
            schema = os.path.join(root, filename.replace('.swagger.', '.schema.'))
            go = schema[:-len('.json')] + '.go'
            # print go
            reqs, resps = detect_objs(swagger)
            if reqs or resps:
                defs = swagger_defs(read_json(swagger)['definitions'])
                # overwrite requests with json schema from *.schema.json
                # to preserve hand written rules
                defs['requests'] = read_json(schema)['definitions']
                schemas = schema_go(reqs, resps, defs)
                if schemas['requests'] or schemas['responses']:
                    write_file(go, render_schema_go(detect_go_pkg(swagger), schemas))


def apply_naming_policy():
    for root, dirnames, filenames in os.walk(REPO_ROOT):
        for filename in fnmatch.filter(filenames, '*.schema.json'):
            schema = os.path.join(root, filename)
            print(schema)
            content = read_json(schema)
            for key, defs in content['definitions'].items():
                if 'properties' in defs:
                    for p, v in defs['properties'].items():
                        if p in [
                            'cluster_name',
                            'namespace', 'name',
                            'bucket_name',
                            'secret_name',
                            'snapshot_name',
                            'auth_secret_name',
                            'cloud_credential'
                        ]:
                            print('====>>>> ' + p)
                            if 'maxLength' not in v:
                                v['maxLength'] = 63
                            if 'pattern' not in v:
                                v['pattern'] = "^[a-z0-9](?:[a-z0-9\-]{0,61}[a-z0-9])?$"
            write_json(content, schema)


def generate_url_summary():
    urls = {}
    for root, dirnames, filenames in os.walk(REPO_ROOT):
        for filename in fnmatch.filter(filenames, '*.proto'):
            proto = os.path.join(root, filename)
            rel_proto = proto[proto.index('github.com/appscode/api') + len('github.com/appscode/api'):]
            with open(proto) as f:
                content = f.read()
                apis = re.findall(URL_REGREX, content, re.MULTILINE)
                for verb, url in apis:
                    # form url key by replacing {} parts with *
                    parts = url.split('/')
                    for idx, part in enumerate(parts):
                        if part.startswith('{') and part.endswith('}'):
                            parts[idx] = '*'
                    url_key = '/'.join(parts)
                    # check that verb is not redefined
                    if url_key in urls:
                        if verb in urls[url_key]:
                            print(verb, 'verb is redefined for URL:', url_key)
                            print('previously defined as ', urls[url_key][verb]['url'], 'in', urls[url_key][verb]['file'])
                            sys.exit(1)
                    else:
                        urls[url_key] = {}
                    urls[url_key][verb] = {
                        'url': url,
                        'file': rel_proto
                    }
    write_json(urls, REPO_ROOT + '/api_urls.json')


if __name__ == "__main__":
    if len(sys.argv) > 1:
        # http://stackoverflow.com/a/834451
        # http://stackoverflow.com/a/817296
        globals()[sys.argv[1]](*sys.argv[2:])
    else:
        generate_json_schema()
        # apply_naming_policy()
        generate_go_schema()
        generate_url_summary()
