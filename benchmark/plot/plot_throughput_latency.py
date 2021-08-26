import matplotlib.pyplot as plt
import sys
import csv

result_file_path = str(sys.argv[1]) 

with open(result_file_path, 'r') as res_file:
    results = list(csv.reader(res_file, delimiter=','))

    float_res = []

    for res in results:
        float_res.append((float(res[0]), float(res[1])))

    float_res.sort(key=lambda x: x[0])

    plt.plot([r[0] for r in float_res], [r[1] for r in float_res], '-')
    plt.show()