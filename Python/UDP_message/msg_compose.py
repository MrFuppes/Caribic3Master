# -*- coding: utf-8 -*-
"""
Created on Mon Jul 27 17:38:58 2020

@author: va6504
"""
import socket
import struct
import time
import zlib



def composeMsg(addr_from: tuple, addr_to: tuple,
               msg_type: int,
               data: bytes,
               fmap: dict,
               timestamp=None):
    """
    Compose a message according to protocol [add specification].
    """
    if timestamp is None:
        timestamp = time.time()

    # pack message content to binary string
    packet = (socket.inet_aton(addr_from[0]) + struct.pack(fmap['addr_from_port'], addr_from[1]) +
              socket.inet_aton(addr_to[0]) + struct.pack(fmap['addr_to_port'], addr_to[1]) +
              b'\x00\x00' + # place holder bytes for packet length
              struct.pack(fmap['ts'], timestamp) +
              struct.pack(fmap['type'], msg_type) +
              data)

    # insert packet length, including checksum bytes
    packet = packet[:12] + struct.pack(fmap['len'], len(packet)+4) + packet[14:]

    # add checksum
    packet += struct.pack(fmap['cs'], zlib.adler32(packet))

    return packet



#------------------------------------------------------------------------------
if __name__ == '__main__':

    # format codes for struct, for each meassage part
    # could be loaded from config file
    fmap = {'addr_from_ip': '!BBBB',
            'addr_from_port': '!H',
            'addr_to_ip': '!BBBB',
            'addr_to_port': '!H',
            'len': '!H',
            'ts': '!d',
            'type': '!B',
            'cs': '!I'}

    # define message content
    addr_from = ('192.168.1.1', 16101)
    addr_to = ('192.168.1.64', 16164)
    msg_type, data = 0, 'SB'.encode('ASCII')

    packet = composeMsg(addr_from, addr_to, msg_type, data, fmap)
    print('bytes in packet ->\n', list(map(hex, packet)))
