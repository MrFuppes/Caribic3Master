# -*- coding: utf-8 -*-
"""
Created on Mon Jul 27 17:38:58 2020

@author: va6504
"""

import socket
import struct
import time
import zlib

def msg_gen(addr_from: tuple, addr_to: tuple, msg_type: int, data: bytes, fmap: dict,
            timestamp=None):
    """
    Generate a message.
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



def msg_check(packet: bytes, addr_from: tuple, addr_to: tuple,
              fmap: dict, idxmap: dict,
              min_pac_len=22):
    """
    Check if a packet contains a valid message.
    """

    # make sure sender + receiver address in message:
    seq = (socket.inet_aton(addr_from[0]) + struct.pack(fmap['addr_to_port'], addr_from[1]) +
           socket.inet_aton(addr_to[0]) + struct.pack(fmap['addr_to_port'], addr_to[1]))

    # make sure address sequence is found
    if seq in packet:
        packet = packet[packet.index(seq):]
        # make sure packet contains at least header
        if len(packet) >= min_pac_len:
            msg_len = struct.unpack(fmap['len'], packet[idxmap['len']])[0]
            # make sure length of packet is sufficient
            if len(packet) >= msg_len:
                # truncate trailing bytes
                if len(packet) > msg_len:
                    packet = packet[:msg_len]
                # make sure checksum is correct:
                if zlib.adler32(packet[:-4]) == struct.unpack(fmap['cs'], packet[idxmap['cs']])[0]:
                    return packet
    return False


#------------------------------------------------------------------------------


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

# index positions for message parts
# could be loaded from config file
idxmap = {'addr_from_ip': slice(0, 4),
          'addr_from_port': slice(4, 6),
          'addr_to_ip': slice(6, 10),
          'addr_to_port': slice(10, 12),
          'len': slice(12, 14),
          'ts': slice(14, 22),
          'type': slice(22, 23),
          'data': slice(23, -4),
          'cs': slice(-4, None)}


# define message content
addr_from = ('192.168.1.1', 16101)
addr_to = ('192.168.1.64', 16164)
msg_type, data = 0, 'SB'.encode('ASCII')

packet = msg_gen(addr_from, addr_to, msg_type, data, fmap)

print('bytes in packet ->\n', list(map(hex, packet)))