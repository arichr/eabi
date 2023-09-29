<p align="center">
    <h1 align="center">Eabi (Encode As Binary)</h1>
    <p align="center">A lightweight alternative to MessagePack.</p>
</p>

Eabi aims to be a faster and lighter alternative to MessagePack, while being completely independent from it.

## Getting started!

> **Warning:**
> Eabi is still under active development and does not yet provide stable releases.
> Consider contributing to make it more user friendly! <img src="https://slackmojis.com/emojis/30273-meow_photo/image/1680406618/meow_photo.png" width="20" alt="" valign="bottom">

Add Eabi as a dependency to your project:
```bash
go get github.com/arichr/eabi
```

### Compiling examples
If you are on Linux, you can use `sh build.sh`.

On other operating systems, the result is no different from building each example yourself:
```sh
go build -o build/eabi cmd/eabi/eabi.go
# and so on...
```

<!--
## Roadmap

For the first release:
* [ ] Support `Map`
* [ ] Support `Uint64`
* [ ] Support `Int16`
* [ ] Support `Uint16`
* [ ] Support `Int8`
* [ ] Support `Uint8`
* [ ] Support `Float64`
* [ ] Support `String`
* [ ] Support slices and arrays
  * [ ] `[]byte` and `[...]byte` as special cases
* [ ] Support structures that implements `Marshaler`

Planned:
* [ ] Implement `typedef`
-->
