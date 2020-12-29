use crossbeam::channel::{self, Receiver};
use std::thread;

fn arrange_data(from: i64, to: i64) -> Receiver<i64> {
    let (sender, receiver) = channel::unbounded();
    let increment = if from > to { -1 } else { 1 };

    thread::spawn(move || {
        let mut value = from;
        while value != to {
            if let Err(unsent_value) = sender.send(value) {
                eprintln!("Failed to send: Channel closed: {}", unsent_value);
            }

            value += increment;
        }
    });

    receiver
}

fn merge(channels: Vec<Receiver<i64>>) -> Receiver<i64> {
    let (sender, receiver) = channel::unbounded();

    for channel in channels {
        // Each thread needs it's own copy of the sender
        let sender = sender.clone();

        // move: channel and sender need to be cloned for each thread
        thread::spawn(move || {
            for value in channel {
                if let Err(unsent_value) = sender.send(value) {
                    eprintln!("Failed to merge: Channel closed: {}", unsent_value);
                }
            }
        });
    }

    receiver
}

fn main() {
    let c1 = arrange_data(1, 10);
    let c2 = arrange_data(7, -3);
    let c3 = arrange_data(100, 110);

    // We won't use the channels in `main()` anymore, so it is safe
    // to move their ownership to `merge` to avoid a copy/clone.
    for value in merge(vec![c1, c2, c3]) {
        println!("{}", value);
    }
}
