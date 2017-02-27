# dnsentinel

This is a small Go program that can handle watching and auto-updating a DNS A record.
Basically I have a few A records that point to residential IPs which can occasionally change.
I dump this compiled binary onto a box sitting behind a NAT (or my linux router that handles that),
and it will update my A records if needed.

Right now it works only with Amazon's Route53 service (that's what I use), but I tried to
write it in such a way that adding functionality for another service provider is not very difficult.

Functionality is declared via a .toml config file, an example is given (sentinel.toml)

In regards to the code, it can be a bit verbose and this whole thing is possibly a bit
over-complicated. As I was writing this I was maybe thinking a bit too much, and
experimenting with some features of Go so maybe it could be simpler but whatever.

Right now it works and I've got it to what I feel like is a decent basis point so
I'm probably not going to spend _too_ much more time with it - but will re-visit it
if I want to add some functionality.
