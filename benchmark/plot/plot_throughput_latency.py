import matplotlib.pyplot as plt
import sys
import csv

result_file_path = str(sys.argv[1]) 


with open(result_file_path, 'r') as res_file:
    results = list(csv.reader(res_file, delimiter=','))

    append_res = []
    read_res = []

    for res in results:
        append_res.append((float(res[0]) + float(res[2]), float(res[1])))
        read_res.append((float(res[2]), float(res[3])))

    append_res.sort(key=lambda x: x[0])
    read_res.sort(key=lambda x: x[0])

    plt.plot([r[0] for r in append_res], [r[1] for r in append_res], label='Append')
    plt.plot([r[0] for r in append_res], [r[1] for r in read_res], label='Read')
    plt.ylabel("Latency ($\mu s)$") 
    plt.xlabel("Throughput")
    plt.title("Append/Read: 95%/5%")
    plt.legend(loc='best')
    plt.show()