# Retrotxt

## Install

There are [numerous download](https://github.com/bengarrett/retrotxtgo/releases/latest/) releases. Otherwise these installation options are available.

### Windows 
### [Scoop](https://scoop.sh/)

```sh
# NOT CURRENTLY WORKING
scoop bucket add retrotxt https://github.com/bengarrett/retrotxtgo.git
scoop install bengarrett/retrotxt
retrotxt -v
```

### macOS 

Unfortunately, newer macOS versions do not permit the running of unsigned terminal applications out of the box. But there is a workaround.

In System Settings, Privacy & Security, Security, toggle **Allow applications downloaded from App store and identified developers**.

Use Finder to extract the download for macOS, [`retrotxtgo_darwin_all.tar.gz`](https://github.com/bengarrett/retrotxtgo/releases/latest/download/retrotxtgo_darwin_all.tar.gz).

Use Finder to select the extracted `retrotxtgo` binary. <kbd>^</kbd> control-click the binary and choose _Open_. macOS will ask if you are sure you want to open it. Confirm by choosing _Open_, which will open a terminal and run the program. After this one-time confirmation, you can run this program within the terminal.

#### Terminal download

```sh
# download the program
curl -OL https://github.com/bengarrett/retrotxtgo/releases/latest/download/retrotxtgo_darwin_all.tar.gz
# decompress and unpack the program
tar -xvzf retrotxtgo_darwin_all.tar.gz
# if needed, follow the instructions above
# then run the command
retrotxtgo --version
```

### Linux
### DEB
```sh
# Debian DEB package for Intel
wget https://github.com/bengarrett/retrotxtgo/releases/latest/download/retrotxtgo.deb
dpkg -i retrotxt.deb
retrotxt -v
```

### RPM
```sh
# Redhat RPM package
wget https://github.com/bengarrett/retrotxtgo/releases/latest/download/retrotxtgo.rpm
rpm -i retrotxt.rpm
retrotxt -v
```

### APK
```sh
# Alpine APK package
wget https://github.com/bengarrett/retrotxtgo/releases/latest/download/retrotxtgo.apk
apk add retrotxt.apk
retrotxt -v
```
