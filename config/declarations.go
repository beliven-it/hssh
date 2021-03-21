package config

// HSSHHostFolderName is the target folder
// to store the list of host files
const HSSHHostFolderName = "config.hssh.d"

// HomePath ....
var HomePath = getHomePath()

// SSHFolderPath ...
var SSHFolderPath = HomePath + "/.ssh"

// SSHConfigFilePath ...
var SSHConfigFilePath = SSHFolderPath + "/config"

// HSSHHostFolderPath ...
var HSSHHostFolderPath = SSHFolderPath + "/" + HSSHHostFolderName

// HSSHConfigFilePath ...
var HSSHConfigFilePath = HomePath + "/.config/hssh/config.yml"

// InitializedFilePath ...
var InitializedFilePath = HomePath + "/.config/hssh/init"
