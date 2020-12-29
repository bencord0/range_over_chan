# Range over chan

A quick experiemnt to remind myself how chan works.

## Usage

### go

    $ go run -v ./

### rust

    $ cargo run

## Observations

The golang version was easier to write (lower cognitive load). The
basic primatives to create a `chan`, pass it around between threads
and iterate values through `range`.

An easy mistake to make is to forget about the WaitGroup and completely
break the program. I deadlocked the goroutines this way. An errgroup
offers a better API for this by scoping the wait group within a block.

Interestingly, the rust version is shorter; probably because I didn't
need to worry about thread synchronisation and the automatically closing
the channel on `Drop`.

It was however harder for me to write since I needed to look up the API
for the mpsc channel.

The easiest mistake to make was when to clone and/or move values into
the threads. This is omitted in the go version thanks to garbage
collection. In lieu of a GC, we have to tell rust which values can be
copied, changed and when (e.g. pre/post `thread::spawn`).

In contrast, for go, we don't need to worry about `Drop` semantics, but
do need to manually close the channel.

Rust can detect if the other side of a channel is closed. I've added
error handling because I could. This is part of the `mpsc::channel` API
and not a language construct.
