# -*- coding: utf-8 -*-
"""
Created on Wed Jul 22 11:37:15 2020

@author: f.obersteiner
"""

# https://de.wikipedia.org/wiki/Liste_der_standardisierten_Ports#Registrierte_Ports:_1024%E2%80%9349151

# Master Computer, sends message
#   from 192.168.232.1:16001
#   to 192.168.232.64:16064
#   listens for reply on 192.168.232.1:16064

from datetime import datetime, timezone
import socket
import struct
import time


N, wait = 5, 1

send_from = ('192.168.232.1', 16001) # the socket from which to send
send_to = ('192.168.232.64', 16064) # the target socket

message = 'this is a test\n'.encode('utf-8')

with socket.socket(family=socket.AF_INET, type=socket.SOCK_DGRAM) as UDPServerSocket:
    UDPServerSocket.bind(send_from)
    UDPServerSocket.settimeout(wait)
    for _ in range(N):
        t = time.time()
        sent = UDPServerSocket.sendto((struct.pack('!d', t) + message), send_to)

        print(f'sent {sent} bytes from {send_from} to {send_to}')
        print(f'message: {(struct.pack("!d", t) + message)}')
        print('awaiting reply...')
        try:
            data, addr = UDPServerSocket.recvfrom(1024)
        except socket.timeout:
            print('timed out.')
        else:
            print(data[:8], data[8:16])
            t_repl = datetime.fromtimestamp(
                struct.unpack('!d', data[:8])[0],
                tz=timezone.utc).isoformat()
            t_sent = datetime.fromtimestamp(
                struct.unpack('!d', data[8:16])[0],
                tz=timezone.utc).isoformat()
            print(data[16:], t_sent, t_repl)
            print(struct.unpack('!d', data[:8])[0]-struct.unpack('!d', data[8:16])[0])
            break

print("done.")