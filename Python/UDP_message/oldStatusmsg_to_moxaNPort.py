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


N, wait = 10, 2
send_from = ('192.168.1.1', 4001)
send_to = ('192.168.127.254', 4001) # moxa nport
status_set = 'MS'

data_rec = b''

with socket.socket(family=socket.AF_INET, type=socket.SOCK_DGRAM) as UDPServerSocket:
    UDPServerSocket.bind(send_from)
    UDPServerSocket.settimeout(wait)
    for _ in range(N):
        # ------
        message = ('\x02' + '216,' +
                   datetime.now(tz=timezone.utc).strftime('%H:%M:%S.%f')[:-4] +
                   ',' + status_set + '\x03')
        checksum = 0
        for el in message[1:]:
            checksum ^= ord(el)
        message = bytes(message, 'ASCII') + struct.pack("B", checksum)
        # ------
        sent = UDPServerSocket.sendto(message, send_to)

        print(f'sent {sent} bytes from {send_from} to {send_to}')
        print(f'message: {message}')
        print('awaiting reply...')
        try:
            data, addr = UDPServerSocket.recvfrom(1024)
        except socket.timeout:
            print('timed out.')
        else:
            print(data)
            data_rec += data
            # break
        time.sleep(wait)

# s = []

# for port in range(16101,16111):
#     server_socket = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
#     server_socket.bind(('192.168.1.1', port))
#     s.append(server_socket)

# for so in s:
#     so.sendto('asdf'.encode('utf-8'), ('192.168.1.64', 16164))
#     so.close()

print("done. received:")
print(data_rec)
