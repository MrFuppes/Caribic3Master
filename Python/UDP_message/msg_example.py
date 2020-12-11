# -*- coding: utf-8 -*-
"""
Created on Tue Jul 21 14:03:15 2020

@author: f.obersteiner
"""
import random
import socket
import struct
import time
import zlib

# define message content
addr_from = ('192.168.232.1', 16001)
addr_to = ('192.168.232.64', 16064)
msg_type, status_str = 0, 'SB'

# format codes for struct, for each meassage part
fmap = {'addr_from_ip': '!BBBB',
        'addr_from_port': '!H',
        'addr_to_ip': '!BBBB',
        'addr_to_port': '!H',
        'len': '!H',
        'ts': '!d',
        'type': '!B',
        'cs': '!I'}

# pack message content to binary string
packet = (socket.inet_aton(addr_from[0]) + struct.pack(fmap['addr_from_port'], addr_from[1]) +
          socket.inet_aton(addr_to[0]) + struct.pack(fmap['addr_to_port'], addr_to[1]) +
          b'\x00\x00' + # place holder bytes for packet length
          struct.pack(fmap['ts'], time.time()) +
          struct.pack(fmap['type'], msg_type) +
          status_str.encode('ASCII'))

# insert packet length, including checksum bytes
packet = packet[:12] + struct.pack(fmap['len'], len(packet)+4) + packet[14:]

# add checksum
packet += struct.pack(fmap['cs'], zlib.adler32(packet))

# add random number of random bytes...
packet = (bytearray(random.getrandbits(8) for _ in range(random.randrange(10))) +
          packet +
          bytearray(random.getrandbits(8) for _ in range(random.randrange(10))))

#------------------------------------------------------------------------------

# assuming 'packet' is received...

# index positions for message parts
idxmap = {'addr_from_ip': slice(0, 4),
          'addr_from_port': slice(4, 6),
          'addr_to_ip': slice(6, 10),
          'addr_to_port': slice(10, 12),
          'len': slice(12, 14),
          'ts': slice(14, 22),
          'type': slice(22, 23),
          'data': slice(23, -4),
          'cs': slice(-4, None)}

# make sure sender + receiver address in message:
checksequence = (socket.inet_aton(addr_from[0]) + struct.pack('!H', addr_from[1]) +
                 socket.inet_aton(addr_to[0]) + struct.pack('!H', addr_to[1]))

assert checksequence in packet, "checksequence not found"

# truncate preceeding bytes
packet = packet[packet.index(checksequence):]

# make sure at least header is complete
assert len(packet) >= 22, "incomplete packet"

msg_len = struct.unpack('!H', packet[idxmap['len']])[0]

# assure correct packet contains full message
assert len(packet) >= msg_len, "incomplete packet"

# truncate message if it has trailing bytes
if len(packet) > msg_len:
    packet = packet[:msg_len]

# assure correct checksum
assert zlib.adler32(packet[:-4]) == struct.unpack(fmap['cs'], packet[idxmap['cs']])[0], "invalid packet checksum"

unpacked = (('.'.join(map(str, struct.unpack(fmap['addr_from_ip'], packet[idxmap['addr_from_ip']]))),
             struct.unpack(fmap['addr_from_port'], packet[idxmap['addr_from_port']])[0]),
            ('.'.join(map(str, struct.unpack(fmap['addr_to_ip'], packet[idxmap['addr_to_ip']]))),
             struct.unpack(fmap['addr_to_port'], packet[idxmap['addr_to_port']])[0]),
            struct.unpack(fmap['len'], packet[idxmap['len']])[0],
            struct.unpack(fmap['ts'], packet[idxmap['ts']])[0],
            struct.unpack(fmap['type'], packet[idxmap['type']])[0],
            ''.join(map(chr, packet[idxmap['data']])))

print("retrieved message:", unpacked)

# print(list(map(hex, packet)))

# size as string would approximately double:
# len(';'.join(('192.168.1.1', '192.168.1.64', '24', '0', '1595337472.0', 'SB')))
# >>> 45