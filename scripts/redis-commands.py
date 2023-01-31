#!/usr/bin/env python

import re
from itertools import groupby
from operator import itemgetter
from collections import OrderedDict
from pprint import pprint
import urllib.request
from operator import attrgetter
import sys
import json


TARGET_URL = 'https://raw.githubusercontent.com/redis/redis/unstable/src/commands.c'

def fetchCommands(url: str):
    redisCommands = urllib.request.urlopen(url)
    commands = redisCommands.read().decode('utf-8')

    implementedCommands = {}

    with open('./src/server/commands.json') as f:
        for cmd in json.loads(f.read()):
            implementedCommands[cmd["name"].lower()] = cmd

    regex = r"{\"([a-z_]+)\",\"([^\"]+)\",\"([^\"]+)\",\"([^\"]+)\"(.*COMMAND_GROUP_([A-Z]+)|)"
    matches = re.finditer(regex, commands, re.MULTILINE)
    commands = []
    for match in matches:
        groups = dict(zip(["name", "description", "complexity", "version", None, "command_group"], match.groups()))
        del(groups[None])

        if groups["name"] in implementedCommands:
            groups["implemented"] = implementedCommands[groups["name"]]

        commands.append(groups)

    commands = sorted(commands, key=lambda c: c['name'])

    return commands

def filterByVersion(commands: list, version: str):
    if not version:
        return commands
    v = version.split('.')
    return filter(lambda c: c["version"].split('.') <= v, commands)

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

def groupByCommandGroup(commands: list):
    group = OrderedDict()
    for cmd in commands:
        v = cmd['command_group']
        if v not in group:
            group[v] = []
        
        group[v].append(cmd)
    return group

def printCommands(commands: dict):
    status_representation = {
        "won't-do": "🚫",
        "implemented": "✅",
        "partially-implemented":  "🏗️ "
    }

    print("This is an automatically generated file! DO NOT EDIT. Run \"rm commands.md && make commands.md\" to regenerate")
    print()
    for version, commands in commands_by_version.items():
        print( "# {0}".format(version))
        for group, cmds in groupByCommandGroup(commands).items():
            print("\t" + group)
            for cmd in cmds:  
                st = status_representation.get(cmd.get("implemented", {}).get("status"), "  ")
                print("\t\t{status} {name}: {description}".format(status=st, **cmd) )


if __name__ == "__main__":
    version = sys.argv[1] if len(sys.argv) >= 2 else ""
    commands = fetchCommands(TARGET_URL)
    commands = filterByVersion(commands, version)
    commands_by_version = sortByVersion(commands)
    printCommands(commands_by_version)