# -*- coding: utf-8 -*-
"""
Created on Fri Sep 25 09:50:54 2020

@author: f.obersteiner
"""
from datetime import datetime, timezone
import socket
import struct
import zlib

from msg_compose import composeMsg

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

#------------------------------------------------------------------------------
def checkMsg(packet: bytes,
             addr_from: tuple, addr_to: tuple,
             fmap: dict,
             idxmap: dict,
             min_pac_len=22,
             verbose=False):
    """
    Check if a packet contains a valid message according to protocol [add specification].
    """
    verboseprint = print if verbose else lambda *a, **k: None
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
                else:
                    verboseprint("invalid checksum")
            else:
                verboseprint(f"invalid packet length, {msg_len} specified vs. {len(packet)} bytes found.")
        else:
            verboseprint("message too short, minimum {min_pac_len} vs. {len(packet)} bytes found.")
    else:
        verboseprint("sender/receiver sequence not found in packet.")

    return False
#------------------------------------------------------------------------------


#------------------------------------------------------------------------------
def parseMsg(packet: bytes,
             fmap: dict,
             idxmap: dict,
             data_fmt=None,
             ts_to_dtobj=False):
    """
    Parse a packet to a dictionary.
    Make sure the packet is a valid message first (run checkMsg).
    """
    content = {'addr_from_ip': '',
                'addr_from_port': -1,
                'addr_to_ip': '',
                'addr_to_port': -1,
                'len': -1,
                'ts': -1,
                'type': -1,
                'data': '',
                'cs': -1}

    for k, fmt in fmap.items():
        content[k] = struct.unpack(fmt, packet[idxmap[k]])
        if len(content[k]) == 1:
            content[k] = content[k][0]

    if data_fmt: # if a format for the data is given, unpack the bytes:
        content['data'] = struct.unpack(data_fmt, packet[idxmap['data']])
    else: # else just leave it as a bytes string.
        content['data'] = packet[idxmap['data']]

    if ts_to_dtobj:
        content['ts'] = datetime.fromtimestamp(content['ts'], tz=timezone.utc)

    return content
#------------------------------------------------------------------------------


#------------------------------------------------------------------------------
if __name__ == '__main__':
    # define message content
    FROM = ('192.168.232.1', 16001)
    TO = ('192.168.232.64', 16064)
    msg_type = 0
    data = 'SB'.encode('ASCII')
    # ...and create a message packet:
    packet = composeMsg(FROM, TO, msg_type, data, fmap)
    print('bytes in packet ->\n', list(map(hex, packet)))

    assert checkMsg(packet, FROM, TO, fmap, idxmap), "invalid message"
    print('message packet passed test :)\n')

    packet = bytes.fromhex('C0A8 0101 3EE5 C0A8 0140 3F24 001E 41D7 DB70 7B6D A254 0077 7466 8B67 0A4A')
    assert checkMsg(packet, FROM, TO, fmap, idxmap), "invalid message"

    parsed = parseMsg(packet, fmap, idxmap, ts_to_dtobj=True)
    for k, v in parsed.items():
        print(f"{k}: {v}")