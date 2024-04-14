package main

type Return int // Return values for LZMA functions
type Action int // Actions for lzma_code

const (
	Ok               Return = 0
	StreamEnd        Return = 1
	NoCheck          Return = 2
	UnsupportedCheck Return = 3
	GetCheck         Return = 4
	MemoryError      Return = 5
	MemoryLimitError Return = 6
	FormatError      Return = 7
	OptionsError     Return = 8
	DataError        Return = 9
	BufferError      Return = 10
	ProgrammingError Return = 11
	SeekNeeded       Return = 12
)

func (ret Return) String() string {
	switch ret {
	case Ok:
		return "Ok"
	case StreamEnd:
		return "StreamEnd"
	case NoCheck:
		return "NoCheck"
	case UnsupportedCheck:
		return "UnsupportedCheck"
	case GetCheck:
		return "GetCheck"
	case MemoryError:
		return "MemoryError"
	case MemoryLimitError:
		return "MemoryLimitError"
	case FormatError:
		return "FormatError"
	case OptionsError:
		return "OptionsError"
	case DataError:
		return "DataError"
	case BufferError:
		return "BufferError"
	case ProgrammingError:
		return "ProgrammingError"
	case SeekNeeded:
		return "SeekNeeded"
	}
	return "Unknown"
}

const (
	Run         Action = 0
	SyncFlush   Action = 1
	FullFlush   Action = 2
	Finish      Action = 3
	FullBarrier Action = 4
)

func (action Action) String() string {
	switch action {
	case Run:
		return "Run"
	case SyncFlush:
		return "SyncFlush"
	case FullFlush:
		return "FullFlush"
	case Finish:
		return "Finish"
	case FullBarrier:
		return "FullBarrier"
	}
	return "Unknown"
}
