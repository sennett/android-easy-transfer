#!/usr/bin/expect

# spawn adb -s RFCW32969NW push  /Users/sennett/Downloads/Chuck\ S01E05\ Chuck\ Versus\ the\ Sizzling\ Shrimp\ \(1080p\ x265\ Joy\).mkv /storage/self/primary/Download/
set program [lrange $argv 0 end]
spawn {*}$program

expect {
    -re "/storage/self" {
        # Process the captured output here
        puts $expect_out(buffer);
        exp_continue
    }
}
