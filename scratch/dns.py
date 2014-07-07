from random import shuffle

from twisted.internet import reactor
from twisted.names import dns, server, client, cache
from twisted.application import service, internet

import services

 
class RiftResolver(client.Resolver):
    def __init__(self, domain, servers):
        client.Resolver.__init__(self, servers=servers)
        self.domain = domain
        self.ttl = 10
 
    def lookupAddress(self, name, timeout=None):
        parts = name.split('.')
        tld = parts[-1]
        thld = parts[0]
        sld = '.'.join(parts[1:-1])
        if sld in self.domain and thld in self.domain[sld]:
            values = self.domain[sld][thld]
            results = [dns.RRHeader(name, dns.CNAME, dns.IN, self.ttl, dns.Record_CNAME(value, self.ttl)) for value in values]
            shuffle(results)
            return [tuple(results), (), ()]
        else:
            return self._lookup(name, dns.IN, dns.A, timeout)
 
 
 
domain = services.build_services() 

upstream_dns = '8.8.8.8'
 
simpledns = RiftResolver(domain, servers=[(upstream_dns, 53)])
 
f = server.DNSServerFactory(clients=[simpledns])
p = dns.DNSDatagramProtocol(f)
f.noisy = p.noisy = False
 
reactor.listenUDP(53, p)
reactor.listenTCP(53, f)
reactor.run()
