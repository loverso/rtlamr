import matplotlib.pyplot as plt
import numpy as np
from scipy import signal
import math

pktlen = 38144 / 4

raw = np.memmap("sample.bin", dtype=np.uint8, offset=(17600<<1)+2048+256, mode='r')

window = raw[:pktlen].copy()
level = 127.4
iq = ((level-(window.astype(np.float64))) / level).view(np.complex128)

fig, subplots = plt.subplots(nrows=2)
fig.set_size_inches(9,9*0.6180339887)

(mag_plot, spec_plot) = subplots

mag = np.abs(iq)

filtered = np.correlate(mag, np.append(np.ones(78), -np.ones(78)))

mag_plot.plot(filtered)
mag_plot.grid(axis='both')
mag_plot.autoscale(tight=True)

quantized = np.digitize(filtered, [0])

spec_plot.plot(quantized)

spec_plot.grid(axis='both')
spec_plot.autoscale(tight=True)
spec_plot.set_ylim(-0.125, 1.125)

plt.savefig('filter.png', dpi=96, transparent=True, bbox_inches="tight")
# plt.show()
