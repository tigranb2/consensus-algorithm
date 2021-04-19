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

    return results[len(results)-1]

def main():
    f = read_file("performance.txt")
    rates = extract_data(f, "Cumulative (sent/received/total):").split()
    tx = rates[2]
    rx = rates[3]
    
    f = read_file("log-1.txt")
    time = extract_data(f, "Time taken:").split()[2]

    data = time + " " + tx + " " + rx + "\n"
    file_name = f"{description}-data.txt"
    f = open(file_name, "a")
    f.write(data)
    
main()