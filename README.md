# Astrotyper

My submission for the Gamedev compo at [Birdie 16](https://www.birdie.org/en/).

## Synopsis

An endless swarm of asteroids are heading for earth and will destroy all
life as we know it. You have been sent out into space to destroy these
asteroids before they hit earth.

Target and destroy asteroids by typing the words next to them. When the
earth's percentage has reached 0% the game is over. For each level the
number of asteroids and their speed is increased and on top of that the
word you have to type to destroy an asteroid gets longer.

## Building

This game is written in [Go](https://golang.org) with
[bindings for SDL2](https://github.com/veandco/go-sdl2).

```
export GOPATH=`pwd`/go
go get github.com/veandco/go-sdl2/sdl
go build
```

This will create a binary called `astrotyper`.

## License

- All source code is licensed under [GPLv3](https://www.gnu.org/licenses/gpl-3.0.en.html).
- All resources in `resources/` are licensed under [CC0](http://creativecommons.org/publicdomain/zero/1.0/).
- The font in `resources/font/` is licensed under [OFL-1.1](https://opensource.org/licenses/OFL-1.1).
