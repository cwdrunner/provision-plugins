# Digital Rebar RAID Configuration Tool

The Digital Rebar RAID configuration tool (also known as drp-raid) is
designed to be an automated wrapper around lower-level RAID
configuration tools.  It is designed to work with the RAID workload in
the Sledgehammer environment.

## Usage

drp-raid can perform a few common tasks:

### Report Current Configuration

`drp-raid` with no arguments reports on the current state of all the
RAID controllers it knows about, including the state of all physical
drives and volumes.  This report is output in the form of a JSON array
of controllers.

### Clear All Configuration

`drp-raid -clear` will clear all configuration on all RAID controllers
attached to the system, including any foreign configuration.

### Show Current Volume Specifications

`drp-raid -volspecs` shows the current list of volumes on all RAID
controllers in the system in a form that can be fed into `drp-raid
-configure` to recreate the volumes as they are.

### Compile Volume Specifications Into Final Format

`drp-raid -compile` takes a list of volume specifications on stdin,
compiles them to a format suitable for the current system (including
resolving what drives to use and what the approximate final size of
each given volume will be), and outputs that list on stdout.  It
assumes that the RAID controllers are unconfigured, and it will change
the current RAID configuration.

### Configure Volumes on Raid Controllers

`drp-raid -configure` takes a list of volume specifications on stdin,
compiles them, and then attempts to create any needed volumes on the
RAID controllers.  Existing volumes will not be modified or deleted.

### Add Volumes to Raid Controllers

`drp-raid -append` attemps to add the list of volume specifications
from stdin to the preexisting volumes.  Unlike configure, you do not
need to pass in volume specifications for existing volumes.

## Supported RAID Controllers

drp-raid supports all RAID controllers supported by current MegaCLI
and StorCLI tooling for Avago-based RAID controllers, along with all
HP SmartArray controllers that are manageable by ssacli.  Support for
other tools and RAID controller types can be added as part of a
consulting or support engagement, depending in your license terms.

## Volume Specifications

Volume specifications (volspecs for short) are what drp-raid uses to
configure RAID volumes, no matter what RAID toolset or controller type
you are using.  A list of volume specifications looks like this:

```json
[
  {
    "RaidLevel": "string",
    "Size": "string",
    "StripeSize": "string",
    "Name": "string",
    "VolumeID": "string",
    "Bootable": false,
    "Type": "string",
    "Protocol": "string",
    "Controller": 0,
    "DiskCount": "string",
    "Disks": [
      {
        "Size": 0,
        "Slot": 0,
        "Enclosure": "string",
        "Type": "string",
        "Protocol": "string",
        "Volume": "string"
      }
    ]
  }
]
```

Each field in a volume specification is defined as follows:

* Controller: an integer that specifies which discovered controller
  the RAID volume should be created on.  drp-raid orders controllers
  based on PCI address.  This value defaults to 0, indicating that the
  volume will be built on the first discovered controller.

* RaidLevel: A string value that specifies what level of RAID to
  build.  RaidLevel must be present in all volspecs, and it has no
  default.  Valid RaidLevels are:

  * "jbod": All disks in this volspec should be turned into JBOD
    volumes.  If a given RAID controller does not support jbod mode, a
    single-disk RAID0 will be created for each disk instead.
  * "concat": The disks should be concatenated.
  * "raid0": Make a RAID0 with the chosen disks.
  * "raid1": Make a RAID1 with the chosen disks.
  * "raid5": Make a RAID5 with the chosen disks.
  * "raid6": Make a RAID6 with the chosen disks.
  * "raid00": Make a RAID00 with the chosen disks.
  * "raid10": Make a RAID10 with the chosen disks.
  * "raid1e": Make a RAID1e with the chosen disks.
  * "raid50": make a RAID50 with the chosen disks.
  * "raid60": make a RAID60 with the chosen disks.
  * "raidS": Works like jbod, but makes raid0 volumes.

  Not all RAID controllers will support the above RAID levels, and
  with the exception of jbod if a RAID level is not supported then
  volume creation will fail.

* Disks: A list of physical disk specifiers that indicate which
  physical disks should be used to create the volume.  When creating a
  volume by specifying individual disks, you are responsible for
  making sure that the choice of disks is sane for the desired
  controller.  If Disks is specified, then DiskCount is ignored and
  drp-raid will not perform more than minimal sanity checking on the
  disks provided.  Each entry in the Disks list must contain the
  following fields:

  * "Enclosure": a string that identifies which enclosure the disk is in.
  * "Slot": a number that identifies which slot in the Enclosure the desired
    physical disk is in.

* DiskCount: The number of disks that should be used to build the
  desired volume.  DiskCount can be one of the following:

  * "min", which indicates that the smallest number of disks that can
    be used for the requested RaidLevel whould be used.
  * "max", which indicates that the largest number of disks with the
    same Type and Protocol should be used.
  * A positive integer.

  If DiskCount is unspecified and Disks is also unspecified, DiskCount
  will default to "min".

* Size: A string value that indicates what the desired total useable
  size of the RAID array should be.  When you let drp-raid decide what
  physical disks to pick for volume creation, it will pick the
  smallest disks that can be used to satisfy the volume creation
  request.  Size can be one of the following values:

  * "min", which will pick the smallest disks that meet the rest of
    the constraints in the volspec.
  * "max", which will pick the largest disks that meet the rest of the
    constraints in the volspec.
  * A string containing a human-readable size (100 MB, 1 TB, 5 PB).

  If Size is unspecified or left blank, it will default to "max".

* Type: A comma-seperated list of disk types that should be tried in
  order when creating a volume.  Currently, individual items can be
  "disk" for rotational disks, and "ssd" for solid-state disks.  If
  unspecified, "disk,ssd" will be used.  All physical disks in a
  created volume will be of the same type, and it is perrmitted to
  have a list with one entry.

* Protocol: A comma-seperated list of low-level protocols that should
  be tried in order when creating a volume.  Currently, individual
  items can be "sata" for disks that communicate using the SATA
  protocol, and "sas" for disks that communicate using the SAS
  protocol.  All physical disks in a created volume will communicate
  using the same protocol, and it is perrmitted to have a list with
  one entry.

* StripeSize: A string containing the human-readable size of each
  individual stripe in the RAID array. It is generally a power of two
  less than or equal to 1 MB.  If unspecified, it defaults to "64 KB".

* Name: The desired name of the volume, if the RAID controller
  supports naming volumes.  Naming a volume is currently unsupported.

* VolumeID: The controller-specific ID of the created volume.

* Bootable: A boolean value indicating whether this volume should be
  the default one the RAID controller will use when booting the
  system.  Defaults to false, and Bootable support is currently not implemented.

## Example Transcript

```
[root@df8-bc-12-41-c4-72 ~]# ./drp-raid -clear
[root@df8-bc-12-41-c4-72 ~]# ./drp-raid
[
    {
      "ID": "0",
      "Driver": "storcli7",
      "PCI": {
        "Bus": 10,
        "Device": 0,
        "Function": 0
      },
      "JBODCapable": false,
      "RaidCapable": true,
      "RaidLevels": [
        "raid0",
        "raid1",
        "raid5",
        "raid6",
        "raid00",
        "raid10",
        "raid50",
        "raid60"
      ],
      "Volumes": [],
      "Disks": [
        {
          "ControllerID": "0",
          "ControllerDriver": "storcli7",
          "VolumeID": "",
          "Enclosure": "252",
          "Size": 499558383616,
          "UsedSize": 0,
          "SectorCount": 975699968,
          "PhysicalSectorSize": 512,
          "LogicalSectorSize": 512,
          "Slot": 0,
          "Protocol": "sata",
          "MediaType": "disk",
          "Status": "Unconfigured(good), Spun Up",
          "Info": {
            "Coerced Size": "465.25 GB [0x3a280000 Sectors]",
            "Connected Port Number": "2(path0)",
            "Device Firmware Level": "1S04",
            "Device Id": "1",
            "Device Speed": "3.0Gb/s",
            "Drive": "Not Certified",
            "Drive has flagged a S.M.A.R.T alert": "No",
            "Drive is formatted for PI information": "No",
            "Drive's NCQ setting": "N/A",
            "Enclosure Device ID": "252",
            "Enclosure position": "0",
            "FDE Capable": "Not Capable",
            "FDE Enable": "Disable",
            "Firmware state": "Unconfigured(good), Spun Up",
            "Foreign State": "None",
            "Inquiry Data": "WD-WMAYP8041031WDC WD5003ABYX-18WERA0                  01.01S04",
            "Last Predictive Failure Event Seq Number": "0",
            "Link Speed": "3.0Gb/s",
            "Locked": "Unlocked",
            "Logical Sector Size": "512",
            "Media Error Count": "0",
            "Media Type": "Hard Disk Device",
            "Needs EKM Attention": "No",
            "Non Coerced Size": "465.261 GB [0x3a286030 Sectors]",
            "Other Error Count": "0",
            "PD Type": "SATA",
            "PI": "No PI",
            "PI Eligibility": "No",
            "Physical Sector Size": "512",
            "Port status": "Active",
            "Port's Linkspeed": "3.0Gb/s",
            "Predictive Failure Count": "0",
            "Raw Size": "465.761 GB [0x3a386030 Sectors]",
            "SAS Address(0)": "0x4433221103000000",
            "Secured": "Unsecured",
            "Sequence Number": "27",
            "Shield Counter": "0",
            "Slot Number": "0",
            "Successful diagnostics completion on": "N/A",
            "WWN": "50014ee0ae7b61c5"
          }
        },
        {
          "ControllerID": "0",
          "ControllerDriver": "storcli7",
          "VolumeID": "",
          "Enclosure": "252",
          "Size": 499558383616,
          "UsedSize": 0,
          "SectorCount": 975699968,
          "PhysicalSectorSize": 512,
          "LogicalSectorSize": 512,
          "Slot": 1,
          "Protocol": "sata",
          "MediaType": "disk",
          "Status": "Unconfigured(good), Spun Up",
          "Info": {
            "Coerced Size": "465.25 GB [0x3a280000 Sectors]",
            "Connected Port Number": "1(path0)",
            "Device Firmware Level": "1S04",
            "Device Id": "2",
            "Device Speed": "3.0Gb/s",
            "Drive": "Not Certified",
            "Drive has flagged a S.M.A.R.T alert": "No",
            "Drive is formatted for PI information": "No",
            "Drive's NCQ setting": "N/A",
            "Enclosure Device ID": "252",
            "Enclosure position": "0",
            "FDE Capable": "Not Capable",
            "FDE Enable": "Disable",
            "Firmware state": "Unconfigured(good), Spun Up",
            "Foreign State": "None",
            "Inquiry Data": "WD-WMAYP7926193WDC WD5003ABYX-18WERA0                  01.01S04",
            "Last Predictive Failure Event Seq Number": "0",
            "Link Speed": "3.0Gb/s",
            "Locked": "Unlocked",
            "Logical Sector Size": "512",
            "Media Error Count": "0",
            "Media Type": "Hard Disk Device",
            "Needs EKM Attention": "No",
            "Non Coerced Size": "465.261 GB [0x3a286030 Sectors]",
            "Other Error Count": "0",
            "PD Type": "SATA",
            "PI": "No PI",
            "PI Eligibility": "No",
            "Physical Sector Size": "512",
            "Port status": "Active",
            "Port's Linkspeed": "3.0Gb/s",
            "Predictive Failure Count": "0",
            "Raw Size": "465.761 GB [0x3a386030 Sectors]",
            "SAS Address(0)": "0x4433221102000000",
            "Secured": "Unsecured",
            "Sequence Number": "27",
            "Shield Counter": "0",
            "Slot Number": "1",
            "Successful diagnostics completion on": "N/A",
            "WWN": "50014ee05925c428"
          }
        },
        {
          "ControllerID": "0",
          "ControllerDriver": "storcli7",
          "VolumeID": "",
          "Enclosure": "252",
          "Size": 499558383616,
          "UsedSize": 0,
          "SectorCount": 975699968,
          "PhysicalSectorSize": 512,
          "LogicalSectorSize": 512,
          "Slot": 2,
          "Protocol": "sata",
          "MediaType": "disk",
          "Status": "Unconfigured(good), Spun Up",
          "Info": {
            "Coerced Size": "465.25 GB [0x3a280000 Sectors]",
            "Connected Port Number": "3(path0)",
            "Device Firmware Level": "1S04",
            "Device Id": "3",
            "Device Speed": "3.0Gb/s",
            "Drive": "Not Certified",
            "Drive has flagged a S.M.A.R.T alert": "No",
            "Drive is formatted for PI information": "No",
            "Drive's NCQ setting": "N/A",
            "Enclosure Device ID": "252",
            "Enclosure position": "0",
            "FDE Capable": "Not Capable",
            "FDE Enable": "Disable",
            "Firmware state": "Unconfigured(good), Spun Up",
            "Foreign State": "None",
            "Inquiry Data": "WD-WMAYP8040536WDC WD5003ABYX-18WERA0                  01.01S04",
            "Last Predictive Failure Event Seq Number": "0",
            "Link Speed": "3.0Gb/s",
            "Locked": "Unlocked",
            "Logical Sector Size": "512",
            "Media Error Count": "0",
            "Media Type": "Hard Disk Device",
            "Needs EKM Attention": "No",
            "Non Coerced Size": "465.261 GB [0x3a286030 Sectors]",
            "Other Error Count": "0",
            "PD Type": "SATA",
            "PI": "No PI",
            "PI Eligibility": "No",
            "Physical Sector Size": "512",
            "Port status": "Active",
            "Port's Linkspeed": "3.0Gb/s",
            "Predictive Failure Count": "0",
            "Raw Size": "465.761 GB [0x3a386030 Sectors]",
            "SAS Address(0)": "0x4433221101000000",
            "Secured": "Unsecured",
            "Sequence Number": "21",
            "Shield Counter": "0",
            "Slot Number": "2",
            "Successful diagnostics completion on": "N/A",
            "WWN": "50014ee0ae7b6632"
          }
        },
        {
          "ControllerID": "0",
          "ControllerDriver": "storcli7",
          "VolumeID": "",
          "Enclosure": "252",
          "Size": 499558383616,
          "UsedSize": 0,
          "SectorCount": 975699968,
          "PhysicalSectorSize": 512,
          "LogicalSectorSize": 512,
          "Slot": 3,
          "Protocol": "sata",
          "MediaType": "disk",
          "Status": "Unconfigured(good), Spun Up",
          "Info": {
            "Coerced Size": "465.25 GB [0x3a280000 Sectors]",
            "Connected Port Number": "0(path0)",
            "Device Firmware Level": "1S04",
            "Device Id": "0",
            "Device Speed": "3.0Gb/s",
            "Drive": "Not Certified",
            "Drive has flagged a S.M.A.R.T alert": "No",
            "Drive is formatted for PI information": "No",
            "Drive's NCQ setting": "N/A",
            "Enclosure Device ID": "252",
            "Enclosure position": "0",
            "Exit Code": "0x00",
            "FDE Capable": "Not Capable",
            "FDE Enable": "Disable",
            "Firmware state": "Unconfigured(good), Spun Up",
            "Foreign State": "None",
            "Inquiry Data": "WD-WMAYP8042658WDC WD5003ABYX-18WERA0                  01.01S04",
            "Last Predictive Failure Event Seq Number": "0",
            "Link Speed": "3.0Gb/s",
            "Locked": "Unlocked",
            "Logical Sector Size": "512",
            "Media Error Count": "0",
            "Media Type": "Hard Disk Device",
            "Needs EKM Attention": "No",
            "Non Coerced Size": "465.261 GB [0x3a286030 Sectors]",
            "Other Error Count": "0",
            "PD Type": "SATA",
            "PI": "No PI",
            "PI Eligibility": "No",
            "Physical Sector Size": "512",
            "Port status": "Active",
            "Port's Linkspeed": "3.0Gb/s",
            "Predictive Failure Count": "0",
            "Raw Size": "465.761 GB [0x3a386030 Sectors]",
            "SAS Address(0)": "0x4433221100000000",
            "Secured": "Unsecured",
            "Sequence Number": "21",
            "Shield Counter": "0",
            "Slot Number": "3",
            "Successful diagnostics completion on": "N/A",
            "WWN": "50014ee003d097f2"
          }
        }
      ],
      "Info": {
        "Access Policy": "Yes",
        "Alarm": "Present",
        "Alarm Control": "Yes",
        "Alarm Disable": "Yes",
        "Allow Boot with Preserved Cache": "No",
        "Allow Ctrl Encryption": "No",
        "Allow HDD SAS/SATA Mix in VD": "Yes",
        "Allow HDD/SSD Mix in VD": "Yes",
        "Allow Mix in Enclosure": "Yes",
        "Allow Mixed Redundancy on Array": "No",
        "Allow SATA in Cluster": "No",
        "Allow SSD SAS/SATA Mix in VD": "Yes",
        "Allowed Device Type": "SAS/SATA Mix",
        "Any Offline VD Cache Preserved": "No",
        "Auto Detect BackPlane Enable": "SGPIO/i2c SEP",
        "Auto Detect BackPlane Enabled": "SGPIO/i2c SEP",
        "Auto Enhanced Import": "Yes",
        "Auto Rebuild": "Enabled",
        "BBU": "Absent",
        "BGI Rate": "30%",
        "BIOS Continue on Error": "Yes",
        "BIOS Enumerate VDs": "Yes",
        "BIOS Version": "3.30.02.2_4.16.08.00_0x06060A05",
        "BOOT Version": "09.250.01.219",
        "Background Rate": "30",
        "Battery FRU": "N/A",
        "Battery Warning": "Disabled",
        "Block SSD Write Disk Cache Change": "No",
        "Boot Block Version": "2.02.00.00-0000",
        "Boot Volume Supported": "NO",
        "BreakMirror RAID Support": "Yes",
        "CC Rate": "Yes",
        "Cache Flush Interval": "4s",
        "Cache When BBU Bad": "Disabled",
        "Cached IO": "No",
        "Check Consistency Rate": "30%",
        "ChipRevision": "B4",
        "Cluster Active": "No",
        "Cluster Disable": "Yes",
        "Cluster Mode": "Disabled",
        "Cluster Permitted": "No",
        "Cluster Support": "No",
        "Coercion Mode": "None",
        "Controller Id": "0000",
        "Critical Disks": "0",
        "Current Size of CacheCade": "0 GB",
        "Current Size of FW Cache": "346 MB",
        "Current Time": "21:21:25 2/14, 2018",
        "Dedicated Hot Spare": "Yes",
        "Default LD PowerSave Policy": "Controller Defined",
        "Default spin down time in minutes": "30",
        "Degraded": "0",
        "Delay Among Spinup Groups": "2s",
        "Delay during POST": "0",
        "Deny CC": "No",
        "Deny Clear": "No",
        "Deny Force Failed": "No",
        "Deny Force Good/Bad": "No",
        "Deny Locate": "No",
        "Deny Missing Replace": "No",
        "Deny SCSI Passthrough": "No",
        "Deny SMP Passthrough": "No",
        "Deny STP Passthrough": "No",
        "Device Id": "0079",
        "Device Interface": "PCIE",
        "Direct PD Mapping": "No",
        "Dirty LED Shows Drive Activity": "No",
        "Disable Copyback": "No",
        "Disable Ctrl-R": "Yes",
        "Disable Join Mirror": "No",
        "Disable Online Controller Reset": "No",
        "Disable Online PFK Change": "No",
        "Disable Puncturing": "No",
        "Disable Spin Down of hot spares": "No",
        "Disk Cache Policy": "Yes",
        "Disks": "4",
        "ECC Bucket Count": "0",
        "Ecc Bucket Leak Rate": "1440 Minutes",
        "Ecc Bucket Size": "15",
        "Enable Copyback on SMART": "No",
        "Enable Copyback to SSD on SMART Error": "Yes",
        "Enable JBOD": "No",
        "Enable LDBBM": "Yes",
        "Enable Led Header": "Yes",
        "Enable SSD Patrol Read": "No",
        "Enable Shield State": "No",
        "Enable Spin Down of UnConfigured Drives": "Yes",
        "Enable Web BIOS": "Yes",
        "EnableCrashDump": "No",
        "EnableLDBBM": "Yes",
        "Exit Code": "0x00",
        "Expose Enclosure Devices": "Enabled",
        "FW Package Build": "12.15.0-0239",
        "FW Version": "2.130.403-4660",
        "Failed Disks": "0",
        "Flash": "Present",
        "Flush Time": "4 seconds",
        "Force Offline": "Yes",
        "Force Online": "Yes",
        "Force Rebuild": "Yes",
        "Foreign Config Import": "Yes",
        "Global Hot Spares": "Yes",
        "Host Interface": "PCIE",
        "Host Request Reordering": "Enabled",
        "IO Policy": "Yes",
        "Interrupt Throttle Active Count": "16",
        "Interrupt Throttle Completion": "50us",
        "LED Show Drive Activity": "Yes",
        "Load Balance Mode": "Auto",
        "Maintain PD Fail History": "Enabled",
        "Max Arms Per VD": "32",
        "Max Arrays": "128",
        "Max Chained Enclosures": "16",
        "Max Configurable CacheCade Size": "0 GB",
        "Max Data Transfer Size": "8192 sectors",
        "Max Drives to Spinup at One Time": "4",
        "Max LD per array": "16",
        "Max Number of VDs": "64",
        "Max Parallel Commands": "1008",
        "Max SGE Count": "80",
        "Max Spans Per VD": "8",
        "Max Strip Size": "1.0 MB",
        "Max Strips PerIO": "42",
        "Maximum number of direct attached drives to spin up in 1 min": "120",
        "Memory": "Present",
        "Memory Correctable Errors": "0",
        "Memory Size": "512MB",
        "Memory Uncorrectable Errors": "0",
        "Mfg. Date": "03/13/12",
        "Min Strip Size": "8 KB",
        "NCQ": "No",
        "NVDATA Version": "2.09.03-0058",
        "NVRAM": "Present",
        "Number of Backend Port": "8",
        "Number of Frontend Port": "0",
        "Offline": "0",
        "On board Expander": "Absent",
        "PFK TrailTime Remaining": "0 days 0 hours",
        "PFK in NVRAM": "No",
        "POST delay": "90 seconds",
        "PR Correct Unconfigured Areas": "Yes",
        "PR Rate": "30%",
        "Patrol Read Rate": "Yes",
        "Phy Polarity": "0",
        "Phy PolaritySplit": "0",
        "Physical Devices": "5",
        "Physical Drive Coercion Mode": "Disabled",
        "Point In Time Progress": "No",
        "Power Saving option": "Don't Auto spin down Configured Drives",
        "Power Savings": "No",
        "PreBoot CLI Enabled": "Yes",
        "Preboot CLI Version": "04.04-020:#%00009",
        "Predictive Fail Poll Interval": "300sec",
        "Product Name": "LSI MegaRAID SAS 9260-8i",
        "RAID Level Supported": "RAID0, RAID1, RAID5, RAID6, RAID00, RAID10, RAID50, RAID60, PRL 11, PRL 11 with spanning, SRL 3 supported, PRL11-RLQ0 DDF layout with no span, PRL11-RLQ0 DDF layout with span",
        "Read Policy": "Yes",
        "Rebuild Rate": "30%",
        "Reconstruct Rate": "Yes",
        "Reconstruction": "Yes",
        "Reconstruction Rate": "30%",
        "Restore Hot Spare on Insertion": "No",
        "Restore HotSpare on Insertion": "Disabled",
        "Revertible Hot Spares": "Yes",
        "Revision No": "79B",
        "Rework Date": "00/00/00",
        "SAS Address": "500605b004872390",
        "SAS Disable": "No",
        "SMART Mode": "Mode 6",
        "Security Key Assigned": "No",
        "Security Key Failed": "No",
        "Security Key Not Backedup": "No",
        "Self Diagnostic": "Yes",
        "Serial Debugger": "Present",
        "Serial No": "SV21106771",
        "Snapshot Enabled": "No",
        "Spanning": "Yes",
        "Spin Down Mode": "None",
        "Spin Down time": "30",
        "Strip Size": "256kB",
        "SubDeviceId": "9261",
        "SubVendorId": "1000",
        "Support Boot Time PFK Change": "No",
        "Support Breakmirror": "No",
        "Support PFK": "Yes",
        "Support PI": "No",
        "Support Security": "No",
        "Support Shield State": "No",
        "Support Temperature": "Yes",
        "Support the OCE without adding drives": "Yes",
        "Supported Drives": "SAS, SATA",
        "T10 Power State": "No",
        "TPM": "Absent",
        "TTY Log In Flash": "No",
        "Temperature sensor for ROC": "Absent",
        "Temperature sensor for controller": "Absent",
        "Time taken to detect CME": "60s",
        "Topology Type": "None",
        "Treat Single span R1E as R10": "No",
        "Un-Certified Hard Disk Drives": "Allow",
        "Upgrade Key": "Absent",
        "Use FDE Only": "No",
        "Use disk activity for locate": "No",
        "Vendor Id": "1000",
        "Virtual Drives": "0",
        "WebBIOS Version": "6.0-54-e_50-Rel",
        "Write Policy": "Yes",
        "ZCR Config": "Unknown",
        "Zero Based Enclosure Enumeration": "No"
      }
    }
  ]

[root@df8-bc-12-41-c4-72 ~]# ./drp-raid -volspecs
[]
[root@df8-bc-12-41-c4-72 ~]# echo '[{"RaidLevel": "raid1","DiskCount": "min"},{"RaidLevel": "raid1","DiskCount": "min"}]' |./drp-raid -compile
[
    {
      "RaidLevel": "raid1",
      "Size": "465.25 GB",
      "StripeSize": "64 KB",
      "Name": "",
      "VolumeID": "",
      "Bootable": false,
      "Type": "disk",
      "Protocol": "sata",
      "Controller": 0,
      "Disks": [
        {
          "Size": 499558383616,
          "Slot": 0,
          "Enclosure": "252",
          "Type": "disk",
          "Protocol": "sata",
          "Volume": ""
        },
        {
          "Size": 499558383616,
          "Slot": 1,
          "Enclosure": "252",
          "Type": "disk",
          "Protocol": "sata",
          "Volume": ""
        }
      ],
      "DiskCount": ""
    },
    {
      "RaidLevel": "raid1",
      "Size": "465.25 GB",
      "StripeSize": "64 KB",
      "Name": "",
      "VolumeID": "",
      "Bootable": false,
      "Type": "disk",
      "Protocol": "sata",
      "Controller": 0,
      "Disks": [
        {
          "Size": 499558383616,
          "Slot": 2,
          "Enclosure": "252",
          "Type": "disk",
          "Protocol": "sata",
          "Volume": ""
        },
        {
          "Size": 499558383616,
          "Slot": 3,
          "Enclosure": "252",
          "Type": "disk",
          "Protocol": "sata",
          "Volume": ""
        }
      ],
      "DiskCount": ""
    }
  ]
[root@df8-bc-12-41-c4-72 ~]# echo '[{"RaidLevel": "raid1","DiskCount": "min"}]' |./drp-raid -configure
2018/02/14 21:22:28 Adapter 0: Created VD 0

Adapter 0: Configured the Adapter!!


Exit Code: 0x00

2018/02/14 21:22:28 Created raid1 on storcli7:0
[root@df8-bc-12-41-c4-72 ~]# echo '[{"RaidLevel": "raid1","DiskCount": "min"}]' |./drp-raid -configure
2018/02/14 21:22:43 All volumes already present, nothing to to
[root@df8-bc-12-41-c4-72 ~]# ./drp-raid -volspecs
[
    {
      "RaidLevel": "raid1",
      "Size": "465.25 GB",
      "StripeSize": "64 KB",
      "Name": "",
      "VolumeID": "0",
      "Bootable": false,
      "Type": "disk",
      "Protocol": "sata",
      "Controller": 0,
      "Disks": [
        {
          "Size": 499558383616,
          "Slot": 0,
          "Enclosure": "252",
          "Type": "disk",
          "Protocol": "sata",
          "Volume": "0"
        },
        {
          "Size": 499558383616,
          "Slot": 1,
          "Enclosure": "252",
          "Type": "disk",
          "Protocol": "sata",
          "Volume": "0"
        }
      ],
      "DiskCount": ""
    }
  ]
[root@df8-bc-12-41-c4-72 ~]# echo '[{"RaidLevel": "raid1","DiskCount": "min"}]' |./drp-raid -append
2018/02/14 21:23:13 Adapter 0: Created VD 1

Adapter 0: Configured the Adapter!!


Exit Code: 0x00

2018/02/14 21:23:13 Created raid1 on storcli7:0
[root@df8-bc-12-41-c4-72 ~]# ./drp-raid -volspecs
[
    {
      "RaidLevel": "raid1",
      "Size": "465.25 GB",
      "StripeSize": "64 KB",
      "Name": "",
      "VolumeID": "0",
      "Bootable": false,
      "Type": "disk",
      "Protocol": "sata",
      "Controller": 0,
      "Disks": [
        {
          "Size": 499558383616,
          "Slot": 0,
          "Enclosure": "252",
          "Type": "disk",
          "Protocol": "sata",
          "Volume": "0"
        },
        {
          "Size": 499558383616,
          "Slot": 1,
          "Enclosure": "252",
          "Type": "disk",
          "Protocol": "sata",
          "Volume": "0"
        }
      ],
      "DiskCount": ""
    },
    {
      "RaidLevel": "raid1",
      "Size": "465.25 GB",
      "StripeSize": "64 KB",
      "Name": "",
      "VolumeID": "1",
      "Bootable": false,
      "Type": "disk",
      "Protocol": "sata",
      "Controller": 0,
      "Disks": [
        {
          "Size": 499558383616,
          "Slot": 2,
          "Enclosure": "252",
          "Type": "disk",
          "Protocol": "sata",
          "Volume": "1"
        },
        {
          "Size": 499558383616,
          "Slot": 3,
          "Enclosure": "252",
          "Type": "disk",
          "Protocol": "sata",
          "Volume": "1"
        }
      ],
      "DiskCount": ""
    }
  ]
[root@df8-bc-12-41-c4-72 ~]# echo '[{"RaidLevel": "raid1","DiskCount": "min"},{"RaidLevel": "raid1","DiskCount": "min"}]' |./drp-raid -configure
2018/02/14 21:23:58 All volumes already present, nothing to to
[root@df8-bc-12-41-c4-72 ~]# ./drp-raid
[
    {
      "ID": "0",
      "Driver": "storcli7",
      "PCI": {
        "Bus": 10,
        "Device": 0,
        "Function": 0
      },
      "JBODCapable": false,
      "RaidCapable": true,
      "RaidLevels": [
        "raid0",
        "raid1",
        "raid5",
        "raid6",
        "raid00",
        "raid10",
        "raid50",
        "raid60"
      ],
      "Volumes": [
        {
          "ControllerID": "0",
          "ControllerDriver": "storcli7",
          "ID": "0",
          "Name": "",
          "Status": "Optimal",
          "RaidLevel": "raid1",
          "Size": 499558383616,
          "StripeSize": 65536,
          "Spans": 1,
          "SpanLength": 2,
          "Disks": [
            {
              "ControllerID": "0",
              "ControllerDriver": "storcli7",
              "VolumeID": "0",
              "Enclosure": "252",
              "Size": 499558383616,
              "UsedSize": 0,
              "SectorCount": 975699968,
              "PhysicalSectorSize": 512,
              "LogicalSectorSize": 512,
              "Slot": 0,
              "Protocol": "sata",
              "MediaType": "disk",
              "Status": "Online, Spun Up",
              "Info": {
                "Coerced Size": "465.25 GB [0x3a280000 Sectors]",
                "Commissioned Spare": "No",
                "Connected Port Number": "2(path0)",
                "Device Firmware Level": "1S04",
                "Device Id": "1",
                "Device Speed": "3.0Gb/s",
                "Drive": "Not Certified",
                "Drive has flagged a S.M.A.R.T alert": "No",
                "Drive is formatted for PI information": "No",
                "Drive's NCQ setting": "N/A",
                "Drive's postion": "DiskGroup: 0, Span: 0, Arm: 0",
                "Emergency Spare": "No",
                "Enclosure Device ID": "252",
                "Enclosure position": "0",
                "FDE Capable": "Not Capable",
                "FDE Enable": "Disable",
                "Firmware state": "Online, Spun Up",
                "Foreign State": "None",
                "Inquiry Data": "WD-WMAYP8041031WDC WD5003ABYX-18WERA0                  01.01S04",
                "Last Predictive Failure Event Seq Number": "0",
                "Link Speed": "3.0Gb/s",
                "Locked": "Unlocked",
                "Logical Sector Size": "512",
                "Media Error Count": "0",
                "Media Type": "Hard Disk Device",
                "Needs EKM Attention": "No",
                "Non Coerced Size": "465.261 GB [0x3a286030 Sectors]",
                "Other Error Count": "0",
                "PD": "1 Information",
                "PD Type": "SATA",
                "PI": "No PI",
                "PI Eligibility": "No",
                "Physical Sector Size": "512",
                "Port status": "Active",
                "Port's Linkspeed": "3.0Gb/s",
                "Predictive Failure Count": "0",
                "Raw Size": "465.761 GB [0x3a386030 Sectors]",
                "SAS Address(0)": "0x4433221103000000",
                "Secured": "Unsecured",
                "Sequence Number": "28",
                "Shield Counter": "0",
                "Slot Number": "0",
                "Successful diagnostics completion on": "N/A",
                "WWN": "50014ee0ae7b61c5"
              }
            },
            {
              "ControllerID": "0",
              "ControllerDriver": "storcli7",
              "VolumeID": "0",
              "Enclosure": "252",
              "Size": 499558383616,
              "UsedSize": 0,
              "SectorCount": 975699968,
              "PhysicalSectorSize": 512,
              "LogicalSectorSize": 512,
              "Slot": 1,
              "Protocol": "sata",
              "MediaType": "disk",
              "Status": "Online, Spun Up",
              "Info": {
                "Coerced Size": "465.25 GB [0x3a280000 Sectors]",
                "Commissioned Spare": "No",
                "Connected Port Number": "1(path0)",
                "Device Firmware Level": "1S04",
                "Device Id": "2",
                "Device Speed": "3.0Gb/s",
                "Drive": "Not Certified",
                "Drive has flagged a S.M.A.R.T alert": "No",
                "Drive is formatted for PI information": "No",
                "Drive's NCQ setting": "N/A",
                "Drive's postion": "DiskGroup: 0, Span: 0, Arm: 1",
                "Emergency Spare": "No",
                "Enclosure Device ID": "252",
                "Enclosure position": "0",
                "FDE Capable": "Not Capable",
                "FDE Enable": "Disable",
                "Firmware state": "Online, Spun Up",
                "Foreign State": "None",
                "Inquiry Data": "WD-WMAYP7926193WDC WD5003ABYX-18WERA0                  01.01S04",
                "Last Predictive Failure Event Seq Number": "0",
                "Link Speed": "3.0Gb/s",
                "Locked": "Unlocked",
                "Logical Sector Size": "512",
                "Media Error Count": "0",
                "Media Type": "Hard Disk Device",
                "Needs EKM Attention": "No",
                "Non Coerced Size": "465.261 GB [0x3a286030 Sectors]",
                "Other Error Count": "0",
                "PD Type": "SATA",
                "PI": "No PI",
                "PI Eligibility": "No",
                "Physical Sector Size": "512",
                "Port status": "Active",
                "Port's Linkspeed": "3.0Gb/s",
                "Predictive Failure Count": "0",
                "Raw Size": "465.761 GB [0x3a386030 Sectors]",
                "SAS Address(0)": "0x4433221102000000",
                "Secured": "Unsecured",
                "Sequence Number": "28",
                "Shield Counter": "0",
                "Slot Number": "1",
                "Successful diagnostics completion on": "N/A",
                "WWN": "50014ee05925c428"
              }
            }
          ],
          "Info": {
            "Bad Blocks Exist": "No",
            "Creation Date": "14-02-2018",
            "Creation Time": "09:22:25 PM",
            "Current Access Policy": "Read/Write",
            "Current Cache Policy": "WriteThrough, ReadAhead, Direct, No Write Cache if Bad BBU",
            "Default Access Policy": "Read/Write",
            "Default Cache Policy": "WriteBack, ReadAhead, Direct, No Write Cache if Bad BBU",
            "Disk Cache Policy": "Disk's Default",
            "Encryption Type": "None",
            "Is VD Cached": "No",
            "Logical Sector Size": "512",
            "Mirror Data": "465.25 GB",
            "Name": "",
            "Number Of Drives": "2",
            "Number of Spans": "1",
            "PD": "0 Information",
            "Physical Sector Size": "512",
            "RAID Level": "Primary-1, Secondary-0, RAID Level Qualifier-0",
            "Size": "465.25 GB",
            "Span Depth": "1",
            "State": "Optimal",
            "Strip Size": "64 KB",
            "VD has Emulated PD": "No",
            "Virtual Drive": "0 (Target Id: 0)"
          }
        },
        {
          "ControllerID": "0",
          "ControllerDriver": "storcli7",
          "ID": "1",
          "Name": "",
          "Status": "Optimal",
          "RaidLevel": "raid1",
          "Size": 499558383616,
          "StripeSize": 65536,
          "Spans": 1,
          "SpanLength": 2,
          "Disks": [
            {
              "ControllerID": "0",
              "ControllerDriver": "storcli7",
              "VolumeID": "1",
              "Enclosure": "252",
              "Size": 499558383616,
              "UsedSize": 0,
              "SectorCount": 975699968,
              "PhysicalSectorSize": 512,
              "LogicalSectorSize": 512,
              "Slot": 2,
              "Protocol": "sata",
              "MediaType": "disk",
              "Status": "Online, Spun Up",
              "Info": {
                "Coerced Size": "465.25 GB [0x3a280000 Sectors]",
                "Commissioned Spare": "No",
                "Connected Port Number": "3(path0)",
                "Device Firmware Level": "1S04",
                "Device Id": "3",
                "Device Speed": "3.0Gb/s",
                "Drive": "Not Certified",
                "Drive has flagged a S.M.A.R.T alert": "No",
                "Drive is formatted for PI information": "No",
                "Drive's NCQ setting": "N/A",
                "Drive's postion": "DiskGroup: 1, Span: 0, Arm: 0",
                "Emergency Spare": "No",
                "Enclosure Device ID": "252",
                "Enclosure position": "0",
                "FDE Capable": "Not Capable",
                "FDE Enable": "Disable",
                "Firmware state": "Online, Spun Up",
                "Foreign State": "None",
                "Inquiry Data": "WD-WMAYP8040536WDC WD5003ABYX-18WERA0                  01.01S04",
                "Last Predictive Failure Event Seq Number": "0",
                "Link Speed": "3.0Gb/s",
                "Locked": "Unlocked",
                "Logical Sector Size": "512",
                "Media Error Count": "0",
                "Media Type": "Hard Disk Device",
                "Needs EKM Attention": "No",
                "Non Coerced Size": "465.261 GB [0x3a286030 Sectors]",
                "Other Error Count": "0",
                "PD": "1 Information",
                "PD Type": "SATA",
                "PI": "No PI",
                "PI Eligibility": "No",
                "Physical Sector Size": "512",
                "Port status": "Active",
                "Port's Linkspeed": "3.0Gb/s",
                "Predictive Failure Count": "0",
                "Raw Size": "465.761 GB [0x3a386030 Sectors]",
                "SAS Address(0)": "0x4433221101000000",
                "Secured": "Unsecured",
                "Sequence Number": "22",
                "Shield Counter": "0",
                "Slot Number": "2",
                "Successful diagnostics completion on": "N/A",
                "WWN": "50014ee0ae7b6632"
              }
            },
            {
              "ControllerID": "0",
              "ControllerDriver": "storcli7",
              "VolumeID": "1",
              "Enclosure": "252",
              "Size": 499558383616,
              "UsedSize": 0,
              "SectorCount": 975699968,
              "PhysicalSectorSize": 512,
              "LogicalSectorSize": 512,
              "Slot": 3,
              "Protocol": "sata",
              "MediaType": "disk",
              "Status": "Online, Spun Up",
              "Info": {
                "Coerced Size": "465.25 GB [0x3a280000 Sectors]",
                "Commissioned Spare": "No",
                "Connected Port Number": "0(path0)",
                "Device Firmware Level": "1S04",
                "Device Id": "0",
                "Device Speed": "3.0Gb/s",
                "Drive": "Not Certified",
                "Drive has flagged a S.M.A.R.T alert": "No",
                "Drive is formatted for PI information": "No",
                "Drive's NCQ setting": "N/A",
                "Drive's postion": "DiskGroup: 1, Span: 0, Arm: 1",
                "Emergency Spare": "No",
                "Enclosure Device ID": "252",
                "Enclosure position": "0",
                "Exit Code": "0x00",
                "FDE Capable": "Not Capable",
                "FDE Enable": "Disable",
                "Firmware state": "Online, Spun Up",
                "Foreign State": "None",
                "Inquiry Data": "WD-WMAYP8042658WDC WD5003ABYX-18WERA0                  01.01S04",
                "Last Predictive Failure Event Seq Number": "0",
                "Link Speed": "3.0Gb/s",
                "Locked": "Unlocked",
                "Logical Sector Size": "512",
                "Media Error Count": "0",
                "Media Type": "Hard Disk Device",
                "Needs EKM Attention": "No",
                "Non Coerced Size": "465.261 GB [0x3a286030 Sectors]",
                "Other Error Count": "0",
                "PD Type": "SATA",
                "PI": "No PI",
                "PI Eligibility": "No",
                "Physical Sector Size": "512",
                "Port status": "Active",
                "Port's Linkspeed": "3.0Gb/s",
                "Predictive Failure Count": "0",
                "Raw Size": "465.761 GB [0x3a386030 Sectors]",
                "SAS Address(0)": "0x4433221100000000",
                "Secured": "Unsecured",
                "Sequence Number": "22",
                "Shield Counter": "0",
                "Slot Number": "3",
                "Successful diagnostics completion on": "N/A",
                "WWN": "50014ee003d097f2"
              }
            }
          ],
          "Info": {
            "Bad Blocks Exist": "No",
            "Creation Date": "14-02-2018",
            "Creation Time": "09:23:09 PM",
            "Current Access Policy": "Read/Write",
            "Current Cache Policy": "WriteThrough, ReadAhead, Direct, No Write Cache if Bad BBU",
            "Default Access Policy": "Read/Write",
            "Default Cache Policy": "WriteBack, ReadAhead, Direct, No Write Cache if Bad BBU",
            "Disk Cache Policy": "Disk's Default",
            "Encryption Type": "None",
            "Is VD Cached": "No",
            "Logical Sector Size": "512",
            "Mirror Data": "465.25 GB",
            "Name": "",
            "Number Of Drives": "2",
            "Number of Spans": "1",
            "PD": "0 Information",
            "Physical Sector Size": "512",
            "RAID Level": "Primary-1, Secondary-0, RAID Level Qualifier-0",
            "Size": "465.25 GB",
            "Span Depth": "1",
            "State": "Optimal",
            "Strip Size": "64 KB",
            "VD has Emulated PD": "No",
            "Virtual Drive": "1 (Target Id: 1)"
          }
        }
      ],
      "Disks": [
        {
          "ControllerID": "0",
          "ControllerDriver": "storcli7",
          "VolumeID": "",
          "Enclosure": "252",
          "Size": 499558383616,
          "UsedSize": 0,
          "SectorCount": 975699968,
          "PhysicalSectorSize": 512,
          "LogicalSectorSize": 512,
          "Slot": 0,
          "Protocol": "sata",
          "MediaType": "disk",
          "Status": "Online, Spun Up",
          "Info": {
            "Coerced Size": "465.25 GB [0x3a280000 Sectors]",
            "Commissioned Spare": "No",
            "Connected Port Number": "2(path0)",
            "Device Firmware Level": "1S04",
            "Device Id": "1",
            "Device Speed": "3.0Gb/s",
            "Drive": "Not Certified",
            "Drive has flagged a S.M.A.R.T alert": "No",
            "Drive is formatted for PI information": "No",
            "Drive's NCQ setting": "N/A",
            "Drive's postion": "DiskGroup: 0, Span: 0, Arm: 0",
            "Emergency Spare": "No",
            "Enclosure Device ID": "252",
            "Enclosure position": "0",
            "FDE Capable": "Not Capable",
            "FDE Enable": "Disable",
            "Firmware state": "Online, Spun Up",
            "Foreign State": "None",
            "Inquiry Data": "WD-WMAYP8041031WDC WD5003ABYX-18WERA0                  01.01S04",
            "Last Predictive Failure Event Seq Number": "0",
            "Link Speed": "3.0Gb/s",
            "Locked": "Unlocked",
            "Logical Sector Size": "512",
            "Media Error Count": "0",
            "Media Type": "Hard Disk Device",
            "Needs EKM Attention": "No",
            "Non Coerced Size": "465.261 GB [0x3a286030 Sectors]",
            "Other Error Count": "0",
            "PD Type": "SATA",
            "PI": "No PI",
            "PI Eligibility": "No",
            "Physical Sector Size": "512",
            "Port status": "Active",
            "Port's Linkspeed": "3.0Gb/s",
            "Predictive Failure Count": "0",
            "Raw Size": "465.761 GB [0x3a386030 Sectors]",
            "SAS Address(0)": "0x4433221103000000",
            "Secured": "Unsecured",
            "Sequence Number": "28",
            "Shield Counter": "0",
            "Slot Number": "0",
            "Successful diagnostics completion on": "N/A",
            "WWN": "50014ee0ae7b61c5"
          }
        },
        {
          "ControllerID": "0",
          "ControllerDriver": "storcli7",
          "VolumeID": "",
          "Enclosure": "252",
          "Size": 499558383616,
          "UsedSize": 0,
          "SectorCount": 975699968,
          "PhysicalSectorSize": 512,
          "LogicalSectorSize": 512,
          "Slot": 1,
          "Protocol": "sata",
          "MediaType": "disk",
          "Status": "Online, Spun Up",
          "Info": {
            "Coerced Size": "465.25 GB [0x3a280000 Sectors]",
            "Commissioned Spare": "No",
            "Connected Port Number": "1(path0)",
            "Device Firmware Level": "1S04",
            "Device Id": "2",
            "Device Speed": "3.0Gb/s",
            "Drive": "Not Certified",
            "Drive has flagged a S.M.A.R.T alert": "No",
            "Drive is formatted for PI information": "No",
            "Drive's NCQ setting": "N/A",
            "Drive's postion": "DiskGroup: 0, Span: 0, Arm: 1",
            "Emergency Spare": "No",
            "Enclosure Device ID": "252",
            "Enclosure position": "0",
            "FDE Capable": "Not Capable",
            "FDE Enable": "Disable",
            "Firmware state": "Online, Spun Up",
            "Foreign State": "None",
            "Inquiry Data": "WD-WMAYP7926193WDC WD5003ABYX-18WERA0                  01.01S04",
            "Last Predictive Failure Event Seq Number": "0",
            "Link Speed": "3.0Gb/s",
            "Locked": "Unlocked",
            "Logical Sector Size": "512",
            "Media Error Count": "0",
            "Media Type": "Hard Disk Device",
            "Needs EKM Attention": "No",
            "Non Coerced Size": "465.261 GB [0x3a286030 Sectors]",
            "Other Error Count": "0",
            "PD Type": "SATA",
            "PI": "No PI",
            "PI Eligibility": "No",
            "Physical Sector Size": "512",
            "Port status": "Active",
            "Port's Linkspeed": "3.0Gb/s",
            "Predictive Failure Count": "0",
            "Raw Size": "465.761 GB [0x3a386030 Sectors]",
            "SAS Address(0)": "0x4433221102000000",
            "Secured": "Unsecured",
            "Sequence Number": "28",
            "Shield Counter": "0",
            "Slot Number": "1",
            "Successful diagnostics completion on": "N/A",
            "WWN": "50014ee05925c428"
          }
        },
        {
          "ControllerID": "0",
          "ControllerDriver": "storcli7",
          "VolumeID": "",
          "Enclosure": "252",
          "Size": 499558383616,
          "UsedSize": 0,
          "SectorCount": 975699968,
          "PhysicalSectorSize": 512,
          "LogicalSectorSize": 512,
          "Slot": 2,
          "Protocol": "sata",
          "MediaType": "disk",
          "Status": "Online, Spun Up",
          "Info": {
            "Coerced Size": "465.25 GB [0x3a280000 Sectors]",
            "Commissioned Spare": "No",
            "Connected Port Number": "3(path0)",
            "Device Firmware Level": "1S04",
            "Device Id": "3",
            "Device Speed": "3.0Gb/s",
            "Drive": "Not Certified",
            "Drive has flagged a S.M.A.R.T alert": "No",
            "Drive is formatted for PI information": "No",
            "Drive's NCQ setting": "N/A",
            "Drive's postion": "DiskGroup: 1, Span: 0, Arm: 0",
            "Emergency Spare": "No",
            "Enclosure Device ID": "252",
            "Enclosure position": "0",
            "FDE Capable": "Not Capable",
            "FDE Enable": "Disable",
            "Firmware state": "Online, Spun Up",
            "Foreign State": "None",
            "Inquiry Data": "WD-WMAYP8040536WDC WD5003ABYX-18WERA0                  01.01S04",
            "Last Predictive Failure Event Seq Number": "0",
            "Link Speed": "3.0Gb/s",
            "Locked": "Unlocked",
            "Logical Sector Size": "512",
            "Media Error Count": "0",
            "Media Type": "Hard Disk Device",
            "Needs EKM Attention": "No",
            "Non Coerced Size": "465.261 GB [0x3a286030 Sectors]",
            "Other Error Count": "0",
            "PD Type": "SATA",
            "PI": "No PI",
            "PI Eligibility": "No",
            "Physical Sector Size": "512",
            "Port status": "Active",
            "Port's Linkspeed": "3.0Gb/s",
            "Predictive Failure Count": "0",
            "Raw Size": "465.761 GB [0x3a386030 Sectors]",
            "SAS Address(0)": "0x4433221101000000",
            "Secured": "Unsecured",
            "Sequence Number": "22",
            "Shield Counter": "0",
            "Slot Number": "2",
            "Successful diagnostics completion on": "N/A",
            "WWN": "50014ee0ae7b6632"
          }
        },
        {
          "ControllerID": "0",
          "ControllerDriver": "storcli7",
          "VolumeID": "",
          "Enclosure": "252",
          "Size": 499558383616,
          "UsedSize": 0,
          "SectorCount": 975699968,
          "PhysicalSectorSize": 512,
          "LogicalSectorSize": 512,
          "Slot": 3,
          "Protocol": "sata",
          "MediaType": "disk",
          "Status": "Online, Spun Up",
          "Info": {
            "Coerced Size": "465.25 GB [0x3a280000 Sectors]",
            "Commissioned Spare": "No",
            "Connected Port Number": "0(path0)",
            "Device Firmware Level": "1S04",
            "Device Id": "0",
            "Device Speed": "3.0Gb/s",
            "Drive": "Not Certified",
            "Drive has flagged a S.M.A.R.T alert": "No",
            "Drive is formatted for PI information": "No",
            "Drive's NCQ setting": "N/A",
            "Drive's postion": "DiskGroup: 1, Span: 0, Arm: 1",
            "Emergency Spare": "No",
            "Enclosure Device ID": "252",
            "Enclosure position": "0",
            "Exit Code": "0x00",
            "FDE Capable": "Not Capable",
            "FDE Enable": "Disable",
            "Firmware state": "Online, Spun Up",
            "Foreign State": "None",
            "Inquiry Data": "WD-WMAYP8042658WDC WD5003ABYX-18WERA0                  01.01S04",
            "Last Predictive Failure Event Seq Number": "0",
            "Link Speed": "3.0Gb/s",
            "Locked": "Unlocked",
            "Logical Sector Size": "512",
            "Media Error Count": "0",
            "Media Type": "Hard Disk Device",
            "Needs EKM Attention": "No",
            "Non Coerced Size": "465.261 GB [0x3a286030 Sectors]",
            "Other Error Count": "0",
            "PD Type": "SATA",
            "PI": "No PI",
            "PI Eligibility": "No",
            "Physical Sector Size": "512",
            "Port status": "Active",
            "Port's Linkspeed": "3.0Gb/s",
            "Predictive Failure Count": "0",
            "Raw Size": "465.761 GB [0x3a386030 Sectors]",
            "SAS Address(0)": "0x4433221100000000",
            "Secured": "Unsecured",
            "Sequence Number": "22",
            "Shield Counter": "0",
            "Slot Number": "3",
            "Successful diagnostics completion on": "N/A",
            "WWN": "50014ee003d097f2"
          }
        }
      ],
      "Info": {
        "Access Policy": "Yes",
        "Alarm": "Present",
        "Alarm Control": "Yes",
        "Alarm Disable": "Yes",
        "Allow Boot with Preserved Cache": "No",
        "Allow Ctrl Encryption": "No",
        "Allow HDD SAS/SATA Mix in VD": "Yes",
        "Allow HDD/SSD Mix in VD": "Yes",
        "Allow Mix in Enclosure": "Yes",
        "Allow Mixed Redundancy on Array": "No",
        "Allow SATA in Cluster": "No",
        "Allow SSD SAS/SATA Mix in VD": "Yes",
        "Allowed Device Type": "SAS/SATA Mix",
        "Any Offline VD Cache Preserved": "No",
        "Auto Detect BackPlane Enable": "SGPIO/i2c SEP",
        "Auto Detect BackPlane Enabled": "SGPIO/i2c SEP",
        "Auto Enhanced Import": "Yes",
        "Auto Rebuild": "Enabled",
        "BBU": "Absent",
        "BGI Rate": "30%",
        "BIOS Continue on Error": "Yes",
        "BIOS Enumerate VDs": "Yes",
        "BIOS Version": "3.30.02.2_4.16.08.00_0x06060A05",
        "BOOT Version": "09.250.01.219",
        "Background Rate": "30",
        "Battery FRU": "N/A",
        "Battery Warning": "Disabled",
        "Block SSD Write Disk Cache Change": "No",
        "Boot Block Version": "2.02.00.00-0000",
        "Boot Volume Supported": "NO",
        "BreakMirror RAID Support": "Yes",
        "CC Rate": "Yes",
        "Cache Flush Interval": "4s",
        "Cache When BBU Bad": "Disabled",
        "Cached IO": "No",
        "Check Consistency Rate": "30%",
        "ChipRevision": "B4",
        "Cluster Active": "No",
        "Cluster Disable": "Yes",
        "Cluster Mode": "Disabled",
        "Cluster Permitted": "No",
        "Cluster Support": "No",
        "Coercion Mode": "None",
        "Controller Id": "0000",
        "Critical Disks": "0",
        "Current Size of CacheCade": "0 GB",
        "Current Size of FW Cache": "346 MB",
        "Current Time": "21:24:30 2/14, 2018",
        "Dedicated Hot Spare": "Yes",
        "Default LD PowerSave Policy": "Controller Defined",
        "Default spin down time in minutes": "30",
        "Degraded": "0",
        "Delay Among Spinup Groups": "2s",
        "Delay during POST": "0",
        "Deny CC": "No",
        "Deny Clear": "No",
        "Deny Force Failed": "No",
        "Deny Force Good/Bad": "No",
        "Deny Locate": "No",
        "Deny Missing Replace": "No",
        "Deny SCSI Passthrough": "No",
        "Deny SMP Passthrough": "No",
        "Deny STP Passthrough": "No",
        "Device Id": "0079",
        "Device Interface": "PCIE",
        "Direct PD Mapping": "No",
        "Dirty LED Shows Drive Activity": "No",
        "Disable Copyback": "No",
        "Disable Ctrl-R": "Yes",
        "Disable Join Mirror": "No",
        "Disable Online Controller Reset": "No",
        "Disable Online PFK Change": "No",
        "Disable Puncturing": "No",
        "Disable Spin Down of hot spares": "No",
        "Disk Cache Policy": "Yes",
        "Disks": "4",
        "ECC Bucket Count": "0",
        "Ecc Bucket Leak Rate": "1440 Minutes",
        "Ecc Bucket Size": "15",
        "Enable Copyback on SMART": "No",
        "Enable Copyback to SSD on SMART Error": "Yes",
        "Enable JBOD": "No",
        "Enable LDBBM": "Yes",
        "Enable Led Header": "Yes",
        "Enable SSD Patrol Read": "No",
        "Enable Shield State": "No",
        "Enable Spin Down of UnConfigured Drives": "Yes",
        "Enable Web BIOS": "Yes",
        "EnableCrashDump": "No",
        "EnableLDBBM": "Yes",
        "Exit Code": "0x00",
        "Expose Enclosure Devices": "Enabled",
        "FW Package Build": "12.15.0-0239",
        "FW Version": "2.130.403-4660",
        "Failed Disks": "0",
        "Flash": "Present",
        "Flush Time": "4 seconds",
        "Force Offline": "Yes",
        "Force Online": "Yes",
        "Force Rebuild": "Yes",
        "Foreign Config Import": "Yes",
        "Global Hot Spares": "Yes",
        "Host Interface": "PCIE",
        "Host Request Reordering": "Enabled",
        "IO Policy": "Yes",
        "Interrupt Throttle Active Count": "16",
        "Interrupt Throttle Completion": "50us",
        "LED Show Drive Activity": "Yes",
        "Load Balance Mode": "Auto",
        "Maintain PD Fail History": "Enabled",
        "Max Arms Per VD": "32",
        "Max Arrays": "128",
        "Max Chained Enclosures": "16",
        "Max Configurable CacheCade Size": "0 GB",
        "Max Data Transfer Size": "8192 sectors",
        "Max Drives to Spinup at One Time": "4",
        "Max LD per array": "16",
        "Max Number of VDs": "64",
        "Max Parallel Commands": "1008",
        "Max SGE Count": "80",
        "Max Spans Per VD": "8",
        "Max Strip Size": "1.0 MB",
        "Max Strips PerIO": "42",
        "Maximum number of direct attached drives to spin up in 1 min": "120",
        "Memory": "Present",
        "Memory Correctable Errors": "0",
        "Memory Size": "512MB",
        "Memory Uncorrectable Errors": "0",
        "Mfg. Date": "03/13/12",
        "Min Strip Size": "8 KB",
        "NCQ": "No",
        "NVDATA Version": "2.09.03-0058",
        "NVRAM": "Present",
        "Number of Backend Port": "8",
        "Number of Frontend Port": "0",
        "Offline": "0",
        "On board Expander": "Absent",
        "PFK TrailTime Remaining": "0 days 0 hours",
        "PFK in NVRAM": "No",
        "POST delay": "90 seconds",
        "PR Correct Unconfigured Areas": "Yes",
        "PR Rate": "30%",
        "Patrol Read Rate": "Yes",
        "Phy Polarity": "0",
        "Phy PolaritySplit": "0",
        "Physical Devices": "5",
        "Physical Drive Coercion Mode": "Disabled",
        "Point In Time Progress": "No",
        "Power Saving option": "Don't Auto spin down Configured Drives",
        "Power Savings": "No",
        "PreBoot CLI Enabled": "Yes",
        "Preboot CLI Version": "04.04-020:#%00009",
        "Predictive Fail Poll Interval": "300sec",
        "Product Name": "LSI MegaRAID SAS 9260-8i",
        "RAID Level Supported": "RAID0, RAID1, RAID5, RAID6, RAID00, RAID10, RAID50, RAID60, PRL 11, PRL 11 with spanning, SRL 3 supported, PRL11-RLQ0 DDF layout with no span, PRL11-RLQ0 DDF layout with span",
        "Read Policy": "Yes",
        "Rebuild Rate": "30%",
        "Reconstruct Rate": "Yes",
        "Reconstruction": "Yes",
        "Reconstruction Rate": "30%",
        "Restore Hot Spare on Insertion": "No",
        "Restore HotSpare on Insertion": "Disabled",
        "Revertible Hot Spares": "Yes",
        "Revision No": "79B",
        "Rework Date": "00/00/00",
        "SAS Address": "500605b004872390",
        "SAS Disable": "No",
        "SMART Mode": "Mode 6",
        "Security Key Assigned": "No",
        "Security Key Failed": "No",
        "Security Key Not Backedup": "No",
        "Self Diagnostic": "Yes",
        "Serial Debugger": "Present",
        "Serial No": "SV21106771",
        "Snapshot Enabled": "No",
        "Spanning": "Yes",
        "Spin Down Mode": "None",
        "Spin Down time": "30",
        "Strip Size": "256kB",
        "SubDeviceId": "9261",
        "SubVendorId": "1000",
        "Support Boot Time PFK Change": "No",
        "Support Breakmirror": "No",
        "Support PFK": "Yes",
        "Support PI": "No",
        "Support Security": "No",
        "Support Shield State": "No",
        "Support Temperature": "Yes",
        "Support the OCE without adding drives": "Yes",
        "Supported Drives": "SAS, SATA",
        "T10 Power State": "No",
        "TPM": "Absent",
        "TTY Log In Flash": "No",
        "Temperature sensor for ROC": "Absent",
        "Temperature sensor for controller": "Absent",
        "Time taken to detect CME": "60s",
        "Topology Type": "None",
        "Treat Single span R1E as R10": "No",
        "Un-Certified Hard Disk Drives": "Allow",
        "Upgrade Key": "Absent",
        "Use FDE Only": "No",
        "Use disk activity for locate": "No",
        "Vendor Id": "1000",
        "Virtual Drives": "2",
        "WebBIOS Version": "6.0-54-e_50-Rel",
        "Write Policy": "Yes",
        "ZCR Config": "Unknown",
        "Zero Based Enclosure Enumeration": "No"
      }
    }
  ]
[root@df8-bc-12-41-c4-72 ~]# ./drp-raid -clear
[root@df8-bc-12-41-c4-72 ~]#
```