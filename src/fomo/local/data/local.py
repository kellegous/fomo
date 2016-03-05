#!/usr/bin/env python

import imp
import optparse
import os
import socket
import sys

def check_module(mod):
	# a local is not required. If you don't specify a local, we will simply discard all streamed data. If you try
	# to recv in a remote that has not corresponding local, what will we do?
	return hasattr(mod, 'remote')

class Stream(object):
	def __init__(self, sock):
		self.__sock = sock

def main():
	parser = optparse.OptionParser()
	opts, args = parser.parse_args()

	if len(args) != 2:
		sys.stderr.write('usage\n')
		return 1

	remote = imp.load_source('remote', args[1])

	sock = socket.socket(socket.AF_UNIX, socket.SOCK_STREAM)
	sock.connect(args[0])

	remote.local(Stream(sock))


if __name__ == '__main__':
	sys.exit(main())