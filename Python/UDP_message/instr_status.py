# -*- coding: utf-8 -*-
"""
Created on Fri Jul 24 13:58:12 2020

@author: va6504
"""

from enum import IntEnum, unique


@unique
class Instr_Status(IntEnum):
    """
    - Enum Class for Instrument Status.
    - Status values do not represent priority of a Status_SET but only the
          order of appearance during operation.
    """
    IN = 0
    WU = 1
    SB = 2
    MS = 3


#------------------------------------------------------------------------------


if __name__ == '__main__':
    for s in Instr_Status:
        print(s)

    print(Instr_Status.WU > Instr_Status.SB)
    # False
    print(Instr_Status.MS > Instr_Status.IN)
    # True
