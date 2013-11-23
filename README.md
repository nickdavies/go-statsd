go-statsd
=========

# This port is very incomplete

    See todo

# Known Incompatibilities

 * The config file in statsd is eval'd as javasript. As such this functionalty will break if you depend on anything other than it being JSON

# TODO

 * Build backends interface
 * Write the rest of process.go
 * work out keyflushInterval variable
 * proper logging
 * real error handling
 * interal stats
 * management server
 * config
 * tests for all things
