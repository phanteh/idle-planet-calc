# Idle Planet Calc

A simple app to help calculate a bill of materials when crafting items / alloys in Idle Planet Miner.

A Windows executable and android apk are available on the build page. Should work on linux / osx / iphone, but un-tested / un-supported.

No unit tests, as that's too much like work.

## Build

Building requires a [Go installation](https://go.dev/doc/install) and the [fyne tool](https://docs.fyne.io/started/)

### Windows

```
fyne package -os windows
```

### Android

Building for android requires [an SDK / NDK](https://developer.android.com/tools/sdkmanager) (also see [Android Studio](https://developer.android.com/studio)) to be configured

This app is not available on the app store (at the moment)
To install via `adb` you need [developer options](https://developer.android.com/studio/debug/dev-options) enabled

```
fyne package -os android
adb install idle_planet_calc.apk
```

## Known Issues

* Bonus math may not be accurate
* Missing possible bonuses

## Planned

* Additional bonuses
* Possibly removing value as it's material focused
* Ore / Alloy / Item images (maybe)
* Replace Dickbutt
* App Store (free, no ads)
