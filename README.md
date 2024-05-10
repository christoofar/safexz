# safexz

A safe Go interface to `liblzma`, the "xz" compression library.

This is a Go package for compression in the xz / lzma format that provides a safer way to call `liblzma`  without the fear of type-safety issues and simplifies the complexity of lzma-compression options into common use cases.

## Installation (UNIX and Mac)

Install the `liblzma-dev` package from your favorite package manager.  For Debian/Ubuntu, for instance:
```bash
sudo apt install liblzma-dev
```

then:
```bash
go get -u github.com/christoofar/safexz
```

At runtime you must ensure that `liblzma.so` (UNIX) `liblzma.dll` (MSYS2) can be properly found by your operating system.
If you're on a Linux distro that has `systemd`, don't worry: you already have `liblzma`.

This package uses CGo.  By default, that means your own project must not (at least directly) disable CGo.   If you're willing to
stick around, my next GitHub project is to demonstrate how you can `CGO_ENABLE=0` but still vendor CGo support back into your program,
which is a technique you can use to make your program hot-swappable between the `musl` C standard library and the GNU C library.  I
use this technique to make tiny CGo-supporting containers!

## Installation on Windows

The easiest way to getting this done is to use [MSYS2](https://www.msys2.org).  You will be compiling your Go program with CGo
using the GCC compiler that is in MinGW64 (which comes with MSYS2).

After installing MSYS2, issue:

```
pacman -S base-devel git mingw-w64-x86_64-gcc mingw-w64-x86_64-lzma
```

Open a MinGW64 command prompt and run `nano /etc/profile` and add the approprite paths to the go compiler at the bottom of this file.  In my case I did this:
```
export GOPATH="/c/Users/chris/go"
export PATH="/go/bin:${PATH}:${GOPATH}/bin"
```

The `/go/bin` folder should be going to your `go.exe` compiler.   What I did to make this easier is symlink the default Go for Windows bin folder to `/go` like this:

```
ln -s /c/Program\ Files/Go /go
```

Restart your machine so all the paths can be found.   Open a MinGW64 command prompt and try `go version` to make sure the Go compiler is reachable.

If you are using Visual Studio Code as your IDE, it can help to make MSYS2 your default terminal program instead of Powershell.  Go to File -> Preferences -> Settings
and in `settings.json` incorporate this into your vscode settings:

```
"terminal.integrated.profiles.windows": {
    "MSYS2": {
        "path": "C:\\msys64\\usr\\bin\\bash.exe",
        "args": ["--login", "-i"],
        "env": {
            "MSYSTEM": "MINGW64",
            "CHERE_INVOKING": "1"
        }
    }
},
"terminal.integrated.defaultProfile.windows": "MSYS2"
```

Restart `vscode`.   Now you should be able to compile your Go program that uses this library.

### Eliminating the MSYS2 dependency for Windows

When you are building your program that uses `safexz` using CGo on Windows, you don't need to install MSYS2 on the target machine.
Compile your program this way:

`go build -ldflags "-extldflags \"-static\""`

This will statically link the lzma dependency from MSYS2 into your binary.

## Usage

```go
// Compressing a file
safexz.CompressFile("data/mydata.dat", "data/mydata.xz")

// Compressing a string
compressedString := safexz.CompressString("Hello World!")

// Direct-reading a compressed `xz` archive into a decompressed slice of bytes
myPicture := safexz.DecompressFileToMemory("images/monalisa.png.xz")

// Raising the compression level and asking for all cores to be
// engaged in compression
myPictureBytes := safexz.CompressFileToMemory("images/monalisa.png", CompressionFullPowerBetter)

// Forwarding an uncompressed data stream into a compressed writer stream,
// compressing as it goes
safexz.CompressStream(networkLogSource, compressedStreamWriter)  // XZWriter is the writer

// Forwarding an XZ-compressed stream into an io.Writer stream,
// decompressing as it goes
safexz.DecompressStream(compressedStream, streamToWriteTo) // XZReader is the reader

// Wrapping an existing reader behind XZReader, which will decompress the xz/lzma stream
// inside
myXZreader := safexz.XZReader.NewReader(myCompressedDataReader)  // type XZReader

// Wrapping an existing writer behind XZWriter, so the underlying writer sees and will
// save or send a valid `xz` data stream
myXZWriter := safexz.XZWriter.NewWriter(file) // type XZWriter

// Compressing a stream with XZWriter while reading its contents
	resp, err := http.Get("https://media.istockphoto.com/id/1453319272/photo/columbus-ohio-usa-skyline-on-the-scioto-river.jpg?s=2048x2048&w=is&k=20&c=tgQ4HAX-dX7A1XTanxHMrkFOg5Fpa2kW87m96JKLcUM=")

	if err != nil {
		t.Error("Error downloading image:", err)
	}
	defer resp.Body.Close()

	// Compress the image
	f, err := os.Create("test.jpg.xz")
	if err != nil {
		t.Error("Error creating compressed file:", err)
	}
	defer f.Close()
	compressedImageWriter := safexz.NewWriter(f)
	_, err = io.Copy(compressedImageWriter, resp.Body)
	if err != nil {
		t.Error("Error compressing image:", err)
	}

	compressedImageWriter.Close()
```

Full API documentation with even more functions than this can be found at [go.dev](https://pkg.go.dev/github.com/christoofar/safexz#pkg-functions)   Please review the Go test cases, as they test over 90% of all the code.

## xz backdoor

Late March 2024 [CVE-2024-3094](https://research.swtch.com/xz-timeline) was issued against the `liblzma` compression library for a [supply chain attack](https://www.crowdstrike.com/cybersecurity-101/cyberattacks/supply-chain-attacks/). That episode began with the attacker gaining maintainer role on the FOSS project in a 2 year campaign to inject a highly-complicated and stealty [backdoor](<https://en.wikipedia.org/wiki/Backdoor_(computing)>) into the software by injecting its own pre-compiled and ready-to-be-linked `.o` file into the build stream.

Since then, `systemd` and OpenSSH, the two primary projects the `xz` backdoor exploited to gain access to `sshd` have since made code changes that remove `liblzma.so` from static linking.

I'm not a firm believer that `dlopen()` is really much of a [cure](https://github.com/golang/go/issues/58548) than it is a quick excuse to ignore your supply chain. And still: dynamic linking in software is frought with other problems, one of which is security.

I am pinning my version of `lzma` to post-backdoor versions as [Lasse Collins](https://tukaani.org/contact.html) refactors out the work of Jia Tan. Similarly, there have been recent commits by Sam James [@thesamesam](https://github.com/thesamesam) done in the spirit of simplifying maintenance of `liblzma` so that the build chain is more compatible with the expectations of modern C developers. The software is healing, the backdoor in `xz` is dead, and by the time you've found this project, new versions of `liblzma` will ready for production use.

## limited but easy interface

Rather than expose the call-chain dependencies of `lzma` directly to you, a simpler Go interface is provided for your integration projects. The interface breaks down your calls into simpler one way `chan` transfers stream to `lzma`.

Direct file and `[]byte`, I/O functions for syntactic sugar are provided, so that you can avoid standing up a stream monitor and input feeder.

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

This work is based off the direct `liblzma` library stubs [published](https://github.com/jamespfennell/xz) by [@jamespfennell](https://github.com/jamespfennell/xz) under the MIT license. `liblzma` is the published work of [Lasse Collin](https://git.tukaani.org/?p=xz.git;a=blob_plain;f=AUTHORS;hb=fcbd0d199933a69713cb293cbd7409a757d854cd) and [many others](https://git.tukaani.org/?p=xz.git;a=blob;f=THANKS;h=7d2d4fe82ad8ab14161d1bacd8ef3437fe51634d;hb=fcbd0d199933a69713cb293cbd7409a757d854cd) and is published under the 0BSD software license.

`safexz` is the published work of Christopher Sawyer and is made available under the BSD3 license.
