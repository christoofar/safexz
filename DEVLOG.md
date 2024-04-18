# Development Log

## Apr 17 2024

- Completed encoding pathway using `*<-chan` and `*chan<-` streaming paths to `liblzma.so`
- Set up Go interface stubs.  First one to be wired up is `CompressFileWithProgress` in `compress.go`
- Set up a pre-defined matrix of compression strategies to simplify the options in `liblzmna.so`, these are:

|Strategy  | Threads  |  `liblzma.so` level |
|----------|----------|---------------------|
|CompressionSimple      | 1 |  4  |
|CompressionSimpleFast  | 1 |  2  |
|CompressionSimpleBetter| 1 |  7  |
|CompressionSimpleMax   | 1 |  9  |
|CompressionMulti*       | ½ vCPUs | 4 |
|CompressionMultiFast   | ½ vCPUs | 2 |
|CompressionMultiBetter | ½ vCPUs | 7 |
|CompressionMultiMax    | ½ vCPUs | 9 |
|CompressionFullPower         | All vCPUs | 4 |
|CompressionFullPowerFast     | All vCPUs | 2 |
|CompressionFullPowerBetter   | All vCPUs | 7 |
|CompressionFullPowerMax      | All vCPUs | 9 |

`* If only 2 cores are available, the Multi option will use both cores just the same as asking for the FullPower strategy.   On a single-core machine, all the options will default to single thread.`

I chose `CompressionMulti` as the default option for unspecified, as it balances for a multi-core environment which most people are running in containers but stays conservative on the memory so only 300MB max is expected to go into reservation when passing in a big file.   If you own a beast of a ThreadRipper you obviously are going to reach for `CompressionFullPowerMax`, while the default setting will still work fine on older Raspberry PI boards.

- To check for memory leaks I wrote a simple loop to bang away at file compression using my copy of the Debian12 installation DVD.  I was surprised that `xz` is a bit non-deterministic in compression results, which is not an expected outcome.   Pass 3 compressed to 2,640,800,208 bytes while Pass 10 compressed to 2,600,824,772.   I checked the progress report calls to make sure I wasn't missing a call to report progress at the tail end of compression.
![image](https://github.com/christoofar/safexz/assets/5059144/0ff252d7-41b6-4d9f-8afc-d781d095d6d4)

Let's decompress the result of Pass 3...
![image](https://github.com/christoofar/safexz/assets/5059144/334658d3-1134-4f14-999b-c68d32c0d598)

Looks like I have a bug somewhere in my stream mechanism.   I'll need to make sure the input bytes are all getting into the encoder and that I am not prematurely closing out of the encoding cycle.

- At least the memory leak checks are passing.
