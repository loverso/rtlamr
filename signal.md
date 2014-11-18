---
layout: page
title: Signal Processing
index: 3
---

{% include image.html path="/assets/signal_flow.png" caption="Overview of signal flow." %}

## Data Acquisition
***

There are two methods to get data from an rtl-sdr dongle, directly with librtlsdr and via tcp with the `rtl_tcp` spectrum server. Using librtlsdr requires the use of cgo which prevents cross-compilation; `rtl_tcp` is used instead. This has the added benefit of allowing the receiver to be somewhere other than the system running `rtlamr`.

## Demodulation
***

The ERT protocol is an on-off keyed manchester-coded signal transmitted at bit-rate of 32.768kbps. On-off keying is a type of amplitude shift keying. Individual symbols are represented as either a carrier of fixed amplitude or no transmission at all.

{% include image.html path="/assets/magnitude.png" caption="<strong>Top:</strong> Inphase component of received signal. <strong>Bottom:</strong> Magnitude of complex signal. <strong>Note:</strong> Signal is truncated show detail." %}

The signal is made up of interleaved in-phase and quadrature samples, 8-bits per component. The amplitude of each sample is:

$$\vert z\vert = \sqrt{\Re(z)^2 + \Im(z)^2}$$

To meet performance requirements the magnitude computation has two implementations. The first uses a pre-computed lookup table which maps all possible 8-bit values to their floating-point squares. Calculating the magnitude using the lookup table then only involves two lookups, one addition and one square-root.

The second implementation is an approximation known as Alpha-Max Beta-Min whose efficiency comes from eliminating the square root operation:

$$
\begin{eqnarray*}
	\alpha &=& \frac{2\cos\frac{\pi}{8}}{1+cos\frac{\pi}{8}} \qquad \beta = \frac{2\sin\frac{\pi}{8}}{1+cos\frac{\pi}{8}} \\ \\
	\vert z\vert &\approx& \alpha\cdot\max(\Re(z),\Im(z)) + \beta\cdot\min(\Re(z),\Im(z))
\end{eqnarray*}
$$