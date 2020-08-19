# -*- coding: utf-8 -*-
"""
Created on Wed Jul 22 11:41:44 2020

@author: va6504
"""

# https://de.wikipedia.org/wiki/Liste_der_standardisierten_Ports#Registrierte_Ports:_1024%E2%80%9349151

# Instrument #64, receives message
#   on 192.168.1.64:16164
#   replies to 192.168.1.1:16101

from datetime import datetime, timezone
import socket
import struct
import time

TIMEOUT = 10
recv_on = ('192.168.1.64', 16164)
repl_on = ('192.168.1.1', 16101)

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


