# Image Conversion

Utilities around dealing with images inside of game dev. Inspired by my hate for TGA.

## Install 

```
go install ./cmd/imgconv
```

## Examples


### TGA to PNG

Convert all TGAs found to PNG, optionally resizing them along in the process

```bash
imgconv ttp -r -s 2048

# Match only on AO textures
imgconv ttp -r -s 1024 *_AO.tga
```

### Resize

Shrink all PNGs found that are above the maximum resolution provided, in this example, 2048x2048.

```
imgconv resize -r -m 2048
```

