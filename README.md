todo
====

minimal command line todo list

todos are stored at ~/.todo

install:

    $ go get github.com/MarcoLucidi01/todo

add todos:

    $ todo fix car blinker
     0 [ ] fix car blinker
    $ todo watch https://www.youtube.com/watch?v=dQw4w9WgXcQ
     0 [ ] fix car blinker
     1 [ ] watch https://www.youtube.com/watch?v=dQw4w9WgXcQ
    $ todo clean dog house
     0 [ ] fix car blinker
     1 [ ] watch https://www.youtube.com/watch?v=dQw4w9WgXcQ
     2 [ ] clean dog house
    $

mark todo as complete:

    $ todo -c 1
     0 [ ] fix car blinker
     2 [ ] clean dog house
    $

show incomplete todos (default when command is a "no print" one):

    $ todo
     0 [ ] fix car blinker
     2 [ ] clean dog house
    $

show also completed todos:

    $ todo -c
     0 [ ] fix car blinker
     1 [x] watch https://www.youtube.com/watch?v=dQw4w9WgXcQ
     2 [ ] clean dog house
    $

complete usage:

    $ todo -h
    usage: todo [-c|-i|-e|-s|-r|-h] [id...] [desc...]
    command line todo list
    todos are stored at /home/marco/.todo
    
      desc...          add new todo
      -c               print also completed todos
      -c id...         mark specified todos as complete
      -i id...         mark specified todos as incomplete
      -e id desc...    edit description of specified todo
      -e id /sub/rep/  replace substring sub with rep in description of specified todo
      -s id id         swap position of specified todos
      -r               remove completed todos
      -r id...         remove specified todos
      -h               show usage message
    
    repo: https://github.com/MarcoLucidi01/todo
    $

inspired by [t](https://github.com/sjl/t)
