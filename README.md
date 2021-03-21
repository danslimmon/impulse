# impulse

**Caveat: this README is aspirational. The code doesn't do anything useful yet.**

Impulse is a workflow, described
[here](https://blog.danslimmon.com/2021/03/15/managing-attention-with-stacks/), for managing your
intentions as a stack. Impulse can be used for lots of purposes:

* As a replacement for traditional to-do lists
* As a lightweight project planning tool
* To work effectively in the face of ADHD and anxiety.

`impulse` is a CLI tool that implements the Impulse workflow.

## The Impulse workflow

Suppose you set out to cook some pasta. You write `cook pasta` on a sheet of paper. You place that
sheet of paper on a desk. Cooking pasta is a multi-step process, and the first step is to put water
in a pot. So you take another sheet of paper and write on it, `put water in pot`. You place this
sheet on top of the `cook pasta` sheet. You now have a stack of papers that looks like this:

        put water in pot
    cook pasta

There is a `cook pasta` task on the stack. Above the `cook pasta` task is a `put water in pot`
task. So, working from the top of the stack down: you want to put water in the pot, and then return
to the cook pasta frame.

So, once you've put water in the pot, you remove that sheet and return to the next task down, `cook
pasta`. You're left with this:

    cook pasta

With Impulse, we always work from the top of the stack. Therefore, what you want to do next is
continue to `cook pasta`.

Are you done cooking pasta? Well no. So you need to put some more sheets onto the stack:

        place pot on burner
        turn on burner
        wait for water to boil
    cook pasta

There are now 4 sheets of paper on the stack. Again: you always work from the top of the stack. So
the next thing you have to do is `place pot on burner`. After doing this and removing the `place pot
on burner` sheet from the stack, you are now on to `turn on burner`. So you turn on the burner,
removing _that_ frame from the stack. Now your task is `wait for water to boil`.

        wait for water to boil
    cook pasta

There are now 2 sheets of paper on the stack. But water takes a while to boil. Maybe while you're
waiting for the water to boil, you decide to check Twitter. So you write `check Twitter` on a new
sheet of paper, and place it on top of the stack:

    
    check Twitter
        wait for water to boil
    cook pasta

There are 3 tasks on the stack. The top sheet, `check Twitter`, has
interrupted the `cook pasta … wait for water to boil` tasks. It is now the thing you're doing. If
you want to be really precise, you can add more sheets for the Twitter task:

        check Twitter notifications
        check Twitter timeline
    check Twitter
        wait for water to boil
    cook pasta

So now the thing you're doing is `check Twitter notifications`. When you finish that, you'll take it
off the stack and proceed to `check Twitter timeline`. And finally you'll remove _that_, and get
back to `check Twitter`. There are no more tasks involved in `check[ing] Twitter`, so nothing more
needs to be added. Instead, `check Twitter` itself is now removed, returning us to `wait for water
to boil`. If the water hasn't boiled yet, we may add another interrupt, such as `text Mom`. If the
water is boiling, then we remove the `wait for water to boil` task and return to `cook pasta`:

    cook pasta

which now necessitates adding new frames onto the stack (namely, `open pasta box` and `pour pasta in
pot`.)

        open pasta box
        pour pasta in pot
    cook pasta

You continue like this until you're done, at which point you move on to whatever task is beneath
`cook pasta` (for example, `eat dinner`).

## The `impulse` tool

**Reminder: this README is aspirational. Don't expect the code to do what I'm describing here yet.**

This repo is home to a CLI tool called `impulse` that implements the Impulse workflow described
above. You use it like so:

    --- Moving the Cursor

    j ↓     move cursor down
    k ↑     move cursor up
    h ←     move cursor to parent
    l →     move cursor to child
    t       move cursor to top

    --- Moving tasks

    J ⇧↓    move task down (among its siblings)
    K ⇧↑    move task up (among its siblings)
    H ⇧←    move task left (make it a child of the task that's currently its grandparent)
    L ⇧→    move task right (make it a child of the sibling directly above it)

    --- Changing tasks

    c       add child task(s)
    s       add sibling task(s)
    d       delete task
    Enter   edit task name

    --- Etc.

    ?       help (this message)
