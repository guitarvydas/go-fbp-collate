This is a quick implementation of Paul Morrison's Collate FBP example (page 91 of "Flow-Based Programming, 2nd Edition").

This example makes every func, except for main and print, into an FBP component, sitting in its own package.  Each component accepts parameters that consist of the input and output ports for the component.  Main "wires up" the components according to the diagram on page 91.

Each FBP component resides in its own sub-directory.

The code for Component.java was used as a guide / spec.  

The point of this exercise was not to build a robust Collate, but as a proof of concept - how this would look in Go.


To use this project, you also need the github components:

ip
collate
readfile
printer

