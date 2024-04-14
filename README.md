# safexz
(in development)  A safe Go interface to `liblzma`, the "xz" compression library.

This is a Go package for compression in the xz / lzma format that provides a safer way to call `liblzma` for common use cases without the fear of type-safety issues and utilizes Go's goroutines to protect your project from [unforseen control hijacks](https://research.swtch.com/xz-timeline) from the 5.6.0 and 5.6.1 versions of `liblzma`

## compatibility with older platforms
`liblzma` is still one of the best compression algorithms for compacting data.  It has backportage to even OS/2 Warp from the mid-1990s and there is an implementation for nearly every 32-bit-or-better processor.

As Go itself is much newer, I should be able to maintain compatibility back to `$TODO: backportage tests :-)`.

## xz backdoor
Late March 2024 [CVE-2024-3094](https://research.swtch.com/xz-timeline) was issued against the `liblzma` compression library for a [supply chain attack](https://www.crowdstrike.com/cybersecurity-101/cyberattacks/supply-chain-attacks/).  That episode began with the attacker gaining maintainer role on the FOSS project in a 2 year campaign in inject a highly-complicated and stealty [backdoor](https://en.wikipedia.org/wiki/Backdoor_(computing)) into the software by injecting its own pre-compiled and ready-to-be-linked `.o` file into the build stream.

Since then, `systemd` and OpenSSH, the two primary projects the `xz` backdoor exploited to gain access to `sshd` have since made code changes that remove `liblzma.so` from static linking.

I'm not a firm believer that `dlopen()` is really much of a [cure](https://github.com/golang/go/issues/58548) than it is a quick excuse to ignore your supply chain.  And still: dynamic linking in software is frought with other problems, one of which is security.

I am pinning my version of `lzma` to versions as [Lasse Collins](https://tukaani.org/contact.html) refactors out the work of Jia Tan.  Similarly, there have been recent commits by Sam James [@thesamesam](https://github.com/thesamesam).   The software is healing, the backdoor in `xz` is dead, and by the time you've found this project it's ready for production use.

## limited but easy interface
Rather than expose the call-chain dependencies of `lzma` directly to you, a simpler Go interface is provided for your integration projects.  The interface breaks down your calls into simpler `chan` transfers to `lzma`.

## safexz wraps the lzma library in isolated go routines
This project's goal is to abide by a chaos/containment theory which believes _unless you wrote the C code, tread carefully_.

To increase the difficulty of a user with a suspicious C library gaining control over your own go program that uses a C library, `safexz`:

- Hides C call access and data structures as internal packages
- No call is made to a C library without boxing it inside 2 layers of goroutines
- All communication of data to and from the C library's functions must take place from a separate goroutine and data may only pass via channels using packed types (fields without pointers) wherever possible.
- Go pointers into `lzma` are not possible
- No one using this library would ever have to manage `liblzma`'s internal state, nor will `liblzma` be able to "see" your data types of your own code except for the []byte stream of data to be compressed/decompressed.

More detail about this technique is [over here](https://gist.github.com/christoofar/880b4bcf3018f4681bb71bfdf1c16a6a).

## credits
This work is based off the direct `liblzma` library stubs [published](https://github.com/jamespfennell/xz) by [@jamespfennell](https://github.com/jamespfennell/xz) under the MIT license.  `liblzma` is the published work of [Lasse Collin](https://git.tukaani.org/?p=xz.git;a=blob_plain;f=AUTHORS;hb=fcbd0d199933a69713cb293cbd7409a757d854cd) and [many others](https://git.tukaani.org/?p=xz.git;a=blob;f=THANKS;h=7d2d4fe82ad8ab14161d1bacd8ef3437fe51634d;hb=fcbd0d199933a69713cb293cbd7409a757d854cd) and is published under the 0BSD software license.

`safexz` is the published work of Christopher Sawyer and is made available under the MIT license.