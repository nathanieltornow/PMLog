import csv
import os
import sys
import matplotlib.pyplot as plt

result_folder = str(sys.argv[1]) 

num_of_clients = [8, 16]
colors = ['red', 'blue', 'green', 'orange']

result_dict = {}


for num in num_of_clients:

    for filename in os.listdir(result_folder):
        if (str(num) + 'client') in filename:

            with open(os.path.join(result_folder, filename), 'r') as csv_file:
                res_list = list(csv.reader(csv_file, delimiter=','))


                data_list = []
                j = 0
                for res in res_list:
                    if j == 0 or j == len(res_list) + 1 or j % 1 == 0:
                        data_list.append((float(res[0]), float(res[1]) * 1000000, float(res[5]) * 1000000))
                    j += 1

                result_dict[num] = data_list

i = 0
for num, res in result_dict.items():
    plt.plot([r[0] for r in res], [r[1] for r in res], '.--',label=('Append, ' + str(num) + ' clients'), color=colors[i])
    plt.plot([r[0] * .95 for r in res], [r[2] for r in res], '.-',label=('Read, ' + str(num) + ' clients'), color=colors[i])
    i += 1

plt.ylabel("Latency ($\mu s)$") 
plt.xlabel("Throughput")
# plt.title("Append/Read: 50%/50%")
plt.legend(loc='best')
plt.show()