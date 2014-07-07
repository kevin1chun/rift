import sys
import json

from twisted.internet import reactor, task
from twisted.internet.protocol import Protocol, Factory

class RProtocol(Protocol):
    def connectionMade(self):
        print('Cxn created with peer: %s' % self.transport.getPeer())

    def dataReceived(self, data):
        if data[0] == 'L':
            print('Sending service descriptors')
            self.transport.write(json.dumps(self.factory.services))
            self.transport.loseConnection()
        elif data[0] == 'R':
            serv_desc = json.loads(data[1:])
            print('Adding service: %s' % serv_desc)
            self.factory.services.append(serv_desc)
        else:
            print('Unknown request: %s' % data)


def build_rproto_factory(services=[]):
    factory = Factory()
    factory.protocol = RProtocol
    factory.services = services
    return factory


if __name__ == '__main__':
    reactor.listenTCP(int(sys.argv[1]), build_rproto_factory())
    reactor.run()

