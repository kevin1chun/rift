import socket
import json

RN = ('172.16.1.52', 8674)

def register(desc):
    sck = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    sck.connect(RN)
    sck.send('%s%s' % ('R', json.dumps(desc)))
    sck.close()
