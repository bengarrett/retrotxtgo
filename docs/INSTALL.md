# Retrotxt

## Install

There are [numerous download](https://github.com/bengarrett/retrotxtgo/releases/latest/) releases. Otherwise these package manager installations are available.

### Windows 
### [Scoop](https://scoop.sh/)

```ps
scoop bucket add retrotxt https://github.com/bengarrett/retrotxtgo.git
scoop install bengarrett/retrotxt
retrotxt -v
```

### macOS 
### [Homebrew](https://brew.sh/)

```sh
brew install bengarrett/retrotxt/retrotxt
retrotxt -v
```

### Linux
### DEB
```sh
# Debian DEB package for Intel
wget https://github.com/bengarrett/retrotxtgo/releases/latest/download/retrotxt_linux.deb
dpkg -i retrotxt.deb
retrotxt -v

# Debian DEB package for the Raspberry Pi & ARM
wget https://github.com/bengarrett/retrotxtgo/releases/latest/download/retrotxt_raspberry_pi.deb
dpkg -i retrotxt_raspberry_pi.deb
retrotxt -v
```

### RPM
```sh
# Redhat RPM package
wget https://github.com/bengarrett/retrotxtgo/releases/latest/download/retrotxt_linux.rpm
rpm -i retrotxt.rpm
retrotxt -v
```

### APK
```sh
# Alpine APK package
wget https://github.com/bengarrett/myip/releases/latest/download/retrotxt.apk
apk add retrotxt.apk
retrotxt -v
```
