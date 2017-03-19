QStr is a library for the Go programming language that handles the strings you often encounter in Quake-based games.
They are usually used to communicate between other players on other servers via chat, but can also be found in player
nicknames or even server names. An identifying characteristic of these strings is their use of short color codes to
alter how they look, much like inline CSS styles do to text on the web. Here is what one looks like:

    ^x444Anti^5body

This is the text "Antibody" with two colors applied: ^x444 and ^5. The first color is a short hexadecimal
representation of the color #444444 (gray). The second, ^5, is a shorthand form for light blue. Each color applies to
all text that follows until another color code is found.

This library aims to make these types of strings easier to display on the web or in 2D graphics. It provides
facilities to strip the string of its color codes for a "stripped" version using the Stripped() method, like so:

    nick := qstr.QStr("^x444Anti^5body").Stripped() // "Antibody"

For display on the web the `HTML()` method can be used. This returns an html.Template-compatible object that wraps
the colorized elements in nested span elements. Using the same example:

    nick := qstr.QStr("^x444Anti^5body").HTML()
    // <span style="color:rgb(127,127,127)">Anti<span style='color:rgb(51,255,255)'>body</span></span>

For the most control and customization the `ColorParts` method can be used. This essentially breaks down the string into
its colorized pieces. Calling this method will give you a slice of the textual components along with their corresponding
RGB values. Here's that in action:

    nick := qstr.QStr("^x444Anti^5body")
    for _, part := range nick.ColorParts() {
        fmt.Printf("part %s has color %+v\n", part.Part, part.Color)
    }
    // part Anti has color {R:0.26666666666666666 G:0.26666666666666666 B:0.26666666666666666}
    // part body has color {R:0.2 G:0.4 B:1}

