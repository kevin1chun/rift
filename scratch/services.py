from schema import SchemaOrg

with open('./providers.txt') as f:
    providers = {t.lower(): p.split(',') for t, p in [l.split() for l in f.readlines()]}


def build_services():
    with open('./schema.json') as f:
        s = SchemaOrg(json_file=f)
    actions = [type.name.lower() for type in s.types if 'Action' in type.ancestry]
    return {
        'org.schema': {k: v for k, v in providers.iteritems() if k in actions}
    }

