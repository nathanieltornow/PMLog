import csv
import os
import sys
import matplotlib.pyplot as plt

# plt.rcParams["figure.figsize"] = (8,5)

result_folders = [str(sys.argv[1]), str(sys.argv[2]), str(sys.argv[3])]


append_descr = ['5% read ratio', '50% read ratio', '95% read ratio']
read_descr = ['5% append ratio', '50% append ratio', '95% append ratio']

l = [.95, .5, .05]

num_of_clients = 8
colors = ['orange', 'purple', 'green']

results_lists = []

for folder in result_folders:
    for filename in os.listdir(folder):
        if (str(num_of_clients) + 'client') in filename:

            with open(os.path.join(folder, filename), 'r') as csv_file:
                res_list = list(csv.reader(csv_file, delimiter=','))


                data_list = []
                j = 0
                for res in res_list:
                    if j == 0 or j == len(res_list) + 1 or j % 2 == 0:
                        data_list.append((float(res[0]), float(res[1]) * 1000000, float(res[5]) * 1000000))
                    j += 1

                results_lists.append(data_list)


for i in range(len(results_lists) - 1):
    plt.plot([r[0] / 1000 for r in results_lists[i]], [r[1] for r in results_lists[i]], 'v-', label=append_descr[i], color=colors[i])
    # plt.plot([r[0] / 1000 for r in results_lists[i]], [r[2] for r in results_lists[i]], 'o--',label=read_descr[i], color=colors[i])


# plt.ylim(0, 900)

plt.ylabel("Latency ($\mu s)$") 
plt.xlabel("KRequests / sec")
# plt.title("Append/Re
# ad: 50%/50%")
plt.legend(loc="lower center", bbox_to_anchor=(0.5, -0.3),  ncol=3)
plt.tight_layout()
plt.show()