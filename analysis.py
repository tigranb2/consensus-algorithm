from sys import argv
description = str(argv[1])

def read_file(file):
    with open(file) as f:
        return f.readlines()

def extract_data(f, string_to_search): 
    results = []
    for line in f:
        if string_to_search in line:
            results.append(line)

    return results

def main():
    f = read_file("performance.txt")
    allRates = extract_data(f, "Cumulative (sent/received/total):")
    rates = allRates[len(allRates)-2].split()
    tx = rates[2]
    rx = rates[3]
    
    f = read_file("log-1.txt")
    allTimes = extract_data(f, "Time taken:")
    time = allTimes[len(allTimes)-1].split()[2]

    data = time + " " + tx + " " + rx + "\n"
    file_name = f"{description}-data.txt"
    f = open(file_name, "a")
    f.write(data)
    
main()