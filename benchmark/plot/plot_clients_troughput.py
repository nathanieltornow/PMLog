import csv
import sys
import numpy as np
import matplotlib.pyplot as plt

result_file_path = str(sys.argv[1]) 


with open(result_file_path, 'r') as res_file:
    results = list(csv.reader(res_file, delimiter=','))

    ind = np.arange(len(results))  # the x locations for the groups
    width = 0.27       # the width of the bars

    val95 = []
    val50 = []


    fig = plt.figure()
    ax = fig.add_subplot(111)

    yvals = [4, 9, 2]
    rects1 = ax.bar(ind, yvals, width, color='r')
    zvals = [1,2,3]
    rects2 = ax.bar(ind+width, zvals, width, color='g')
    kvals = [11,12,13]
    rects3 = ax.bar(ind+width*2, kvals, width, color='b')

    ax.set_ylabel('Scores')
    ax.set_xticks(ind+width)
    ax.set_xticklabels( ('2011-Jan-4', '2011-Jan-5', '2011-Jan-6') )
    ax.legend( (rects1[0], rects2[0], rects3[0]), ('y', 'z', 'k') )

    def autolabel(rects):
        for rect in rects:
            h = rect.get_height()
            ax.text(rect.get_x()+rect.get_width()/2., 1.05*h, '%d'%int(h),
                    ha='center', va='bottom')

    autolabel(rects1)
    autolabel(rects2)
    autolabel(rects3)

    plt.show()