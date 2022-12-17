#!/usr/bin/env python

import re
from itertools import groupby
from collections import OrderedDict
from pprint import pprint
import urllib.request
from operator import attrgetter



TARGET_URL = 'https://raw.githubusercontent.com/redis/redis/unstable/src/commands.c'


def fetchCommands(url: str):
    redisCommands = urllib.request.urlopen(url)

    regex = r"{\"([a-z]+)\",\"([^\"]+)\",\"([^\"]+)\",\"([^\"]+)\""
    matches = re.finditer(regex, redisCommands.read().decode('utf-8'), re.MULTILINE)
    commands = []
    for matchNum, match in enumerate(matches, start=1):
        groups = dict(zip(["name", "description", "complexity", "version"], match.groups()))
        commands.append(groups)

    commands = sorted(commands, key=lambda c: c['name'])

    return commands

def sortByVersion(commands: list):
    commands_by_version = OrderedDict()
    for cmd in commands:
        v = cmd['version']
        if v not in commands_by_version:
            commands_by_version[v] = []

        commands_by_version[v].append(cmd)

    for key in sorted(commands_by_version):
        commands_by_version.move_to_end(key)

    return commands_by_version

def printCommands(commands: dict):
    for version, commands in commands_by_version.items():
        print( "# {0}".format(version))
        for cmd in commands:
            print("\t{name}: {description}".format(**cmd) )


if __name__ == "__main__":
    commands = fetchCommands(TARGET_URL)
    commands_by_version = sortByVersion(commands)
    printCommands(commands_by_version)