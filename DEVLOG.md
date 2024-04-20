# Development Log

## Why Aren't You Supporting Multi-threaded decompression?

If you're looking at the functions I put in `decompression.go`, I've skipped on multi-threaded decompression.   In the decompression scenario it (yet again) comes down to the working storage in RAM that will determine the decompression speed and this time output I/O will play a bigger factor as bytes in the working area need to be cleared away to make room for the compressed data coming in.

For `safexz` I have set a hard maximum area of `250<<20`, or 250MB of decompression working storage.  For the original Raspberry PI 1 which has 512MB of working storage, `128MB` will be selected instead.  This is sufficient without causing too much headache.

On the off-chance that you are using TinyGo and working on ridiculously constrained machine, I have come up with this solve, which isn't one of my prouder moments:

![image](https://github.com/christoofar/safexz/assets/5059144/af67591d-9981-463d-86ae-2547bdf7755c)


On normal VMs `liblzma` will get a very ample working area and you'll see nice I/O coming out of the streamer.  If you're trying to re-create a Commodore128 in TinyGo you'll at least get... working storage.  Of 64KB.


## Apr 20 2024

Made a quick prover to see how the different canned compression strategies pan out.  This is over a puny 6-core dev VM so it will naturally be slow, thus a better comparison of the multithreaded behaviors of `liblzma.so` can be spotted.  Compressing the King James Bible yields these results (The results as they come out are unsorted, so I've resorted them here):

```
christoofar@pop-os:~/src/safexz/cmd/speedtest$ ./speedtest -i ../../test/canterbury-corpus/large/bible.txt 
Starting compression with CompressionSimpleFast...
Starting compression with CompressionSimple...
Starting compression with CompressionSimpleBetter...
Starting compression with CompressionSimpleMax...
Compression complete.  Moving on to CompressionMultiFast...
Starting compression with CompressionMulti...
Starting compression with CompressionMultiBetter...
Starting compression with CompressionMultiMax...
Compression complete.  Moving on to CompressionFullPowerFast...
Starting compression with CompressionFullPower...
Starting compression with CompressionFullPowerBetter...
Starting compression with CompressionFullPowerMax...
Compression complete.
Compression Results:
Algorithm                      :                 Time : Size
---------                      :                 ---- : ----
CompressionSimpleFast          :         476.479748ms : 1085880 bytes
CompressionSimple              :         1.722374928s : 944900 bytes
CompressionSimpleBetter        :         1.791940736s : 885192 bytes
CompressionSimpleMax           :         1.822931302s : 885192 bytes
CompressionMultiFast           :         2.284920555s : 1085880 bytes
CompressionMulti               :         1.402525781s : 944900 bytes
CompressionMultiBetter         :          1.84167011s : 885192 bytes
CompressionMultiMax            :          1.80917974s : 885192 bytes
CompressionFullPowerFast       :         2.264714487s : 1085880 bytes
CompressionFullPower           :         1.362964771s : 944900 bytes
CompressionFullPowerBetter     :         1.818037158s : 885192 bytes
CompressionFullPowerMax        :         1.819513768s : 885192 bytes
```

So there's no savings to be had at all going with the `Max` option when it comes to raw text, as that just burns CPU.  At least when it comes to the common text case sizes at least.   The most interesting result is the `CompressionSimpleFast` option beat everything else on time.   When you think about it, it makes sense.

Single-stream compression algorithms don't lend themselves well to multiprocessing because of a basic way multiprocessing on single-tasks works called `segmentation`.   Segmentation is when you break up an unworked dataset like this:

```
+--------------------+-------------------+-------------------+
+    Data Part 1     +    Data Part 2    +     Data Part 3   +
++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
```
Then you would assign coroutines or whole OS process threads to work on all three parts, then stitch the results back into one post-process dataset.

You can segment the data again into sub-sub segments like this if you have lots of resources:
```
+--------------------+-------------------+-------------------+
+  Block 1 + Block 2 + Block 3 + Block 4 + Block 5 + Block 6 +
++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
```
Segmentation eliminates the problem of shared cooperative memory and the need to lock sections of memory to prevent collisions, but now you've created areas of the dataset where repetitive data that crosses a block boundary might escape the attention of the compression routine each thread is running:
```
+--------------------+-------------------+-------------------+
+ ----> [ Repetitive data ] Data Part 2  +     Data Part 3   +
++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
```
The simplist answer to this conundrum is to load more of the data to be compressed into working storage (RAM), so that expensive cleanup passes and/or multi-segment sized sweeper passes can examine the finished areas of the compression stream in the working storage and correct areas that escaped attention of the standard compression scheme working at the smallest segment size.

Whether running these cleanup passes is worth it entirely depends on the nature of the data underneath.  Large datasets with huge repeating page blocks will certainly pass across the compression window if it's large enough, but the compressed result will likely compress much further if the repeating patterns are seen again much later in the datastream, which would get pickup up by a pass using a larger inspection size.

And this is what the higher settings of `liblzma.so` essentially, more-or-less, do.   More threads don't make compression necessarily faster, but more RAM certainly will.  It does in a dramatic way.  RAM speed and the amount of it you have by and large will dominate the time spent compressing, less so on the underlying I/O media speed or the number of cores you throw at it.

So for big data you're best off giving a few extra gigs of RAM to `safexz` and run it with `CompressionFullPowerBetter`.  As you can see, an 8GB VM is not going to handle compressing `debian.iso` that well when you have the compression level cranked up.  This has VSCode running in debug mode sucking up about 4.5GB of RAM and `safexz` in the `CompressionFullPowerMax` setting has pulled down an extra 1.5GB of RAM (the `RES` column) and sent the Linux swap system into overdrive.
![image](https://github.com/christoofar/safexz/assets/5059144/9dbdff68-2496-4519-82a5-246fe4a9832f)

At the other end of the spectrum are small machines, like the [Pi Zero](https://www.canakit.com/raspberry-pi-zero.html) or something even smaller than that.   The working storage for a simple Go program using `CompressionSimpleFast` is only 43K:

![image](https://github.com/christoofar/safexz/assets/5059144/9a366287-9f27-4699-89bc-f812f82f2b4c)

`liblzma.so` gives you tremendous flexibility where you can compress data from the smallest computers to something as monsterous as an IBM z/Series mainframe.





## Apr 19 2024

So, I finally figured out why I was getting such non-deterministic results.  This is what I did to fix it. 
![image](https://github.com/christoofar/safexz/assets/5059144/f4136c9c-3742-4262-9a26-03eacd338ac0)

`readbuf` is dirty after the read.  I thought that this was a simple reusable type but it's got some distinct behavior when used with `file.Read()` because of the syscall that occurs.   To fix this, I pull out what was read into a clean byte slice and send that into the `chan` for processing.

That results in a clean byte-for-byte accounted-for `diff`:

![image](https://github.com/christoofar/safexz/assets/5059144/98292003-274d-4375-bb88-6e413aa6726a)

Now I can move on to the multi-threaded decompression.

## Apr 18 2024

So, all options using Simple (single-threading) produce a good result.  The multi-threading ones do not.  Probably another signal that I need to pick up from `liblzma.so` to know that all the threads underneath in the innermost goroutine have completed.  I'll be hunting around for some multithread examples in C to see if the calling pattern is different.

*Eureka!* I found what was wrong with the multithreaded compress options.  Turns out that I goofed and did not check `.avail_in` on the stream before pushing data.  Apparently this issue doesn't turn up frequently enough for me to see it in single-threaded mode but it will come up in multi-threaded.   `.avail_in` tells you that there are bytes waiting to be drained into the memory area where lzma is working.   You can try to fill up to `MAX_BUFFER_SIZE` but it's easier to wait for it to clear to zero in a cycle and on the next cycle it's likely for the drain to occur.  If you just set all the bytes for the cycle then the bytes waiting to be drained will be destroyed.

Now, testing again with compressing the Debian12 DVD I still get two different compression result sizes.
![image](https://github.com/christoofar/safexz/assets/5059144/ef8fdbec-957e-4d2c-8023-c338f656b996)

But this time the byte counts matches the origin.
![image](https://github.com/christoofar/safexz/assets/5059144/4120ddad-c460-4481-8841-a3909c0d1b82)

Is LZMA still not deterministic?  That's wild.  Can't be right.  Let's check that all the bytes exactly match.
![image](https://github.com/christoofar/safexz/assets/5059144/cc43aa82-99b2-47d5-ab38-31d8063906da)

You know... I think maybe `MAX_BUFFER_SIZE` really should be down in the low count, like `1024`.  Anything higher might break some limits or plausables. 

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
