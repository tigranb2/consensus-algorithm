# consensus-algorithm
Compilation of ACR, BACR, MP3, and BMP3 algorithms.

## Installation
Execute line by line:
```shell
sudo su 
cd /usr/local/go/src # install directory 
git clone https://github.com/tigranb2/consensus-algorithm.git
cd consensus-algorithm && . setup.sh # installs dependencies (Go, Python, Mininet)
```

## Usage
### Configuration:
To configure the algorithm, edit config.yaml and config.json.  
    
Open config.yaml to edit the number of Mininet nodes and loss rate:
```shell
nano config.yaml
```
You will see:
```yaml
network:
  topo:
    class: SingleSwitchTopo
    args:
      - 9 # set number of Mininet hosts. Should be the same as number of nodes

  link:
    bw: null    # e.g., 100 (Mb/s)
    delay: 10  # e.g., 1s, 1ms
    jitter: 3  # e.g., 1s, 1ms
    loss: null    # set loss rate (%) (set it to null for MP3 and BMP3!)
...
```
**_For the MP3 and BMP3 algorithms, ALWAYS keep loss rate as "null"._**
      
      
Open config.json to set node count and fault count:
```shell
nano config.json
```
The first number in the first line is the node count, the second number of the first line is the fault count:
```json
100 49 //# of nodes, # of faulty nodes
1 10.0.0.1 8001
2 10.0.0.2 8002
...
```

### Running:
Execute:
```shell
. run.sh {algorithm} {num_of_nodes} {broadcast_frequency OR timeout} {loss_rate} {description}
# algorithm is the name of the algorithm you wish to run (e.g. ACR, BMP3, etc.)
# num_of_nodes is the number of nodes (same as number in config.json and config.yaml)
# the third argument is the broadcast frequency for ACR & BACR, and it is the timeout for MP3 & BMP3
# loss_rate is the percent of messages lost (1.5 would be 1.5%)
# description should describe the test (e.g. ACR_scalability-100). There should be no spaces. The file where perfromance data is stored is {description}-data.txt
```

After running, the collected data will be printed on screen and saved to a file (see above). Copy the 5 lines for test.
