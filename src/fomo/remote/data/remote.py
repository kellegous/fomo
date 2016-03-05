#!/usr/bin/env python

import imp
import optparse
import os
import socket
import sys

def ParseAddr(str):
	host, port = str.split(':')
	return host, int(port)

class Stream(object):
	def __init__(self, sock):
		self.__sock = sock

def main():
	root = os.path.abspath(os.path.dirname(__file__))

	parser = optparse.OptionParser()
	opts, args = parser.parse_args()

	if len(args) != 2:
		sys.stderr.write('usage: %s addr task.py' % os.argv[0])
		return 1

	host, port = ParseAddr(args[0])

	remote = imp.load_source('remote', os.path.join(root, args[1]))

	sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
	sock.connect((host, port))

	remote.remote(Stream(sock))

if __name__ == '__main__':
	sys.exit(main())