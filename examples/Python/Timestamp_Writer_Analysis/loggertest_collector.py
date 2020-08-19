# -*- coding: utf-8 -*-
"""
Created on Fri May  8 11:30:48 2020

@author: F. Obersteiner, f/obersteiner//kit/edu
"""
from datetime import datetime
from pathlib import Path
import pandas as pd


src = Path('F:/try_05')
src = Path('D:/logtest/try_05')
for f in src.rglob('*.txt'):
    print(f)

df = pd.concat(map(lambda f: pd.read_csv(f, sep='\t', parse_dates=['t_log[UTC]'], date_parser=datetime.fromisoformat), src.rglob('*.txt')))

df['dt'] = df['t_log[UTC]'].diff().dt.total_seconds() - 1
print(f"dt min, max: {df['dt'].min()}, {df['dt'].max()}")
print(f"at {df['t_log[UTC]'].loc[df['dt'] == df['dt'].min()].iloc[0]}, "
      f"{df['t_log[UTC]'].loc[df['dt'] == df['dt'].max()].iloc[0]}")



p0 = df.plot(x='t_log[UTC]', y='dt')
p0.set_ylabel("diff vs. reference 1 second")
p0.plot(df['t_log[UTC]'], df['PC-NTP[s]'], marker='*', color='r', label="PC-NTP")
p0.plot(df['t_log[UTC]'], df['NTPoffset[s]'], marker='o', color='b', label="NTP offset")
p0.legend()


v = list(map(datetime.timestamp, df['t_log[UTC]']))
df['v_0'] = v
df['v_0'] = df['v_0'] - v[0]
p1 = df.plot(x='t_log[UTC]', y='v_0')



p2 = df['dt'].plot.box()


# with pd.HDFStore('D:/KIT/A350_Changeover/MasterComputer/logger_eval/03-200514/loggeddata.h5') as store:
#     store['data'] = df