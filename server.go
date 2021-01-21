package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/Telmate/proxmox-api-go/proxmox"
)

func main() {
	insecurevar := true
	debugvar := false
	taskTimeoutvar, _ := strconv.Atoi(os.Getenv("TIMEOUT"))
	var insecure *bool
	insecure = &insecurevar        //, false, "TLS insecure mode")
	proxmox.Debug = &debugvar      //, false, "debug mode")
	taskTimeout := &taskTimeoutvar //, 300, "api task timeout in seconds")
	// fvmid := -1                    // "custom vmid (instead of auto)"
	// flag.Parse()
	tlsconf := &tls.Config{InsecureSkipVerify: true}
	if !*insecure {
		tlsconf = nil
	}
	c, _ := proxmox.NewClient(os.Getenv("PM_API_URL"), nil, tlsconf, *taskTimeout)
	err := c.Login(os.Getenv("PM_USER"), os.Getenv("PM_PASS"), os.Getenv("PM_OTP"))
	if err != nil {
		log.Fatal(err)
	}
	// vmid := *fvmid
	// if vmid < 0 {
	// 	if len(flag.Args()) > 1 {
	// 		vmid, err = strconv.Atoi(flag.Args()[len(flag.Args())-1])
	// 		if err != nil {
	// 			vmid = 0
	// 		}
	// 	} else if flag.Args()[0] == "idstatus" {
	// 		vmid = 0
	// 	}
	// }

	var jbody interface{}
	var vmr *proxmox.VmRef

	// if len(flag.Args()) == 0 {
	// 	fmt.Printf("Missing action, try start|stop vmid\n")
	// 	os.Exit(0)
	// }
	action := "getConfig"
	// action = "getNodes"
	// action = "createQemu"
	vmid := 108
	vmid = 108
	node := "node1"

	switch action {
	case "start":
		vmr = proxmox.NewVmRef(vmid)
		jbody, _ = c.StartVm(vmr)

	case "stop":

		vmr = proxmox.NewVmRef(vmid)
		jbody, _ = c.StopVm(vmr)

	case "destroy":
		// vmr = proxmox.NewVmRef(vmid)
		// jbody, err = c.StopVm(vmr)
		// failError(err)
		// jbody, _ = c.DeleteVm(vmr)

	case "getNodes":
		nodeList, err := c.GetNodeList()
		failError(err)
		cj, err := json.MarshalIndent(nodeList, "", "  ")
		log.Println(string(cj))
		fmt.Println(string(cj))

	case "getConfig":
		vmr = proxmox.NewVmRef(vmid)
		c.CheckVmRef(vmr)
		vmType := vmr.GetVmType()
		var config interface{}
		var err error
		if vmType == "qemu" {
			config, err = proxmox.NewConfigQemuFromApi(vmr, c)
		} else if vmType == "lxc" {
			config, err = proxmox.NewConfigLxcFromApi(vmr, c)
		}
		failError(err)
		cj, err := json.MarshalIndent(config, "", "  ")
		log.Println(string(cj))
		fmt.Println(string(cj))

	case "getNetworkInterfaces":
		vmr = proxmox.NewVmRef(vmid)
		c.CheckVmRef(vmr)
		networkInterfaces, err := c.GetVmAgentNetworkInterfaces(vmr)
		failError(err)

		networkInterfaceJson, err := json.Marshal(networkInterfaces)
		fmt.Println(string(networkInterfaceJson))

	case "createQemuJson":
		config, err := proxmox.NewConfigQemuFromJson(os.Stdin)
		failError(err)
		vmr = proxmox.NewVmRef(vmid)
		vmr.SetNode(flag.Args()[2])
		failError(config.CreateVm(vmr, c))
		log.Println("Complete")

	case "createQemu":
		/**
		{
			"vmid": 0,
			"name": "mail.dextion.com",
			"desc": "",
			"bios": "seabios",
			"onboot": true,
			"agent": 0,
			"memory": 2048,
			"balloon": 0,
			"os": "l26",
			"cores": 2,
			"sockets": 1,
			"vcpus": 0,
			"cpu": "host",
			"numa": false,
			"kvm": true,
			"hotplug": "network,disk,usb",
			"iso": "local:iso/ubuntu-20.04.1-live-server-amd64.iso",
			"fullclone": null,
			"boot": "cdn",
			"bootdisk": "scsi0",
			"scsihw": "virtio-scsi-pci",
			"disk": {
				"0": {
				"file": "vm-100-disk-0",
				"format": "raw",
				"size": "150G",
				"slot": 0,
				"storage": "datadisk",
				"storage_type": "lvm",
				"type": "scsi",
				"volume": "datadisk:vm-100-disk-0"
				}
			},
			"unused_disk": {},
			"network": {
				"0": {
				"bridge": "vmbr0",
				"firewall": true,
				"id": 0,
				"macaddr": "9A:2B:2C:6F:D5:A1",
				"model": "virtio"
				}
			},
			"tags": "",
			"diskGB": 0,
			"storage": "",
			"storageType": "",
			"nic": "",
			"bridge": "",
			"vlan": -1,
			"mac": "",
			"ciuser": "",
			"cipassword": "",
			"cicustom": "",
			"searchdomain": "",
			"nameserver": "",
			"sshkeys": "",
			"ipconfig0": "",
			"ipconfig1": "",
			"ipconfig2": ""
		}
		**/
		config := proxmox.ConfigQemu{
			VmID:        108,
			Name:        "Test",
			Description: "description test",
			Bios:        "seabios",
			Onboot:      true,
			Agent:       0,
			Memory:      2048,
			Balloon:     0,
			QemuOs:      "l26",
			QemuCores:   2,
			QemuSockets: 1,
			QemuVcpus:   0,
			QemuCpu:     "host",
			QemuNuma:    false,
			QemuKVM:     true,
			Hotplug:     "network,disk,usb",
			QemuIso:     "local:iso/ubuntu-20.04.1-live-server-amd64.iso",
			FullClone:   nil,
			Boot:        "cdn",
			BootDisk:    "scsi0",
			Scsihw:      "virtio-scsi-pci",
			QemuDisks: proxmox.QemuDevices{
				0: {
					"format":       "raw",
					"size":         "150G",
					"storage":      "datadisk",
					"storage_type": "lvm",
					"type":         "scsi",
				},
			},
			QemuNetworks: proxmox.QemuDevices{
				0: {
					"bridge":   "vmbr0",
					"firewall": true,
					"model":    "virtio",
				},
			},
			// QemuUnusedDisks:config.QemuDevices{},
			// QemuVga: config.QemuDevice{},
			// QemuSerials: config.QemuSerials{},
			HaState: "",
			Tags:    "",
			// CIuser:     "root",
			// CIpassword: "roots",
			// Ipconfig0:"",
		}
		vmr = proxmox.NewVmRef(vmid)
		vmr.SetNode(node)
		failError(config.CreateVm(vmr, c))
		log.Println("Completed")

	case "createLxc":
		config, err := proxmox.NewConfigLxcFromJson(os.Stdin)
		failError(err)
		vmr = proxmox.NewVmRef(vmid)
		vmr.SetNode(flag.Args()[2])
		failError(config.CreateLxc(vmr, c))
		log.Println("Complete")

	case "installQemu":
		config, err := proxmox.NewConfigQemuFromJson(os.Stdin)
		failError(err)
		if vmid > 0 {
			vmr = proxmox.NewVmRef(vmid)
		} else {
			nextid, err := c.GetNextID(0)
			failError(err)
			vmr = proxmox.NewVmRef(nextid)
		}
		vmr.SetNode(flag.Args()[1])
		log.Print("Creating node: ")
		log.Println(vmr)
		failError(config.CreateVm(vmr, c))
		_, err = c.StartVm(vmr)
		failError(err)
		sshPort, err := proxmox.SshForwardUsernet(vmr, c)
		failError(err)
		log.Println("Waiting for CDRom install shutdown (at least 5 minutes)")
		failError(proxmox.WaitForShutdown(vmr, c))
		log.Println("Restarting")
		_, err = c.StartVm(vmr)
		failError(err)
		sshPort, err = proxmox.SshForwardUsernet(vmr, c)
		failError(err)
		log.Println("SSH Portforward on:" + sshPort)
		log.Println("Complete")

	case "idstatus":
		maxid, err := proxmox.MaxVmId(c)
		failError(err)
		nextid, err := c.GetNextID(vmid)
		failError(err)
		log.Println("---")
		log.Printf("MaxID: %d\n", maxid)
		log.Printf("NextID: %d\n", nextid)
		log.Println("---")

	case "cloneQemu":
		config, err := proxmox.NewConfigQemuFromJson(os.Stdin)
		failError(err)
		log.Println("Looking for template: " + flag.Args()[1])
		sourceVmr, err := c.GetVmRefByName(flag.Args()[1])
		failError(err)
		if sourceVmr == nil {
			log.Fatal("Can't find template")
			return
		}
		if vmid == 0 {
			vmid, err = c.GetNextID(0)
		}
		vmr = proxmox.NewVmRef(vmid)
		vmr.SetNode(flag.Args()[2])
		log.Print("Creating node: ")
		log.Println(vmr)
		failError(config.CloneVm(sourceVmr, vmr, c))
		failError(config.UpdateConfig(vmr, c))
		log.Println("Complete")

	case "createQemuSnapshot":
		sourceVmr, err := c.GetVmRefByName(flag.Args()[1])
		jbody, err = c.CreateQemuSnapshot(sourceVmr, flag.Args()[2])
		failError(err)

	case "deleteQemuSnapshot":
		sourceVmr, err := c.GetVmRefByName(flag.Args()[1])
		jbody, err = c.DeleteQemuSnapshot(sourceVmr, flag.Args()[2])
		failError(err)

	case "listQemuSnapshot":
		sourceVmr, err := c.GetVmRefByName(flag.Args()[1])
		jbody, _, err = c.ListQemuSnapshot(sourceVmr)
		if rec, ok := jbody.(map[string]interface{}); ok {
			temp := rec["data"].([]interface{})
			for _, val := range temp {
				snapshotName := val.(map[string]interface{})
				if snapshotName["name"] != "current" {
					fmt.Println(snapshotName["name"])
				}
			}
		} else {
			fmt.Printf("record not a map[string]interface{}: %v\n", jbody)
		}
		failError(err)

	case "rollbackQemu":
		sourceVmr, err := c.GetVmRefByName(flag.Args()[1])
		jbody, err = c.RollbackQemuVm(sourceVmr, flag.Args()[2])
		failError(err)

	case "sshforward":
		vmr = proxmox.NewVmRef(vmid)
		sshPort, err := proxmox.SshForwardUsernet(vmr, c)
		failError(err)
		log.Println("SSH Portforward on:" + sshPort)

	case "sshbackward":
		vmr = proxmox.NewVmRef(vmid)
		err = proxmox.RemoveSshForwardUsernet(vmr, c)
		failError(err)
		log.Println("SSH Portforward off")

	case "sendstring":
		vmr = proxmox.NewVmRef(vmid)
		err = proxmox.SendKeysString(vmr, c, flag.Args()[2])
		failError(err)
		log.Println("Keys sent")

	case "nextid":
		id, err := c.GetNextID(0)
		failError(err)
		log.Printf("Getting Next Free ID: %d\n", id)

	case "checkid":
		i, err := strconv.Atoi(flag.Args()[1])
		failError(err)
		id, err := c.VMIdExists(i)
		failError(err)
		log.Printf("Selected ID is free: %d\n", id)

	case "migrate":
		vmr := proxmox.NewVmRef(vmid)
		c.GetVmInfo(vmr)
		args := flag.Args()
		if len(args) <= 1 {
			fmt.Printf("Missing target node\n")
			os.Exit(1)
		}
		_, err := c.MigrateNode(vmr, args[1], true)

		if err != nil {
			log.Printf("Error to move %+v\n", err)
			os.Exit(1)
		}
		log.Printf("VM %d is moved on %s\n", vmid, args[1])

	default:
		fmt.Printf("unknown action, try start|stop vmid\n")
	}
	if jbody != nil {
		log.Println(jbody)
	}
	//log.Println(vmr)
}

func failError(err error) {
	if err != nil {
		log.Fatal(err)
	}
	return
}
