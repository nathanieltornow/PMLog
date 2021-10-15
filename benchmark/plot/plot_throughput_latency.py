from os import read
import matplotlib.pyplot as plt
import sys
import csv

result_file_path = str(sys.argv[1]) 


with open(result_file_path, 'r') as res_file:
    results = list(csv.reader(res_file, delimiter=','))

    append_res = []
    read_res = []
    append99 = []
    read99 = []
    append95 = []
    read95 = []
    appendmedian = []
    readmedian = []

    for res in results:
        append_res.append((float(res[0]), float(res[1]) * 1000000))
        appendmedian.append((float(res[0]), float(res[2]) * 1000000))
        append99.append((float(res[0]), float(res[3]) * 1000000))
        append95.append((float(res[0]), float(res[4]) * 1000000))

        read_res.append((float(res[0]), float(res[5]) * 1000000))
        readmedian.append((float(res[0]), float(res[6]) * 1000000))
        read99.append((float(res[0]), float(res[7]) * 1000000))
        read95.append((float(res[0]), float(res[8]) * 1000000))

    plt.plot([r[0] for r in append_res], [r[1] for r in append_res], '--',label='Append')
    # plt.plot([r[0] for r in appendmedian], [r[1] for r in appendmedian], '--',label='Append Meadian')
    # plt.plot([r[0] for r in append99], [r[1] for r in append99], '--',label='Append 99%')
    # plt.plot([r[0] for r in append95], [r[1] for r in append95], '--', label='Append 95%')
    plt.plot([r[0] for r in read_res], [r[1] for r in read_res], label='Read')
    # plt.plot([r[0] for r in readmedian], [r[1] for r in readmedian], label='Read Median')
    # plt.plot([r[0] for r in read99], [r[1] for r in read99], label='Read 99%')
    # plt.plot([r[0] for r in read95], [r[1] for r in read95], label='Read 95%')
    plt.ylabel("Latency ($\mu s)$") 
    plt.xlabel("Throughput")
    plt.title("Append/Read: 50%/50%")
    plt.legend(loc='best')
    plt.show()