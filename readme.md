# Arkantos, a text editor in the making.

## For whom it concerns: Not usable at the moment in a realistic way.

### TODO: Add syntax highlighting and better command support. Add plugin support. (Very busy with uni rn, maybe later in time I will get back to this)

<img src="Arkantos.png">

## Explanation of the project:

#### This is a learning project to learn GO, since I have never built a text editor some of choices made may be questionable.

### In order to use to use it:

-       Pass in buffers names, and it will open them all.
-       Don't pass anything and open a empty buffer.

#### Commands:

-       CTRL + s/S : saves file

#### Normal mode:

-       hjkl: Movement keys (like vim)
-       w/Q + ENTER: save file
-       w/W + q/Q + ENTER: save and quit
-       q/Q + ENTER: quits
-       i/I: enters insertion mode
-       TAB: Cycles buffers

#### Insertion mode:

-       ESC: enters normal mode (behaves like a normal text editor)

-       TODO: copy and paste

#### Sources:

<a>https://github.com/gen2brain/raylib-go</a>

