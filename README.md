# Carpool Creator!

**Carpool Creator** is a service meant to allow users to input a number of addresses and persons, and easily map out an optimal route (least distance total, not currently achieved with algorithm as we use something similar to an MST). I initially built it when I realized planning ride sheets for my student organization was a big mess, and could be optimized with the use of Google Distance Matrix API after I learned a bit with graph theory and algorithms. It's also been a bit of a way for me to experiment with server-side development in Go, which has been a blast to learn thus far; I also have a separate project for a client side application that actually structures a JSON to communicate with this service, in a separate repository. Feel free to check it out!

To use the application, properly install Go (https://golang.org/). Then, navigate to the "main" folder, and run "go build" from the command line, and a file main.exe should be built. You can directly run this and your server will be up; I'll add details as to how to use your own API key later :)

There is also an authentication module I was playing with, and plan to separate from the project entirely as a new project. It's a generic, multipurpose authentication tool that returns temporary session cookies to the user on successful login, and generates tables in a postgresql database for authentication.
