# consensus-algorithm
Compilation of ACR, BACR, MP3, and BMP3 algorithms.

## Installation
Execute line by line:
```shell
sudo su 
mkdir /usr/local/go && mkdir /usr/local/go/src
cd /usr/local/go/src # install directory 
git clone https://github.com/tigranb2/consensus-algorithm.git
cd consensus-algorithm && . setup.sh # installs dependencies (Go, Python, Mininet)
```

## Usage
### Configuration:
To configure the algorithm, edit config.yaml and config.json.  
    
Open config.yaml to edit the number of Mininet hosts and loss rate:
```shell
nano config.yaml
```
You will see:
```yaml
network:
  topo:
    class: SingleSwitchTopo
    args:
      - 100 # set number of Mininet hosts. Should be the same as number of nodes

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
. run.sh {algorithm} {num_of_nodes} {broadcast_frequency OR timeout} {loss_rate} {description} # ex: . run.sh ACR 300 10 0.5 ACR_scalability_300
# algorithm is the name of the algorithm you wish to run (e.g. ACR, BMP3, etc.)
# num_of_nodes is the number of nodes (same as number in config.json and config.yaml)
# the third argument is the broadcast frequency for ACR & BACR, but the timeout for MP3 & BMP3
# loss_rate is the percent of messages lost (1.5 would be 1.5%)
# description should describe the test (e.g. ACR_scalability_100). There should be no spaces. The file where perfromance data is stored is {description}-data.txt
```

After running, the collected data will be printed on screen and saved to a file (see above). Copy the 5 lines for each test.    
Example (your numbers will be much larger... this is only 10 nodes):
```shell
...
*** Done # copy the next 5 lines
219.939342ms 89.2KB 101KB 
184.599021ms 102KB 105KB
237.078078ms 94.3KB 105KB
210.949478ms 99.1KB 107KB
206.415603ms 97.9KB 106KB
```
