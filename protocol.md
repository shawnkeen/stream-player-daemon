# The Music Stream Player Daemon Protocol

This document describes the command and control protocol for the Music Stream Player Daemon.

## Overview and Basics

The MSPD protocol uses line-based text for commands and data exchange.  Connections follow a client/server setup via TCP.  All communication is initiated by the client.

When the client connects to the server, the server will send a single line:

    OK MSPD version

where `version` is the protocol version.

The client may now send commands.  Each command consists of a single line, terminated by the new line character `\n`.  The server then responds with one or more lines.  The response ends with a line giving the completion status of the command.  This line contains a [return code](#return-codes) and an error message.  The line starts with either `OK` or `ERROR`, as a visual cue as well as a method to mark the last line of the response.  No other line of the server's response may start with `OK` or `ERROR`.  Example:

    ERROR 202 'my fav' is not a valid station ID

Each part of the error line is separated by a single space.  The error message may include spaces itself.  It is terminated by the new line terminating the response.

### Return Codes
Return codes are positive integers ranging from 0 to 999.

* **0-100** ID of currently playing station (0 for no playback) or the currently set volume in percent.
* **200**  general protocol error
* **201**  unknown command
* **202**  invalid argument
* **203**  incorrect syntax
* **300**  internal server error
* **301**  requested property not available
* **302**  number out of range

### Ending Connections
The protocol does not provide a separate command to end the connection.  If the client wishes to abort the connection, this should be done at TCP level, by closing the socket.  Both server and client should always be prepared that a connection is abruptly ended.  If the last command or response is incomplete, it is to be ignored.  Timeouts should be used and set to appropriate values.

### Concepts and Nomenclature
The operation of the Music Stream Player Daemon involves a number of concepts.  This section gives an overview and some definitions.

##### Stations  
The main purpose of the MSPD is to play audio streams, available over the network.  These streams may be thought of as *stations*, similar to radio stations.  The most significant property of a *station* is a *url* that points to the stream.  A *station* also has unique *ID*.  Currently, this *ID* is a number in the range between 1 and 99 and corresponds to the position of the *station* in the list of stations.  The client may present a different order of stations to the user.  The currently used backend uses two names for each station: a *display name*, as well as a *short name*.  The *short name* is usually only used in the configuration file and consists of a single word.  The *display name* can be any string of characters and is the name usually presented to the user.

##### Volume
Another important aspect is *volume* control.  A client must be able to request and change the current volume setting.  The concrete implementation of volume control on the server side is handled in the server's configuration and is not controlled by the client.  Also, there is only one volume control available over the protocol, even if the underlying hard ware / sound system provides several channels and controls.

##### Playback Status
The *status* describes what is currently going on with the MSPD.  This includes whether MSPD is actually playing a *station* or not, the currently set *volume*, and additional information about the stream.  This additional stream information consists of the last *tag* that was embedded in the stream, typically denoting the playing song.  In the future, this might include further stream properties, like audio compression rate, sample rate, etc.

*version 1.x.x*

### Versioning

The version string is structured into three parts.  Each part consists of a number and the parts are separated by a dot (`.`).

The first part is the base protocol version.  This refers to the basic mechanics of the protocol.  Should the protocol ever include, e.g., byte transfers, as opposed to line-based text, this change would be reflected in the first part of the version string.

The second part covers changes in capabilities, like available commands and corresponding responses.

The third part is reserved for any changes that do not affect compatibility.

## Commands
First, we give a brief overview of the available commands.

* [**stations**](#stations) Request a list of available stations.
* [**play**](#play) Play the given station.
* [**next**](#next) Play the next station in the list.
* [**prev**](#prev) Play the previous station in the list.
* [**stop**](#stop) Stop all playback.
* [**current**](#current) Requests information about the currently playing station.
* [**status**](#status) Give a full status report.
* [**status**](#status) Return full status report.
* [**volume**](#volume) Set the volume.
* [**help**](#help) Request a list of available commands.

### stations
A parameterless command that requests a list of all stations.

In the response, there is one line for each available station.  The line begins with the number (ID) of the station, followed by a single space.  The rest of the string is the display name of the station.

*version x.1.x*

### play
Requests the server to set the current station.

The command takes a single parameter.  The parameter is the number (ID) of the station to be played.  In the special case that the number is `0`, playback is stopped.  This is equal to the [`stop`](#stop) command.

If an unknown station number is given, the error code `301` is returned.

*version x.1.x*

### next
Requests the server to play the next station in the list of stations.  

If the current station is the last in the list, the first station will be played.  This command does not take any parameters.

Returns error code `301` if no station is playing.

*version x.1.x*

### prev
Requests the server to play the previous station in the list stations.

Works analogously to  [`next`](#next).

*version x.1.x*

### stop
Requests the server to stop all playback.

This command does not take any parameters.  No error is returned, if no station was playing when issuing the command.  In this case, playback status does not change.

*version x.1.x*

### current
Request information about the currently playing station.  If no station is playing, the result is simply the status code for successful command execution and a return code of `0`.  If a station is playing, the first line consists of the song/programme information, and the second line is the normal return code and the station ID.

*version x.1.x*

### status
Return a full status report.

Results in several lines being sent: *station*, *url*, *id*, *tag*, and *volume*.  Each line starts with the name of the information printed, followed by a colon and a single space.  The rest of the line is the corresponding information.

Returns an error, if status information could not be gathered.

### volume
Changes the playback volume.

The command takes two arguments.  The second argument is a number between `0`and `100`, denoting a percentage. The first argument is any of `set`, `inc`, `dec` or `get`.

If `set` is used, the second parameter is used as the absolute percentage of volume. E.g.,

    volume set 80

sets the volume to 80 percent.  If the specified volume exceeds 100, no error is returned.  The volume is set to 100.

To increase or decrease the volume by a given percentage, `inc` and `dec` are used, respectively. E.g.,

    volume inc 10

increases the volume by 10 percent points.  If, e.g., the volume was set to `70`, the new volume would be `80`.  If the command would increase the volume beyond 100 percent, it would not result in an error.  In such a case, the volume is set to 100 percent.  The `dec` parameter works analogously.

To request the currently set volume, `get` is used.  This subcommand requires a second argument for compatibility with the other subcommands of `volume`.  However, this seconds argument is ignored.  E.g.,

   volume get 0

returns the currently set volume in the return code.

*version x.1.x*

### help
Requests a brief overview over the supported commands.  The response consists of a line per command, containing only its name.

*version x.1.x*