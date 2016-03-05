import socket

def remote(stream):
	print('remote: %s' % socket.gethostname())
	# stream.send(fomo.File('/etc/passwd'))

	# for i in range(100):
	# 	stream.send(i)

def local(stream):
	print('local: %s' % socket.gethostname())
	# for host, obj in stream:
	# 	obj.write_to(host)