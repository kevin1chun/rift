"""
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
"""

from urllib2 import urlopen
import json

class SchemaOrg:
    def __init__(self, json_file=None, json_url='http://schema.rdfs.org/all.json'):
        j = json_file or urlopen(json_url)
        self.json_schema = json.load(j)

    @property
    def types(self):
        props = dict({})
        for prop, definition in self.json_schema['properties'].items():
            for t in definition['domains']:
                if not t in props:
                    props[t] = []
                props[t].append({ 'name': prop, 'type': definition['ranges'] })

        return [SchemaOrgType(name, definition['ancestors'], props[name] if name in props else []) for name, definition in self.json_schema['types'].items()]

class SchemaOrgType:
    def __init__(self, name, ancestry, props):
        self.name = name
        self.ancestry = ancestry
        self.props = props

    def __str__(self):
        return '"%s" (%s): %s' % (self.name, self.ancestry, self.props)
