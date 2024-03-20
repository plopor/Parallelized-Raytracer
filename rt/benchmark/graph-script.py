import matplotlib.pyplot as plt
import sys

threads = [2, 4, 6, 8, 12]
inputs = ['nosteal&noskew', 'steal&noskew', 'steal&skew', 'nosteal&skew']

data = {'nosteal&noskew': {'threads': threads, 'speedup': [0] * 5},
        'steal&noskew': {'threads': threads, 'speedup': [0] * 5},
        'steal&skew': {'threads': threads, 'speedup': [0] * 5},
        'nosteal&skew': {'threads': threads, 'speedup': [0] * 5}}

batchNum = sys.argv[1]

with open(f'./slurm/out/{batchNum}.slurm1.stdout', 'r') as file:
    avgs = []
    lines = []
    for line in file:
        lines.append(float(line))
        if len(lines) == 5:
            avgs.append(sum(lines)/5.0)
            lines = []

seq_avgs = avgs[:1]
avgs = avgs[1:]

for i in range(4):
    thread_avgs = avgs[:5]
    avgs = avgs[5:]
    for j in range(5):
        data[inputs[i]]['speedup'][j] = seq_avgs[0] / thread_avgs[j]

for inputs, plots in data.items():
    plt.plot(plots['threads'], plots['speedup'], label=f"{inputs}")
plt.title(f"Speedup Graph")
plt.ylabel("Speedup")
plt.xlabel("Number of Threads")
plt.legend()
plt.savefig(f"{batchNum}_speedup.png")
plt.close()