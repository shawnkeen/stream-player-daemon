# stream-player-daemon
Simple server to control playing music streams over the network. 
This is a first attempt to write something useful in Golang.

The server uses a simple text protocol, so you can use netcat to try it out.

Playing a stream is done by running a separate command in the background. 
The main purpose is to not only control the streams but also get the currently 
played song title. The command playing the stream should extract the the song 
title tag from the stream and write it to a file. The server reads that file
and provides this information to a client.

Frankly, I have never actually used this, since I 
run [MPD](https://www.musicpd.org/) on my Pi. But it was a fun way to get acquainted 
with Golang.
