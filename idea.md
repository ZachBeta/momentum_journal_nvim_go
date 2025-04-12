the approach so far feels more like chasing typescript issues

I'm feeling a drive to pull back in scope

```
gh repo create
? What would you like to do? Create a new repository on GitHub from scratch
? Repository name momentum_journal_nvim_go
? Repository owner ZachBeta
? Description creating a minimalist tool written in golang that allows an artist to follow "the artist's way" journaling, with agentic support, editing markdown files, and
```

but that approach might throw us, because as soon as I started brainstorming it, I started to realize that focusing on markdown, or vim, or a terminal, it's got a barrier to interaction

I want something like google docs, but without google docs, all local, an have an agentic flow

I've tried interacting with LLMs in a few different ways
pure chat bot in ollama run
cursor
chatgpt voice
text messages
atting discord bots

git repositories

I want to follow "the artist's way" and am struggling with pen on paper - I want to get there, but I need to reprogram my brain away from wanting to dopamine seek first thing when I wake, or first thing when I'm tired, or uncomfortable, or overwhelmed, etc etc

So the idea is to build this interface

single player
multiplayer

phone screen
laptop screen

sharing?
burning

privacy

reading it, not reading it

reviewing it?

we're talking about text transformation
idea transformation

how do we capture ideas, read ideas, write ideas, share ideas, collaborate on ideas

on the one hand I want to follow existing paradigms
on the other, I want something simple but fluid, as open as possible to capturing an idea, no matter how complex it is, or easy or difficult
just capture it
lower the barrier

and then how do I take that concept, and prototype something where I can open up my laptop anywhere, or maybe even my phone, and start journaling, capturing thoughts, and doing so in a way that's facilitated by an agent that I can slowly train to have a personality that helps bring out my creativity

the places I probably type the most text on any given day are cursor and discord, occasionally vim in a terminal

Google docs has been in the mix a lot as well

various other messaging apps

but it boils down to me writing in conversation, to track ideas to myself, to write code that can execute

My consumption habits?
lately it's reddit posts - a particular algorithm
and youtube - another particular algorithm

I have noticed that algorithms of consumption have consumed enough data themselves that now they are becoming algorithms of creation, of production. They are now capable of interacting with humans to create new things together, rather than coerce the humans to consume more content or products, or funnel all meaningful social interactions through systems that are powered by advertising incentives

chat gpt just rolled out a more expansive means of handling memories
cursor and windsurf have their own mechanisms of memory

arguably interacting with an LLM that writes and reads files is interacting with a concept of memory

I tried a neovim plugin that can interact with an agent it quickly became too large in scope and complexity to iterate quickly on
I tried to do that in this repo today with various tech stacks of handling a browser based interaction with an LLM

maybe the scope is too large

I need a page with two large full screen text boxes
One is freely editable content
One is more conversational - only add to the latest message
They need to be able to reference eachother

So we have a large text area that is freely edited on the left
And we have a chat window on the right

That is the core behavior of cursor. Everything after that is convenience and QOL for software development, but at the core, we're editing a large file, passing that to the conversation context, and then being able to have the conversation context modify it

The ux matters.
A free text area feels so very different from a series of messages
Wanting to edit it freely
to move around
to format

a file editor is basically free edit
a terminal is basically a conversation

If the agent makes an edit to a file token by token, I can't view that in vim. I can look at that if there's a google docs type system
A raw text field would get messy quickly
it's almost like I want to have the text area be tokenized underneath, then that tree/array/whatever data is updated

So a few simple plugins with vim, tmux, ollama might work after all
The web based version will quickly experience scope creep

In fact all the documents here are barely related to this new concept
In fact the neovim plugin app is arguably closer in terms of data flows
Streaming might get weird

I think we might be able to set up some kind of special mode of neovim to have a buffer that can stream content from the LLM conversation window. In fact a simple wrapper around the ux of `ollama run gemma3:12b` is probably enough. Mostly I want vim commands to be able to sample content from the conversation buffer into the compose buffer

so we open a markdown file in nvim left pane buffer
we have a pane open on the right that is a streaming LLm conversation that can contextualize the compose buffer - quite literally by adding the entire document, or some assessed context window of the compose buffer, and the conversation buffer. The conversation will certainly have a rolling context window as the first version, future versions can handle how to pull appropriately from conversation history
the compose windowpane can pass the entire document, or if the token count grows too large, the tail of it
And at that, how do we then take those two buffers, and combine them into a single context. the first version can fill the first half of context from the tail of compose, and the second half of context from the tail of conversation

I worry that making nvim plugins for this, we might end up leaning too hard on lua, which feels like it might not be the best language and tech stack for this?
but it also would make it easier to exist as a pure nvim plugin

this is probably plenty to use in a fresh repo
