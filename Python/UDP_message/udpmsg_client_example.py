# -*- coding: utf-8 -*-
"""
Created on Wed Jul 22 11:41:44 2020

@author: f.obersteiner
"""

# https://de.wikipedia.org/wiki/Liste_der_standardisierten_Ports#Registrierte_Ports:_1024%E2%80%9349151

# Instrument #64, receives message
#           on 192.168.232.64:16064
#   replies to 192.168.232.1:16064

from datetime import datetime, timezone
import socket
import struct
import time

TIMEOUT = 5
recv_on = ('192.168.232.64', 16064)
repl_on = ('192.168.232.1', 16064)
reply = "message received.".encode('utf-8')


# Create a UDP socket
with socket.socket(socket.AF_INET, socket.SOCK_DGRAM) as sock:
    sock.bind(recv_on)
    sock.settimeout(TIMEOUT)
    print(f"listening on {recv_on}")
    data, addr = sock.recvfrom(1024) # buffer size is 1024 bytes
    print(f"received message {data} from {addr}")
    t = time.time()
    print(f"sending reply on {datetime.fromtimestamp(t, tz=timezone.utc).isoformat()}")
    reply = struct.pack('!d', t) + data[:8] + reply
    sent = sock.sendto(reply, repl_on)

print("done.")