#!/usr/bin/env python

import optparse
import os
import socket
import sys

def ParseAddr(str):
	host, port = str.split(':')
	return host, int(port)

def main():
	parser = optparse.OptionParser()
	opts, args = parser.parse_args()

	if len(args) != 2:
		sys.stderr.write('usage: %s addr task.py' % os.argv[0])
		return 1

	host, port = ParseAddr(args[0])

	sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
	sock.connect((host, port))
	print('connected')

if __name__ == '__main__':
	sys.exit(main())