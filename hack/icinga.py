#!/usr/bin/env python

from collections import OrderedDict
import fnmatch
import json
import os
import re
import subprocess
import sys
from os.path import expandvars
import urllib2
from xml.dom.minidom import parse, parseString
import tempfile
import shutil
import lxml.html
from lxml.cssselect import CSSSelector
import HTMLParser
from os.path import expanduser


def parse_md():
    p = re.compile(ur'###.*?</a>[ ](?P<cmd>[a-zA-Z0-9_-]+)\n\n(?P<desc>.*?)\n\n.*?---\n(?P<vars>.*?)\n\n', re.DOTALL)
    q = re.compile(ur'Examples:\n(?P<examples>.*?)Send email', re.DOTALL)
    r = re.compile(ur'(-(?P<short>\w),\s+)?--(?P<long>[a-zA-Z-_]+)=(?P<type>\S+)')

    link = "https://raw.githubusercontent.com/Icinga/icinga2/master/doc/10-icinga-template-library.md"
    resp = urllib2.urlopen(link).read()
    resp = HTMLParser.HTMLParser().unescape(resp)
    plugins = []
    for m in re.finditer(p, resp):
        g = m.groupdict()
        # print '%02d-%02d: %s' % (m.start(), m.end(), g['cmd'])
        plugin = {
            'name': g['cmd'],
            'description': g['desc'].replace('\n', ' '),
            'vars': []
        }
        desc = g['vars'].replace('.**', '')
        desc = desc.replace('**', '')
        for attr in desc.split('\n'):
            x = attr.split('|')
            if len(x) == 2:
                # print x[0].strip()
                # print x[1].strip()
                attrDesc = x[1].strip()
                optional = attrDesc.startswith('Optional ')
                if optional:
                    attrDesc = attrDesc[len('Optional '):].strip()
                plugin['vars'].append({
                    'name': x[0].strip().replace('\\_', '_'),
                    'optional': optional,
                    'description': attrDesc.replace('\n', ' ')
                })
        try:
            manpage_url = 'https://www.monitoring-plugins.org/doc/man/check_{0}.html'.format(plugin['name'])
            print 'manpage_url =================', manpage_url
            resp = urllib2.urlopen(manpage_url).read()
            tree = lxml.html.fromstring(resp)
            sel = CSSSelector('pre code')
            results = sel(tree)
            if len(results) > 0:
                match = lxml.html.tostring(sel(tree)[0])
                for line in match.split('\n'):
                    line = line.strip()
                    if line.startswith('-'):
                        m3 = re.search(r, line)
                        if m3:
                            g3 = m3.groupdict()
                            print '-----------------------------------', line, g3['long'], g3['long'].replace('-', ''), g3['long'].replace('-', '_')
                            for v in plugin['vars']:
                                if g3['long'] in v['name'] or g3['long'].replace('-', '') in v['name'] or g3['long'].replace('-', '_') in v['name']:
                                    if 'type' in g3:
                                        t = HTMLParser.HTMLParser().unescape(g3['type']).strip().strip("'")
                                        if t == 'INTEGER' or t == 'DOUBLE' or t == 'STRING':
                                            v['type'] = t
                                        elif t == 'SECONDS':
                                            v['type'] = 'INTEGER'
                                            v['format'] = t
                                        elif t == 'INTEGER[,INTEGER]':
                                            v['type'] = 'INTEGER'
                                        elif '|' in t:
                                            v['type'] = 'STRING'
                                            v['format'] = 'enum'
                                            v['values'] = t.split('|')
                                        else:
                                            v['type'] = 'STRING'
                                            v['format'] = t
                                    v['flag'] = {
                                        'long':  g3['long']
                                    }
                                    if 'short' in g3:
                                        v['flag']['short'] = g3['short']
                                    print '==========', v
                m2 = re.search(q, match)
                if m2:
                    g2 = m2.groupdict()
                    plugin['examples'] = g2['examples'].strip()

            plugins.append(plugin)
        except urllib2.HTTPError, e:
            pass
            # print e.code
            # print e.msg
            # print e.headers
            # print e.fp.read()

    schema = {
        'plugin': plugins
    }
    with open(expanduser("~") + '/go/src/go.searchlight.dev/icinga-operator/data/files/icinga.gen.json', 'w') as f:
        return json.dump(schema, f, sort_keys=True, indent=2, separators=(',', ': '))


if __name__ == "__main__":
    parse_md()
