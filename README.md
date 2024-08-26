# USP
	                                                                                     
![image](https://github.com/user-attachments/assets/ae04ddd3-9a81-4e27-be60-8d4e98d4e605)

## Features
This Go program establishes persistence on a Linux system by creating a udev rule that triggers the execution of a specified payload (binary or script). It offers two trigger options:

- **USB Persistence:** The payload is executed whenever a USB device is inserted.
- **Boot Persistence:** The payload is executed during system boot, leveraging the `/dev/random` device.

Additionally, it provides a cleanup option to remove the established persistence.

## Usage

1. **Compile the Go script:** `go build main.go -o usp`

2. **Run the script with root privileges (sudo):** `sudo ./usp`


You can customize the behavior using the following flags:

-   `-f <filename>`: Specify the path where the payload will be written (default: `/usr/local/bin/persistence`).
-   `-p <payload>`: Specify the path to the payload file (binary or script) that will be executed. This is a required flag.
-   `-r <rulesname>`: Specify the name of the udev rules file (default: `75-persistence.rules`).
-   `-usb`: Enable USB persistence.
-   `-random`: Enable boot persistence using `/dev/random`.
-   `-c`: Cleanup persistence, removing the payload file and udev rule.


## Example
- The following uses the USB persistence method to run `example.sh` everytime a USB device is connected.
```bash
sudo ./usp -p ./example.sh 
```

- The following uses the "random" persistence method to run `my_backdoor_binary` everytime `/dev/random` is loaded (such as at boot). It is installed at `/bin/ripgrep`. (Masquerading as the `ripgrep` binary). Additionally, the rules file is created in `/etc/udev/rules.d/123-notsektchy.rules`

```bash
sudo ./usp -random  -f /bin/ripgrep -p my_backdoor_binary -r 123-notsketchy.rules
```



## References && Additional Reading 
- https://www.aon.com/en/insights/cyber-labs/unveiling-sedexp
- https://ch4ik0.github.io/en/posts/leveraging-Linux-udev-for-persistence/
- https://opensource.com/article/18/11/udev
