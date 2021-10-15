import csv
import os

def calc_average(csv_file_path):
    with open(csv_file_path, 'r') as csv_file:
        res_list = list(csv.reader(csv_file, delimiter=','))
        
        sum_quad = [0, 0, 0, 0]
        for res in res_list:
            sum_quad[0] += float(res[0])
            sum_quad[1] += float(res[1])
            sum_quad[2] += float(res[2])
            sum_quad[3] += float(res[3])
        
        for i in range(len(sum_quad)):
            sum_quad[i] /= len(res_list)

        return sum_quad


def aggregate_averages(result_folder, out_csv):
    with open(out_csv, 'w', newline='') as csv_file:
        writer = csv.writer(csv_file, delimiter=',')
        for filename in os.listdir(result_folder):
            average = calc_average(os.path.join(result_folder, filename))
            writer.writerow(average)


def main():
    aggregate_averages("results/mem_log_go/95w_5r", "results/mem_log_go/95w_5r_avg.csv")


if __name__ == "__main__":
    main()