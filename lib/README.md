# lib/

Native build tree for `golibjpeg`, aligned with [pylibjpeg-libjpeg](https://github.com/pydicom/pylibjpeg-libjpeg).

```
lib/
├── libjpeg/      # git submodule → thorfdbg/libjpeg
├── interface/    # decode + streamhook (from pylibjpeg-libjpeg/lib/interface)
├── capi/         # C ABI exported for purego (golibjpeg.h)
└── CMakeLists.txt
```

## Build

CI (`build-libs.yml`) builds all platform libraries on push to `main` when `lib/**` changes,
then commits artifacts to `native/libs/`.

Optional local build (debugging only):

```bash
git submodule update --init --recursive lib/libjpeg
make build-native
```

## Reference

Implementation and tests are tracked against `ref/pylibjpeg-libjpeg/` (read-only submodule).
