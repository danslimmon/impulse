# impulse

`impulse` is a thought tool that manages your intentions as a **stack**. I come to the term
**stack** from computer science, but you don't need to worry about that. You can just think of a
literal stack. A stack of papers, for instance. (And if you're thinking of the computer science
thing, don't expect a totally direct analog here.)

Suppose you set out to cook some pasta. You write `cook pasta` on a sheet of paper. You place that
sheet of paper on a desk. Cooking pasta is a multi-step process, and the first step is to put water
in a pot. So you take another sheet of paper and write on it, `put water in pot`. You place this
sheet on top of the `cook pasta` sheet. You now have a stack of papers that looks like this:

         put water in pot
       / 
    川 cook pasta

There is a `cook pasta` frame on the stack. Above the `cook pasta` frame is a `put water in pot`
rame. So, working from the top of the stack down: you want to put water in the pot, and then return
to the cook pasta frame.

So, once you've put water in the pot, you come back to the next frame down, `cook pasta`. This is
called **popping** the frame `put water in pot`. You're left with this:

    川 cook pasta
       \
        put water in pot

`cook pasta` is now at the base of the stack (as indicated by the `川`). `put water in pot` has
already been done, so it appears _below_ `cook pasta`. It is of historical interest only. What
matter is what we want to do next. And, as always, we work from the top of the stack. What we want
to do next, then, is continue to `cook pasta`.

Are you done cooking pasta? Well no. So you need to `push` some more frames onto the stack:

        
        place pot on burner
        | 
        turn on burner
        |
        wait for water to boil
       /
    川 cook pasta
       \
        put water in pot

There are now 4 frames (sheets of paper) on the stack. You always work from the top of the stack, so
the next thing you have to do is `place pot on burner`. After doing this and **popping** `place pot
on burner` off the stack, you are now on to `turn on burner`. So you turn on the burner, popping
that frame off the stack. Now your task is `wait for water to boil`.

        wait for water to boil
       /
    川 cook pasta
       \
        put water in pot
        |
        place pot on burner
        |
        turn on burner

(Remember, anything below `川` is history, not part of the stack. The stack starts at `川` and goes
up.)

There are now 2 frames (sheets of paper) on the stack. But water takes a while to boil. Maybe while
you're waiting for the water to boil, you decide to check your Twitter. So you write `check Twitter`
on a new sheet of paper, and place it on top of the stack:

    
    口 check Twitter
        wait for water to boil
       /
    川 cook pasta
       \
        put water in pot
        |
        place pot on burner
        |
        turn on burner

There are 3 frames (sheets of paper) on the stack. The most recent frame, `check Twitter`, has
**interrupted** the `cook pasta … wait for water to boil` frame. It is now the thing you're doing.
Since it's not the _main_ thing you're doing, it's called an **interrupt** and it's indicated by
`口`.  It can **push** its own frames onto the stack:

        check Twitter notifications
        |
        check Twitter timeline
       /
    口 check Twitter

        wait for water to boil
       /
    川 cook pasta
       \
        put water in pot
        |
        place pot on burner
        |
        turn on burner

So now the thing you're doing is `check Twitter notifications`. When you finish that, you'll **pop**
it off the stack and proceed to `check Twitter timeline`. And finally you'll **pop** _that_, and get
back to `check Twitter`. There are no more tasks involved in `check[ing] Twitter`, so nothing more
needs to be pushed. Instead, `check Twitter` itself is now **popped**, returning us to `wait for
water to boil`. If the water hasn't boiled yet, we may add another interrupt, such as `text Mom`. If
the water is boiling, then we **pop** `wait for water to boil` and return to `cook pasta`,

    川 cook pasta
       \
        put water in pot
        |
        place pot on burner
        |
        turn on burner
        |
        wait for water to boil

which now necessitates **pushing** new frames onto the stack (namely, `pour pasta in pot` and `open
pasta box`, pushed in that order so that `open pasta box` ends up at the top of the stack):

        open pasta box
        |
        pour pasta in pot
       /
    cook pasta
    川
       \
        put water in pot
        |
        place pot on burner
        |
        turn on burner
        |
        wait for water to boil

## aspirational

all this is subject to the principle that frames are immutable once they have been popped.

### semantics of 川

- when you pop the frame that is at `川`, it gets saved to a database with its descendants and their
    accompanying docs if any
- the database of frames is searchable
- you can look through a list of frames that have recently occupied `川`
- you can take a frame from the history (or search) and move it to `川`. any changes you
    subsequently make will apply to the frame you grabbed, and when it comes back off of `川` it
    will be archived with those changes.
- you can take a frame from the history (or search) and _duplicate_ it to `川`. the original frame
    will stay archived, and this new frame will be saved as its own thing.

### workflow

- there is a cursor (indicated by unique text color) that indicates what frame is **selected**.
- the cursor can be moved with keystrokes:

    j   down
    k   up
    h   go to parent
    l   go to top

- frames can be manipulated with keystrokes:

    p   pop       top frame moves below its parent, new top frame becomes whatever was underneath
                  the popped frame.
    u   push      new frame is created as a child of the current frame
    i   insert    new frame is created as a sibling of the current frame. it starts directly above
                  the parent frame, but after entering its name, you can move it with j and k
                  among its siblings, then hit enter to commit the position and return to the top
                  frame.
    I             same as `i`, but the new frame starts directly below the top frame.
    !   interrupt new frame is created as stack interrupt. it preempts `川`. when the interrupt
                  frame is later popped, control returns to `川`.

### frame context

- a given frame may have a "context", which is a markdown blob associated with that frame.
- when a frame is archived, its context gets saved with it
- with a frame selected, you can hit a keystroke to open its context in your text editor. the h1 of
    the document is the name of the frame. when you save it and exit, everything below the header
    gets saved as the frame context. the header gets discarded.
- frame contexts are part of the searchable material in the database.

### frame templates

- the user can take a frame and convert it to a **frame template**.
- at any time, a frame template can be pushed onto the stack. this means we take the template and
    its descendants and throw them all onto the stack as a new frame. (one thing i'm thinking about
    for this is zero hour)
- maybe by convertion templates have all-caps names? so if you're pushing and being prompted for a
    name, and you enter the all caps name of a template, then that template gets pushed?

## mvp

everything in memory, not saved to disk. one only manipulates the top frame. and its siblings.

don't worry about what happens to frames after they're popped. they go away forever. then we can
have a simple terminal interface that doesn't use ncurses. commands are received with
bufio.Reader(stdin).ReadByte().

show stack
p
u
i (inserts directly above parent)
!

    $ impulse
    川: [cook pasta\n]

    川 cook pasta

    [u]: [boil water\n]

       / boil water
    川 cook pasta

    [u]: [fill pot with water\n]

           / fill pot with water
       / boil water
    川 cook pasta
    
    [i]: [place pot on burner\n]

           / fill pot with water
           / place pot on burner
       / boil water
    川 cook pasta

    [i]: [turn on burner\n]

           / fill pot with water
           / place pot on burner
           / turn on burner
       / boil water
    川 cook pasta

    [p]

           / place pot on burner
           / turn on burner
       / boil water
    川 cook pasta

    [p]

           / turn on burner
       / boil water
    川 cook pasta

    [p]

       / boil water
    川 cook pasta

    [u]: [wait for water to boil\n]

           / wait for water to boil
       / boil water
    川 cook pasta

    [!]: [check twitter\n]

    口 check twitter
           / wait for water to boil
       / boil water
    川 cook pasta

    [p]

           / wait for water to boil
       / boil water
    川 cook pasta
