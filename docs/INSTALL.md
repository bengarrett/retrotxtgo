# Retrotxt

## Install

There are [numerous download](https://github.com/bengarrett/retrotxtgo/releases/latest/) releases. Otherwise these installation options are available.

### Windows 

The [download for Windows](https://github.com/bengarrett/retrotxtgo/releases/latest/download/retrotxtgo_windows_amd64.zip) is for Windows 10 or newer, for use with PowerShell or the Command Prompt. 

When unzipped, the program is a single file application that doesn't require installation. It can be placed into a directory of your choice.

Users of WSL (Windows Subsystem for Linux) should instead use the Linux instructions further in this document.

### macOS 

Unfortunately, newer macOS versions do not permit the running of unsigned terminal applications out of the box. But there is a workaround.

In System Settings, Privacy & Security, Security, toggle **Allow applications downloaded from App store and identified developers**.

1. Use Finder to extract the download for macOS, [`retrotxtgo_darwin_all.tar.gz`](https://github.com/bengarrett/retrotxtgo/releases/latest/download/retrotxtgo_darwin_all.tar.gz).
1. Use Finder to select the extracted `retrotxtgo` binary. 
1. <kbd>^</kbd> control-click the binary and choose _Open_. 
1. macOS will ask if you are sure you want to open it. 
1. Confirm by choosing _Open_, which will open a terminal and run the program. 
1. After this one-time confirmation, you can run this program within the terminal.

```sh
# A common location for terminal programs on macOS
/usr/local/bin
```

#### Terminal download

```sh
# download the program
curl -OL https://github.com/bengarrett/retrotxtgo/releases/latest/download/retrotxtgo_darwin_all.tar.gz

# decompress and unpack the program
tar -xvzf retrotxtgo_darwin_all.tar.gz

# if needed, follow the itemized instructions above
# then run the command
retrotxt --version
```

### Linux
### DEB
```sh
# Debian DEB package
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
