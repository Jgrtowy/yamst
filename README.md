# 🛠️ YAMST - Yet Another Minecraft Server Tools

### So what is this?

YAMST is a collection of tools to help install, manage and host a Minecraft server. It is written in Go and is designed to be cross-platform.
It's pretty useless right now, but I'm working on it to make it as good as it can be.

### What can it do right now?

-   Install a Minecraft server (Vanilla and Paper currently)
-   Apply default configuration
-   Cache server jars
-   Download latest release and snapshot

### Planned features

-   Install and manage plugins
-   Install and manage worlds
-   Install and manage mods
-   Add more server types (Spigot, Forge, Fabric, etc.)
-   Simple GUI
-   Use proper Java version

## How do I use it?

### Installation

### Windows

1. Download the latest .exe from the releases page

2. Move the .exe to a folder of your choice

3. Add the folder to your PATH

4. Open a command prompt and run `yamst -h` to verify the installation and see the available commands

### Linux

1. Clone the repository

2. Run `go build` in the root directory

3. Move the executable to a folder of your choice

4. Add the folder to your PATH

5. Open a command prompt and run `yamst -h` to verify the installation and see the available commands

### Contributing

1. Fork the repository

2. Make your changes

3. Run `go build` in the root directory

4. Test your changes

5. Create a pull request

6. Wait for me to review it

### Usage

1. Create a folder for your server

2. Open a command prompt in the folder

3. Run `yamst -i` to install the latest version of vanilla server
