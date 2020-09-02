# -*- coding: utf-8 -*-
"""
Created on Wed Jul 22 11:37:15 2020

@author: va6504
"""

# https://de.wikipedia.org/wiki/Liste_der_standardisierten_Ports#Registrierte_Ports:_1024%E2%80%9349151

# Master Computer, sends message
#   from 192.168.1.1:16100
#   to 192.168.1.64:16164
#   listens for reply on 192.168.1.1:16101

from datetime import datetime, timezone
import socket
import struct
import time


N, wait = 5, 1

send_from = ('192.168.1.1', 16101) # the socket from which to send
send_to = ('192.168.1.64', 16164) # the target socket


send_from = ('192.168.1.1', 4001)
send_to = ('192.168.127.254', 4001)

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

# s = []

# for port in range(16101,16111):
#     server_socket = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
#     server_socket.bind(('192.168.1.1', port))
#     s.append(server_socket)

# for so in s:
#     so.sendto('asdf'.encode('utf-8'), ('192.168.1.64', 16164))
#     so.close()

print("done.")
