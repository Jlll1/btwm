# btwm

btwm is a simple, (kinda) tiling window manager. It's meant to be as small as possible, do as little as possible and be as easy to understand as possible.

## Using btwm
btwm opens each window full-screen and puts it on the top of your view-stack
Each window is assigned the first free tag, so the first application you open is tagged `1`, the next is tagged `2` and so forth.
You can focus applications by using the `FocusTagCommand` (by default Super+[Tag])

btwm also provides *Split Mode* for times when you want to have more than one application visible on the screen. Split mode displays two applications next to each other, each taking 50% of the width of the screen.
You can enter split mode by using the `SplitOnTagCommand` (by default Super+Shift+[Tag])

To browse other keybindings view [/configuration/keybindings.go](/configuration/keybindings.go)

## Configuring keybindings

The keybindings can be found in [/configuration/keybindings.go](/configuration/keybindings.go)
You can assign commands to specific key combinations using the `wm.BindKey()` function.
The list of available commands can be found in [/wm/commands.go](/wm/commands.go)

