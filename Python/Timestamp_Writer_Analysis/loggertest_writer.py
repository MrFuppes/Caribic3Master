# -*- coding: utf-8 -*-
"""
Created on Fri May  8 09:09:21 2020

@author: F. Obersteiner, f/obersteiner//kit/edu
"""
import os
import time
import shutil
import argparse
from pathlib import Path
from random import random
from datetime import datetime, timezone
import ntplib

# setup argparser
parser = argparse.ArgumentParser(description='Test the logger.')
parser.add_argument('--verbose', type=bool, default=False,
                    help='print stuff to the console')
parser.add_argument('iterations', metavar='N', type=int, nargs=1,
                    help='how many seconds to run')

# get the arguments that the script was called with
args = parser.parse_args()

# create function so that if --verbose=True, it prints stuff to the console
verboseprint = print if args.verbose else lambda *a, **k: None

# where to store the logfiles
# TODO: make this a parser arg?
dst = Path('D:/logtest/try_05')

# check if dir exists, create if not
if not os.path.isdir(dst):
    dst.mkdir(parents=True, exist_ok=True)
    verboseprint(f"created logfile directory: {str(dst)}")

# create dir for each day; will need this again so def a function...
def make_day_dir(dst):
    dst_day = dst/datetime.utcnow().strftime('%Y-%m-%d')
    # check if dir exists, create if not
    if not os.path.isdir(dst_day):
        dst_day.mkdir(parents=True, exist_ok=True)
        verboseprint(f"created directory for new day: {str(dst_day)}")
    return dst_day

# a function to get an independent timestamp:
def get_NTPtime(server='de.pool.ntp.org'):
    c = ntplib.NTPClient()
    try:
        r = c.request(server, version=3)
    except Exception as e:
        print(f"NTP error: {e}")
        return None, None
    return datetime.fromtimestamp(r.tx_time, tz=timezone.utc), r.offset


# inital day folder
dst_day = make_day_dir(dst)

# initial now
now = datetime.utcnow()

# create default column header
sep = '\t'
lineend = '\n'
n_randoms = 30
header = sep.join(['t_log[UTC]', 'NTP[UTC]', 'NTPoffset[s]', 'PC-NTP[s]'] + [f'v_{i}' for i in range(n_randoms)]) + lineend

# logger loop
for i in range(args.iterations[0]):

    verboseprint('iteration', i, end='\r')

    # check if date has changed; update dst_day if so
    if now.date() < datetime.utcnow().date():
        dst_day = make_day_dir(dst)

    # update now
    now = datetime.now(tz=timezone.utc)

    # defaults for reference time:
    ntptime, ntpoffset, pc_ntp = 'NaT', 'NaN', 'NaN'

    # generate the line to write to the logfile
    line = sep.join([now.isoformat(), ntptime, ntpoffset, pc_ntp] + [f'{random():.6f}' for _ in range(n_randoms)]) + lineend

    # create file for each 15 min
    fname = dst_day/now.replace(minute=divmod(now.minute, 15)[0]*15,
                                second=0,
                                microsecond=0).strftime('%Y-%m-%d_%H-%M.txt')

    # check each iteration if the current logfile exists, create if not
    if not os.path.exists(fname):
        # check available disk space; exit if less than 10%
        total, _, free = shutil.disk_usage(dst.drive)
        if free/total < 0.1:
            print("error: disk low on memory, exiting logger.")
            break
        # implicit else:
        ntptime, ntpoffset = get_NTPtime()
        if ntptime:
            pc_ntp = f"{(now-ntptime).total_seconds():.8f}"
            ntptime, ntpoffset = ntptime.isoformat(), f"{ntpoffset:.8f}"
        else:
            ntptime, ntpoffset, pc_ntp = 'NaT', 'NaN', 'NaN'
        with open(fname, 'w') as fobj:
            fobj.write(header)
            fobj.write(sep.join([now.isoformat(),
                                 ntptime,
                                 ntpoffset,
                                 pc_ntp] + [f'{random():.6f}' for _ in range(n_randoms)]) + lineend)
        verboseprint(f"created new logfile: {fname.name}")
    else:
        with open(fname, 'a') as fobj:
            fobj.write(line)

    time.sleep(1 - (time.time() % 1))

verboseprint('logger iterations done.')
