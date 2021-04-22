from sys import modules
from sys import argv
from os import system
from functools import partial
from time import time, sleep

from mininet.net import Mininet
from mininet.node import CPULimitedHost
from mininet.link import TCLink
from mininet.util import dumpNodeConnections
from mininet.log import setLogLevel
from mininet.cli import CLI

from topos import *
from config import conf

algorithm = str(argv[1])
num_of_nodes = int(argv[2])
delay = int(argv[3])
lossRate = float(argv[4])

def get_topology():
    privateDirs = []
    host = custom(CPULimitedHost, cpu=.003, privateDirs=privateDirs)
    try:
        topo_cls = getattr(modules[__name__], conf["topo"]["class"])
        topo_obj = topo_cls(*conf['topo']["args"], **conf['topo']["kwargs"])
        net = Mininet(topo=topo_obj, host=host, link=TCLink)
        return topo_obj, net
    except Exception as e:
        print("Specified topology not found: ", e)
        exit(0)

def main():
    def command(host, cmd, print=True):
        if print:
            hs[host - 1].cmdPrint(cmd)
        else:
            hs[host - 1].cmd(cmd)

    system('sudo mn --clean')
    setLogLevel('info')

    # reads YAML configs and creates the network
    topo, net = get_topology()
    net.start()

    hs = topo.hosts(sort=True)
    hs = [net.getNodeByName(h) for h in hs]

    command(1, "iftop -t > performance.txt &")
    
    for i in range(2, num_of_nodes + 1):
        cmd = f"./{algorithm} {i} {delay} {lossRate} > log-%s.txt &" % i
        command(i, cmd)
    
    cmd = f"./{algorithm} 1 {delay} {lossRate} | tee log-1.txt"
    command(1, cmd)

    sleep(20)
    
    # stop the network
    net.stop()


if __name__ == '__main__':
    main()
