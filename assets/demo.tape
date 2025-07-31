# Golden Test Library Demo
# VHS documentation: https://github.com/charmbracelet/vhs

Output assets/demo.gif

Set Shell "bash"
Set FontSize 15
Set Width 1200
Set Height 700
Set Theme "Dracula"

Type "# Golden Test Library - Modern Go Testing 🚀"
Sleep 800ms
Enter
Enter

Type "# Let's run a test that will show a beautiful diff"
Sleep 800ms
Enter
Type "go test ./examples"
Enter
Sleep 3s

Type "# 🔍 See that beautiful colored diff output!"
Sleep 1s
Enter
Type "# Golden shows exactly what changed with colors"
Sleep 1s
Enter
Enter

Type "# Now let's update the golden file with environment variable"
Sleep 800ms
Enter
Type "GOLDEN_UPDATE=true go test ./examples"
Enter
Sleep 2s

Type "# ✨ Golden automatically updated the file!"
Sleep 800ms
Enter
Type "# One API, Any Data, Beautiful Diffs - That's Golden! 🌟"
Sleep 1s
