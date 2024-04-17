package safexz

// CompressionStrategy is an enum type for the compression strategy to use for encoding.
type CompressionStrategy int

// CompressionSimple uses a single thread with not that much demand on memory.
const CompressionSimple CompressionStrategy = 0

// CompressionSimpleFast uses a single thread with a faster compression speed at the
// expense of a larger result.
const CompressionSimpleFast CompressionStrategy = 1

// CompressionSimpleBetter uses a single thread with a better compression ratio than "Better",
// but uses more memory.
const CompressionSimpleBetter CompressionStrategy = 2

// CompressionSimpleMax uses a single thread with the best compression ratio, but uses the most
// memory.  In memory-constrained environments this will spill your process into the swap area
// and slow down processing greatly.
const CompressionSimpleMax CompressionStrategy = 3

// CompressionMulti (default) uses multiple threads (half the number of available cores) to compress data.
// This is the best option in terms of speed and memory consumption without crowding out other processing
// on the system.
const CompressionMulti CompressionStrategy = 4

// CompressionMultiFast uses multiple threads (half the number of available cores) to compress data but
// with a faster compression speed at the expense of a larger result.
const CompressionMultiFast CompressionStrategy = 5

// CompressionMultiBetter uses multiple threads (half the number of available cores) to compress data but
// with a better compression ratio than the default, but uses more memory.
const CompressionMultiBetter CompressionStrategy = 6

// CompressionMultiMax uses multiple threads (half the number of available cores) to compress data but
// with its higher compression ratio comes the most memory usage.  This option is not recommended for
// systems with constrained memory resources.
const CompressionMultiMax CompressionStrategy = 7

// CompressionFullPower uses all available cores to compress data.  This is a faster option but will crowd
// out other processing on the system.
const CompressionFullPower CompressionStrategy = 8

// CompressionFullPowerFast uses all available cores to compress data with a faster compression speed at
// the expense of a larger result, but less memory usage.  This is the fastest option for large files.
const CompressionFullPowerFast CompressionStrategy = 9

// CompressionFullPowerBetter uses all available cores to compress data with a better compression ratio
// than the default, but uses more memory.  This option is not recommended for systems with constrained
// memory resources.
const CompressionFullPowerBetter CompressionStrategy = 10

// CompressionFullPowerMax uses all available cores to compress data with the best compression ratio, but
// uses a large amount of memory. This option is not recommended for systems with constrained memory
// resources, and on large files will crowd out other processing on the system.
const CompressionFullPowerMax CompressionStrategy = 11
