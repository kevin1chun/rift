import random

from twisted.internet import reactor, defer
from twisted.names import client, dns, error, server, resolve

import services

class RiftResolver(object):
    def __init__(self, domain):
        self.domain = domain


    def _get_answers(self, name):
        parts = name.split('.')
        tld = parts[-1]
        thld = parts[0]
        sld = '.'.join(parts[1:-1])
        if sld in self.domain and thld in self.domain[sld]:
            return self.domain[sld][thld]
        else:
            return None


    def query(self, query, timeout):
        name = query.name.name
        results = self._get_answers(name)
        if results:
            print('In my domain: %s' % name)
            random.shuffle(results)
            forward = (dns.RRHeader(name=results[0], type=dns.CNAME, payload=dns.Record_CNAME(name=results[0])),),(),() 
            return defer.succeed(forward)
        return defer.fail(error.DomainError())


if __name__ == '__main__':
    backup_resolvers = [('8.8.8.8', 53), ('8.8.4.4', 53)]
    backup = client.Resolver(servers=backup_resolvers)

    domain = services.build_services()

    print('Domain discovered, starting DNS...')

    factory = server.DNSServerFactory(
        clients=[resolve.ResolverChain([RiftResolver(domain), backup])]
    )

    protocol = dns.DNSDatagramProtocol(controller=factory)

    factory.noisy = protocol.noisy = False

    reactor.listenUDP(53, protocol)
    reactor.listenTCP(53, factory)

    reactor.run()

