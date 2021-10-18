from os import read
import matplotlib.pyplot as plt
import sys
import csv
plt.rcParams["figure.figsize"] = (5,3)

colors = ['blue', 'green', 'red']
l = 0

for file in sys.argv[1:]:
    with open(str(file), 'r') as res_file:
        results = list(csv.reader(res_file, delimiter=','))

        append_res = []
        read_res = []
        append99 = []
        read99 = []
        append95 = []
        read95 = []
        appendmedian = []
        readmedian = []
        i = 0
        for res in results:
            if i % 2 == 0 or i == 0 or i == len(res) - 1:
                append_res.append((float(res[0]), float(res[1]) * 1000000))
                appendmedian.append((float(res[0]), float(res[2]) * 1000000))
                append99.append((float(res[0]), float(res[3]) * 1000000))
                append95.append((float(res[0]), float(res[4]) * 1000000))

                read_res.append((float(res[0]), float(res[5]) * 1000000))
                readmedian.append((float(res[0]), float(res[6]) * 1000000))
                read99.append((float(res[0]), float(res[7]) * 1000000))
                read95.append((float(res[0]), float(res[8]) * 1000000))
            i+=1

        plt.plot([r[0]*3 / 1000 for r in append_res], [r[1] for r in append_res],'o-', markersize=4, label='Append - Average', color=colors[l])
        # plt.plot([r[0] / 1000 for r in appendmedian], [r[1] for r in appendmedian],'o-', markersize=4, label='Append - Median', color='blue')
        # plt.plot([r[0] / 1000 for r in append99], [r[1] for r in append99], 'o--', markersize=4, label='Append 99%', color='blue')
        # plt.plot([r[0] / 1000 for r in append95], [r[1] for r in append95], 'o--', markeredgecolor='black', label='Append - 95%')
        plt.plot([r[0]*3 / 1000 for r in read_res], [r[1] for r in read_res], 'v--', markersize=4, label='Read - Average', color=colors[l])
        # plt.plot([r[0] / 1000 for r in readmedian], [r[1] for r in readmedian], 'v-', markersize=4, label='Read - Median', color='green')
        # plt.plot([r[0] / 1000 for r in read99], [r[1] for r in read99], 'v--', markersize=4, label='Read - 99%', color='green')
        # plt.plot([r[0] / 1000 for r in read95], [r[1] for r in read95], label='Read - 95%')
    l+=1

plt.ylabel("Latency ($\mu s)$") 
plt.xlabel("KRequests / sec")
plt.tight_layout()
# plt.ylim(0, 1500)
plt.legend(loc='best')
plt.show()